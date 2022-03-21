/*
 * Copyright 2021 Nikki Nikkhoui <nnikkhoui@wikimedia.org>, Frances Goodwin <fgoodwin@wikimedia.org>
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
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type httperr struct {
	Status int
	Detail string
}

func TestInvalidMethod(t *testing.T) {

	req, err := http.NewRequest("GET", "/v0/echo", nil)

	require.NoError(t, err)

	rr := httptest.NewRecorder()

	handler := EchoHandler{}
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code, "Incorrect status code thrown")

	res := rr.Result().Body

	body, _ := ioutil.ReadAll(res)

	e := httperr{}

	err = json.Unmarshal(body, &e)

	require.NoError(t, err)

	assert.Equal(t, e.Status, http.StatusMethodNotAllowed, "Incorrect status code delivered to client")
	assert.Equal(t, e.Detail, "Method Not Allowed", "Incorrect detail delivered to client")

}

func TestEmptyBody(t *testing.T) {

	req, err := http.NewRequest("POST", "/v0/echo", strings.NewReader(""))

	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := EchoHandler{}

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Incorrect status code thrown")

	res := rr.Result().Body

	body, _ := ioutil.ReadAll(res)

	e := httperr{}

	err = json.Unmarshal(body, &e)

	require.NoError(t, err)

	assert.Equal(t, e.Status, http.StatusBadRequest, "Incorrect status code delivered to client")
	assert.Equal(t, e.Detail, "Body must include `message`", "Incorrect detail delivered to client")
}

func TestInvalidJSON(t *testing.T) {

	req, err := http.NewRequest("POST", "/v0/echo", strings.NewReader("bad req body"))

	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := EchoHandler{}

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Incorrect status code returned")

	res := rr.Result().Body

	body, _ := ioutil.ReadAll(res)

	e := httperr{}

	err = json.Unmarshal(body, &e)

	require.NoError(t, err)

	assert.Equal(t, e.Status, http.StatusBadRequest, "Incorrect status code delivered to client")
	assert.Equal(t, e.Detail, "Invalid JSON", "Incorrect detail delivered to client")
}

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

func TestFailedBody(t *testing.T) {

	req, err := http.NewRequest("POST", "/v0/echo", errReader(0))

	require.NoError(t, err)

	rr := httptest.NewRecorder()

	handler := EchoHandler{}

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Incorrect status code returned")

	res := rr.Result().Body

	body, _ := ioutil.ReadAll(res)

	e := httperr{}

	err = json.Unmarshal(body, &e)

	require.NoError(t, err)

	assert.Equal(t, e.Status, http.StatusBadRequest, "Incorrect status code delivered to client")
	assert.Equal(t, e.Detail, "Unable to read response body", "Incorrect detail delivered to client")
}

func TestEcho(t *testing.T) {
	msg, _ := json.Marshal(Echo{Message: "echo"})

	req, err := http.NewRequest("POST", "/v0/echo", bytes.NewBuffer(msg))

	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := EchoHandler{}

	handler.ServeHTTP(rr, req)

	res := rr.Result().Body

	body, _ := ioutil.ReadAll(res)

	e := Echo{}

	err = json.Unmarshal(body, &e)

	require.NoError(t, err)

	timestamp, errc := time.Parse(time.RFC3339, e.Timestamp)

	assert.NoError(t, errc, "Response body does not contain timestamp")
	assert.NotEmpty(t, timestamp, "No timestamp found")
	assert.Equal(t, e.Message, "echo", "Echoed message did not match sent message")
}
