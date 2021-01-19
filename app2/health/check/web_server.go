package check

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

type WebServer struct {
	URL string
}

func (ws WebServer) Check() error {
	client := http.Client{
		Timeout: time.Second * time.Duration(5),
	}
	request, _ := http.NewRequest(http.MethodGet, ws.URL, nil)

	res, err := client.Do(request)
	if err != nil {
		log.Printf("WebService health check failed: %w", err)
		return err
	}

	if res.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("Expected: %d, Received: %d", http.StatusOK, res.StatusCode))
		log.Printf("WebService health check failed: %w", err)
		return err
	}

	return nil
}
