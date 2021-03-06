// Copyright 2018 The Kubeflow Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"

	"github.com/onrik/logrus/filename"
	log "github.com/sirupsen/logrus"

	"github.com/NTHU-LSALAB/DRAGON/cmd/DRAGON/app"
	"github.com/NTHU-LSALAB/DRAGON/cmd/DRAGON/app/options"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func init() {
	// Add filename as one of the fields of the structured log message.
	filenameHook := filename.NewHook()
	filenameHook.Field = "filename"
	log.AddHook(filenameHook)
}

func startMonitoring(monitoringPort int) {
	go func() {
		log.Infof("Setting up client for monitoring on port: %s", strconv.Itoa(monitoringPort))
		http.Handle("/metrics", promhttp.Handler())
		err := http.ListenAndServe(fmt.Sprintf(":%s", strconv.Itoa(monitoringPort)), nil)
		if err != nil {
			log.Error("Monitoring endpoint setup failure.")
		}
	}()
}

func main() {
	s := options.NewServerOption()
	s.AddFlags(flag.CommandLine)

	flag.Parse()

	if s.JSONLogFormat {
		// Output logs in a json format so that it can be parsed by services like Stackdriver.
		log.SetFormatter(&log.JSONFormatter{})
	}

	startMonitoring(s.MonitoringPort)

	if err := app.Run(s); err != nil {
		log.Fatalf("%v\n", err)
	}

}
