package main

import (
	"os"

	"github.com/wcharczuk/go-chart"
)

func main() {
	graph := chart.Chart{
		Series: []chart.Series{
			chart.ContinuousSeries{
				XValues: []float64{1.0, 2.0, 3.0, 4.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0},
			},
		},
	}

	f, _ := os.Create("chart.png")
	defer f.Close()
	graph.Render(chart.PNG, f)

}
