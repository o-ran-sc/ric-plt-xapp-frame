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
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"sync"
)

//-----------------------------------------------------------------------------
// Alias
//-----------------------------------------------------------------------------
type CounterOpts prometheus.Opts
type Counter prometheus.Counter
type Gauge prometheus.Gauge

type CounterVec struct {
	Vec    *prometheus.CounterVec
	Opts   CounterOpts
	Labels []string
}

type GaugeVec struct {
	Vec    *prometheus.GaugeVec
	Opts   CounterOpts
	Labels []string
}

func strSliceCompare(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

//-----------------------------------------------------------------------------
//
//-----------------------------------------------------------------------------

type MetricGroupsCache struct {
	sync.RWMutex //This is for map locking
	counters     map[string]Counter
	gauges       map[string]Gauge
}

func (met *MetricGroupsCache) CIs(metric string) bool {
	met.RLock()
	defer met.RUnlock()
	_, ok := met.counters[metric]
	return ok
}

func (met *MetricGroupsCache) CGet(metric string) Counter {
	met.RLock()
	defer met.RUnlock()
	return met.counters[metric]
}

func (met *MetricGroupsCache) CInc(metric string) {
	met.RLock()
	defer met.RUnlock()
	met.counters[metric].Inc()
}

func (met *MetricGroupsCache) CAdd(metric string, val float64) {
	met.RLock()
	defer met.RUnlock()
	met.counters[metric].Add(val)
}

func (met *MetricGroupsCache) GIs(metric string) bool {
	met.RLock()
	defer met.RUnlock()
	_, ok := met.gauges[metric]
	return ok
}

func (met *MetricGroupsCache) GGet(metric string) Gauge {
	met.RLock()
	defer met.RUnlock()
	return met.gauges[metric]
}

func (met *MetricGroupsCache) GSet(metric string, val float64) {
	met.RLock()
	defer met.RUnlock()
	met.gauges[metric].Set(val)
}

func (met *MetricGroupsCache) GAdd(metric string, val float64) {
	met.RLock()
	defer met.RUnlock()
	met.gauges[metric].Add(val)
}

func (met *MetricGroupsCache) GInc(metric string) {
	met.RLock()
	defer met.RUnlock()
	met.gauges[metric].Inc()
}

func (met *MetricGroupsCache) GDec(metric string) {
	met.RLock()
	defer met.RUnlock()
	met.gauges[metric].Dec()
}

func (met *MetricGroupsCache) CombineCounterGroupsWithPrefix(prefix string, srcs ...map[string]Counter) {
	met.Lock()
	defer met.Unlock()
	for _, src := range srcs {
		for k, v := range src {
			met.counters[prefix+k] = v
		}
	}
}

func (met *MetricGroupsCache) CombineCounterGroups(srcs ...map[string]Counter) {
	met.Lock()
	defer met.Unlock()
	for _, src := range srcs {
		for k, v := range src {
			met.counters[k] = v
		}
	}
}

func (met *MetricGroupsCache) CombineGaugeGroupsWithPrefix(prefix string, srcs ...map[string]Gauge) {
	met.Lock()
	defer met.Unlock()
	for _, src := range srcs {
		for k, v := range src {
			met.gauges[prefix+k] = v
		}
	}
}

func (met *MetricGroupsCache) CombineGaugeGroups(srcs ...map[string]Gauge) {
	met.Lock()
	defer met.Unlock()
	for _, src := range srcs {
		for k, v := range src {
			met.gauges[k] = v
		}
	}
}

func NewMetricGroupsCache() *MetricGroupsCache {
	entry := &MetricGroupsCache{}
	entry.counters = make(map[string]Counter)
	entry.gauges = make(map[string]Gauge)
	return entry
}

//-----------------------------------------------------------------------------
// All counters/gauges registered via Metrics instances:
// Counter names are build from: namespace, subsystem, metric and possible labels
//-----------------------------------------------------------------------------
var globalLock sync.Mutex
var cache_allcounters map[string]Counter
var cache_allgauges map[string]Gauge
var cache_allcountervects map[string]CounterVec
var cache_allgaugevects map[string]GaugeVec

func init() {
	cache_allcounters = make(map[string]Counter)
	cache_allgauges = make(map[string]Gauge)
	cache_allcountervects = make(map[string]CounterVec)
	cache_allgaugevects = make(map[string]GaugeVec)
}

//-----------------------------------------------------------------------------
//
//-----------------------------------------------------------------------------
type Metrics struct {
	Namespace string
}

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
 * Helpers
 */
func (m *Metrics) getFullName(opts prometheus.Opts, labels []string) string {
	labelname := ""
	for _, lbl := range labels {
		if len(labelname) == 0 {
			labelname += lbl
		} else {
			labelname += "_" + lbl
		}
	}
	return fmt.Sprintf("%s_%s_%s_%s", opts.Namespace, opts.Subsystem, opts.Name, labelname)
}

//
//
//
func (m *Metrics) RegisterCounter(opts CounterOpts, subsytem string) Counter {
	globalLock.Lock()
	defer globalLock.Unlock()
	opts.Namespace = m.Namespace
	opts.Subsystem = subsytem
	id := m.getFullName(prometheus.Opts(opts), []string{})
	if _, ok := cache_allcounters[id]; !ok {
		Logger.Info("Register new counter with opts: %v", opts)
		cache_allcounters[id] = promauto.NewCounter(prometheus.CounterOpts(opts))
	}
	return cache_allcounters[id]
}

//
//
//
func (m *Metrics) RegisterCounterGroup(optsgroup []CounterOpts, subsytem string) map[string]Counter {
	c := make(map[string]Counter)
	for _, opts := range optsgroup {
		c[opts.Name] = m.RegisterCounter(opts, subsytem)
	}
	return c
}

//
//
//
func (m *Metrics) RegisterLabeledCounter(opts CounterOpts, labelNames []string, labelValues []string, subsytem string) Counter {
	globalLock.Lock()
	defer globalLock.Unlock()
	opts.Namespace = m.Namespace
	opts.Subsystem = subsytem
	vecid := m.getFullName(prometheus.Opts(opts), []string{})
	if _, ok := cache_allcountervects[vecid]; !ok {
		Logger.Info("Register new counter vector with opts: %v labelNames: %v", opts, labelNames)
		entry := CounterVec{}
		entry.Opts = opts
		entry.Labels = labelNames
		entry.Vec = promauto.NewCounterVec(prometheus.CounterOpts(entry.Opts), entry.Labels)
		cache_allcountervects[vecid] = entry
	}
	entry := cache_allcountervects[vecid]
	if strSliceCompare(entry.Labels, labelNames) == false {
		Logger.Warn("id:%s cached counter vec labels dont match %v != %v", vecid, entry.Labels, labelNames)
	}

	valid := m.getFullName(prometheus.Opts(entry.Opts), labelValues)
	if _, ok := cache_allcounters[valid]; !ok {
		Logger.Info("Register new counter from vector with opts: %v labelValues: %v", entry.Opts, labelValues)
		cache_allcounters[valid] = entry.Vec.WithLabelValues(labelValues...)
	}
	return cache_allcounters[valid]

}

//
//
//
func (m *Metrics) RegisterLabeledCounterGroup(optsgroup []CounterOpts, labelNames []string, labelValues []string, subsytem string) map[string]Counter {
	c := make(map[string]Counter)
	for _, opts := range optsgroup {
		c[opts.Name] = m.RegisterLabeledCounter(opts, labelNames, labelValues, subsytem)
	}
	return c
}

//
//
//
func (m *Metrics) RegisterGauge(opts CounterOpts, subsytem string) Gauge {
	globalLock.Lock()
	defer globalLock.Unlock()
	opts.Namespace = m.Namespace
	opts.Subsystem = subsytem
	id := m.getFullName(prometheus.Opts(opts), []string{})
	if _, ok := cache_allgauges[id]; !ok {
		Logger.Info("Register new gauge with opts: %v", opts)
		cache_allgauges[id] = promauto.NewGauge(prometheus.GaugeOpts(opts))
	}
	return cache_allgauges[id]
}

//
//
//
func (m *Metrics) RegisterGaugeGroup(optsgroup []CounterOpts, subsytem string) map[string]Gauge {
	c := make(map[string]Gauge)
	for _, opts := range optsgroup {
		c[opts.Name] = m.RegisterGauge(opts, subsytem)
	}
	return c
}

//
//
//
func (m *Metrics) RegisterLabeledGauge(opt CounterOpts, labelNames []string, labelValues []string, subsytem string) Gauge {
	globalLock.Lock()
	defer globalLock.Unlock()
	opt.Namespace = m.Namespace
	opt.Subsystem = subsytem
	vecid := m.getFullName(prometheus.Opts(opt), []string{})
	if _, ok := cache_allgaugevects[vecid]; !ok {
		Logger.Info("Register new gauge vector with opt: %v labelNames: %v", opt, labelNames)
		entry := GaugeVec{}
		entry.Opts = opt
		entry.Labels = labelNames
		entry.Vec = promauto.NewGaugeVec(prometheus.GaugeOpts(entry.Opts), entry.Labels)
		cache_allgaugevects[vecid] = entry
	}
	entry := cache_allgaugevects[vecid]
	if strSliceCompare(entry.Labels, labelNames) == false {
		Logger.Warn("id:%s cached gauge vec labels dont match %v != %v", vecid, entry.Labels, labelNames)
	}
	valid := m.getFullName(prometheus.Opts(entry.Opts), labelValues)
	if _, ok := cache_allgauges[valid]; !ok {
		Logger.Info("Register new gauge from vector with opts: %v labelValues: %v", entry.Opts, labelValues)
		cache_allgauges[valid] = entry.Vec.WithLabelValues(labelValues...)
	}
	return cache_allgauges[valid]

}

//
//
//
func (m *Metrics) RegisterLabeledGaugeGroup(opts []CounterOpts, labelNames []string, labelValues []string, subsytem string) map[string]Gauge {
	c := make(map[string]Gauge)
	for _, opt := range opts {
		c[opt.Name] = m.RegisterLabeledGauge(opt, labelNames, labelValues, subsytem)
	}
	return c
}

/*
 * Handling counter vectors
 *
 * Examples:

  //---------
	vec := Metric.RegisterCounterVec(
		CounterOpts{Name: "counter0", Help: "counter0"},
		[]string{"host"},
		"SUBSYSTEM")

	stat:=Metric.GetCounterFromVect([]string{"localhost:8888"},vec)
	stat.Inc()

  //---------
	vec := Metric.RegisterCounterVecGroup(
		[]CounterOpts{
			{Name: "counter1", Help: "counter1"},
			{Name: "counter2", Help: "counter2"},
		},
		[]string{"host"},
		"SUBSYSTEM")

	stats:=Metric.GetCounterGroupFromVects([]string{"localhost:8888"}, vec)
	stats["counter1"].Inc()
*/

// Deprecated: Use RegisterLabeledCounter
func (m *Metrics) RegisterCounterVec(opts CounterOpts, labelNames []string, subsytem string) CounterVec {
	globalLock.Lock()
	defer globalLock.Unlock()
	opts.Namespace = m.Namespace
	opts.Subsystem = subsytem
	id := m.getFullName(prometheus.Opts(opts), []string{})
	if _, ok := cache_allcountervects[id]; !ok {
		Logger.Info("Register new counter vector with opts: %v labelNames: %v", opts, labelNames)
		entry := CounterVec{}
		entry.Opts = opts
		entry.Labels = labelNames
		entry.Vec = promauto.NewCounterVec(prometheus.CounterOpts(entry.Opts), entry.Labels)
		cache_allcountervects[id] = entry
	}
	entry := cache_allcountervects[id]
	if strSliceCompare(entry.Labels, labelNames) == false {
		Logger.Warn("id:%s cached counter vec labels dont match %v != %v", id, entry.Labels, labelNames)
	}
	return entry
}

// Deprecated: Use RegisterLabeledCounterGroup
func (m *Metrics) RegisterCounterVecGroup(optsgroup []CounterOpts, labelNames []string, subsytem string) map[string]CounterVec {
	c := make(map[string]CounterVec)
	for _, opts := range optsgroup {
		c[opts.Name] = m.RegisterCounterVec(opts, labelNames, subsytem)
	}
	return c
}

// Deprecated: Use RegisterLabeledCounter
func (m *Metrics) GetCounterFromVect(labelValues []string, vec CounterVec) (c Counter) {
	globalLock.Lock()
	defer globalLock.Unlock()
	id := m.getFullName(prometheus.Opts(vec.Opts), labelValues)
	if _, ok := cache_allcounters[id]; !ok {
		Logger.Info("Register new counter from vector with opts: %v labelValues: %v", vec.Opts, labelValues)
		cache_allcounters[id] = vec.Vec.WithLabelValues(labelValues...)
	}
	return cache_allcounters[id]
}

// Deprecated: Use RegisterLabeledCounterGroup
func (m *Metrics) GetCounterGroupFromVects(labelValues []string, vects ...map[string]CounterVec) map[string]Counter {
	c := make(map[string]Counter)
	for _, vect := range vects {
		for name, vec := range vect {
			c[name] = m.GetCounterFromVect(labelValues, vec)
		}
	}
	return c
}

// Deprecated: Use RegisterLabeledCounterGroup
func (m *Metrics) GetCounterGroupFromVectsWithPrefix(prefix string, labelValues []string, vects ...map[string]CounterVec) map[string]Counter {
	c := make(map[string]Counter)
	for _, vect := range vects {
		for name, vec := range vect {
			c[prefix+name] = m.GetCounterFromVect(labelValues, vec)
		}
	}
	return c
}

/*
 * Handling gauge vectors
 *
 * Examples:

  //---------
	vec := Metric.RegisterGaugeVec(
		CounterOpts{Name: "gauge0", Help: "gauge0"},
		[]string{"host"},
		"SUBSYSTEM")

	stat:=Metric.GetGaugeFromVect([]string{"localhost:8888"},vec)
	stat.Inc()

  //---------
	vecgrp := Metric.RegisterGaugeVecGroup(
		[]CounterOpts{
			{Name: "gauge1", Help: "gauge1"},
			{Name: "gauge2", Help: "gauge2"},
		},
		[]string{"host"},
		"SUBSYSTEM")

	stats:=Metric.GetGaugeGroupFromVects([]string{"localhost:8888"},vecgrp)
	stats["gauge1"].Inc()
*/

// Deprecated: Use RegisterLabeledGauge
func (m *Metrics) RegisterGaugeVec(opt CounterOpts, labelNames []string, subsytem string) GaugeVec {
	globalLock.Lock()
	defer globalLock.Unlock()
	opt.Namespace = m.Namespace
	opt.Subsystem = subsytem
	id := m.getFullName(prometheus.Opts(opt), []string{})
	if _, ok := cache_allgaugevects[id]; !ok {
		Logger.Info("Register new gauge vector with opt: %v labelNames: %v", opt, labelNames)
		entry := GaugeVec{}
		entry.Opts = opt
		entry.Labels = labelNames
		entry.Vec = promauto.NewGaugeVec(prometheus.GaugeOpts(entry.Opts), entry.Labels)
		cache_allgaugevects[id] = entry
	}
	entry := cache_allgaugevects[id]
	if strSliceCompare(entry.Labels, labelNames) == false {
		Logger.Warn("id:%s cached gauge vec labels dont match %v != %v", id, entry.Labels, labelNames)
	}
	return entry
}

// Deprecated: Use RegisterLabeledGaugeGroup
func (m *Metrics) RegisterGaugeVecGroup(opts []CounterOpts, labelNames []string, subsytem string) map[string]GaugeVec {
	c := make(map[string]GaugeVec)
	for _, opt := range opts {
		c[opt.Name] = m.RegisterGaugeVec(opt, labelNames, subsytem)
	}
	return c
}

// Deprecated: Use RegisterLabeledGauge
func (m *Metrics) GetGaugeFromVect(labelValues []string, vec GaugeVec) Gauge {
	globalLock.Lock()
	defer globalLock.Unlock()
	id := m.getFullName(prometheus.Opts(vec.Opts), labelValues)
	if _, ok := cache_allgauges[id]; !ok {
		Logger.Info("Register new gauge from vector with opts: %v labelValues: %v", vec.Opts, labelValues)
		cache_allgauges[id] = vec.Vec.WithLabelValues(labelValues...)
	}
	return cache_allgauges[id]
}

// Deprecated: Use RegisterLabeledGaugeGroup
func (m *Metrics) GetGaugeGroupFromVects(labelValues []string, vects ...map[string]GaugeVec) map[string]Gauge {
	c := make(map[string]Gauge)
	for _, vect := range vects {
		for name, vec := range vect {
			c[name] = m.GetGaugeFromVect(labelValues, vec)
		}
	}
	return c
}

// Deprecated: Use RegisterLabeledGaugeGroup
func (m *Metrics) GetGaugeGroupFromVectsWithPrefix(prefix string, labelValues []string, vects ...map[string]GaugeVec) map[string]Gauge {
	c := make(map[string]Gauge)
	for _, vect := range vects {
		for name, vec := range vect {
			c[prefix+name] = m.GetGaugeFromVect(labelValues, vec)
		}
	}
	return c
}
