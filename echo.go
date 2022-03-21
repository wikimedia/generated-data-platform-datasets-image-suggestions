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
	"io/ioutil"
	"net/http"
	"time"

	"github.com/eevans/servicelib-golang/logger"
	"schneider.vip/problem"
)

// Echo represents the JSON object sent as the body of the POST, as well as the
// the one that is echoed back in the response.
type Echo struct {
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

// EchoHandler is an http.Handler that implements an echo service endpoint.
type EchoHandler struct {
	Logger *logger.Logger
}

func (s *EchoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var echoBody []byte
	var echoResponse = Echo{}
	var err error
	var response []byte

	if r.Method != "POST" {
		problem.New(
			problem.Status(http.StatusMethodNotAllowed),
			problem.Detail("Method Not Allowed"),
		).WriteTo(w)
		return
	}

	if echoBody, err = ioutil.ReadAll(r.Body); err != nil {
		problem.New(
			problem.Status(http.StatusBadRequest),
			problem.Detail("Unable to read response body"),
		).WriteTo(w)
		return
	}

	if len(echoBody) < 1 {
		problem.New(
			problem.Status(http.StatusBadRequest),
			problem.Detail("Body must include `message`"),
		).WriteTo(w)
		return
	}

	if err = json.Unmarshal(echoBody, &echoResponse); err != nil {
		problem.New(
			problem.Status(http.StatusBadRequest),
			problem.Detail("Invalid JSON"),
		).WriteTo(w)
		return
	}

	echoResponse.Timestamp = time.Now().Format(time.RFC3339)

	if response, err = json.MarshalIndent(echoResponse, "", "  "); err != nil {
		s.Logger.Request(r).Log(logger.ERROR, "Unable to echo: %s", err)
		problem.New(
			problem.Status(http.StatusInternalServerError),
			problem.Detail("Unable to echo"),
		).WriteTo(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
