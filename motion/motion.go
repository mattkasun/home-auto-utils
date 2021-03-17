package main

import (
	"fmt"
	"log"
	"log/syslog"
	"math"
	"os"
	"time"

	redistime "github.com/RedisTimeSeries/redistimeseries-go"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var host string = "10.0.0.3:6379"
var name string = "time"
var broker string = "tcp://pinode:1883"
var mqttClientID = "nusak"

var handler mqtt.MessageHandler = func(client mqtt.Client, message mqtt.Message) {
	sysLog, err := syslog.Dial("", "localhost", syslog.LOG_INFO|syslog.LOG_DAEMON, "motion-to-redis")
	if err != nil {
		log.Fatal(err)
	}
	//print/log received topic/message
	sysLog.Emerg(fmt.Sprint("message recevied", message.Topic(), string(message.Payload())))
	//connect to redis timeseries
	redis := redistime.NewClient(host, name, nil)
	options := redistime.DefaultCreateOptions
	options.DuplicatePolicy = redistime.LastDuplicatePolicy
	//parse topic and payload
	topic := message.Topic()
	//convert payload to unix timestamp
	payload := string(message.Payload())
	time, err := time.Parse("Mon Jan 2 15:04:05 MST 2006", payload)
	if err != nil {
		sysLog.Err(fmt.Sprint("error converting timestamp", err))
	}
	timestamp := time.Unix()
	//convert timestamp to hour bucket
	min_seconds := int64(math.Mod(float64(timestamp), 3600))
	timestamp = timestamp - min_seconds
	//get the last update from redis
	last, errors := redis.Get(topic)
	if errors != nil {
		if errors.Error() == "ERR TSDB: the key does not exist" { //new key
			sysLog.Err(fmt.Sprint("key does not exist"))
			_, err = redis.IncrBy(topic, timestamp, 1, options)
			if err != nil {
				sysLog.Err(fmt.Sprint("error incrementing ", topic, err))
			}
			return
		} else {
			sysLog.Emerg(fmt.Sprinti("unable to get redis key", errors.Error()))
			return
		}
	}
	if last.Timestamp == timestamp { // we are in the same hour bucket
		_, err = redis.IncrBy(topic, timestamp, 1, options)
		if err != nil {
			sysLog.Err(fmt.Sprint("error incrementing ", topic, err))
		}
		return
	}
	// start a new bucket
	_, err = redis.AddWithOptions(topic, timestamp, 1, options)
	if err != nil {
		sysLog.Err(fmt.Sprint("error adding ", topic, err))
	}
}

func main() {
	sysLog, err := syslog.Dial("", "localhost", syslog.LOG_WARNING|syslog.LOG_DAEMON, "motion-to-redis")
	if err != nil {
		log.Fatal(err)
	}
	//connect to MQTT broker and set handler
	options := mqtt.NewClientOptions().AddBroker(broker).SetClientID(mqttClientID)
	options.SetDefaultPublishHandler(handler)

	c := mqtt.NewClient(options)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		sysLog.Emerg(fmt.Sprint("no connection to mqtt broker", token.Error()))
		os.Exit(1)
	}
	//subscribe to motion events
	if token := c.Subscribe("motion/#", 0, nil); token.Wait() && token.Error() != nil {
		sysLog.Emerg(fmt.Sprint("unable to subscribe", token.Error()))
		os.Exit(1)
	}
	//loop forever
	for {
	}

	c.Disconnect(10)
}
