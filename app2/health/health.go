package health

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type (
	ApplicationDetailed struct {
		Name         string        `json:"name"`
		Version      string        `json:"version"`
		Status       string  `json:"status"`
		Integrations []Integration `json:"integrations"`
	}

	Integration struct {
		Name         string  `json:"name"`
		Port         string  `json:"port"`
		Host         string  `json:"host"`
		ResponseTime float64 `json:"response_time"`
		Status       string  `json:"status"`
		Error        string  `json:"error,omitempty"`
	}

	Health struct {
		checks map[string]Config
	}
)

func (h *Health) Handler() http.Handler {
	return http.HandlerFunc(h.HandlerFunc)
}

func (h *Health) Check() ApplicationDetailed {
	app := ApplicationDetailed{
		Name:         "App2",
		Version:      "1.0.0",
		Status: "OK",
		Integrations: make([]Integration, 0),
	}

	for _, service := range h.checks {
		integration := Integration{
			Name:   service.Name,
			Port:   service.Port,
			Host:   service.Host,
			Status: "OK",
		}

		start := time.Now()
		err := service.Check()
		if err != nil {
			integration.Error = err.Error()
			integration.Status = "NOK"
			app.Status = "NOK"
		}
		integration.ResponseTime = time.Now().Sub(start).Seconds()

		app.Integrations = append(app.Integrations, integration)
	}

	return app
}

func (h *Health) HandlerFunc(w http.ResponseWriter, r *http.Request) {
	c := h.Check()

	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(c)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	code := http.StatusOK
	if c.Status == "NOK" {
		code = http.StatusServiceUnavailable
	}

	w.WriteHeader(code)
	w.Write(data)
}

type (
	Checker func() error

	Config struct {
		Name  string
		Port  string
		Host  string
		Check Checker
	}

	Option func(*Health) error
)

func New(opts ...Option) (*Health, error) {
	h := &Health{
		checks: make(map[string]Config),
	}

	for _, opt := range opts {
		if err := opt(h); err != nil {
			return nil, err
		}
	}

	return h, nil
}

func Register(c Config) Option {
	return func(h *Health) error {
		if c.Name == "" {
			return errors.New("health check must have a name to be registered")
		}

		h.checks[c.Name] = c
		return nil
	}
}
