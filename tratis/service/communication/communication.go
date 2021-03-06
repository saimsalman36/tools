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

package communication

/*
This package allows one to query the jaeger UI.
*/

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func FindNumServices(toolAddr string, toolPortNum string) []byte {
	pageAddress := fmt.Sprintf("http://%s:%s/jaeger/api/services",
		toolAddr, toolPortNum)

	resp, err := http.Get(pageAddress)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return body
}

func ExtractTraces(toolAddr string, toolPortNum string,
	appEntryPoint string, numTraces int,
	startTime string, duration string) []byte {

	pageAddress := ""

	if startTime != "" {
		start, err := time.Parse(time.RFC3339, startTime)
		if err != nil {
			log.Fatalln(err)
		}

		d, err := time.ParseDuration(duration)
		if err != nil {
			log.Fatalln(err)
		}

		end := start.Add(d)

		pageAddress = fmt.Sprintf("http://%s:%s/jaeger/api/traces?service=%s&limit=%d&lookback=custom&start=%d&end=%d",
			toolAddr, toolPortNum, appEntryPoint, numTraces, start.Unix()*1000000, end.Unix()*1000000)
	} else {
		pageAddress = fmt.Sprintf("http://%s:%s/jaeger/api/traces?service=%s&limit=%d",
			toolAddr, toolPortNum, appEntryPoint, numTraces)
	}

	resp, err := http.Get(pageAddress)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return (body)
}
