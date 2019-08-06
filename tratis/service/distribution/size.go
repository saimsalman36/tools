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
	"strconv"

	"gonum.org/v1/gonum/stat"
	"istio.io/tools/tratis/service/graph"
)

type MessageSize struct {
	ResponseSize string `json:"responseSize"`
	RequestSize  string `json:"requestSize"`
}

type MessageSizeInfo struct {
	MsgInfo       []MessageSize `json:"size"`
	OperationName string        `json:"operationName,omitempty"`
}

type CombinedSizeInformation struct {
	Size [][]float64 `json:"sizes"`
	// Sums          []float64   `json:sum"`
	Average  []float64 `json:"average"`
	Variance []float64 `json:"variance"`
	Median   []float64 `json:"median"`
	// Min          []uint64   `json:"minimum"`
	// Max          []uint64   `json:"maximum"`
	OperationName string `json:"operationName"`
}

func (data CombinedSizeInformation) GetDistributionData() [][]float64 {
	return data.Size
}

func (data CombinedSizeInformation) GetOperation() string {
	return data.OperationName
}

func (data CombinedSizeInformation) Convert() HasDistributionData {
	return HasDistributionData(data)
}

func CombineSizeInformation(data [][]MessageSizeInfo, isResponse bool) []CombinedSizeInformation {
	ret := make([]CombinedSizeInformation, 0)

	for _, trace := range data {
		for idx, span := range trace {
			if len(ret) < idx+1 {
				ret = append(ret, CombinedSizeInformation{})
			}
			ret[idx].OperationName = span.OperationName

			if len(ret[idx].Size) == 0 {
				ret[idx].Size = make([][]float64, len(span.MsgInfo))
			}

			for sizeIndex, size := range span.MsgInfo {
				if isResponse {
					value, _ := strconv.Atoi(size.ResponseSize)
					ret[idx].Size[sizeIndex] =
						append(ret[idx].Size[sizeIndex], float64(value))
				} else {
					value, _ := strconv.Atoi(size.RequestSize)
					ret[idx].Size[sizeIndex] =
						append(ret[idx].Size[sizeIndex], float64(value))
				}
			}
		}
	}

	for idx, op := range ret {
		for _, slice := range op.Size {
			ret[idx].Average = append(ret[idx].Average, stat.Mean(slice, nil))
			ret[idx].Variance = append(ret[idx].Variance, stat.Variance(slice, nil))

			sort.Float64s(slice)
			ret[idx].Median = append(ret[idx].Median, stat.Quantile(0.5, stat.Empirical, slice, nil))
		}
	}

	return ret
}

func ExtractSizeInformation(g *graph.Graph) []MessageSizeInfo {
	ret := make([]MessageSizeInfo, 0)
	ExtractSizeInformationWrapper(g.Root, &ret)
	return ret
}

func ExtractSizeInformationWrapper(n *graph.Node, t *[]MessageSizeInfo) {
	if n == nil {
		return
	}

	for _, child := range *n.Children {
		ExtractSizeInformationWrapper(&child, t)
	}

	if n.Data.RequestType == "inbound" {
		sizeData := make([]MessageSize, 0)
		newSize := MessageSize{n.Data.ResponseSize, n.Data.RequestSize}
		sizeData = append(sizeData, newSize)
		*t = append(*t, MessageSizeInfo{sizeData, n.Data.OperationName})
	}
}
