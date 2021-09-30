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

var mCVect CounterVec
var mGVect GaugeVec

var mCGroupVect map[string]CounterVec
var mGGroupVect map[string]GaugeVec

func TestMetricSetup(t *testing.T) {
	mCVect = Metric.RegisterCounterVec(CounterOpts{Name: "counter1", Help: "counter1"}, []string{"name", "event"}, "SUBSYSTEM0")

	mCGroupVect = Metric.RegisterCounterVecGroup(
		[]CounterOpts{
			{Name: "counter1", Help: "counter1"},
		},
		[]string{"name", "event"},
		"SUBSYSTEM1")

	mGVect = Metric.RegisterGaugeVec(CounterOpts{Name: "gauge1", Help: "gauge1"}, []string{"name", "event"}, "SUBSYSTEM0")

	mGGroupVect = Metric.RegisterGaugeVecGroup(
		[]CounterOpts{
			{Name: "gauge1", Help: "gauge1"},
		},
		[]string{"name", "event"},
		"SUBSYSTEM1")

	tmpCVect := Metric.RegisterCounterVec(CounterOpts{Name: "counter1", Help: "counter1"}, []string{"name", "event"}, "SUBSYSTEM0")

	if tmpCVect.Vec != mCVect.Vec {
		t.Errorf("tmpCVect not same than mCVect. cache not working?")
	}

	tmpGVect := Metric.RegisterGaugeVec(CounterOpts{Name: "gauge1", Help: "gauge1"}, []string{"name", "event"}, "SUBSYSTEM0")

	if tmpGVect.Vec != mGVect.Vec {
		t.Errorf("tmpGVect not same than mGVect. cache not working?")
	}

	Metric.RegisterCounterVec(CounterOpts{Name: "counter1", Help: "counter1"}, []string{"name", "eventMismatch"}, "SUBSYSTEM0")
	Metric.RegisterGaugeVec(CounterOpts{Name: "gauge1", Help: "gauge1"}, []string{"name", "eventMismatch"}, "SUBSYSTEM0")

}

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

func TestMetricCounterVector(t *testing.T) {
	//
	//
	c_1_1 := Metric.GetCounterFromVect([]string{"name1", "event1"}, mCVect)
	c_1_1.Inc()
	c_1_2 := Metric.GetCounterFromVect([]string{"name1", "event1"}, mCVect)
	c_1_2.Inc()
	if c_1_1 != c_1_2 {
		t.Errorf("c_1_1 not same than c_1_2. cache not working?")
	}
	//
	//
	c_2_1 := Metric.GetCounterFromVect([]string{"name1", "event2"}, mCVect)
	c_2_1.Inc()
	c_2_2 := Metric.GetCounterFromVect([]string{"name1", "event2"}, mCVect)
	c_2_2.Inc()
	if c_2_1 != c_2_2 {
		t.Errorf("c_2_1 not same than c_2_2. cache not working?")
	}
	if c_1_1 == c_2_1 {
		t.Errorf("c_1_1 same than c_2_1. what?")
	}

}

func TestMetricCounterGroupVector(t *testing.T) {
	//
	//
	c_grp1 := Metric.GetCounterGroupFromVects([]string{"name1", "event1"}, mCGroupVect)
	if _, ok := c_grp1["counter1"]; ok == false {
		t.Errorf("c_grp1 counter1 not exists")
	}
	c_grp1["counter1"].Inc()

	//
	//
	c_grp2 := Metric.GetCounterGroupFromVects([]string{"name1", "event2"}, mCGroupVect)
	if _, ok := c_grp2["counter1"]; ok == false {
		t.Errorf("c_grp2 counter1 not exists")
	}
	c_grp2["counter1"].Inc()
}

func TestMetricGaugeVector(t *testing.T) {
	//
	//
	c_1_1 := Metric.GetGaugeFromVect([]string{"name1", "event1"}, mGVect)
	c_1_1.Inc()
	c_1_2 := Metric.GetGaugeFromVect([]string{"name1", "event1"}, mGVect)
	c_1_2.Inc()
	if c_1_1 != c_1_2 {
		t.Errorf("c_1_1 not same than c_1_2. cache not working?")
	}
	//
	//
	c_2_1 := Metric.GetGaugeFromVect([]string{"name1", "event2"}, mGVect)
	c_2_1.Inc()
	c_2_2 := Metric.GetGaugeFromVect([]string{"name1", "event2"}, mGVect)
	c_2_2.Inc()
	if c_2_1 != c_2_2 {
		t.Errorf("c_2_1 not same than c_2_2. cache not working?")
	}
	if c_1_1 == c_2_1 {
		t.Errorf("c_1_1 same than c_2_1. what?")
	}
}

func TestMetricGaugeGroupVector(t *testing.T) {
	//
	//
	g_grp1 := Metric.GetGaugeGroupFromVects([]string{"name1", "event1"}, mGGroupVect)
	if _, ok := g_grp1["gauge1"]; ok == false {
		t.Errorf("g_grp1 gauge1 not exists")
	}
	g_grp1["gauge1"].Inc()

	//
	//
	g_grp2 := Metric.GetGaugeGroupFromVects([]string{"name1", "event2"}, mGGroupVect)
	if _, ok := g_grp2["gauge1"]; ok == false {
		t.Errorf("g_grp2 gauge1 not exists")
	}
	g_grp2["gauge1"].Inc()
}

func TestMetricCounterGroupVectorPrefix(t *testing.T) {
	//
	//
	c_grp1 := Metric.GetCounterGroupFromVectsWithPrefix("event1_", []string{"name1", "event1"}, mCGroupVect)
	if _, ok := c_grp1["event1_counter1"]; ok == false {
		t.Errorf("c_grp1 event1_counter1 not exists")
	}
	c_grp1["event1_counter1"].Inc()

	//
	//
	c_grp2 := Metric.GetCounterGroupFromVectsWithPrefix("event2_", []string{"name1", "event2"}, mCGroupVect)
	if _, ok := c_grp2["event2_counter1"]; ok == false {
		t.Errorf("c_grp2 event2_counter1 not exists")
	}
	c_grp2["event2_counter1"].Inc()

	//
	//
	m_grp := NewMetricGroupsCache()
	m_grp.CombineCounterGroups(c_grp1, c_grp2)

	//
	//
	if m_grp.CIs("event1_counter1") == false {
		t.Errorf("m_grp event1_counter1 not exists")
	}
	m_grp.CInc("event1_counter1")

	//
	//
	if m_grp.CIs("event2_counter1") == false {
		t.Errorf("m_grp event2_counter1 not exists")
	}

	m_grp.CAdd("event2_counter1", 1)
	m_grp.CGet("event2_counter1")
}

func TestMetricGaugeGroupVectorPrefix(t *testing.T) {
	//
	//
	g_grp1 := Metric.GetGaugeGroupFromVectsWithPrefix("event1_", []string{"name1", "event1"}, mGGroupVect)
	if _, ok := g_grp1["event1_gauge1"]; ok == false {
		t.Errorf("g_grp1 event1_gauge1 not exists")
	}
	g_grp1["event1_gauge1"].Inc()

	//
	//
	g_grp2 := Metric.GetGaugeGroupFromVectsWithPrefix("event2_", []string{"name1", "event2"}, mGGroupVect)
	if _, ok := g_grp2["event2_gauge1"]; ok == false {
		t.Errorf("g_grp2 event2_gauge1 not exists")
	}
	g_grp2["event2_gauge1"].Inc()

	m_grp := NewMetricGroupsCache()
	m_grp.CombineGaugeGroups(g_grp1, g_grp2)

	//
	//
	if m_grp.GIs("event1_gauge1") == false {
		t.Errorf("m_grp event1_gauge1 not exists")
	}
	m_grp.GInc("event1_gauge1")

	//
	//
	if m_grp.GIs("event2_gauge1") == false {
		t.Errorf("m_grp event2_gauge1 not exists")
	}
	m_grp.GInc("event2_gauge1")

	m_grp.GGet("event2_gauge1")
	m_grp.GDec("event2_gauge1")
	m_grp.GSet("event2_gauge1", 1)
}

func TestMetricGroupCache(t *testing.T) {
	//
	//
	c_grp1 := Metric.GetCounterGroupFromVectsWithPrefix("event1_", []string{"name1", "event1"}, mCGroupVect)
	if _, ok := c_grp1["event1_counter1"]; ok == false {
		t.Errorf("c_grp1 event1_counter1 not exists")
	}
	c_grp1["event1_counter1"].Inc()

	//
	//
	c_grp2 := Metric.GetCounterGroupFromVects([]string{"name1", "event2"}, mCGroupVect)
	if _, ok := c_grp2["counter1"]; ok == false {
		t.Errorf("c_grp2 counter1 not exists")
	}
	c_grp2["counter1"].Inc()

	//
	//
	g_grp1 := Metric.GetGaugeGroupFromVectsWithPrefix("event1_", []string{"name1", "event1"}, mGGroupVect)
	if _, ok := g_grp1["event1_gauge1"]; ok == false {
		t.Errorf("g_grp1 event1_gauge1 not exists")
	}
	g_grp1["event1_gauge1"].Inc()

	//
	//
	g_grp2 := Metric.GetGaugeGroupFromVects([]string{"name1", "event2"}, mGGroupVect)
	if _, ok := g_grp2["gauge1"]; ok == false {
		t.Errorf("g_grp2 gauge1 not exists")
	}
	g_grp2["gauge1"].Inc()

	//
	//
	m_grp := NewMetricGroupsCache()
	m_grp.CombineCounterGroups(c_grp1)
	m_grp.CombineCounterGroupsWithPrefix("event2_", c_grp2)
	m_grp.CombineGaugeGroups(g_grp1)
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
