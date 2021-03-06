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

package distribution

import (
	"sort"

	"gonum.org/v1/gonum/stat"
	"istio.io/tools/tratis/service/graph"
)

type Time struct {
	StartTime uint64 `json:"startTime"`
	EndTime   uint64 `json:"endTime"`
	Duration  uint64 `json:"duration"`
}

type TimeInformation struct {
	TimeData      []Time `json:"time"`
	OperationName string `json:"operationName,omitempty"`
}

type CombinedTimeInformation struct {
	Duration      [][]float64 `json:"durations"`
	Average       []float64   `json:"average"`
	Variance      []float64   `json:"variance"`
	Median        []float64   `json:"median"`
	OperationName string      `json:"operationName"`
}

func (data CombinedTimeInformation) GetDistributionData() [][]float64 {
	return data.Duration
}

func (data CombinedTimeInformation) GetOperation() string {
	return data.OperationName
}

func (data CombinedTimeInformation) Convert() HasDistributionData {
	return HasDistributionData(data)
}

func CombineTimeInformation(data [][]TimeInformation) []CombinedTimeInformation {
	ret := make([]CombinedTimeInformation, 0)

	for _, trace := range data {
		for idx, span := range trace {
			if len(ret) < idx+1 {
				ret = append(ret, CombinedTimeInformation{})
			}
			ret[idx].OperationName = span.OperationName

			if len(ret[idx].Duration) == 0 {
				ret[idx].Duration = make([][]float64, len(span.TimeData))
			}

			for timeIndex, duration := range span.TimeData {
				ret[idx].Duration[timeIndex] =
					append(ret[idx].Duration[timeIndex], float64(duration.Duration))
			}
		}
	}

	for idx, op := range ret {
		for _, slice := range op.Duration {
			ret[idx].Average = append(ret[idx].Average, stat.Mean(slice, nil))
			ret[idx].Variance = append(ret[idx].Variance, stat.Variance(slice, nil))

			sort.Float64s(slice)
			ret[idx].Median = append(ret[idx].Median, stat.Quantile(0.5, stat.Empirical, slice, nil))
		}
	}

	return ret
}

func ExtractTimeInformation(g *graph.Graph) []TimeInformation {
	ret := make([]TimeInformation, 0)
	ExtractTimeInformationWrapper(g.Root, &ret)
	return ret
}

func ExtractTimeInformationWrapper(n *graph.Node, t *[]TimeInformation) {
	if n == nil {
		return
	}

	nodeStartTime := n.Data.StartTime
	nodeEndTime := n.Data.StartTime + n.Data.Duration

	timeData := make([]Time, 0)

	for _, child := range *n.Children {
		d := child.Data.StartTime - nodeStartTime
		newTime := Time{nodeStartTime, child.Data.StartTime, d}
		timeData = append(timeData, newTime)
		nodeStartTime = child.Data.StartTime + child.Data.Duration

		ExtractTimeInformationWrapper(&child, t)
	}

	d := nodeEndTime - nodeStartTime
	newTime := Time{nodeStartTime, nodeEndTime, d}
	timeData = append(timeData, newTime)

	if n.Data.RequestType == "inbound" {
		*t = append(*t, TimeInformation{timeData, n.Data.OperationName})
	}
}
