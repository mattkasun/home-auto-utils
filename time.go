package main

import (
	"fmt"
	"time"

	redistime "github.com/RedisTimeSeries/redistimeseries-go"
)

func main() {

	client := redistime.NewClient("10.0.0.3:6379", "time", nil)
	//key := "download"

	download, err := client.RangeWithOptions("download", redistime.TimeRangeMinimum, redistime.TimeRangeFull, redistime.DefaultRangeOptions)

	if err != nil {
		fmt.Println("error retrieving data", err)
	}
	var timestamps []int64
	var values []float64
	for _, counts := range download {
		timestamps = append(timestamps, counts.Timestamp/1000)
		values = append(values, counts.Value/100000)
		//fmt.Println(counts.Timestamp)
		//time1 := time.Unix(0, counts.Timestamp/1000)
		//fmt.Println(time1, counts.Value)
		//time2 := time.Unix(counts.Timestamp/1000, 0)
		//fmt.Println(time2, counts.Value)
		fmt.Println(counts.Timestamp/1000, counts.Value/100000)
		for i := 1; i < 1000000; i = i * 10 {
			date := time.Date(0, 0, 0, 0, 0, 0, int(counts.Timestamp/int64(i)), time.UTC)
			date2 := time.Unix(counts.Timestamp/int64(i), 0)
			fmt.Println(i, counts.Timestamp/int64(i), date, date2)
		}
		for i := 1; i < 1000000; i = i * 10 {
			date := time.Date(0, 0, 0, 0, 0, 0, int(counts.Timestamp*int64(i)), time.UTC)
			fmt.Println(i, counts.Timestamp*int64(i), date)
		}
	}
	fmt.Println(timestamps)
	fmt.Println(values)
}
