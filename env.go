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
	"os"

	log "gerrit.wikimedia.org/r/mediawiki/services/servicelib-golang/logger"
)

// MergeEnvironment takes override values supplied in the environment, and merges them with a Config struct.
func MergeEnvironment(config *Config, logger *log.Logger) {

	// Cassandra credentials
	cassandraUsername, userOk := os.LookupEnv("CASSANDRA_USERNAME")
	cassandraPassword, passOk := os.LookupEnv("CASSANDRA_PASSWORD")

	if userOk && !passOk {
		logger.Warning("CASSANDRA_USERNAME env var provided but CASSANDRA_PASSWORD unset, using values from configuration file.")
	}

	if passOk && !userOk {
		logger.Warning("CASSANDRA_PASSWORD env var provided but CASSANDRA_USERNAME unset, using values from configuration file.")
	}

	if userOk && passOk {
		config.Cassandra.Authentication.Username = cassandraUsername
		config.Cassandra.Authentication.Password = cassandraPassword
		logger.Info("Cassandra login credentials configured from environment")
	}
}
