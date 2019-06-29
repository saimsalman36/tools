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

package graph

import (
	"encoding/json"
	jaeger "github.com/jaegertracing/jaeger/model/json"
	"sort"
)

type NodeData struct {
	SpanID        jaeger.SpanID     `json:"spanID,omitempty"`
	OperationName string            `json:"operationName,omitempty"`
	StartTime     uint64            `json:"startTime,omitempty"`
	Duration      uint64            `json:"duration,omitempty"`
	Tags          []jaeger.KeyValue `json:"tags"`
}

type Node struct {
	Data     NodeData `json:"data"`
	Children *[]Node  `json:"children"`
}

type Graph struct {
	Root *Node `'json:"root"`
}

func (g *Graph) ExtractGraphData() []byte {
	bytes, _ := json.Marshal(g.Root)

	return bytes
}

func GenerateGraph(data []jaeger.Span) *Graph {
	for _, v := range data {
		if len(v.References) == 0 {
			childrenList := make([]Node, 0, 0)
			d := NodeData{v.SpanID, v.OperationName,
				v.StartTime, v.Duration, v.Tags}
			root := Node{d, &childrenList}
			UpdateChildren(data, &childrenList, v.SpanID)
			return &Graph{&root}
		}
	}

	return &Graph{&Node{NodeData{}, nil}}
}

func UpdateChildren(data []jaeger.Span, children *[]Node, SpanID jaeger.SpanID) {
	for _, v := range data {
		if len(v.References) == 0 {
			continue
		}

		ref := v.References[0]
		if ref.RefType == jaeger.ChildOf && ref.SpanID == SpanID {
			childrenList := make([]Node, 0, 0)
			d := NodeData{v.SpanID, v.OperationName,
				v.StartTime, v.Duration, v.Tags}
			*children = append(*children, Node{d, &childrenList})

			UpdateChildren(data, &childrenList, v.SpanID)

			sort.Slice(childrenList, func(i, j int) bool {
				return (childrenList[i].Data.StartTime <
					childrenList[j].Data.StartTime)
			})
		}
	}
}