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

type ControllerStrategy interface {
	ShouldNotify(event string) bool
	Decorate(event string) (string, error)
}

type TemperatureController struct {
	stopChan <-chan string
	watcher  watcher.Watcher
	notifier notifier.Notifier
	strategy ControllerStrategy
}

func (ctrl *TemperatureController) Start() {
	ctrl.watcher.Watch()
	for {
		select {
		case <-ctrl.stopChan:
			return
		case event := <-ctrl.watcher.NextEvent():
			if ok := ctrl.strategy.ShouldNotify(event); ok {
				transformedEvent, err := ctrl.strategy.Decorate(event)
				if err != nil {
					log.Printf("Warning: could not transform event %s: %s\n", event, err)
					continue
				}
				err = ctrl.notifier.Notify(transformedEvent)
				if err != nil {
					log.Printf("Error during notify: %s", err)
				}
			}
		default:
			continue
		}
	}

}

type defaultControllerStrategy struct {
}

func (strat *defaultControllerStrategy) ShouldNotify(event string) bool {
	return true
}

func (strat *defaultControllerStrategy) Decorate(event string) (string, error) {
	return event, nil
}

type thresholdingControllerStrategy struct {
	defaultControllerStrategy
	threshold         float32
	thresholdExceeded bool
}

func (strat *thresholdingControllerStrategy) ShouldNotify(event string) bool {
	value, err := strconv.ParseFloat(event, 32)
	if err != nil {
		log.Printf("Warning: could not process event %s. Will not notify.", event)
		return false
	}
	val := float32(value)

	// Notify only if we have just exceeded the thresholds
	if val > strat.threshold && !strat.thresholdExceeded {
		strat.thresholdExceeded = true
		return true
	}

	if val <= strat.threshold {
		strat.thresholdExceeded = false
	}

	return false
}

type celsiusControllerStrategy struct {
	thresholdingControllerStrategy
}

func (strat *celsiusControllerStrategy) Decorate(event string) (string, error) {
	value, err := strconv.ParseFloat(event, 32)
	if err != nil {
		return "", err
	}
	val := float32(value)
	tempVal := goDS18B20.Temperature(val)
	return fmt.Sprintf("%.03f C", tempVal.Celsius()), nil
}

func newCelsiusControllerStrategy(threshold float32) ControllerStrategy {
	return &celsiusControllerStrategy{thresholdingControllerStrategy{threshold: threshold}}
}

type fahrenheitControllerStrategy struct {
	thresholdingControllerStrategy
}

func (strat *fahrenheitControllerStrategy) Decorate(event string) (string, error) {
	value, err := strconv.ParseFloat(event, 32)
	if err != nil {
		return "", err
	}
	val := float32(value)
	tempVal := goDS18B20.Temperature(val)
	return fmt.Sprintf("%.03f F", tempVal.Fahrenheit()), nil
}

func newFahrenheitControllerStrategy(threshold float32) ControllerStrategy {
	return &fahrenheitControllerStrategy{thresholdingControllerStrategy{threshold: threshold}}
}

func NewCelsiusController(stopChan <-chan string, w watcher.Watcher, n notifier.Notifier, threshold float32) Controller {
	return &TemperatureController{
		stopChan: stopChan,
		watcher:  w,
		notifier: n,
		strategy: newCelsiusControllerStrategy(threshold),
	}
}

func NewFahrenheitController(stopChan <-chan string, w watcher.Watcher, n notifier.Notifier, threshold float32) Controller {
	return &TemperatureController{
		stopChan: stopChan,
		watcher:  w,
		notifier: n,
		strategy: newFahrenheitControllerStrategy(threshold),
	}
}
