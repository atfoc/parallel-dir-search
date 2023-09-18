package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"html/template"
	"io"
	"math"
	"net/http"
	"os"
)

var templates *template.Template

func main() {

	templates = template.Must(template.New("template").Funcs(
		template.FuncMap{
			"NewUUID": func() string {
				return uuid.New().String()
			},
		},
	).ParseGlob("templates/*.gohtml"))

	http.HandleFunc("/", index)

	if err := http.ListenAndServe("localhost:8080", nil); err != nil {
		panic(errors.Join(errors.New("failed to start server"), err))
	}
}

type DatasetCollection struct {
	labels []float64
	values map[string][]*string
	isInit bool
}

func (c *DatasetCollection) Add(data map[string]string, time float64) {
	c.panicIfNotInit()
	c.labels = append(c.labels, time)

	for key, value := range data {
		c.values[key] = append(c.values[key], &value)
	}

	for key, value := range c.values {
		for len(value) < len(c.labels) {
			value = append(value, nil)
		}
		c.values[key] = value
	}
}

func (c *DatasetCollection) GetDatasets() []ChartDataset {
	result := make([]ChartDataset, 0)
	for key, value := range c.values {
		result = append(result, ChartDataset{
			Label: key,
			Data:  buildData(c.labels, value),
		})
	}

	return result
}

func (c *DatasetCollection) panicIfNotInit() {
	if !c.isInit {
		panic(errors.New("can not use DatasetCollection when it it not init"))
	}
}

func buildData(labels []float64, values []*string) []ChartData {
	result := make([]ChartData, 0)
	for i, value := range values {
		result = append(result, ChartData{
			X: fmt.Sprintf("%f", labels[i]),
			Y: value,
		})
	}

	return result
}

func GetCharts() []Chart {
	open, err := os.Open("metric.json")
	if err != nil {
		panic(err)
	}
	defer open.Close()

	all, err := io.ReadAll(open)
	if err != nil {
		panic(all)
	}

	var data []any
	if err = json.Unmarshal(all, &data); err != nil {
		panic(err)
	}

	//result := make(map[string]ChartDataset)
	result := GetFirstChart(data)
	result1 := GetSecondChart(data)
	result2 := GetThirdChart(data)

	//add(&result, processOneTime(lastTimeArray), fmt.Sprintf("%f", lastTime))

	return []Chart{
		&LineChart{Datasets: result.GetDatasets()},
		&LineChart{Datasets: result1.GetDatasets()},
		&LineChart{Datasets: result2.GetDatasets()},
	}

}

func GetSecondChart(data []any) DatasetCollection {
	return IterateOverData(data, 1, processOneTime1)
}

func GetThirdChart(data []any) DatasetCollection {
	return IterateOverData(data, 1, processOneTime2)
}

func IterateOverData(data []any, amountOfTimes int, process func([]map[string]any) map[string]string) DatasetCollection {
	result := NewDatasetCollection()
	var lastTime float64
	var lastTimeArray []map[string]any
	var timesAdded int

	for _, el := range data {
		metric := el.(map[string]any)
		if lastTime != metric["time"].(float64) {
			timesAdded += 1
			if timesAdded == amountOfTimes {
				if len(lastTimeArray) != 0 {
					result.Add(process(lastTimeArray), lastTime)
				}
				lastTimeArray = make([]map[string]any, 0)
				lastTime = metric["time"].(float64)
				timesAdded = 0
			}

		}

		lastTimeArray = append(lastTimeArray, metric)
	}
	return result
}

func GetFirstChart(data []any) DatasetCollection {
	return IterateOverData(data, 50, processOneTime)
}

func NewDatasetCollection() DatasetCollection {
	return DatasetCollection{
		labels: make([]float64, 0),
		values: make(map[string][]*string),
		isInit: true,
	}
}

func processOneTime(array []map[string]any) map[string]string {
	tmp := make(map[string]float64)

	for _, el := range array {
		if el["metricName"] == "dir_processed_count" {
			//tmp[el["workerIndex"].(string)] = tmp[el["workerIndex"].(string)] + el["value"].(float64)
			tmp["rate"] = tmp["rate"] + el["value"].(float64)
		}
	}

	result := make(map[string]string)
	for key, el := range tmp {
		result[key] = fmt.Sprintf("%f", el/(array[len(array)-1]["time"].(float64)-array[0]["time"].(float64)))
	}

	return result
}

func processOneTime1(array []map[string]any) map[string]string {
	tmp := make(map[string]float64)

	for _, el := range array {
		if el["metricName"] == "dir_to_process_queue_size" {
			tmp["size"] = math.Max(tmp["size"], el["value"].(float64))
		}
	}

	result := make(map[string]string)
	for key, el := range tmp {
		result[key] = fmt.Sprintf("%f", el)
	}

	return result
}

func processOneTime2(array []map[string]any) map[string]string {
	tmp := make(map[string]float64)

	for _, el := range array {
		if el["metricName"] == "dir_read_time" {
			tmp["time"] = el["value"].(float64)
		}
	}

	result := make(map[string]string)
	for key, el := range tmp {
		result[key] = fmt.Sprintf("%f", el)
	}

	return result
}

type LineChart struct {
	Datasets []ChartDataset
}

type Chart interface {
	GetChartType() string
}

type ChartDataset struct {
	Label string
	Data  []ChartData
}

type ChartData struct {
	X string
	Y *string
}

func (c *LineChart) GetChartType() string {
	return "line"
}

func index(writer http.ResponseWriter, request *http.Request) {

	chartData := GetCharts()
	if err := templates.ExecuteTemplate(writer, "index.gohtml", chartData); err != nil {
		writer.Write([]byte(err.Error()))
	}
}
