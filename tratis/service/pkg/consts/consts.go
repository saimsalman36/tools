// Copyright 2019 Istio Authors
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
	// Traces Limit
	NumTraces = 20000

	MinNumTraces = 1

	DateFiltering = true
	StartTime     = "2019-08-05T21:55:51.534891301Z"
	EndTime       = "2019-08-05T21:59:51.534891301Z"

	// Distribution Fitting File Path
	DistFilePath = "Distribution"
	// Distribution Fitting Function Name
	DistFittingFuncName = "BestFitDistribution"

	// Tracing Tool IP Address
	TracingToolAddress = "localhost"
	// Tarcing Tool Port Number
	TracingToolPortNumber = "15032"
	// Output can be either `RAW` or `DIST`
	OutputType = "RAW"
)

var (
	TracingToolEntryPoints = []string{"productpage.service-graph"}
)
