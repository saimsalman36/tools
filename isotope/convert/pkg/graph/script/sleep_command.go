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

package script

import (
	"encoding/json"
	"log"
	"time"

	"github.com/jmcvetta/randutil"
	"gonum.org/v1/gonum/stat/distuv"
)

type CommandType string

const (
	Static       CommandType = "static"
	Histogram    CommandType = "histogram"
	Distribution CommandType = "dist"
)

type SleepCommandStatic struct {
	Time time.Duration `json:"time"`
}

func (c SleepCommandStatic) Duration() time.Duration {
	return c.Time
}

type SleepCommandDistribution struct {
	DistName string `json:"name"`
	Dist interface {
		Rand() float64
	}
}

func (c SleepCommandDistribution) Duration() time.Duration {
	return time.Duration(c.Dist.Rand() * 10e8)
}

type SleepCommandHistogram struct {
	Histogram []randutil.Choice
}

func (c SleepCommandHistogram) Duration() time.Duration {
	result, err := randutil.WeightedChoice(c.Histogram)
	if err != nil {
		panic(err)
	}

	return (result.Item.(time.Duration))
}

type SleepCommandWrapper struct {
	Type CommandType `json:"Type"`
	Data json.RawMessage
}

type SleepCommand struct {
	Type CommandType
	Data interface {
		Duration() time.Duration
	}
}

// UnmarshalJSON converts a JSON object to a SleepCommand.
func (c *SleepCommand) UnmarshalJSON(b []byte) (err error) {
	var command SleepCommandWrapper
	err = json.Unmarshal(b, &command)
	if err != nil {
		return
	}

	switch command.Type {
	case Static:
		var cmd map[string]interface{}
		err = json.Unmarshal(command.Data, &cmd)

		if err != nil {
			return
		}

		var staticCmd SleepCommandStatic
		staticCmd.Time, err = time.ParseDuration(cmd["time"].(string))

		if err != nil {
			return
		}

		*c = SleepCommand{command.Type, staticCmd}

	case Distribution:
		var cmd map[string]interface{}
		err = json.Unmarshal(command.Data, &cmd)

		if err != nil {
			return
		}

		var distCmd SleepCommandDistribution
		var dist interface {
			Rand() float64
		}

		switch cmd["name"] {
		case "normal":
			distData := cmd["dist"].(map[string]interface{})
			dist = distuv.Normal{
				Mu:    distData["mean"].(float64),
				Sigma: distData["sigma"].(float64),
			}
		case "lognormal":
			distData := cmd["dist"].(map[string]interface{})
			dist = distuv.LogNormal{
				Mu:    distData["mean"].(float64),
				Sigma: distData["sigma"].(float64),
			}
		case "beta":
			distData := cmd["dist"].(map[string]interface{})
			dist = distuv.Beta{
				Alpha: distData["alpha"].(float64),
				Beta:  distData["beta"].(float64),
			}
		case "chi-squared":
			distData := cmd["dist"].(map[string]interface{})
			dist = distuv.ChiSquared{
				K: distData["k"].(float64),
			}
		case "exp":
			distData := cmd["dist"].(map[string]interface{})
			dist = distuv.Exponential{
				Rate: distData["rate"].(float64),
			}
		case "f":
			distData := cmd["dist"].(map[string]interface{})
			dist = distuv.F{
				D1: distData["d1"].(float64),
				D2: distData["d2"].(float64),
			}
		case "gamma":
			distData := cmd["dist"].(map[string]interface{})
			dist = distuv.Gamma{
				Alpha: distData["alpha"].(float64),
				Beta:  distData["beta"].(float64),
			}
		case "gumbel-right":
			distData := cmd["dist"].(map[string]interface{})
			dist = distuv.GumbelRight{
				Mu:   distData["mu"].(float64),
				Beta: distData["beta"].(float64),
			}
		case "inverse-gamma":
			distData := cmd["dist"].(map[string]interface{})
			dist = distuv.InverseGamma{
				Alpha: distData["alpha"].(float64),
				Beta:  distData["beta"].(float64),
			}
		case "laplace":
			distData := cmd["dist"].(map[string]interface{})
			dist = distuv.Laplace{
				Mu:    distData["mu"].(float64),
				Scale: distData["scale"].(float64),
			}
		case "pareto":
			distData := cmd["dist"].(map[string]interface{})
			dist = distuv.Pareto{
				Xm:    distData["xm"].(float64),
				Alpha: distData["alpha"].(float64),
			}
		case "studentst":
			distData := cmd["dist"].(map[string]interface{})
			dist = distuv.StudentsT{
				Mu:    distData["mu"].(float64),
				Sigma: distData["sigma"].(float64),
				Nu:    distData["nu"].(float64),
			}
		case "weibull":
			distData := cmd["dist"].(map[string]interface{})
			dist = distuv.Weibull{
				K:      distData["k"].(float64),
				Lambda: distData["lambda"].(float64),
			}
		}

		distCmd.DistName = cmd["name"].(string)
		distCmd.Dist = dist
		*c = SleepCommand{command.Type, distCmd}

	case Histogram:
		var probDistribution map[string]int
		err = json.Unmarshal(command.Data, &probDistribution)

		if err != nil {
			return
		}

		ret := make([]randutil.Choice, 0, len(probDistribution))
		totalPercentage := 0

		for timeString, percentage := range probDistribution {
			duration, err := time.ParseDuration(timeString)

			if err != nil {
				return err
			}

			ret = append(ret, randutil.Choice{percentage, duration})
			totalPercentage += percentage
		}

		if totalPercentage != 100 {
			log.Fatalf("Total Percentage does not equal 100.")
		}

		var HistCmd SleepCommandHistogram
		HistCmd.Histogram = ret
		*c = SleepCommand{command.Type, HistCmd}

	}
	return
}

func (c SleepCommand) String() string {
	res, err := json.Marshal(c)

	if err != nil {
		log.Fatalf("%s", err)
	}

	return string(res)
}
