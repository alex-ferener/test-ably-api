package main

import (
	"fmt"
	"strings"
	"sync"
	"text/tabwriter"
	"time"
)

type Timing struct {
	sync.RWMutex
	values map[string]map[string][]int64
}

func newTiming() *Timing {
	return &Timing{
		values: make(map[string]map[string][]int64),
	}
}

func (t *Timing) Put(instanceId1, instanceId2 string, value int64) {
	// sort the map in order to build the pairs
	// eg: instance1 -> instance2 == instance2 -> instance1
	if instanceId1 > instanceId2 {
		instanceId1, instanceId2 = instanceId2, instanceId1
	}

	t.Lock()
	defer t.Unlock()

	_, ok := t.values[instanceId1]
	if !ok {
		t.values[instanceId1] = map[string][]int64{
			instanceId2: {},
		}
	}

	_, ok = t.values[instanceId1][instanceId2]
	if !ok {
		t.values[instanceId1][instanceId2] = []int64{}
	}

	t.values[instanceId1][instanceId2] = append(t.values[instanceId1][instanceId2], value)
}

func (t *Timing) String() string {
	tbl := new(strings.Builder)
	w := new(tabwriter.Writer)
	w.Init(tbl, 8, 8, 1, '\t', 0)

	tableFormat := "%s\t%s\t%s\t%s\t%s\t\n"
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, tableFormat, "InstanceId1", "InstanceId2", "Min (ms)", "Max (ms)", "Avg (ms)")
	fmt.Fprintf(w, tableFormat, "-----------", "-----------", "--------", "--------", "--------")

	for instanceId1, innerMap := range t.values {
		for instanceId2, timings := range innerMap {
			min, max, avg := MinMaxAvg(timings)
			fmt.Fprintf(w, tableFormat, instanceId1, instanceId2, formatTiming(min), formatTiming(max), formatTiming(int64(avg)))
		}
	}
	fmt.Fprintf(w, "\n")
	w.Flush()

	return tbl.String()
}

func formatTiming(timing int64) string {
	return fmt.Sprintf("%d", time.Duration(timing)/time.Millisecond)
}
