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

func (met *MetricGroupsCache) CombineCounterGroups(srcs ...map[string]Counter) {
	met.Lock()
	defer met.Unlock()
	for _, src := range srcs {
		for k, v := range src {
			met.counters[k] = v
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

func init() {
	cache_allcounters = make(map[string]Counter)
	cache_allgauges = make(map[string]Gauge)
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

/*
 * Handling counters
 */
func (m *Metrics) registerCounter(opts CounterOpts) Counter {
	Logger.Info("Register new counter with opts: %v", opts)
	return promauto.NewCounter(prometheus.CounterOpts(opts))
}

func (m *Metrics) RegisterCounterGroup(opts []CounterOpts, subsytem string) (c map[string]Counter) {
	globalLock.Lock()
	defer globalLock.Unlock()
	c = make(map[string]Counter)
	for _, opt := range opts {
		opt.Namespace = m.Namespace
		opt.Subsystem = subsytem

		id := m.getFullName(prometheus.Opts(opt), []string{})
		if _, ok := cache_allcounters[id]; !ok {
			cache_allcounters[id] = m.registerCounter(opt)
		}

		c[opt.Name] = cache_allcounters[id]
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
	globalLock.Lock()
	defer globalLock.Unlock()
	c = make(map[string]Gauge)
	for _, opt := range opts {
		opt.Namespace = m.Namespace
		opt.Subsystem = subsytem

		id := m.getFullName(prometheus.Opts(opt), []string{})
		if _, ok := cache_allgauges[id]; !ok {
			cache_allgauges[id] = m.registerGauge(opt)
		}

		c[opt.Name] = cache_allgauges[id]
	}

	return
}

/*
 * Handling counter vectors
 *
 * Example:

	vec := Metric.RegisterCounterVecGroup(
		[]CounterOpts{
			{Name: "counter1", Help: "counter1"},
			{Name: "counter2", Help: "counter2"},
		},
		[]string{"host"},
		"SUBSYSTEM")

	stat:=Metric.GetCounterGroupFromVects([]string{"localhost:8888"}, vec)

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

func (m *Metrics) GetCounterGroupFromVectsWithPrefix(prefix string, labels []string, vects ...map[string]CounterVec) (c map[string]Counter) {
	globalLock.Lock()
	defer globalLock.Unlock()
	c = make(map[string]Counter)
	for _, vec := range vects {
		for name, opt := range vec {

			id := m.getFullName(prometheus.Opts(opt.Opts), labels)
			if _, ok := cache_allcounters[id]; !ok {
				Logger.Info("Register new counter from vector with opts: %v labels: %v prefix: %s", opt.Opts, labels, prefix)
				cache_allcounters[id] = opt.Vec.WithLabelValues(labels...)
			}
			c[prefix+name] = cache_allcounters[id]
		}
	}
	return
}

func (m *Metrics) GetCounterGroupFromVects(labels []string, vects ...map[string]CounterVec) (c map[string]Counter) {
	return m.GetCounterGroupFromVectsWithPrefix("", labels, vects...)
}

/*
 * Handling gauge vectors
 *
 * Example:

	vec := Metric.RegisterGaugeVecGroup(
		[]CounterOpts{
			{Name: "gauge1", Help: "gauge1"},
			{Name: "gauge2", Help: "gauge2"},
		},
		[]string{"host"},
		"SUBSYSTEM")

	stat:=Metric.GetGaugeGroupFromVects([]string{"localhost:8888"},vec)

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

func (m *Metrics) GetGaugeGroupFromVectsWithPrefix(prefix string, labels []string, vects ...map[string]GaugeVec) (c map[string]Gauge) {
	globalLock.Lock()
	defer globalLock.Unlock()
	c = make(map[string]Gauge)
	for _, vec := range vects {
		for name, opt := range vec {

			id := m.getFullName(prometheus.Opts(opt.Opts), labels)
			if _, ok := cache_allgauges[id]; !ok {
				Logger.Info("Register new gauge from vector with opts: %v labels: %v prefix: %s", opt.Opts, labels, prefix)
				cache_allgauges[id] = opt.Vec.WithLabelValues(labels...)
			}
			c[prefix+name] = cache_allgauges[id]
		}
	}
	return
}

func (m *Metrics) GetGaugeGroupFromVects(labels []string, vects ...map[string]GaugeVec) (c map[string]Gauge) {
	return m.GetGaugeGroupFromVectsWithPrefix("", labels, vects...)
}
