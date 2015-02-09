package health

import (
	"flag"
	"github.com/MarshalMediaGroup/goshin"
	"github.com/MarshalMediaGroup/goshin/apps/health/checks"
	"strings"
)

type Health struct {
	*goshin.Goshin

	Ifaces       map[string]bool
	IgnoreIfaces map[string]bool
}

func New() *Health {
	app := &Health{}
	app.Goshin = goshin.NewGoshin()
	app.Configure()

	app.AddCheck("cpu", checks.NewCPUTime())
	app.AddCheck("memory", checks.NewMemoryUsage())
	app.AddCheck("load", checks.NewLoadAverage())
	app.AddCheck("net", checks.NewNetStats(app.Ifaces, app.IgnoreIfaces))
	app.AddCheck("iostat", checks.NewIOStat(app.Interval))
	return app
}

func (app *Health) Configure() {
	var (
		ifacesPtr         = flag.String("interfaces", "", "Interfaces to monitor")
		ignoreIfacesPtr   = flag.String("ignore-interfaces", "lo", "Interfaces to ignore (default: lo)")
		cpuWarningPtr     = flag.Float64("cpu-warning", 0.9, "CPU warning threshold (fraction of total jiffies")
		cpuCriticalPtr    = flag.Float64("cpu-critical", 0.95, "CPU critical threshold (fraction of total jiffies")
		loadWarningPtr    = flag.Float64("load-warning", 3, "Load warning threshold (load average / core")
		loadCriticalPtr   = flag.Float64("load-critical", 8, "Load critical threshold (load average / core)")
		memoryWarningPtr  = flag.Float64("memory-warning", 0.85, "Memory warning threshold (fraction of RAM)")
		memoryCriticalPtr = flag.Float64("memory-critical", 0.95, "Memory critical threshold (fraction of RAM)")
		checksPtr         = flag.String("checks", "cpu,load,memory,net,iostat", "A list of checks to run")
	)
	app.Goshin.Configure()
	app.ExtractEnabledChecks(*checksPtr)

	ifaces := make(map[string]bool)

	if len(*ifacesPtr) != 0 {
		for _, iface := range strings.Split(*ifacesPtr, ",") {
			ifaces[iface] = true
		}
	}
	app.Ifaces = ifaces

	ignoreIfaces := make(map[string]bool)

	if len(*ignoreIfacesPtr) != 0 {
		for _, ignoreIface := range strings.Split(*ignoreIfacesPtr, ",") {
			ignoreIfaces[ignoreIface] = true
		}
	}
	app.IgnoreIfaces = ignoreIfaces

	cpuThreshold := goshin.NewThreshold()
	cpuThreshold.Critical = *cpuCriticalPtr
	cpuThreshold.Warning = *cpuWarningPtr

	app.Thresholds["cpu"] = cpuThreshold

	loadThreshold := goshin.NewThreshold()
	loadThreshold.Critical = *loadCriticalPtr
	loadThreshold.Warning = *loadWarningPtr

	app.Thresholds["load"] = loadThreshold

	memoryThreshold := goshin.NewThreshold()
	memoryThreshold.Critical = *memoryCriticalPtr
	memoryThreshold.Warning = *memoryWarningPtr

	app.Thresholds["memory"] = memoryThreshold

}
