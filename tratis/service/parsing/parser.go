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

package parser

import (
	"encoding/json"
	"log"

	jaeger "github.com/jaegertracing/jaeger/model/json"

	comm "istio.io/tools/tratis/service/communication"
	"istio.io/tools/tratis/service/pkg/consts"
)

type TraceData struct {
	Traces []jaeger.Trace `json:"data"`
	Total  int            `json:"total"`
	Limit  int            `json:"limit"`
	Offset int            `json:"offset"`
}

func ParseJSON(toolName string, startTime string, duration string) (appTrace TraceData,
	err error) {

	ret := TraceData{}

	if toolName == "jaeger" {
		for _, EntryService := range consts.TracingToolEntryPoints {
			temp, err := ParseJaeger(comm.ExtractTraces(consts.TracingToolAddress,
				consts.TracingToolPortNumber, EntryService,
				consts.NumTraces, startTime, duration))

			if err != nil {
				return temp, err
			}

			ret.Traces = append(ret.Traces, temp.Traces...)
			ret.Total += temp.Total
			ret.Limit += temp.Limit
			ret.Offset += temp.Offset
		}

		return ret, nil
	}

	log.Fatalf(`tracing tool "%s" is not correctly supported`, toolName)
	return
}

func ParseJaeger(data []byte) (appTrace TraceData, err error) {
	t := TraceData{}
	err = json.Unmarshal(data, &t)
	if err != nil {
		return
	}

	return t, nil
}
