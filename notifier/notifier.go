package notifier

import (
	"fmt"
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
	resp, err := http.Get(path.Join(notifier.endpoint, event))
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
