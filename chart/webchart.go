package main

import (
	"fmt"
	"net/http"
	"time"

	redistime "github.com/RedisTimeSeries/redistimeseries-go"
	chart "github.com/wcharczuk/go-chart/v2"
)

func getData(key string) ([]time.Time, []float64) {
	client := redistime.NewClient("10.0.0.3:6379", "time", nil)
	download, err := client.RangeWithOptions(key, redistime.TimeRangeMinimum, redistime.TimeRangeFull, redistime.DefaultRangeOptions)
	if err != nil {
		fmt.Println("error retrieving data", err)
	}

	var timestamps []time.Time
	var values []float64
	for _, counts := range download {
		timestamp := time.Unix(counts.Timestamp/1000, 0)
		//timestamp := time.Date(0, 0, 0, 0, 0, 0, int(counts.Timestamp/1000), time.Local)
		timestamps = append(timestamps, timestamp)
		values = append(values, counts.Value/100000)
	}
	return timestamps, values
}

func drawChart(res http.ResponseWriter, req *http.Request) {
	/*
	   This is an example of using the `TimeSeries` to automatically coerce time.Time values into a continuous xrange.
	   Note: chart.TimeSeries implements `ValueFormatterProvider` and as a result gives the XAxis the appropriate formatter to use for the ticks.
	*/
	timestamps, download := getData("download")
	timestamps, upload := getData("upload")
	//timestamps, ping := getData("ping")

	graph := chart.Chart{
		Series: []chart.Series{
			chart.TimeSeries{
				Name:    "upload",
				XValues: timestamps,
				YValues: upload,
			},
			chart.TimeSeries{
				Name:    "download",
				XValues: timestamps,
				YValues: download,
			},
			//chart.TimeSeries{
			//	Name:    "ping",
			//	XValues: timestamps,
			//	YValues: ping,
			//},
		},
	}

	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}

	res.Header().Set("Content-Type", "image/png")
	graph.Render(chart.PNG, res)
}

func drawCustomChart(res http.ResponseWriter, req *http.Request) {
	/*
	   This is basically the other timeseries example, except we switch to hour intervals and specify a different formatter from default for the xaxis tick labels.
	*/
	graph := chart.Chart{
		XAxis: chart.XAxis{
			ValueFormatter: chart.TimeHourValueFormatter,
		},
		Series: []chart.Series{
			chart.TimeSeries{
				XValues: []time.Time{
					time.Now().Add(-10 * time.Hour),
					time.Now().Add(-9 * time.Hour),
					time.Now().Add(-8 * time.Hour),
					time.Now().Add(-7 * time.Hour),
					time.Now().Add(-6 * time.Hour),
					time.Now().Add(-5 * time.Hour),
					time.Now().Add(-4 * time.Hour),
					time.Now().Add(-3 * time.Hour),
					time.Now().Add(-2 * time.Hour),
					time.Now().Add(-1 * time.Hour),
					time.Now(),
				},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0},
			},
		},
	}

	res.Header().Set("Content-Type", "image/png")
	graph.Render(chart.PNG, res)
}

func main() {
	http.HandleFunc("/", drawChart)
	http.HandleFunc("/favicon.ico", func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte{})
	})
	http.HandleFunc("/custom", drawCustomChart)
	http.ListenAndServe(":8080", nil)
}
