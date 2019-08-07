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

package graph

import (
	"encoding/json"
	"errors"
	"log"
	"sort"
	"strings"
	// "fmt"

	jaeger "github.com/jaegertracing/jaeger/model/json"
)

type NodeData struct {
	ServiceName   string        `json:"service,omitempty"`
	SpanID        jaeger.SpanID `json:"-"`
	OperationName string        `json:"operationName"`
	StartTime     uint64        `json:"-"`
	Duration      uint64        `json:"-"`
	RequestType   string        `json:"requestType"`
	NodeID        string        `json:"nodeID"`
	RequestSize   string        `json:"-"`
	ResponseSize  string        `json:"-"`
}

type Node struct {
	Data     NodeData `json:"data"`
	Children *[]Node  `json:"children"`
}

type Graph struct {
	Root *Node `'json:"root"`
}

func (g *Graph) ExtractGraphData() []byte {
	bytes, err := json.Marshal(g.Root)

	if err != nil {
		log.Fatalf(`graph structure is improper`)
	}

	return bytes
}

func CombineGraphs(g1 *Graph, g2 *Graph) *Graph {
	node := _CombineGraphHelper(g1.Root, g2.Root)
	return &Graph{node}
}

func _CombineGraphHelper(node1 *Node, node2 *Node) *Node {
	ret := &Node{}
	if node1.Data.OperationName == node2.Data.OperationName &&
		node1.Data.RequestType == node2.Data.RequestType &&
		node1.Data.ServiceName == node2.Data.ServiceName {

		ret.Data.OperationName = node1.Data.OperationName
		ret.Data.RequestType = node1.Data.RequestType
		ret.Data.ServiceName = node1.Data.ServiceName

		*ret.Children = append(*node1.Children, *node2.Children...)
	}

	return ret
}

func findTag(tags []jaeger.KeyValue, key string) jaeger.KeyValue {
	for _, tag := range tags {
		if tag.Key == key {
			return tag
		}
	}

	return jaeger.KeyValue{}
}

func CompGraph(g1 *Graph, g2 *Graph) bool {
	return _CompGraphHelper(g1.Root, g2.Root)
}

func _CompGraphHelper(node1 *Node, node2 *Node) bool {
	if node1 == nil && node2 == nil {
		return true
	} else if node1 == nil || node2 == nil {
		return false
	}

	ret := true

	if node1.Data.OperationName == node2.Data.OperationName &&
		node1.Data.RequestType == node2.Data.RequestType &&
		node1.Data.ServiceName == node2.Data.ServiceName &&
		node1.Data.NodeID == node2.Data.NodeID &&
		len(*node1.Children) == len(*node2.Children) {
		for i := 0; i < len(*node1.Children); i++ {
			ret = ret && _CompGraphHelper(&(*node1.Children)[i], &(*node2.Children)[i])
		}
	} else {
		return false
	}

	return ret
}

// Root span has no references.
func findRootSpan(spans []jaeger.Span) (jaeger.Span, error) {
	for _, span := range spans {
		if len(span.References) == 0 {
			return span, nil
		}
	}

	return jaeger.Span{}, errors.New("root span not present in spans")
}

func findTags(tags []jaeger.KeyValue) (reqType string,
	nodeID string,
	respSize string,
	reqSize string) {

	/*
		Tag Examples:

		{Key: upstream_cluster
		 Type: string
		 Value: inbound|9080|http|productpage.default.svc.cluster.local
		}

		{Key: upstream_cluster
		 Type: string
		 Value: inbound|9080|http|reviews.default.svc.cluster.local
		}

		{Key: upstream_cluster
		 Type: string
		 Value: outbound|9080||ratings.default.svc.cluster.local
		}
	*/

	tag := findTag(tags, "upstream_cluster")

	if tag.Value == nil {
		reqType = ""
	} else {
		reqType = strings.Split(tag.Value.(string), "|")[0]
	}

	tag = findTag(tags, "node_id")
	if tag.Value == nil {
		nodeID = ""
	} else {
		nodeID = tag.Value.(string)
	}

	tag = findTag(tags, "response_size")
	if tag.Value == nil {
		respSize = ""
	} else {
		respSize = tag.Value.(string)
	}

	tag = findTag(tags, "request_size")
	if tag.Value == nil {
		reqSize = ""
	} else {
		reqSize = tag.Value.(string)
	}

	return
}

func GenerateGraph(spans []jaeger.Span,
	processes map[jaeger.ProcessID]jaeger.Process) (*Graph, error) {
	rootSpan, err := findRootSpan(spans)
	if err != nil {
		return nil, err
	}

	reqType, nodeID, respSize, reqSize := findTags(rootSpan.Tags)
	svcName := processes[rootSpan.ProcessID].ServiceName

	d := NodeData{svcName, rootSpan.SpanID, rootSpan.OperationName,
		rootSpan.StartTime, rootSpan.Duration, reqType, nodeID,
		reqSize, respSize}
	childrenList := UpdateChildren(spans, rootSpan.SpanID, processes)
	root := Node{d, &childrenList}
	return &Graph{&root}, nil
}

func UpdateChildren(spans []jaeger.Span, spanID jaeger.SpanID,
	processes map[jaeger.ProcessID]jaeger.Process) []Node {
	children := make([]Node, 0)

	for _, span := range spans {
		if len(span.References) == 0 {
			continue
		}

		ref := span.References[0]
		if ref.RefType == jaeger.ChildOf && ref.SpanID == spanID {
			reqType, nodeID, respSize, reqSize := findTags(span.Tags)
			svcName := processes[span.ProcessID].ServiceName

			d := NodeData{svcName, span.SpanID, span.OperationName,
				span.StartTime, span.Duration, reqType, nodeID,
				reqSize, respSize}

			nodeChildren := UpdateChildren(spans, span.SpanID, processes)
			children = append(children, Node{d, &nodeChildren})
		}
	}

	sort.Slice(children, func(i, j int) bool {
		return (children[i].Data.StartTime <
			children[j].Data.StartTime)
	})

	return children
}
