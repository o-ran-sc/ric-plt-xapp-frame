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
type CounterOpts prometheus.Opts
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

/*
 * Handling counters
 */
func (m *Metrics) registerCounter(opts CounterOpts) Counter {
	Logger.Info("Register new counter with opts: %v", opts)
	return promauto.NewCounter(prometheus.CounterOpts(opts))
}

func (m *Metrics) RegisterCounterGroup(opts []CounterOpts, subsytem string) (c map[string]Counter) {
	c = make(map[string]Counter)
	for _, opt := range opts {
		opt.Namespace = m.Namespace
		opt.Subsystem = subsytem
		c[opt.Name] = m.registerCounter(opt)
	}

	return
}

/*
 * Handling gauges
 */
func (m *Metrics) registerGauge(opts CounterOpts) Gauge {
	Logger.Info("Register new gauge with opts: %v", opts)
	return promauto.NewGauge(prometheus.GaugeOpts(opts))
}

func (m *Metrics) RegisterGaugeGroup(opts []CounterOpts, subsytem string) (c map[string]Gauge) {
	c = make(map[string]Gauge)
	for _, opt := range opts {
		opt.Namespace = m.Namespace
		opt.Subsystem = subsytem
		c[opt.Name] = m.registerGauge(opt)
	}

	return
}

/*
 * Handling counter vectors
 *
 * Example:

	vec := xapp.Metric.RegisterCounterVecGroup(
		[]xapp.CounterOpts{
			{Name: "counter1", Help: "counter1"},
			{Name: "counter2", Help: "counter2"},
		},
		[]string{"host"},
		"SUBSYSTEM")

	stat:=xapp.Metric.GetCounterGroupFromVec(vec, []string{"localhost:8888"})

*/
type CounterVec struct {
	Vec  *prometheus.CounterVec
	Opts CounterOpts
}

func (m *Metrics) registerCounterVec(opts CounterOpts, labelNames []string) *prometheus.CounterVec {
	Logger.Info("Register new counter vector with opts: %v labelNames: %v", opts, labelNames)

	return promauto.NewCounterVec(prometheus.CounterOpts(opts), labelNames)
}

func (m *Metrics) RegisterCounterVecGroup(opts []CounterOpts, labelNames []string, subsytem string) (c map[string]CounterVec) {
	c = make(map[string]CounterVec)
	for _, opt := range opts {
		entry := CounterVec{}
		entry.Opts = opt
		entry.Opts.Namespace = m.Namespace
		entry.Opts.Subsystem = subsytem
		entry.Vec = m.registerCounterVec(entry.Opts, labelNames)
		c[opt.Name] = entry
	}
	return
}

func (m *Metrics) GetCounterGroupFromVec(vec map[string]CounterVec, labels []string) (c map[string]Counter) {
	c = make(map[string]Counter)
	for name, opt := range vec {
		c[name] = opt.Vec.WithLabelValues(labels...)
		Logger.Info("Register new counter for vector with opts: %v labels: %v", opt.Opts, labels)
	}
	return
}

/*
 * Handling gauge vectors
 *
 * Example:

	vec := xapp.Metric.RegisterGaugeVecGroup(
		[]xapp.CounterOpts{
			{Name: "gauge1", Help: "gauge1"},
			{Name: "gauge2", Help: "gauge2"},
		},
		[]string{"host"},
		"SUBSYSTEM")

	stat:=xapp.Metric.GetGaugeGroupFromVec(vec, []string{"localhost:8888"})

*/
type GaugeVec struct {
	Vec  *prometheus.GaugeVec
	Opts CounterOpts
}

func (m *Metrics) registerGaugeVec(opts CounterOpts, labelNames []string) *prometheus.GaugeVec {
	Logger.Info("Register new gauge vector with opts: %v labelNames: %v", opts, labelNames)

	return promauto.NewGaugeVec(prometheus.GaugeOpts(opts), labelNames)
}

func (m *Metrics) RegisterGaugeVecGroup(opts []CounterOpts, labelNames []string, subsytem string) (c map[string]GaugeVec) {
	c = make(map[string]GaugeVec)
	for _, opt := range opts {
		entry := GaugeVec{}
		entry.Opts = opt
		entry.Opts.Namespace = m.Namespace
		entry.Opts.Subsystem = subsytem
		entry.Vec = m.registerGaugeVec(entry.Opts, labelNames)
		c[opt.Name] = entry

	}
	return
}

func (m *Metrics) GetGaugeGroupFromVec(vec map[string]GaugeVec, labels []string) (c map[string]Gauge) {
	c = make(map[string]Gauge)
	for name, opt := range vec {
		c[name] = opt.Vec.WithLabelValues(labels...)
		Logger.Info("Register new gauge for vector with opts: %v labels: %v", opt.Opts, labels)
	}
	return
}
