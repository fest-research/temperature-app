package notifier

import (
	"fmt"
	"log"
	"net/http"
	"path"
)

const (
	Fahrenheit = "Fahrenheit"
	Celsius    = "Celsius"
)

type Notifier interface {
	Notify(event string) error
}

type UrlPathNotifier struct {
	endpoint string
}

func (notifier *UrlPathNotifier) Notify(event string) error {
	urlPath := path.Join(notifier.endpoint, event)
	log.Printf("GET request to: %s", urlPath)
	resp, err := http.Get(urlPath)
	if err != nil {
		return fmt.Errorf("Failed to notify %s  with event %s: %s", notifier.endpoint, event, err)
	}
	if resp.StatusCode == http.StatusOK {
		return fmt.Errorf("Unexpected status during notify: %s", resp.StatusCode)
	}
	return nil
}

func NewUrlPathNotifier(endpoint string) Notifier {
	return &UrlPathNotifier{endpoint: endpoint}
}
