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

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	jaeger "github.com/jaegertracing/jaeger/model/json"
	"istio.io/tools/tratis/service/output"
	parser "istio.io/tools/tratis/service/parsing"
)

func unique(data []jaeger.Trace) []jaeger.Trace {
	ret := make([]jaeger.Trace, 0)
	set := make(map[jaeger.TraceID]bool)

	for _, trace := range data {
		_, ok := set[trace.TraceID]

		if !ok {
			set[trace.TraceID] = true
			ret = append(ret, trace)
		}
	}

	return ret
}

func main() {
	fmt.Println("Starting Tratis ...")

	start := flag.String("start", "",
		"Start Time (RFC3339 Format)\nExample: -start 2019-08-05T21:55:51.534891301Z")
	duration := flag.String("duration", "60s",
		"Duration\nExample: -duration 3h")
	tool := flag.String("tool", "jaeger",
		"Tool Name\nExample: -tool jaeger/zipkin")
	tracesFile := flag.String("traces", "Traces.json",
		"File to Record Raw Data (Traces)\nExample: -traces Traces.json")
	resultFile := flag.String("results", "Results.json",
		"File to Record Results\nExample: -results Output.json")

	flag.Parse()

	if len(os.Args) < 5 {
		log.Fatalf(`Input Arguments not correctly provided: go run main.go 
			-start=<START_TIME> 
			-duration=<DURATION (minutes)>
			-tool=<TOOL_NAME> 
			-trace=<TRACE_FILE_NAME.json> 
			-results=<RESULTS_FILE_NAME.json>`)
	}

	TracingToolName := *tool
	traceFileName := *tracesFile
	jsonFileName := *resultFile

	fmt.Println("Generating Traces ...")

	data, err := parser.ParseJSON(TracingToolName, *start, *duration)
	if err != nil {
		log.Fatalf(`Connection between "%s" and tratis is broken`,
			TracingToolName)
	}

	fmt.Println(len(data.Traces), " traces generated")

	fmt.Println("Removing Duplicate Traces ...")

	data.Traces = unique(data.Traces)

	fmt.Println(len(data.Traces), " unique traces generated")

	fmt.Println("Writing Traces to File ...")

	f, err := os.Create(traceFileName)
	if err != nil {
		log.Fatalf("Unable to create %s: %v", traceFileName, err)
	}
	bytes, _ := json.MarshalIndent(data, "", "  ")
	n, err := f.Write(bytes)
	if err != nil {
		log.Fatalf("Unable to write json to %s: %v", jsonFileName, err)
	}

	_, _ = fmt.Fprintf(os.Stderr,
		"Successfully wrote %d bytes of Json data to %s\n", n, traceFileName)

	results := output.GenerateOutput(data)

	f, err = os.Create(jsonFileName)
	if err != nil {
		log.Fatalf("Unable to create %s: %v", jsonFileName, err)
	}
	n, err = f.Write(append(results, '\n'))
	if err != nil {
		log.Fatalf("Unable to write json to %s: %v", jsonFileName, err)
	}
	_, _ = fmt.Fprintf(os.Stderr,
		"Successfully wrote %d bytes of Json data to %s\n", n, jsonFileName)
}
