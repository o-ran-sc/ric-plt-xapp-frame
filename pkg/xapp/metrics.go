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
	"sync"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//-----------------------------------------------------------------------------
//
//-----------------------------------------------------------------------------
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

//-----------------------------------------------------------------------------
//
//-----------------------------------------------------------------------------
type MetricGroupsCacheCounterRegisterer interface {
	RegisterCounter(CounterOpts) Counter
}

type MetricGroupsCacheCounterRegistererFunc func(CounterOpts) Counter

func (fn MetricGroupsCacheCounterRegistererFunc) RegisterCounter(copts CounterOpts) Counter {
	return fn(copts)
}

//-----------------------------------------------------------------------------
//
//-----------------------------------------------------------------------------
type MetricGroupsCacheGaugeRegisterer interface {
	RegisterGauge(CounterOpts) Gauge
}

type MetricGroupsCacheGaugeRegistererFunc func(CounterOpts) Gauge

func (fn MetricGroupsCacheGaugeRegistererFunc) RegisterGauge(copts CounterOpts) Gauge {
	return fn(copts)
}

//-----------------------------------------------------------------------------
//
//-----------------------------------------------------------------------------

type MetricGroupsCache struct {
	sync.RWMutex //This is for map locking
	counters     map[string]Counter
	gauges       map[string]Gauge
	regcnt       MetricGroupsCacheCounterRegisterer
	reggau       MetricGroupsCacheGaugeRegisterer
}

func (met *MetricGroupsCache) Registerer(regcnt MetricGroupsCacheCounterRegisterer, reggau MetricGroupsCacheGaugeRegisterer) {
	met.regcnt = regcnt
	met.reggau = reggau
}

func (met *MetricGroupsCache) cReg(metric string) Counter {
	if met.regcnt != nil {
		cntr := met.regcnt.RegisterCounter(CounterOpts{Name: metric, Help: "Amount of " + metric + "(auto)"})
		met.counters[metric] = cntr
		return cntr
	}
	return nil
}
func (met *MetricGroupsCache) gReg(metric string) Gauge {
	if met.reggau != nil {
		gaug := met.reggau.RegisterGauge(CounterOpts{Name: metric, Help: "Amount of " + metric + "(auto)"})
		met.gauges[metric] = gaug
		return gaug
	}
	return nil
}

func (met *MetricGroupsCache) CIs(metric string) bool {
	met.Lock()
	defer met.Unlock()
	_, ok := met.counters[metric]
	return ok
}

func (met *MetricGroupsCache) CGet(metric string) Counter {
	met.Lock()
	defer met.Unlock()
	cntr, ok := met.counters[metric]
	if !ok {
		cntr = met.cReg(metric)
	}
	return cntr
}

func (met *MetricGroupsCache) CInc(metric string) {
	met.Lock()
	defer met.Unlock()
	cntr, ok := met.counters[metric]
	if !ok {
		cntr = met.cReg(metric)
	}
	cntr.Inc()
}

func (met *MetricGroupsCache) CAdd(metric string, val float64) {
	met.Lock()
	defer met.Unlock()
	cntr, ok := met.counters[metric]
	if !ok {
		cntr = met.cReg(metric)
	}
	cntr.Add(val)
}

func (met *MetricGroupsCache) GIs(metric string) bool {
	met.Lock()
	defer met.Unlock()
	_, ok := met.gauges[metric]
	return ok
}

func (met *MetricGroupsCache) GGet(metric string) Gauge {
	met.Lock()
	defer met.Unlock()
	gaug, ok := met.gauges[metric]
	if !ok {
		gaug = met.gReg(metric)
	}
	return gaug
}

func (met *MetricGroupsCache) GSet(metric string, val float64) {
	met.Lock()
	defer met.Unlock()
	gaug, ok := met.gauges[metric]
	if !ok {
		gaug = met.gReg(metric)
	}
	gaug.Set(val)
}

func (met *MetricGroupsCache) GAdd(metric string, val float64) {
	met.Lock()
	defer met.Unlock()
	gaug, ok := met.gauges[metric]
	if !ok {
		gaug = met.gReg(metric)
	}
	gaug.Add(val)
}

func (met *MetricGroupsCache) GInc(metric string) {
	met.Lock()
	defer met.Unlock()
	gaug, ok := met.gauges[metric]
	if !ok {
		gaug = met.gReg(metric)
	}
	gaug.Inc()
}

func (met *MetricGroupsCache) GDec(metric string) {
	met.Lock()
	defer met.Unlock()
	gaug, ok := met.gauges[metric]
	if !ok {
		gaug = met.gReg(metric)
	}
	gaug.Dec()
}

func (met *MetricGroupsCache) combineCounterGroupsWithPrefix(prefix string, srcs ...map[string]Counter) {
	for _, src := range srcs {
		for k, v := range src {
			met.counters[prefix+k] = v
		}
	}
}

func (met *MetricGroupsCache) CombineCounterGroupsWithPrefix(prefix string, srcs ...map[string]Counter) {
	met.Lock()
	defer met.Unlock()
	met.combineCounterGroupsWithPrefix(prefix, srcs...)
}

func (met *MetricGroupsCache) CombineCounterGroups(srcs ...map[string]Counter) {
	met.Lock()
	defer met.Unlock()
	met.combineCounterGroupsWithPrefix("", srcs...)
}

func (met *MetricGroupsCache) combineGaugeGroupsWithPrefix(prefix string, srcs ...map[string]Gauge) {
	for _, src := range srcs {
		for k, v := range src {
			met.gauges[prefix+k] = v
		}
	}
}

func (met *MetricGroupsCache) CombineGaugeGroupsWithPrefix(prefix string, srcs ...map[string]Gauge) {
	met.Lock()
	defer met.Unlock()
	met.combineGaugeGroupsWithPrefix(prefix, srcs...)
}

func (met *MetricGroupsCache) CombineGaugeGroups(srcs ...map[string]Gauge) {
	met.Lock()
	defer met.Unlock()
	met.combineGaugeGroupsWithPrefix("", srcs...)
}

func NewMetricGroupsCache() *MetricGroupsCache {
	entry := &MetricGroupsCache{}
	entry.counters = make(map[string]Counter)
	entry.gauges = make(map[string]Gauge)
	entry.regcnt = nil
	entry.reggau = nil
	return entry
}

func NewMetricGroupsCacheWithRegisterers(regcnt MetricGroupsCacheCounterRegisterer, reggau MetricGroupsCacheGaugeRegisterer) *MetricGroupsCache {
	entry := NewMetricGroupsCache()
	entry.regcnt = regcnt
	entry.reggau = reggau
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
	if _, ok := cache_allcountervects[id]; ok {
		Logger.Warn("Register new counter with opts: %v, name conflicts existing counter vector", opts)
		return nil
	}
	if _, ok := cache_allcounters[id]; !ok {
		Logger.Debug("Register new counter with opts: %v", opts)
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
	if _, ok := cache_allcounters[vecid]; ok {
		Logger.Warn("Register new counter vector with opts: %v labelNames: %v, name conflicts existing counter", opts, labelNames)
		return nil
	}
	if _, ok := cache_allcountervects[vecid]; !ok {
		Logger.Debug("Register new counter vector with opts: %v labelNames: %v", opts, labelNames)
		entry := CounterVec{}
		entry.Opts = opts
		entry.Labels = labelNames
		entry.Vec = promauto.NewCounterVec(prometheus.CounterOpts(entry.Opts), entry.Labels)
		cache_allcountervects[vecid] = entry
	}
	entry := cache_allcountervects[vecid]
	if strSliceCompare(entry.Labels, labelNames) == false {
		Logger.Warn("id:%s cached counter vec labels dont match %v != %v", vecid, entry.Labels, labelNames)
		return nil
	}
	valid := m.getFullName(prometheus.Opts(entry.Opts), labelValues)
	if _, ok := cache_allcounters[valid]; !ok {
		Logger.Debug("Register new counter from vector with opts: %v labelValues: %v", entry.Opts, labelValues)
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
	if _, ok := cache_allgaugevects[id]; ok {
		Logger.Warn("Register new gauge with opts: %v, name conflicts existing gauge vector", opts)
		return nil
	}
	if _, ok := cache_allgauges[id]; !ok {
		Logger.Debug("Register new gauge with opts: %v", opts)
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
func (m *Metrics) RegisterLabeledGauge(opts CounterOpts, labelNames []string, labelValues []string, subsytem string) Gauge {
	globalLock.Lock()
	defer globalLock.Unlock()
	opts.Namespace = m.Namespace
	opts.Subsystem = subsytem
	vecid := m.getFullName(prometheus.Opts(opts), []string{})
	if _, ok := cache_allgauges[vecid]; ok {
		Logger.Warn("Register new gauge vector with opts: %v labelNames: %v, name conflicts existing counter", opts, labelNames)
		return nil
	}
	if _, ok := cache_allgaugevects[vecid]; !ok {
		Logger.Debug("Register new gauge vector with opts: %v labelNames: %v", opts, labelNames)
		entry := GaugeVec{}
		entry.Opts = opts
		entry.Labels = labelNames
		entry.Vec = promauto.NewGaugeVec(prometheus.GaugeOpts(entry.Opts), entry.Labels)
		cache_allgaugevects[vecid] = entry
	}
	entry := cache_allgaugevects[vecid]
	if strSliceCompare(entry.Labels, labelNames) == false {
		Logger.Warn("id:%s cached gauge vec labels dont match %v != %v", vecid, entry.Labels, labelNames)
		return nil
	}
	valid := m.getFullName(prometheus.Opts(entry.Opts), labelValues)
	if _, ok := cache_allgauges[valid]; !ok {
		Logger.Debug("Register new gauge from vector with opts: %v labelValues: %v", entry.Opts, labelValues)
		cache_allgauges[valid] = entry.Vec.WithLabelValues(labelValues...)
	}
	return cache_allgauges[valid]
}

//
//
//
func (m *Metrics) RegisterLabeledGaugeGroup(optsgroup []CounterOpts, labelNames []string, labelValues []string, subsytem string) map[string]Gauge {
	c := make(map[string]Gauge)
	for _, opts := range optsgroup {
		c[opts.Name] = m.RegisterLabeledGauge(opts, labelNames, labelValues, subsytem)
	}
	return c
}
