package main

import (
	"errors"
	"github.com/google/uuid"
	"html/template"
	"net/http"
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
	Y string
}

func (c *LineChart) GetChartType() string {
	return "line"
}

func index(writer http.ResponseWriter, request *http.Request) {

	chartData := []LineChart{{
		Datasets: []ChartDataset{
			{Label: "test",
				Data: []ChartData{
					{X: "0", Y: "0"},
					{X: "1", Y: "1"},
					{X: "2", Y: "2"},
					{X: "3", Y: "1"},
				}},
			{Label: "test1",
				Data: []ChartData{
					{X: "0", Y: "7"},
					{X: "1", Y: "8"},
					{X: "2", Y: "9"},
					{X: "3", Y: "15"},
				}},
		},
	},
		{
			Datasets: []ChartDataset{
				{Label: "test",
					Data: []ChartData{
						{X: "0", Y: "0"},
						{X: "1", Y: "1"},
						{X: "2", Y: "2"},
						{X: "3", Y: "1"},
					}},
				{Label: "test1",
					Data: []ChartData{
						{X: "0", Y: "7"},
						{X: "1", Y: "8"},
						{X: "2", Y: "9"},
						{X: "3", Y: "15"},
					}},
			},
		},
	}
	if err := templates.ExecuteTemplate(writer, "index.gohtml", chartData); err != nil {
		writer.Write([]byte(err.Error()))
	}
}
