/*
 * Copyright 2022 Eric Evans <eevans@wikimedia.org> and Wikimedia Foundation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"

	log "gerrit.wikimedia.org/r/mediawiki/services/servicelib-golang/logger"
	"github.com/gocql/gocql"
	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	http_gateway "gitlab.wikimedia.org/repos/generated-data-platform/cassandra-http-gateway"
)

var (
	// These values are assigned at build using `-ldflags` (see: Makefile)
	buildDate = "unknown"
	buildHost = "unknown"
	version   = "unknown"
)

var (
	reqCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Count of HTTP requests processed, partitioned by status code and HTTP method.",
		},
		[]string{"code", "method"},
	)

	durationHisto = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "A histogram of latencies for requests, partitioned by status code and HTTP method.",
			Buckets: []float64{.001, .0025, .0050, .01, .025, .050, .10, .25, .50, 1},
		},
		[]string{"code", "method"},
	)
	promBuildInfoGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name:        "image_suggestions_build_info",
			Help:        "Build information",
			ConstLabels: map[string]string{"version": version, "build_date": buildDate, "build_host": buildHost, "go_version": runtime.Version()},
		})
)

func init() {
	prometheus.MustRegister(reqCounter, durationHisto, promBuildInfoGauge)
	promBuildInfoGauge.Set(1)
}

// Entrypoint for our service
func main() {
	var confFile = flag.String("config", "./config.yaml", "Path to the configuration file")

	var config *Config
	var err error
	var logger *log.Logger

	flag.Parse()

	if config, err = ReadConfig(*confFile); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if logger, err = log.NewLogger(os.Stdout, config.ServiceName, config.LogLevel); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to initialize the logger: %s", err)
		os.Exit(1)
	}

	// Allow overriding config using environment variables
	MergeEnvironment(config, logger)

	logger.Info(
		"Initializing service %s (Version: %s, Go: %s, Build host: %s, Timestamp: %s",
		config.ServiceName,
		version,
		runtime.Version(),
		buildHost,
		buildDate,
	)

	logger.Info("Connecting to Cassandra database: %s (port %d)", strings.Join(config.Cassandra.Hosts, ","), config.Cassandra.Port)
	logger.Debug("Cassandra: configured for consistency level '%s'", strings.ToLower(config.Cassandra.Consistency))
	logger.Debug("Cassandra: configured for local datacenter '%s'", config.Cassandra.LocalDC)

	var session *gocql.Session

	if session, err = newCassandraSession(config); err != nil {
		logger.Error("Failed to create Cassandra session: %s", err)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	router := httprouter.New()
	builder := http_gateway.SelectBuilder.Logger(logger).Session(session).CounterVec(reqCounter).HistogramVec(durationHisto)

	router.GET("/public/image_suggestions/suggestions/:wiki/:page_id", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		builder.
			From("image_suggestions", "suggestions").
			Bind(ps).
			Build().
			Handle(w, r)
	})

	router.GET("/private/image_suggestions/feedback/:wiki/:page_id", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		builder.
			From("image_suggestions", "feedback").
			Bind(ps).
			Build().
			Handle(w, r)
	})

	router.GET("/private/image_suggestions/title_cache/:wiki/:title", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		builder.
			From("image_suggestions", "title_cache").
			Bind(ps).
			Build().
			Handle(w, r)
	})

	router.GET("/private/image_suggestions/instanceof_cache/:wiki/:page_id", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		builder.
			From("image_suggestions", "instanceof_cache").
			Bind(ps).
			Build().
			Handle(w, r)
	})

	router.Handler("GET", "/healthz", &HealthzHandler{NewHealthz(version, buildDate, buildHost)})
	router.Handler("GET", "/metrics", promhttp.Handler())
	router.HandlerFunc("GET", "/openapi", openAPIHandlerFunc)

	http.ListenAndServe(fmt.Sprintf("%s:%d", config.Address, config.Port), router)
}
