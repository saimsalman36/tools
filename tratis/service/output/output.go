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

package output

import (
	"encoding/json"

	dist "istio.io/tools/tratis/service/distribution"
	"istio.io/tools/tratis/service/graph"
	parser "istio.io/tools/tratis/service/parsing"
	"istio.io/tools/tratis/service/pkg/consts"

	"fmt"
)

type BasicData struct {
	NumTraces int         `json:"NumTraces"`
	CallGraph *graph.Node `json:"Graph"`
}

type RawOutput struct {
	BasicData
	TimeInformation         []dist.HasDistributionData `json:"-"`
	RequestSizeInformation  []dist.HasDistributionData `json:"-"`
	ResponseSizeInformation []dist.HasDistributionData `json:"-"`
}

type DistributionOutput struct {
	BasicData
	TimeInformation         []dist.TotalDistributions `json:"TimeInformation"`
	RequestSizeInformation  []dist.TotalDistributions `json:"RequestSizeInformation"`
	ResponseSizeInformation []dist.TotalDistributions `json:"ResponseSizeInformation"`
}

func addNewGraph(data *[][]*graph.Graph, g *graph.Graph) int {
	ret := -1
	for idx, graphs := range *data {
		if graph.CompGraph(graphs[0], g) {
			ret = idx
			break
		}
	}

	if ret == -1 {
		temp := make([]*graph.Graph, 0)
		temp = append(temp, g)
		*data = append(*data, temp)

		return (len(*data) - 1)
	}
	(*data)[ret] = append((*data)[ret], g)
	return ret
}

func GenerateOutput(data parser.TraceData) []byte {
	traces := data.Traces

	fmt.Printf("Processing %d Traces\n", len(traces))

	d := make([][][]dist.TimeInformation, 0)
	s := make([][][]dist.MessageSizeInfo, 0)

	graphs := make([][]*graph.Graph, 0)

	fmt.Println("Generating Time Information ...")

	for _, trace := range traces {
		g, err := graph.GenerateGraph(trace.Spans, trace.Processes)

		if err != nil {
			continue
		}

		idx := addNewGraph(&graphs, g)

		traceInformation := dist.ExtractTimeInformation(g)
		sizeInformation := dist.ExtractSizeInformation(g)

		if len(d) < idx+1 {
			d = append(d, make([][]dist.TimeInformation, 0))
		}

		if len(s) < idx+1 {
			s = append(s, make([][]dist.MessageSizeInfo, 0))
		}

		d[idx] = append(d[idx], traceInformation)
		s[idx] = append(s[idx], sizeInformation)
	}

	fmt.Println("Combining Results + Distribution Fitting ...")

	ret := make([]RawOutput, 0)

	for idx := range graphs {

		if len(graphs[idx]) > consts.MinNumTraces {
			fmt.Println("Processing Graph Number: ", idx)

			var temp RawOutput

			temp.NumTraces = len(graphs[idx])
			temp.CallGraph = graphs[idx][0].Root

			fmt.Println("Extra Time Information")

			timeResults := dist.CombineTimeInformation(d[idx])
			temp.TimeInformation = dist.ConvertTimeInfo(timeResults)

			fmt.Println("Extracting Response Size Information")

			sizeResults := dist.ConvertSizeInfo(dist.CombineSizeInformation(s[idx], true))
			temp.ResponseSizeInformation = sizeResults

			fmt.Println("Extracting Request Size Information")

			sizeResults = dist.ConvertSizeInfo(dist.CombineSizeInformation(s[idx], false))
			temp.RequestSizeInformation = sizeResults

			ret = append(ret, temp)
		} else {
			fmt.Println("Skipping Graph Number: ", idx)
			fmt.Println("")
		}
	}

	bytes, _ := json.MarshalIndent(ret, "", "  ")
	return bytes
}
