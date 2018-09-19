package generic

import (
	"time"
	"sync"
	"math/rand"
	log "github.com/sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "generic" // For Prometheus metrics.
)

// Exporter implements the prometheus.Exporter interface.
type Exporter struct {
	duration            prometheus.Gauge
	scrapeErrors        prometheus.Gauge
	totalScrapes        prometheus.Counter
	genericMetrics      map[string]*prometheus.GaugeVec
	metricsMtx          sync.RWMutex
	sync.RWMutex
}

// ScrapeResult is the metric data model.
type ScrapeResult struct {
	Name        string
	Value       float64
	Type        string
	Description string
}

var (
	labelNames = []string{"type", "description"}
)

func newGenericMetric(name, help string) *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace:   namespace,
			Name:        name,
			Help:        help,
		},
		labelNames,
	)
}

// NewExporter returns a new generic exporter.
func NewExporter() (*Exporter, error) {
	e := Exporter{
		duration: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "scrape_duration_seconds",
			Help:      "The scrape duration.",
		}),
		totalScrapes: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "scrapes_total",
			Help:      "Total scrapes.",
		}),
		scrapeErrors: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "scrape_error",
			Help:      "The scrape error status.",
		}),
		genericMetrics: map[string]*prometheus.GaugeVec{
			"current_value":  newGenericMetric("current_value", "Current value."),
		},
	}

	return &e, nil
}

func (e *Exporter) initMetric(name, help string) {
	e.genericMetrics = map[string]*prometheus.GaugeVec{}
	e.genericMetrics[name] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      name,
		Help:      help,
	}, labelNames)
}

// Describe output metrics metadata. It implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range e.genericMetrics {
		m.Describe(ch)
	}
	ch <- e.duration.Desc()
	ch <- e.totalScrapes.Desc()
	ch <- e.scrapeErrors.Desc()
}

// Collect process. It implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	scrapes := make(chan ScrapeResult)

	e.Lock()
	defer e.Unlock()

	e.resetMetrics()
	go e.scrape(scrapes)
	e.setMetrics(scrapes)

	ch <- e.duration
	ch <- e.totalScrapes
	ch <- e.scrapeErrors

	e.collectMetrics(ch)
}

func (e *Exporter) resetMetrics() {
	for _, m := range e.genericMetrics {
		m.Reset()
	}
}

func (e *Exporter) collectMetrics(metrics chan<- prometheus.Metric) {
	for _, m := range e.genericMetrics {
		m.Collect(metrics)
	}
}

func (e *Exporter) scrape(scrapes chan<- ScrapeResult) {

	defer close(scrapes)
	now := time.Now().UnixNano()
	e.totalScrapes.Inc()
	errorCount := 0

	// example metric 
	_type := "type"
	_desc := "description"
	_value := rand.Float64()
	
	log.Debugf("Creating new metric: current_value{type=%s, description=%s} = %v.", _type, _desc, _value)
	scrapes <- ScrapeResult{
		Name:        "current_value",
		Value:       _value,
		Type:        _type,
		Description: _desc,
	}

	e.scrapeErrors.Set(float64(errorCount))
	e.duration.Set(float64(time.Now().UnixNano()-now) / 1000000000)
}

func (e *Exporter) setMetrics(scrapes <-chan ScrapeResult) {
	log.Debug("set metrics")
	for scr := range scrapes {
		name := scr.Name
		if _, ok := e.genericMetrics[name]; !ok {
			e.metricsMtx.Lock()
			e.genericMetrics[name] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      name,
			}, labelNames)
			e.metricsMtx.Unlock()
		}
		var labels prometheus.Labels = map[string]string{
			"type":       scr.Type,
			"description": scr.Description,
		}
		e.genericMetrics[name].With(labels).Set(float64(scr.Value))
	}
}