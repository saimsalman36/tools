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

package process

import (
	"encoding/json"
	"fmt"

	"istio.io/tools/tratis/service/parsing/pkg/tag"
)

// UnmarshalJSON converts b to a Service, applying the default values from
// DefaultService.
func (process *Process) UnmarshalJSON(b []byte) (err error) {
	var processData unmarshableProcess
	err = json.Unmarshal(b, &processData)
	if err != nil {
		fmt.Println(err)
		return
	}

	process.ServiceName = processData.ServiceName
	process.Tags = make(tag.Tag)

	for _, t := range processData.Tags {
		process.Tags[t["key"].(string)] = t["value"].(string)
	}

	return
}

type unmarshableProcess struct {
	ServiceName string    `json:"serviceName"`
	Tags        []tag.Tag `json:"tags"`
}
