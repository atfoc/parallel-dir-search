package main

import (
	"main/metriccollection"
	"os"
	"strconv"
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

func GetMetricExporter(fileName string) (metriccollection.MetricExporter, error) {
	open, err := os.Create(fileName)
	if err != nil {
		return metriccollection.MetricExporter{}, err
	}

	return metriccollection.NewMetricExporter(open, &metricCollection), nil
}
