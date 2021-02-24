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

var mCVect map[string]CounterVec
var mGVect map[string]GaugeVec
var mGGroup map[string]Gauge

func TestMetricSetup(t *testing.T) {
	mCVect = Metric.RegisterCounterVecGroup(
		[]CounterOpts{
			{Name: "counter1", Help: "counter1"},
		},
		[]string{"name", "event"},
		"SUBSYSTEM")

	mGVect = Metric.RegisterGaugeVecGroup(
		[]CounterOpts{
			{Name: "counter2", Help: "counter2"},
		},
		[]string{"name", "event"},
		"SUBSYSTEM")

	mGGroup = Metric.RegisterGaugeGroup(
		[]CounterOpts{
			{Name: "counter3", Help: "counter3"},
		},
		"SUBSYSTEM2")
}

func TestMetricCounter(t *testing.T) {
	var TestCounterOpts = []CounterOpts{
		{Name: "Blaah1", Help: "Blaah1"},
		{Name: "Blaah2", Help: "Blaah2"},
		{Name: "Blaah3", Help: "Blaah3"},
		{Name: "Blaah4", Help: "Blaah4"},
	}

	ret1 := Metric.RegisterCounterGroup(TestCounterOpts, "TestMetricCounter")

	if len(ret1) == 0 {
		t.Errorf("ret1 counter group is empty")
	}

	ret2 := Metric.RegisterCounterGroup(TestCounterOpts, "TestMetricCounter")

	if len(ret2) == 0 {
		t.Errorf("ret2 counter group is empty")
	}

	if len(ret1) != len(ret2) {
		t.Errorf("ret1 len %d differs from ret2 len %d", len(ret1), len(ret2))
	}
}

func TestMetricCounterVector(t *testing.T) {
	//
	//
	c_grp1 := Metric.GetCounterGroupFromVects([]string{"name1", "event1"}, mCVect)
	if _, ok := c_grp1["counter1"]; ok == false {
		t.Errorf("c_grp1 counter1 not exists")
	}
	c_grp1["counter1"].Inc()

	//
	//
	c_grp2 := Metric.GetCounterGroupFromVects([]string{"name1", "event2"}, mCVect)
	if _, ok := c_grp2["counter1"]; ok == false {
		t.Errorf("c_grp2 counter1 not exists")
	}
	c_grp2["counter1"].Inc()
}

func TestMetricGaugeVector(t *testing.T) {
	//
	//
	g_grp1 := Metric.GetGaugeGroupFromVects([]string{"name1", "event1"}, mGVect)
	if _, ok := g_grp1["counter2"]; ok == false {
		t.Errorf("g_grp1 counter2 not exists")
	}
	g_grp1["counter2"].Inc()

	//
	//
	g_grp2 := Metric.GetGaugeGroupFromVects([]string{"name1", "event2"}, mGVect)
	if _, ok := g_grp2["counter2"]; ok == false {
		t.Errorf("g_grp2 counter2 not exists")
	}
	g_grp2["counter2"].Inc()
}

func TestMetricCounterVectorPrefix(t *testing.T) {
	//
	//
	c_grp1 := Metric.GetCounterGroupFromVectsWithPrefix("event1_", []string{"name1", "event1"}, mCVect)
	if _, ok := c_grp1["event1_counter1"]; ok == false {
		t.Errorf("c_grp1 event1_counter1 not exists")
	}
	c_grp1["event1_counter1"].Inc()

	//
	//
	c_grp2 := Metric.GetCounterGroupFromVectsWithPrefix("event2_", []string{"name1", "event2"}, mCVect)
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

func TestMetricGaugeVectorPrefix(t *testing.T) {
	//
	//
	g_grp1 := Metric.GetGaugeGroupFromVectsWithPrefix("event1_", []string{"name1", "event1"}, mGVect)
	if _, ok := g_grp1["event1_counter2"]; ok == false {
		t.Errorf("g_grp1 event1_counter2 not exists")
	}
	g_grp1["event1_counter2"].Inc()

	//
	//
	g_grp2 := Metric.GetGaugeGroupFromVectsWithPrefix("event2_", []string{"name1", "event2"}, mGVect)
	if _, ok := g_grp2["event2_counter2"]; ok == false {
		t.Errorf("g_grp2 event2_counter2 not exists")
	}
	g_grp2["event2_counter2"].Inc()

	m_grp := NewMetricGroupsCache()
	m_grp.CombineGaugeGroups(g_grp1, g_grp2)

	//
	//
	if m_grp.GIs("event1_counter2") == false {
		t.Errorf("m_grp event1_counter2 not exists")
	}
	m_grp.GInc("event1_counter2")

	//
	//
	if m_grp.GIs("event2_counter2") == false {
		t.Errorf("m_grp event2_counter2 not exists")
	}
	m_grp.GInc("event2_counter2")

	m_grp.GGet("event2_counter2")
	m_grp.GDec("event2_counter2")
	m_grp.GSet("event2_counter2", 1)
}

func TestMetricGroupCache(t *testing.T) {
	//
	//
	c_grp1 := Metric.GetCounterGroupFromVectsWithPrefix("event1_", []string{"name1", "event1"}, mCVect)
	if _, ok := c_grp1["event1_counter1"]; ok == false {
		t.Errorf("c_grp1 event1_counter1 not exists")
	}
	c_grp1["event1_counter1"].Inc()

	//
	//
	c_grp2 := Metric.GetCounterGroupFromVectsWithPrefix("event2_", []string{"name1", "event2"}, mCVect)
	if _, ok := c_grp2["event2_counter1"]; ok == false {
		t.Errorf("c_grp2 event2_counter1 not exists")
	}
	c_grp2["event2_counter1"].Inc()

	//
	//
	g_grp1 := Metric.GetGaugeGroupFromVectsWithPrefix("event1_", []string{"name1", "event1"}, mGVect)
	if _, ok := g_grp1["event1_counter2"]; ok == false {
		t.Errorf("g_grp1 event1_counter2 not exists")
	}
	g_grp1["event1_counter2"].Inc()

	//
	//
	g_grp2 := Metric.GetGaugeGroupFromVectsWithPrefix("event2_", []string{"name1", "event2"}, mGVect)
	if _, ok := g_grp2["event2_counter2"]; ok == false {
		t.Errorf("g_grp2 event2_counter2 not exists")
	}
	g_grp2["event2_counter2"].Inc()

	//
	//
	m_grp := NewMetricGroupsCache()
	m_grp.CombineCounterGroups(c_grp1)
	m_grp.CombineCounterGroups(c_grp2)
	m_grp.CombineGaugeGroups(g_grp1)
	m_grp.CombineGaugeGroups(g_grp2)

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

	if m_grp.GIs("event1_counter2") == false {
		t.Errorf("m_grp.Gauges event1_counter2 not exists")
	}
	m_grp.GInc("event1_counter2")

	if m_grp.GIs("event2_counter2") == false {
		t.Errorf("m_grp.Gauges event2_counter2 not exists")
	}
	m_grp.GInc("event2_counter2")

	m_grp.CAdd("event2_counter1", 1)
	m_grp.CGet("event2_counter1")
	m_grp.GGet("event2_counter2")
	m_grp.GDec("event2_counter2")
	m_grp.GSet("event2_counter2", 1)
}
