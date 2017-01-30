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

package sensor

import (
	"fmt"

	"log"

	DS18B20 "github.com/traetox/goDS18B20"
)

// TODO add doc
type DS18B20Reader struct{}

// TODO add doc
func (ds DS18B20Reader) ReadFromSensor() (*DS18B20.Temperature, error) {
	log.Print("Reading from sensor DS18B20.")
	id, err := ds.GetSlaveID()
	if err != nil {
		return nil, err
	}

	probe, err := DS18B20.NewProbe(id)
	if err != nil {
		return nil, err
	}

	err = probe.Update()
	if err != nil {
		return nil, err
	}

	temp, err := probe.Temperature()
	if err != nil {
		return nil, err
	}
	return &temp, nil
}

func (DS18B20Reader) GetSlaveID() (string, error) {
	slaves, err := DS18B20.Slaves()
	if err != nil {
		return "", err
	}

	if len(slaves) == 0 {
		return "", fmt.Errorf("%s", "Could not find any DS18B20 sensors.")
	}

	if len(slaves) > 1 {
		log.Print("Found more than 1 DS18B20 sensor, taking first one.")
	}

	log.Printf("Found slave with id: %s", slaves[0])
	return slaves[0], nil
}
