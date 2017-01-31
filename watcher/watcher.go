package watcher

import (
	"github.com/fest-research/temperature-app/sensor/ds18b20"
	"log"
	"time"
)

type Watcher interface {
	Watch()
	NextEvent() <-chan string
}

type TemperatureWatcher struct {
	reader        sensor.DS18B20Reader
	watchInterval time.Duration
	resultChan    chan string
}

func (tempWatcher *TemperatureWatcher) Watch() {
	go func() {
		for x := range time.Tick(tempWatcher.watchInterval) {
			// TODO: produce an event
			//reading, err := tempWatcher.reader.ReadFromSensor()
			//if err != nil {
			//	log.Printf("Warning: can not read from sensor: %s, time: %s\n", err, x)
			//	continue
			//}
			log.Printf("Sending reading at %s", x)
			//tempWatcher.resultChan <- reading.String()
			tempWatcher.resultChan <- "34.6"
		}
	}()
}

func (tempWatcher *TemperatureWatcher) NextEvent() <-chan string {
	return tempWatcher.resultChan
}

func NewTemperatureWatcher(watchInterval time.Duration, reader sensor.DS18B20Reader, resultChan chan string) Watcher {
	return &TemperatureWatcher{reader: reader, watchInterval: watchInterval, resultChan: resultChan}
}
