package metriccollection

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"
)

type MetricCollection struct {
	metrics map[string]metricInstance
	mutex   sync.Mutex
}

type metricInstance struct {
	metricName  string
	labels      []string
	labelValues []string
	value       int
}

type Metric struct {
	Name   string
	Labels []string
}

func (c *MetricCollection) Inc(metric Metric, labelValues []string, amount int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// todo handle error if not all values have been provided
	key := buildKey(metric, labelValues)
	oldMetricValue, ok := c.getMetric(key)
	if !ok {
		oldMetricValue = newMetricInstance(metric, labelValues)
	}

	oldMetricValue.value += amount
	c.setMetric(key, oldMetricValue)
}

func (c *MetricCollection) Dec(metric Metric, labelValues []string, amount int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// todo handle error if not all values have been provided
	key := buildKey(metric, labelValues)
	oldMetricValue, ok := c.getMetric(key)
	if !ok {
		oldMetricValue = newMetricInstance(metric, labelValues)
	}

	oldMetricValue.value -= amount
	c.setMetric(key, oldMetricValue)
}

func (c *MetricCollection) Set(metric Metric, labelValues []string, amount int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// todo handle error if not all values have been provided
	key := buildKey(metric, labelValues)
	oldMetricValue, ok := c.getMetric(key)
	if !ok {
		oldMetricValue = newMetricInstance(metric, labelValues)
	}

	oldMetricValue.value = amount
	c.setMetric(key, oldMetricValue)
}

func newMetricInstance(metric Metric, labelValues []string) metricInstance {
	return metricInstance{
		metricName:  metric.Name,
		labels:      metric.Labels,
		labelValues: labelValues,
		value:       0,
	}
}

func buildKey(metric Metric, values []string) string {
	// todo handle errors
	builder := strings.Builder{}

	builder.WriteString(metric.Name)
	builder.WriteString("_")
	for i, labelName := range metric.Labels {
		builder.WriteString(labelName)
		builder.WriteString("_")
		builder.WriteString(values[i])
	}

	return builder.String()
}

func (c *MetricCollection) getMetric(key string) (metricInstance, bool) {
	if c.metrics == nil {
		return metricInstance{}, false
	}

	value, ok := c.metrics[key]

	return value, ok
}

func (c *MetricCollection) setMetric(key string, value metricInstance) {
	if c.metrics == nil {
		c.metrics = make(map[string]metricInstance)
	}

	c.metrics[key] = value
}

func (c *MetricCollection) getAllMetrics() []metricInstance {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	result := make([]metricInstance, 0)

	for _, value := range c.metrics {
		result = append(result, value)
	}

	return result
}

type MetricExporter struct {
	metricCollection *MetricCollection
	output           io.WriteCloser
	stopChanel       chan interface{}
	doneChanel       chan interface{}
}

func NewMetricExporter(output io.WriteCloser, collection *MetricCollection) MetricExporter {
	return MetricExporter{
		metricCollection: collection,
		output:           output,
		stopChanel:       make(chan interface{}),
		doneChanel:       make(chan interface{}),
	}
}

func (e *MetricExporter) StartExporting() error {
	e.panicIfNotInit()

	_, err := fmt.Fprint(e.output, "[")
	if err != nil {
		return err
	}

	go e.export()

	return nil
}

func (e *MetricExporter) FinishExporting() {
	e.panicIfNotInit()
	close(e.stopChanel)
	<-e.doneChanel
}

func (e *MetricExporter) panicIfNotInit() {
	if e.output == nil {
		panic(errors.New("can not use MetricExporter when it is not init"))
	}
}

func (e *MetricExporter) export() {
	hadOutput := false
	for {
		if hadOutput {
			if _, err := fmt.Fprint(e.output, ","); err != nil {
				panic(err)
			}
		}
		hadOutput = e.writeAllMetricsToOutput(time.Now().UnixMilli())

		select {
		case <-e.stopChanel:
			close(e.doneChanel)
			if _, err := fmt.Fprint(e.output, "]"); err != nil {
				panic(err)
			}
			return
		default:
			time.Sleep(time.Millisecond * 5)
		}
	}
}

func (e *MetricExporter) writeAllMetricsToOutput(currentTime int64) bool {
	metrics := e.metricCollection.getAllMetrics()
	hadOutput := false
	for i, metric := range metrics {
		hadOutput = true
		e.writeMetricToOutput(metric, currentTime)
		if i != len(metrics)-1 {
			if _, err := fmt.Fprint(e.output, ","); err != nil {
				panic(err)
			}
		}
	}

	return hadOutput
}

func (e *MetricExporter) writeMetricToOutput(metric metricInstance, currentTime int64) {
	_, err := fmt.Fprintf(e.output, "{\"time\": %d,\"metricName\": \"%s\", \"value\": %d", currentTime, metric.metricName, metric.value)
	if err != nil {
		panic(err)
	}

	for i, value := range metric.labelValues {
		_, err := fmt.Fprintf(e.output, ",\"%s\": \"%s\"", metric.labels[i], value)
		if err != nil {
			panic(err)
		}
	}

	if _, err := fmt.Fprint(e.output, "}"); err != nil {
		panic(err)
	}
}
