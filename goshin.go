package goshin

import (
	"flag"
	"fmt"
	"github.com/bigdatadev/goryman"
	"github.com/vharitonsky/iniflags"
	"os"
	"strings"
	"time"
)

type Check interface {
	Collect(queue chan *Metric)
}

type Metric struct {
	Service, Description, State string
	Value                       interface{}
}

func NewMetric() *Metric {
	return &Metric{State: "ok"}
}

type Threshold struct {
	Warning, Critical float64
}

func NewThreshold() *Threshold {
	return &Threshold{}
}

type Goshin struct {
	Address         string
	EventHost       string
	Interval        time.Duration
	Tag             []string
	Ttl             float32
	Thresholds      map[string]*Threshold
	EnabledChecks   map[string]bool
	AvailableChecks map[string]Check
}

func NewGoshin() *Goshin {
	return &Goshin{
		Thresholds:      make(map[string]*Threshold),
		AvailableChecks: make(map[string]Check),
	}
}

func (g *Goshin) Configure() {
	var (
		hostname, _ = os.Hostname()

		hostPtr      = flag.String("host", "localhost", "Riemann host")
		portPtr      = flag.Int("port", 5555, "Riemann port")
		eventHostPtr = flag.String("event-host", hostname, "Event hostname")
		intervalPtr  = flag.Int("interval", 5, "Seconds between updates")
		tagPtr       = flag.String("tag", "", "Tag to add to events")
		ttlPtr       = flag.Float64("ttl", 10, "TTL for events")
	)
	iniflags.Parse()

	g.Address = fmt.Sprintf("%s:%d", *hostPtr, *portPtr)
	g.EventHost = *eventHostPtr
	g.Interval = time.Second * time.Duration(*intervalPtr)

	if len(*tagPtr) != 0 {
		g.Tag = strings.Split(*tagPtr, ",")
	}

	g.Ttl = float32(*ttlPtr)

}

func (g *Goshin) Start() {
	fmt.Print("Gare aux goriiillllleeeees!\n\n\n")
	fmt.Printf("Goshin will report each %s\n", g.Interval)

	// channel size has to be large enough
	// to allow Goshin send all metrics to Riemann
	// in g.Interval
	var collectQueue chan *Metric = make(chan *Metric, 1000)

	ticker := time.NewTicker(g.Interval)

	for t := range ticker.C {
		fmt.Println("Tick at ", t)

		for name, check := range g.AvailableChecks {
			_, ok := g.EnabledChecks[name]
			if ok {
				check.Collect(collectQueue)
			}
		}
		go g.Report(collectQueue)
	}
}

func (g *Goshin) EnforceState(metric *Metric) {

	threshold, present := g.Thresholds[metric.Service]

	if present {
		value := metric.Value

		// TODO threshold checking
		// only for int and float type
		switch {
		case value.(float64) > threshold.Critical:
			metric.State = "critical"
		case value.(float64) > threshold.Warning:
			metric.State = "warning"
		default:
			metric.State = "ok"
		}
	}
}

func (g *Goshin) Report(reportQueue chan *Metric) {

	c := goryman.NewGorymanClient(g.Address)
	err := c.Connect()

	if err != nil {
		fmt.Println("Can not connect to host")
	} else {

		more := true

		for more {
			select {
			case metric := <-reportQueue:
				g.EnforceState(metric)

				err := c.SendEvent(&goryman.Event{
					Metric:      metric.Value,
					Ttl:         g.Ttl,
					Service:     metric.Service,
					Description: metric.Description,
					Tags:        g.Tag,
					Host:        g.EventHost,
					State:       metric.State})

				if err != nil {
					fmt.Println("something does wrong:", err)
				}
			default:
				more = false
			}
		}
	}

	defer c.Close()
}

func (g *Goshin) AddCheck(name string, check Check) {
	g.AvailableChecks[name] = check
}

func (g *Goshin) ExtractEnabledChecks(checksPtr string) {
	// TODO make automatic construction of checks flag based on available checks
	checks := make(map[string]bool)

	if len(checksPtr) != 0 {
		for _, check := range strings.Split(checksPtr, ",") {
			checks[check] = true
		}
	}
	g.EnabledChecks = checks
}
