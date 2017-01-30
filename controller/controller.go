package controller

import (
	"fmt"
	"github.com/fest-research/temperature-app/notifier"
	"github.com/fest-research/temperature-app/watcher"
	"github.com/traetox/goDS18B20"
	"log"
	"strconv"
)

type Controller interface {
	Start()
}

type rawTemperatureController struct {
	stopChan <-chan string
	watcher  watcher.Watcher
	notifier notifier.Notifier
}

func (ctrl *rawTemperatureController) Start() {
	ctrl.watcher.Watch()
	for {
		select {
		case <-ctrl.stopChan:
			return
		case event := <-ctrl.watcher.NextEvent():
			if ok := ctrl.shouldNotify(event); ok {
				fmt.Printf("Attempting to notify for event: %s", event)
				transformedEvent, err := ctrl.decorate(event)
				if err != nil {
					log.Printf("Warning: could not transform event %s: %s\n", event, err)
					continue
				}
				ctrl.notifier.Notify(transformedEvent)
			}
		}
	}
}

func (ctrl *rawTemperatureController) shouldNotify(event string) bool {
	return true
}

func (ctrl *rawTemperatureController) decorate(event string) (string, error) {
	return event, nil
}

func newRawTemperatureController(stopChan <-chan string, watcher watcher.Watcher,
	notifier notifier.Notifier) *rawTemperatureController {
	return &rawTemperatureController{stopChan: stopChan, watcher: watcher, notifier: notifier}
}

type thresholdingController struct {
	*rawTemperatureController
	threshold         float32
	thresholdExceeded bool
}

func (ctrl *thresholdingController) shouldNotify(event string) bool {
	value, err := strconv.ParseFloat(event, 32)
	if err != nil {
		log.Printf("Warning: could not process event %s. Will not notify.", event)
		return false
	}
	val := float32(value)

	// Notify only if we have just exceeded the thresholds
	if val > ctrl.threshold && !ctrl.thresholdExceeded {
		ctrl.thresholdExceeded = true
		return true
	}

	if val <= ctrl.threshold {
		ctrl.thresholdExceeded = false
	}

	return false
}

func newThresholdingController(stopChan <-chan string, watcher watcher.Watcher,
	notifier notifier.Notifier, threshold float32) *thresholdingController {
	return &thresholdingController{
		rawTemperatureController: newRawTemperatureController(stopChan, watcher, notifier),
		threshold:                threshold,
		thresholdExceeded:        false,
	}
}

type CelsiusController struct {
	*thresholdingController
	threshold float32
}

func (celController *CelsiusController) decorate(event string) (string, error) {
	value, err := strconv.ParseFloat(event, 32)
	if err != nil {
		return "", err
	}
	val := float32(value)
	tempVal := goDS18B20.Temperature(val)
	return fmt.Sprintf("%.03f C", tempVal.Celsius()), nil
}

func NewCelsiusController(stopChan <-chan string, watcher watcher.Watcher,
	notifier notifier.Notifier, threshold float32) Controller {
	return &CelsiusController{
		thresholdingController: newThresholdingController(stopChan, watcher, notifier, threshold),
		threshold:              threshold,
	}
}

type FahrenheitController struct {
	*thresholdingController
	threshold float32
}

func (ctrl *FahrenheitController) shouldNotify(event string) bool {
	value, err := strconv.ParseFloat(event, 32)
	if err != nil {
		fmt.Printf("Warning: could not process event %s. Will not notify.", event)
		return false
	}
	val := float32(value)
	return val > ctrl.threshold
}

func (celController *FahrenheitController) decorate(event string) (string, error) {
	value, err := strconv.ParseFloat(event, 32)
	if err != nil {
		return "", err
	}
	val := float32(value)
	tempVal := goDS18B20.Temperature(val)
	return fmt.Sprintf("%.03f C", tempVal.Fahrenheit()), nil
}

func NewFahrenheitController(stopChan <-chan string, watcher watcher.Watcher,
	notifier notifier.Notifier, threshold float32) Controller {
	return &CelsiusController{
		thresholdingController: newThresholdingController(stopChan, watcher, notifier, threshold),
	}
}
