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
	c_grp := Metric.CombineCounterGroups(c_grp1, c_grp2)

	//
	//
	if _, ok := c_grp["event1_counter1"]; ok == false {
		t.Errorf("c_grp event1_counter1 not exists")
	}
	c_grp["event1_counter1"].Inc()

	//
	//
	if _, ok := c_grp["event2_counter1"]; ok == false {
		t.Errorf("c_grp event2_counter1 not exists")
	}
	c_grp["event2_counter1"].Inc()
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

	//
	//
	g_grp := Metric.CombineGaugeGroups(g_grp1, g_grp2)

	//
	//
	if _, ok := g_grp["event1_counter2"]; ok == false {
		t.Errorf("g_grp event1_counter2 not exists")
	}
	g_grp["event1_counter2"].Inc()

	//
	//
	if _, ok := g_grp["event2_counter2"]; ok == false {
		t.Errorf("g_grp event2_counter2 not exists")
	}
	g_grp["event2_counter2"].Inc()
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
	cacheid := "CACHEID"
	entry := Metric.GroupCacheGet(cacheid)
	if entry == nil {
		Metric.GroupCacheAddCounters(cacheid, c_grp1)
		Metric.GroupCacheAddCounters(cacheid, c_grp2)
		Metric.GroupCacheAddGauges(cacheid, g_grp1)
		Metric.GroupCacheAddGauges(cacheid, g_grp2)
		entry = Metric.GroupCacheGet(cacheid)
	}

	if entry == nil {
		t.Errorf("Cache failed")
	}

	if _, ok := entry.Counters["event1_counter1"]; ok == false {
		t.Errorf("entry.Counters event1_counter1 not exists")
	}
	entry.Counters["event1_counter1"].Inc()

	if _, ok := entry.Counters["event2_counter1"]; ok == false {
		t.Errorf("entry.Counters event2_counter1 not exists")
	}
	entry.Counters["event2_counter1"].Inc()

	if _, ok := entry.Gauges["event1_counter2"]; ok == false {
		t.Errorf("entry.Gauges event1_counter2 not exists")
	}
	entry.Gauges["event1_counter2"].Inc()

	if _, ok := entry.Gauges["event2_counter2"]; ok == false {
		t.Errorf("entry.Gauges event2_counter2 not exists")
	}
	entry.Gauges["event2_counter2"].Inc()

}
