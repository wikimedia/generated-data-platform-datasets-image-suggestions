/*
 * Copyright 2021 Nikki Nikkhoui <nnikkhoui@wikimedia.org> and Wikimedia Foundation
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
	"encoding/json"
	"net/http"
	"runtime"
)

// Healthz represents the JSON object sent in the body of a `/healthz` response.
type Healthz struct {
	Version   string `json:"version"`
	BuildDate string `json:"build_date"`
	BuildHost string `json:"build_host"`
	GoVersion string `json:"go_version"`
}

// NewHealthz initializes and returns a new Healthz.
func NewHealthz(version, date, host string) *Healthz {
	return &Healthz{
		Version:   version,
		BuildDate: date,
		BuildHost: host,
		GoVersion: runtime.Version(),
	}
}

// HealthzHandler is an http.Handler that implements a readiness endpoint.
type HealthzHandler struct {
	healthz *Healthz
}

// Serves Healthz struct data as JSON
func (s *HealthzHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var response []byte
	var err error

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if response, err = json.MarshalIndent(s.healthz, "", "  "); err != nil {
		w.Write([]byte(`{}`))
		return
	}
	w.Write(response)
}
