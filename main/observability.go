package main

import (
	"errors"
	"fmt"
	"io"
	"main/metriccollection"
	"os"
	"strconv"
	"time"
)

var metricCollection = metriccollection.MetricCollection{}

var dirProcessedMetric = metriccollection.Metric{
	Name:   "dir_processed_count",
	Labels: []string{"workerIndex"},
}

var dirToProcessQueueSize = metriccollection.Metric{
	Name: "dir_to_process_queue_size",
}

func ObserveDirProcessed(workerIndex int) {
	metricCollection.Inc(dirProcessedMetric, []string{strconv.Itoa(workerIndex)}, 1)
}

func ObserveAddDirToQueue() {
	metricCollection.Dec(dirToProcessQueueSize, nil, 1)
}

func ObserveDirRemovedFromQueue() {
	metricCollection.Dec(dirToProcessQueueSize, nil, 1)
}

type MetricExporter struct {
	metricCollection *metriccollection.MetricCollection
	output           io.WriteCloser
	stopChanel       chan interface{}
	doneChanel       chan interface{}
}

func NewMetricExporter(output io.WriteCloser, collection *metriccollection.MetricCollection) MetricExporter {
	return MetricExporter{
		metricCollection: collection,
		output:           output,
		stopChanel:       make(chan interface{}),
		doneChanel:       make(chan interface{}),
	}
}

func (e *MetricExporter) StartExporting() {
	e.panicIfNotInit()

	go e.export()
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
	for {
		_, err := fmt.Fprintln(e.output, "cao")
		if err != nil {
			panic(err)
		}

		select {
		case <-e.stopChanel:
			close(e.doneChanel)
			return
		default:
			time.Sleep(time.Second)
			continue
		}

	}
}

func GetMetricExporter() (MetricExporter, error) {
	open, err := os.Create("metric.json")
	if err != nil {
		return MetricExporter{}, err
	}

	return NewMetricExporter(open, &metricCollection), nil
}
