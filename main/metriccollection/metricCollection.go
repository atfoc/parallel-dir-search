package metriccollection

import (
	"strings"
)

type MetricCollection struct {
	metrics map[string]metricInstance
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
	c.init()

	// todo handle error if not all values have been provided
	// todo make it thread safe
	key := buildKey(metric, labelValues)
	oldMetricValue, ok := c.metrics[key]
	if !ok {
		oldMetricValue = newMetricInstance(metric, labelValues)
	}

	oldMetricValue.value += amount
	c.metrics[key] = oldMetricValue
}

func (c *MetricCollection) Dec(metric Metric, labelValues []string, amount int) {
	c.init()

	// todo handle error if not all values have been provided
	// todo make it thread safe
	key := buildKey(metric, labelValues)
	oldMetricValue, ok := c.metrics[key]
	if !ok {
		oldMetricValue = newMetricInstance(metric, labelValues)
	}

	oldMetricValue.value -= amount
	c.metrics[key] = oldMetricValue
}

func (c *MetricCollection) Set(metric Metric, labelValues []string, amount int) {
	c.init()

	// todo handle error if not all values have been provided
	// todo make it thread safe
	key := buildKey(metric, labelValues)
	oldMetricValue, ok := c.metrics[key]
	if !ok {
		oldMetricValue = newMetricInstance(metric, labelValues)
	}

	oldMetricValue.value = amount
	c.metrics[key] = oldMetricValue
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

func (c *MetricCollection) init() {
	if c.metrics == nil {
		c.metrics = make(map[string]metricInstance)
	}
}
