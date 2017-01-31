// The MIT License
//
// Copyright (c) 2016 Sebastian Florek
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/fest-research/temperature-app/controller"
	"github.com/fest-research/temperature-app/notifier"
	"github.com/fest-research/temperature-app/sensor/ds18b20"
	"github.com/fest-research/temperature-app/watcher"
	"github.com/spf13/pflag"
)

var (
	remoteEndpoint = pflag.String("remote-endpoint",
		"http://104.155.11.172:8080/api/v1/proxy/namespaces/default/services/demo/push",
		"Address of the remove server to which measurements are sent.")
	threshold = pflag.Float32("temperature-threshold", 32.5,
		"Events are produced when this temperature threshold is exceeded.")
	units = pflag.String("temperature-units", notifier.Celsius,
		"What units should the temperature readings be in. Default is Celsius.")
)

func main() {
	// Set logging out to standard console out
	log.SetOutput(os.Stdout)

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	log.Printf("Sending to remote endpoint: %s", *remoteEndpoint)

	watcherChan := make(chan string)
	defer close(watcherChan)
	stopChan := make(chan string)
	defer close(stopChan)

	temperatureWatcher := watcher.NewTemperatureWatcher(time.Second*2, sensor.DS18B20Reader{}, watcherChan)
	urlPathNotifier := notifier.NewUrlPathNotifier(*remoteEndpoint)

	ctrl := controller.Controller(nil)
	if *units == notifier.Celsius {
		ctrl = controller.NewCelsiusController(stopChan, temperatureWatcher, urlPathNotifier, *threshold)
	} else if *units == notifier.Fahrenheit {
		ctrl = controller.NewFahrenheitController(stopChan, temperatureWatcher, urlPathNotifier, *threshold)
	}
	ctrl.Start()
}
