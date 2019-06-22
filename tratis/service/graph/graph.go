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
	"fmt"
	"istio.io/tools/tratis/service/parsing/pkg/span"
	"sort"
)

type Node struct {
	children *[]Node
	Data     span.Span
}

type Graph struct {
	root *Node
}

func PrintGraph(g *Graph, size int) [][]span.Span {
	if g == nil {
		return nil
	}
	ret := make([][]span.Span, size)

	for i := 0; i < size; i++ {
		ret[i] = make([]span.Span, 0)
	}

	PrintGraphHelper(g.root, &ret, 0)

	for i := 0; i < size; i++ {
		for j := 0; j < len(ret[i]); j++ {
			fmt.Print(ret[i][j].OperationName)
			fmt.Print(" ")
			fmt.Print(ret[i][j].SpanID)
			fmt.Print("::")
		}
		fmt.Println(" ")
	}

	return ret
}

func PrintGraphHelper(root *Node, ret *[][]span.Span, level int) {
	if root == nil {
		return
	} else {
		(*ret)[level] = append((*ret)[level], root.Data)

		for _, child := range *root.children {
			PrintGraphHelper(&child, ret, level+1)

		}
	}
}

func GenerateGraph(data map[string]span.Span) *Graph {
	for _, v := range data {
		if len(v.References) == 0 {
			childrenList := make([]Node, 0, 0)
			root := Node{&childrenList, v}
			UpdateChildren(data, &childrenList, v.SpanID)
			return &Graph{&root}
		}
	}

	return &Graph{&Node{nil, span.Span{}}}
}

func UpdateChildren(data map[string]span.Span, children *[]Node, SpanID string) {
	for _, v := range data {
		if len(v.References) == 0 {
			continue
		}

		ref := v.References[0]
		if ref.RefType == "CHILD_OF" && ref.SpanID == SpanID {
			childrenList := make([]Node, 0, 0)
			*children = append(*children, Node{&childrenList, v})

			UpdateChildren(data, &childrenList, v.SpanID)

			sort.Slice(childrenList, func(i, j int) bool {
				return (childrenList[i].Data.StartTime >
					childrenList[j].Data.StartTime)
			})
		}
	}
}
