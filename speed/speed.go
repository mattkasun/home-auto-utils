package main

import (
	"fmt"
	"log"
	"time"

	redistime "github.com/RedisTimeSeries/redistimeseries-go"
	"github.com/showwin/speedtest-go/speedtest"
)

func main() {
	var host string = "10.0.0.3:6379"
	var name string = "time"
	//log := log.Default()

	user, _ := speedtest.FetchUserInfo()

	serverList, _ := speedtest.FetchServerList(user)
	targets, _ := serverList.FindServer([]int{})
	for _, s := range targets {
		s.PingTest()
		s.DownloadTest(false)
		s.UploadTest(false)
		fmt.Printf("Latency: %s, Download: %f, Upload: %f\n", s.Latency, s.DLSpeed, s.ULSpeed)

		time := time.Now().Unix()
		redis := redistime.NewClient(host, name, nil)
		_, err := redis.Add("Upload", time, s.ULSpeed)
		if err != nil {
			log.Print("could not add record to Upload")
		}
		_, err = redis.Add("Downoad", time, s.DLSpeed)
		if err != nil {
			log.Print("could not add record to Download")
		}
		_, err = redis.Add("Ping", time, s.Latency.Seconds())
		if err != nil {
			log.Print("could not add record to Ping")
		}
	}
}
