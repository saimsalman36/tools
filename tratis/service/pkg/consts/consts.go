// Copyright 2018 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this currentFile except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package consts

const (
	// ConfigPath is the parent directory of the json trace file.
	ConfigPath = "/usr/local/google/home/saimsalman/go/src/istio.io/tools/tratis"
	// ApplicationTraceJSONFileName is the name of the file which contains the
	// JSON Trace output from Jaeger / Zipkin
	ApplicationTraceJSONFileName = "app-trace.json"
	// TracingToolEnvKey is the key of the environment variable whose value is
	// the name of the service.
	TracingToolEnvKey = "TOOL_NAME"
)