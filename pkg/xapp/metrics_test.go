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
	"testing"
)

func TestMetricCounter(t *testing.T) {
	TestCounterOpt := CounterOpts{Name: "CounterBlaah1", Help: "CounterBlaah1"}
	ret1 := Metric.RegisterCounter(TestCounterOpt, "TestMetricCounter")
	ret1.Inc()
	ret2 := Metric.RegisterCounter(TestCounterOpt, "TestMetricCounter")
	ret2.Inc()
	if ret1 != ret2 {
		t.Errorf("ret1 not same than ret2. cache not working?")
	}
}

func TestMetricCounterGroup(t *testing.T) {
	var TestCounterOpts = []CounterOpts{
		{Name: "CounterBlaah1", Help: "CounterBlaah1"},
		{Name: "CounterBlaah2", Help: "CounterBlaah2"},
		{Name: "CounterBlaah3", Help: "CounterBlaah3"},
		{Name: "CounterBlaah4", Help: "CounterBlaah4"},
	}

	ret1 := Metric.RegisterCounterGroup(TestCounterOpts, "TestMetricCounterGroup")

	if len(ret1) == 0 {
		t.Errorf("ret1 counter group is empty")
	}

	ret1["CounterBlaah1"].Inc()
	ret1["CounterBlaah2"].Inc()
	ret1["CounterBlaah3"].Inc()
	ret1["CounterBlaah4"].Inc()

	ret2 := Metric.RegisterCounterGroup(TestCounterOpts, "TestMetricCounterGroup")

	if len(ret2) == 0 {
		t.Errorf("ret2 counter group is empty")
	}

	ret2["CounterBlaah1"].Inc()
	ret2["CounterBlaah2"].Inc()
	ret2["CounterBlaah3"].Inc()
	ret2["CounterBlaah4"].Inc()

	if len(ret1) != len(ret2) {
		t.Errorf("ret1 len %d differs from ret2 len %d", len(ret1), len(ret2))
	}
}

func TestMetricLabeledCounter(t *testing.T) {
	//
	//
	c_1_1 := Metric.RegisterLabeledCounter(
		CounterOpts{Name: "counter1", Help: "counter1"},
		[]string{"name", "event"},
		[]string{"name1", "event1"},
		"SUBSYSTEML0")

	c_1_2 := Metric.RegisterLabeledCounter(
		CounterOpts{Name: "counter1", Help: "counter1"},
		[]string{"name", "event"},
		[]string{"name1", "event1"},
		"SUBSYSTEML0")

	c_1_1.Inc()
	c_1_2.Inc()
	if c_1_1 != c_1_2 {
		t.Errorf("c_1_1 not same than c_1_2. cache not working?")
	}

	//
	//
	c_2_1 := Metric.RegisterLabeledCounter(
		CounterOpts{Name: "counter1", Help: "counter1"},
		[]string{"name", "event"},
		[]string{"name1", "event2"},
		"SUBSYSTEML0")

	c_2_2 := Metric.RegisterLabeledCounter(
		CounterOpts{Name: "counter1", Help: "counter1"},
		[]string{"name", "event"},
		[]string{"name1", "event2"},
		"SUBSYSTEML0")

	c_2_1.Inc()
	c_2_2.Inc()
	if c_2_1 != c_2_2 {
		t.Errorf("c_2_1 not same than c_2_2. cache not working?")
	}

	if c_1_1 == c_2_1 {
		t.Errorf("c_1_1 same than c_2_1. what?")
	}
	if c_1_2 == c_2_2 {
		t.Errorf("c_1_2 same than c_2_2. what?")
	}

}

func TestMetricLabeledCounterMissmatch(t *testing.T) {
	Metric.RegisterLabeledCounter(
		CounterOpts{Name: "counter1", Help: "counter1"},
		[]string{"name", "event"},
		[]string{"name1", "event1"},
		"SUBSYSTEMLERRMISSMATCH")

	ret := Metric.RegisterLabeledCounter(
		CounterOpts{Name: "counter1", Help: "counter1"},
		[]string{"name", "eventmiss"},
		[]string{"name1", "event1"},
		"SUBSYSTEMLERRMISSMATCH")

	if ret != nil {
		t.Errorf("Returned counter even its labels are mismatching")
	}

	ret = Metric.RegisterLabeledCounter(
		CounterOpts{Name: "counter1", Help: "counter1"},
		[]string{"name"},
		[]string{"name1"},
		"SUBSYSTEMLERRMISSMATCH")

	if ret != nil {
		t.Errorf("Returned counter even its labels are mismatching")
	}

}

func TestMetricLabeledCounterWrongOrder(t *testing.T) {
	Metric.RegisterLabeledCounter(
		CounterOpts{Name: "counter1", Help: "counter1"},
		[]string{"name", "event"},
		[]string{"name1", "event1"},
		"SUBSYSTEMLERRWRONGORDER")

	ret := Metric.RegisterLabeledCounter(
		CounterOpts{Name: "counter1", Help: "counter1"},
		[]string{"event", "name"},
		[]string{"name1", "event1"},
		"SUBSYSTEMLERRWRONGORDER")

	if ret != nil {
		t.Errorf("Returned counter even its labels order is wrong")
	}
}

func TestMetricLabeledCounterCounterNameExists(t *testing.T) {
	Metric.RegisterCounter(
		CounterOpts{Name: "counter1", Help: "counter1"},
		"SUBSYSTEMLERRNAMEEXISTS")

	ret := Metric.RegisterLabeledCounter(
		CounterOpts{Name: "counter1", Help: "counter1"},
		[]string{"name", "event"},
		[]string{"name1", "event1"},
		"SUBSYSTEMLERRNAMEEXISTS")

	if ret != nil {
		t.Errorf("Returned labeled counter even its name conflicts with existing counter name")
	}
}

func TestMetricCounterLabeledCounterNameExists(t *testing.T) {
	Metric.RegisterLabeledCounter(
		CounterOpts{Name: "counter2", Help: "counter2"},
		[]string{"name", "event"},
		[]string{"name1", "event1"},
		"SUBSYSTEMLERRNAMEEXISTS")

	ret := Metric.RegisterCounter(
		CounterOpts{Name: "counter2", Help: "counter2"},
		"SUBSYSTEMLERRNAMEEXISTS")

	if ret != nil {
		t.Errorf("Returned counter even its name conflicts with existing labeled counter name")
	}
}

func TestMetricLabeledCounterGroup(t *testing.T) {
	//
	//
	c_grp1 := Metric.RegisterLabeledCounterGroup(
		[]CounterOpts{{Name: "counter1", Help: "counter1"}},
		[]string{"name", "event"},
		[]string{"name1", "event1"},
		"SUBSYSTEML1")

	if _, ok := c_grp1["counter1"]; ok == false {
		t.Errorf("c_grp1 counter1 not exists")
	}
	c_grp1["counter1"].Inc()

	//
	//
	c_grp2 := Metric.RegisterLabeledCounterGroup(
		[]CounterOpts{{Name: "counter1", Help: "counter1"}},
		[]string{"name", "event"},
		[]string{"name1", "event2"},
		"SUBSYSTEML1")

	if _, ok := c_grp2["counter1"]; ok == false {
		t.Errorf("c_grp2 counter1 not exists")
	}
	c_grp2["counter1"].Inc()
}

func TestMetricGauge(t *testing.T) {
	TestGaugeOpts := CounterOpts{Name: "GaugeBlaah1", Help: "GaugeBlaah1"}
	ret1 := Metric.RegisterGauge(TestGaugeOpts, "TestMetricGauge")
	ret1.Inc()
	ret2 := Metric.RegisterGauge(TestGaugeOpts, "TestMetricGauge")
	ret2.Inc()
	if ret1 != ret2 {
		t.Errorf("ret1 not same than ret2. cache not working?")
	}
}

func TestMetricGaugeGroup(t *testing.T) {
	var TestGaugeOpts = []CounterOpts{
		{Name: "GaugeBlaah1", Help: "GaugeBlaah1"},
		{Name: "GaugeBlaah2", Help: "GaugeBlaah2"},
		{Name: "GaugeBlaah3", Help: "GaugeBlaah3"},
		{Name: "GaugeBlaah4", Help: "GaugeBlaah4"},
	}

	ret1 := Metric.RegisterGaugeGroup(TestGaugeOpts, "TestMetricGaugeGroup")

	if len(ret1) == 0 {
		t.Errorf("ret1 gauge group is empty")
	}

	ret1["GaugeBlaah1"].Inc()
	ret1["GaugeBlaah2"].Inc()
	ret1["GaugeBlaah3"].Inc()
	ret1["GaugeBlaah4"].Inc()

	ret2 := Metric.RegisterGaugeGroup(TestGaugeOpts, "TestMetricGaugeGroup")

	if len(ret2) == 0 {
		t.Errorf("ret2 gauge group is empty")
	}

	ret2["GaugeBlaah1"].Inc()
	ret2["GaugeBlaah2"].Inc()
	ret2["GaugeBlaah3"].Inc()
	ret2["GaugeBlaah4"].Inc()

	if len(ret1) != len(ret2) {
		t.Errorf("ret1 len %d differs from ret2 len %d", len(ret1), len(ret2))
	}
}

func TestMetricLabeledGauge(t *testing.T) {
	//
	//
	c_1_1 := Metric.RegisterLabeledGauge(
		CounterOpts{Name: "gauge1", Help: "gauge1"},
		[]string{"name", "event"},
		[]string{"name1", "event1"},
		"SUBSYSTEML0")

	c_1_2 := Metric.RegisterLabeledGauge(
		CounterOpts{Name: "gauge1", Help: "gauge1"},
		[]string{"name", "event"},
		[]string{"name1", "event1"},
		"SUBSYSTEML0")

	c_1_1.Inc()
	c_1_2.Inc()
	if c_1_1 != c_1_2 {
		t.Errorf("c_1_1 not same than c_1_2. cache not working?")
	}

	//
	//
	c_2_1 := Metric.RegisterLabeledGauge(
		CounterOpts{Name: "gauge1", Help: "gauge1"},
		[]string{"name", "event"},
		[]string{"name1", "event2"},
		"SUBSYSTEML0")

	c_2_2 := Metric.RegisterLabeledGauge(
		CounterOpts{Name: "gauge1", Help: "gauge1"},
		[]string{"name", "event"},
		[]string{"name1", "event2"},
		"SUBSYSTEML0")

	c_2_1.Inc()
	c_2_2.Inc()
	if c_2_1 != c_2_2 {
		t.Errorf("c_2_1 not same than c_2_2. cache not working?")
	}

	if c_1_1 == c_2_1 {
		t.Errorf("c_1_1 same than c_2_1. what?")
	}
	if c_1_2 == c_2_2 {
		t.Errorf("c_1_2 same than c_2_2. what?")
	}

}

func TestMetricLabeledGaugeMissmatch(t *testing.T) {
	Metric.RegisterLabeledGauge(
		CounterOpts{Name: "gauge1", Help: "gauge1"},
		[]string{"name", "event"},
		[]string{"name1", "event1"},
		"SUBSYSTEMLERRMISSMATCH")

	ret := Metric.RegisterLabeledGauge(
		CounterOpts{Name: "gauge1", Help: "gauge1"},
		[]string{"name", "eventmiss"},
		[]string{"name1", "event1"},
		"SUBSYSTEMLERRMISSMATCH")

	if ret != nil {
		t.Errorf("Returned gauge even its labels are mismatching")
	}

	ret = Metric.RegisterLabeledGauge(
		CounterOpts{Name: "gauge1", Help: "gauge1"},
		[]string{"name"},
		[]string{"name1"},
		"SUBSYSTEMLERRMISSMATCH")

	if ret != nil {
		t.Errorf("Returned gauge even its labels are mismatching")
	}

}

func TestMetricLabeledGaugeWrongOrder(t *testing.T) {
	Metric.RegisterLabeledGauge(
		CounterOpts{Name: "gauge1", Help: "gauge1"},
		[]string{"name", "event"},
		[]string{"name1", "event1"},
		"SUBSYSTEMLERRWRONGORDER")

	ret := Metric.RegisterLabeledGauge(
		CounterOpts{Name: "gauge1", Help: "gauge1"},
		[]string{"event", "name"},
		[]string{"name1", "event1"},
		"SUBSYSTEMLERRWRONGORDER")

	if ret != nil {
		t.Errorf("Returned gauge even its labels order is wrong")
	}

}

func TestMetricLabeledGaugeGaugeNameExists(t *testing.T) {
	Metric.RegisterGauge(
		CounterOpts{Name: "gauge1", Help: "gauge1"},
		"SUBSYSTEMLERRNAMEEXISTS")

	ret := Metric.RegisterLabeledGauge(
		CounterOpts{Name: "gauge1", Help: "gauge1"},
		[]string{"name", "event"},
		[]string{"name1", "event1"},
		"SUBSYSTEMLERRNAMEEXISTS")

	if ret != nil {
		t.Errorf("Returned labeled gauge even its name conflicts with existing gauge name")
	}
}

func TestMetricGaugeLabeledGaugeNameExists(t *testing.T) {
	Metric.RegisterLabeledGauge(
		CounterOpts{Name: "gauge2", Help: "gauge2"},
		[]string{"name", "event"},
		[]string{"name1", "event1"},
		"SUBSYSTEMLERRNAMEEXISTS")

	ret := Metric.RegisterGauge(
		CounterOpts{Name: "gauge2", Help: "gauge2"},
		"SUBSYSTEMLERRNAMEEXISTS")

	if ret != nil {
		t.Errorf("Returned gauge even its name conflicts with existing labeled gauge name")
	}
}

func TestMetricLabeledGaugeGroup(t *testing.T) {
	//
	//
	g_grp1 := Metric.RegisterLabeledGaugeGroup(
		[]CounterOpts{{Name: "gauge1", Help: "gauge1"}},
		[]string{"name", "event"},
		[]string{"name1", "event1"},
		"SUBSYSTEML1")

	if _, ok := g_grp1["gauge1"]; ok == false {
		t.Errorf("g_grp1 gauge1 not exists")
	}
	g_grp1["gauge1"].Inc()

	//
	//
	g_grp2 := Metric.RegisterLabeledGaugeGroup(
		[]CounterOpts{{Name: "gauge1", Help: "gauge1"}},
		[]string{"name", "event"},
		[]string{"name1", "event2"},
		"SUBSYSTEML1")

	if _, ok := g_grp2["gauge1"]; ok == false {
		t.Errorf("g_grp2 gauge1 not exists")
	}
	g_grp2["gauge1"].Inc()
}

func TestMetricGroupCache(t *testing.T) {
	//
	//
	c_grp1 := Metric.RegisterLabeledCounterGroup(
		[]CounterOpts{{Name: "counter1", Help: "counter1"}},
		[]string{"name", "event"},
		[]string{"name1", "event1"},
		"SUBSYSTEML1")
	if _, ok := c_grp1["counter1"]; ok == false {
		t.Errorf("c_grp1 counter1 not exists")
	}
	c_grp1["counter1"].Inc()

	//
	//
	c_grp2 := Metric.RegisterLabeledCounterGroup(
		[]CounterOpts{{Name: "counter1", Help: "counter1"}},
		[]string{"name", "event"},
		[]string{"name1", "event2"},
		"SUBSYSTEML1")
	if _, ok := c_grp2["counter1"]; ok == false {
		t.Errorf("c_grp2 counter1 not exists")
	}
	c_grp2["counter1"].Inc()

	//
	//
	g_grp1 := Metric.RegisterLabeledGaugeGroup(
		[]CounterOpts{{Name: "gauge1", Help: "gauge1"}},
		[]string{"name", "event"},
		[]string{"name1", "event1"},
		"SUBSYSTEML1")
	if _, ok := g_grp1["gauge1"]; ok == false {
		t.Errorf("g_grp1 gauge1 not exists")
	}
	g_grp1["gauge1"].Inc()

	//
	//
	g_grp2 := Metric.RegisterLabeledGaugeGroup(
		[]CounterOpts{{Name: "gauge1", Help: "gauge1"}},
		[]string{"name", "event"},
		[]string{"name1", "event2"},
		"SUBSYSTEML1")
	if _, ok := g_grp2["gauge1"]; ok == false {
		t.Errorf("g_grp2 gauge1 not exists")
	}
	g_grp2["gauge1"].Inc()

	//
	//
	m_grp := NewMetricGroupsCache()
	m_grp.CombineCounterGroupsWithPrefix("event1_", c_grp1)
	m_grp.CombineCounterGroupsWithPrefix("event2_", c_grp2)
	m_grp.CombineGaugeGroupsWithPrefix("event1_", g_grp1)
	m_grp.CombineGaugeGroupsWithPrefix("event2_", g_grp2)

	if m_grp == nil {
		t.Errorf("Cache failed")
	}

	if m_grp.CIs("event1_counter1") == false {
		t.Errorf("m_grp.Counters event1_counter1 not exists")
	}
	m_grp.CInc("event1_counter1")

	if m_grp.CIs("event2_counter1") == false {
		t.Errorf("m_grp.Counters event2_counter1 not exists")
	}
	m_grp.CInc("event2_counter1")

	if m_grp.GIs("event1_gauge1") == false {
		t.Errorf("m_grp.Gauges event1_gauge1 not exists")
	}
	m_grp.GInc("event1_gauge1")

	if m_grp.GIs("event2_gauge1") == false {
		t.Errorf("m_grp.Gauges event2_gauge1 not exists")
	}
	m_grp.GInc("event2_gauge1")

	m_grp.CAdd("event2_counter1", 1)
	m_grp.CGet("event2_counter1")
	m_grp.GGet("event2_gauge1")
	m_grp.GDec("event2_gauge1")
	m_grp.GSet("event2_gauge1", 1)
}

type registerer struct{}

func (met *registerer) RegisterCounter(opts CounterOpts) Counter {
	return Metric.RegisterLabeledCounter(
		opts,
		[]string{"host", "interface"},
		[]string{"testhost", "testinterface"},
		"SUBSYSTEMAUTO")
}

func (met *registerer) RegisterGauge(opts CounterOpts) Gauge {
	return Metric.RegisterLabeledGauge(
		opts,
		[]string{"host", "interface"},
		[]string{"testhost", "testinterface"},
		"SUBSYSTEMAUTO")
}

func TestMetricCounterAutoCGetNoReg(t *testing.T) {
	m_grp := NewMetricGroupsCache()
	m_grp.CGet("cautotest1")
}

func TestMetricCounterAutoCGetFunc(t *testing.T) {
	m_grp := NewMetricGroupsCache()
	m_reg := &registerer{}
	m_grp.Registerer(MetricGroupsCacheCounterRegistererFunc(m_reg.RegisterCounter), nil)
	m_grp.CGet("cautotest1")
}

func TestMetricCounterAutoCGet(t *testing.T) {
	m_grp := NewMetricGroupsCacheWithRegisterers(&registerer{}, nil)
	m_grp.CGet("cautotest1")
}

func TestMetricCounterAutoCInc(t *testing.T) {
	m_grp := NewMetricGroupsCache()
	m_grp.Registerer(&registerer{}, nil)
	m_grp.CInc("cautotest1")
}

func TestMetricCounterAutoCAdd(t *testing.T) {
	m_grp := NewMetricGroupsCache()
	m_grp.Registerer(&registerer{}, nil)
	m_grp.CAdd("cautotest1", float64(10))
}

func TestMetricCounterAutoGGetNoReg(t *testing.T) {
	m_grp := NewMetricGroupsCache()
	m_grp.GGet("gautotest1")
}

func TestMetricCounterAutoGGetFunc(t *testing.T) {
	m_grp := NewMetricGroupsCache()
	m_reg := &registerer{}
	m_grp.Registerer(nil, MetricGroupsCacheGaugeRegistererFunc(m_reg.RegisterGauge))
	m_grp.GGet("gautotest1")
}

func TestMetricCounterAutoGGet(t *testing.T) {
	m_grp := NewMetricGroupsCacheWithRegisterers(nil, &registerer{})
	m_grp.GGet("gautotest1")
}

func TestMetricCounterAutoGInc(t *testing.T) {
	m_grp := NewMetricGroupsCache()
	m_grp.Registerer(nil, &registerer{})
	m_grp.GInc("gautotest1")
}

func TestMetricCounterAutoGSet(t *testing.T) {
	m_grp := NewMetricGroupsCache()
	m_grp.Registerer(nil, &registerer{})
	m_grp.GSet("gautotest1", float64(10))
}

func TestMetricCounterAutoGAdd(t *testing.T) {
	m_grp := NewMetricGroupsCache()
	m_grp.Registerer(nil, &registerer{})
	m_grp.GAdd("gautotest1", float64(10))
}

func TestMetricCounterAutoGDec(t *testing.T) {
	m_grp := NewMetricGroupsCache()
	m_grp.Registerer(nil, &registerer{})
	m_grp.GDec("gautotest1")
}
