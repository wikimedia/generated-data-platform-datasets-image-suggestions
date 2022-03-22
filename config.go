/*
 * Copyright 2021 Nikki Nikkhoui <nnikkhoui@wikimedia.org>, Eric Evans <eevans@wikimedia.org>,
 * and Wikimedia Foundation
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
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/gocql/gocql"
	yaml "gopkg.in/yaml.v2"
)

// Config represents an application-wide configuration.
type Config struct {
	ServiceName string    `yaml:"service_name"`
	BaseURI     string    `yaml:"base_uri"`
	Address     string    `yaml:"listen_address"`
	Port        int       `yaml:"listen_port"`
	LogLevel    string    `yaml:"log_level"`
	Cassandra   cassandra `yaml:"cassandra"`
}

type cassandra struct {
	Port        int      `yaml:"port"`
	Consistency string   `yaml:"consistency"`
	Hosts       []string `yaml:"hosts"`
}

// NewConfig returns a new Config from YAML serialized as bytes.
func NewConfig(data []byte) (*Config, error) {
	// Populate a new Config with sane defaults
	config := Config{
		ServiceName: "service-scaffold-golang",
		BaseURI:     "",
		Address:     "localhost",
		Port:        8080,
		LogLevel:    "info",

		Cassandra: cassandra{
			Port:        9042,
			Consistency: "quorum",
			Hosts:       []string{"localhost"},
		},
	}
	err := yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return validate(&config)
}

// Returns a new Config from a YAML file.
func ReadConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return NewConfig(data)
}

// validateLogLevel ensures a valid log level
func validateLogLevel(config *Config) error {
	switch strings.ToUpper(config.LogLevel) {
	case "DEBUG", "INFO", "WARNING", "ERROR", "FATAL":
		return nil
	}
	return fmt.Errorf("Unsupported log level: %s", config.LogLevel)
}

func validateCassandraConsistency(c cassandra) error {
	switch strings.ToLower(c.Consistency) {
	case "any", "one", "two", "three", "quorum", "all", "localquorum", "eachquorum", "localone":
		return nil
	}
	return fmt.Errorf("Unsupported consistency level: %s", c.Consistency)
}

func validate(config *Config) (*Config, error) {
	var err error

	if !strings.HasPrefix(config.BaseURI, "/") {
		config.BaseURI = "/" + config.BaseURI
	}

	// Validate log level
	if err = validateLogLevel(config); err != nil {
		return nil, err
	}

	// Validate consistency level
	if _, err = goCQLConsistency(config.Cassandra.Consistency); err != nil {
		return nil, err
	}

	return config, nil
}

func goCQLConsistency(c string) (gocql.Consistency, error) {
	switch strings.ToLower(c) {
	case "any":
		return gocql.Any, nil
	case "one":
		return gocql.One, nil
	case "two":
		return gocql.Two, nil
	case "three":
		return gocql.Three, nil
	case "quorum":
		return gocql.Quorum, nil
	case "all":
		return gocql.All, nil
	case "localquorum":
		return gocql.LocalQuorum, nil
	case "eachquorum":
		return gocql.EachQuorum, nil
	case "localone":
		return gocql.LocalOne, nil
	default:
		return 0, fmt.Errorf("Unrecognized Cassandra consistency level")
	}
}