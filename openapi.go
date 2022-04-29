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

import "net/http"

const openapi = `
openapi: 3.0.2

info:
  title: image-suggestions-dataset
  description: |
      [HTTP data gateway service](https://www.mediawiki.org/wiki/Platform_Engineering_Team/Data_Value_Stream/Data_Gateway#Image_Suggestions)
      for Image Suggestions.
  version: "1.0.0"

paths:
  /healthz:
    get:
      description: Establishes service readiness
      responses:
        200:
          description: Success
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Healthz'
      # x-amples is a sequence of request/response pairs which can be issued to
      # test service availability (for example by using
      # https://gerrit.wikimedia.org/r/admin/projects/operations/software/service-checker
      # to generate Icinga alerts).
      x-amples:
        - title: Retrieve readiness data
          request:
            headers:
              Accept: application/json
          response:
            status: 200
            headers:
              Content-Type: application/json
      # Enable/disable service monitoring based on x-amples.
      x-monitor: true

components:
  schemas:
    Healthz:
      type: object
      properties:
        version:
          type: string
        build_date:
          type: string
        build_host:
          type: string
        go_version:
          type: string
`

func openAPIHandlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-yaml")
	w.Write([]byte(openapi))
}
