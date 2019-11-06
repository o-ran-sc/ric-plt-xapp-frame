/*
==================================================================================
  Copyright (c) 2019 AT&T Intellectual Property.
  Copyright (c) 2019 Nokia

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
==================================================================================
*/

package xapp

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	Namespace string
}

// Alias
type CounterOpts prometheus.CounterOpts
type Counter prometheus.Counter
type Gauge prometheus.Gauge

func NewMetrics(url, namespace string, r *mux.Router) *Metrics {
	if url == "" {
		url = "/ric/v1/metrics"
	}
	if namespace == "" {
		namespace = "ricxapp"
	}

	Logger.Info("Serving metrics on: url=%s namespace=%s", url, namespace)

	// Expose 'metrics' endpoint with standard golang metrics used by prometheus
	r.Handle(url, promhttp.Handler())

	return &Metrics{Namespace: namespace}
}

func (m *Metrics) RegisterCounter(opts CounterOpts) Counter {
	Logger.Info("Register new counter with opts: %v", opts)

	return promauto.NewCounter(prometheus.CounterOpts(opts))
}

func (m *Metrics) RegisterCounterGroup(opts []CounterOpts, subsytem string) (c map[string]Counter) {
	c = make(map[string]Counter)
	for _, opt := range opts {
		opt.Namespace = m.Namespace
		opt.Subsystem = subsytem
		c[opt.Name] = m.RegisterCounter(opt)
	}

	return
}

func (m *Metrics) RegisterGauge(opts CounterOpts) Gauge {
	Logger.Info("Register new gauge with opts: %v", opts)

	return promauto.NewGauge(prometheus.GaugeOpts(opts))
}

func (m *Metrics) RegisterGaugeGroup(opts []CounterOpts, subsytem string) (c map[string]Gauge) {
	c = make(map[string]Gauge)
	for _, opt := range opts {
		opt.Namespace = m.Namespace
		opt.Subsystem = subsytem
		c[opt.Name] = m.RegisterGauge(opt)
	}

	return
}
