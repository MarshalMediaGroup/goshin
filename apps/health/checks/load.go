package checks

import (
	"fmt"
	"github.com/MarshalMediaGroup/goshin"
	linuxproc "github.com/pariviere/goprocinfo/linux"
)

type LoadAverage struct {
	last1m, last5m, last15m float64
}

func (l *LoadAverage) Usage() float64 {
	loadAverage, _ := linuxproc.ReadLoadAvg("/proc/loadavg")
	cpuInfo, _ := linuxproc.ReadCPUInfo("/proc/cpuinfo")

	l.last1m = loadAverage.Last1Min / float64(cpuInfo.NumCore())
	l.last5m = loadAverage.Last5Min / float64(cpuInfo.NumCore())
	l.last15m = loadAverage.Last15Min / float64(cpuInfo.NumCore())

	return l.last1m
}

func (l *LoadAverage) Ranking() string {
	return fmt.Sprintf("1-minute load average/core is %f", l.last1m)
}

func (l *LoadAverage) Collect(queue chan *goshin.Metric) {

	metric := goshin.NewMetric()

	metric.Service = "load"

	metric.Value = l.Usage()
	metric.Description = l.Ranking()

	queue <- metric
}

func NewLoadAverage() *LoadAverage {
	return &LoadAverage{}
}
