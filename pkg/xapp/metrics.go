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
	Counters map[string]Counter
	Gauges   map[string]Gauge
}

func (met *MetricGroupsCache) CInc(metric string) {
	met.Counters[metric].Inc()
}

func (met *MetricGroupsCache) CAdd(metric string, val float64) {
	met.Counters[metric].Add(val)
}

func (met *MetricGroupsCache) GSet(metric string, val float64) {
	met.Gauges[metric].Set(val)
}

//-----------------------------------------------------------------------------
//
//-----------------------------------------------------------------------------
type Metrics struct {
	Namespace            string
	MetricGroupsCacheMap map[string]*MetricGroupsCache
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

	return &Metrics{Namespace: namespace, MetricGroupsCacheMap: make(map[string]*MetricGroupsCache)}
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
	c = make(map[string]Counter)
	for _, vec := range vects {
		for name, opt := range vec {
			c[prefix+name] = opt.Vec.WithLabelValues(labels...)
			Logger.Info("Register new counter for vector with opts: %v labels: %v", opt.Opts, labels)
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
	c = make(map[string]Gauge)
	for _, vec := range vects {
		for name, opt := range vec {
			c[prefix+name] = opt.Vec.WithLabelValues(labels...)
			Logger.Info("Register new gauge for vector with opts: %v labels: %v", opt.Opts, labels)
		}
	}
	return
}

func (m *Metrics) GetGaugeGroupFromVects(labels []string, vects ...map[string]GaugeVec) (c map[string]Gauge) {
	return m.GetGaugeGroupFromVectsWithPrefix("", labels, vects...)

}

/*
 *
 */
func (m *Metrics) CombineCounterGroups(srcs ...map[string]Counter) map[string]Counter {
	trg := make(map[string]Counter)
	for _, src := range srcs {
		for k, v := range src {
			trg[k] = v
		}
	}
	return trg
}

func (m *Metrics) CombineGaugeGroups(srcs ...map[string]Gauge) map[string]Gauge {
	trg := make(map[string]Gauge)
	for _, src := range srcs {
		for k, v := range src {
			trg[k] = v
		}
	}
	return trg
}

/*
 *
 */
func (m *Metrics) GroupCacheGet(id string) *MetricGroupsCache {
	entry, ok := m.MetricGroupsCacheMap[id]
	if ok == false {
		return nil
	}
	return entry
}

func (m *Metrics) GroupCacheAddCounters(id string, vals map[string]Counter) {
	entry, ok := m.MetricGroupsCacheMap[id]
	if ok == false {
		entry = &MetricGroupsCache{}
		m.MetricGroupsCacheMap[id] = entry
	}
	m.MetricGroupsCacheMap[id].Counters = m.CombineCounterGroups(m.MetricGroupsCacheMap[id].Counters, vals)
}

func (m *Metrics) GroupCacheAddGauges(id string, vals map[string]Gauge) {
	entry, ok := m.MetricGroupsCacheMap[id]
	if ok == false {
		entry = &MetricGroupsCache{}
		m.MetricGroupsCacheMap[id] = entry
	}
	m.MetricGroupsCacheMap[id].Gauges = m.CombineGaugeGroups(m.MetricGroupsCacheMap[id].Gauges, vals)
}
