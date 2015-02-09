package checks

import (
	"fmt"
	"github.com/MarshalMediaGroup/goshin"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Device struct {
	Name    string
	Metrics map[string]float64
}

type IOStat struct {
	stats    chan []Device
	interval time.Duration
}

func (m *IOStat) parse(output string) ([]Device, error) {
	lines := strings.Split(output, "\n")
	headersLine := lines[2]
	deviceLines := lines[3 : len(lines)-2]
	devices := make([]Device, 0)
	headers := strings.Fields(headersLine)[2:]

	for _, line := range deviceLines {
		parts := strings.Fields(line)
		name := parts[0]
		values := parts[2:]
		device := Device{
			Name:    name,
			Metrics: make(map[string]float64),
		}
		for i, header := range headers {
			value, err := strconv.ParseFloat(strings.Replace(values[i], ",", ".", -1), 64)
			if err != nil {
				return nil, err
			}
			device.Metrics[header] = value
		}
		devices = append(devices, device)
	}
	return devices, nil
}

func (m *IOStat) readOnce() {
	out, err := exec.Command("sh", "-c", fmt.Sprintf("iostat -xymd %d 1", m.interval/time.Second)).Output()
	if err != nil {
		log.Print(err.Error())
		return
	}
	devices, err := m.parse(string(out))
	if err != nil {
		log.Print(err.Error())
		return
	}
	m.stats <- devices
}

func (m *IOStat) readStats() {
	m.readOnce()
	ticker := time.NewTicker(m.interval)
	for _ = range ticker.C {
		m.readOnce()
	}
}

func (m *IOStat) Collect(queue chan *goshin.Metric) {
	devices := <-m.stats
	for _, device := range devices {
		for metricName, metricValue := range device.Metrics {
			metric := goshin.NewMetric()
			metric.Service = fmt.Sprintf("iostat.device.%s.metric.%s", device.Name, metricName)
			metric.Value = metricValue
			queue <- metric
		}
	}
}

func NewIOStat(interval time.Duration) *IOStat {
	stat := &IOStat{
		stats:    make(chan []Device, 1),
		interval: interval,
	}
	go stat.readStats()
	return stat
}
