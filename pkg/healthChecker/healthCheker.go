package healthChecker

import (
	"bytes"
	"encoding/json"
	apiRepo "github.com/rzabhd80/healthCheck/api/healthCheckApi/repository"
	"github.com/rzabhd80/healthCheck/helpers"
	"github.com/rzabhd80/healthCheck/models"
	healthCheckRepo "github.com/rzabhd80/healthCheck/pkg/healthChecker/repository"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type HealthChecker struct {
	apiRepo         apiRepo.APIRepository
	healthCheckRepo healthCheckRepo.HealthCheckRepository
	slackWebhookURL string
	checkInterval   time.Duration
	stopChan        chan struct{}
	wg              sync.WaitGroup
}

func NewHealthChecker(apiRepo apiRepo.APIRepository, slackWebhookURL string, checkInterval time.Duration,
	healthCheckRepo healthCheckRepo.HealthCheckRepository) *HealthChecker {
	return &HealthChecker{
		apiRepo:         apiRepo,
		healthCheckRepo: healthCheckRepo,
		slackWebhookURL: slackWebhookURL,
		checkInterval:   checkInterval,
		stopChan:        make(chan struct{}),
	}
}

func (h *HealthChecker) Start(poolSize int) {
	h.wg.Add(1)
	go func() {
		defer h.wg.Done()
		ticker := time.NewTicker(h.checkInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				h.runHealthChecks(poolSize)
			case <-h.stopChan:
				return
			}
		}
	}()
}

func (h *HealthChecker) Stop() {
	close(h.stopChan)
	h.wg.Wait()
}

func (h *HealthChecker) runHealthChecks(poolSize int) {
	apis, err := h.apiRepo.GetAllActive()
	if err != nil {
		log.Printf("Failed to get active APIs: %v", err)
		return
	}

	sem := make(chan struct{}, poolSize)
	var wg sync.WaitGroup

	for _, api := range apis {
		wg.Add(1)
		sem <- struct{}{}

		go func(api models.API) {
			defer wg.Done()
			h.checkAPIHealth(api)
			<-sem
		}(api)
	}

	wg.Wait()
	return
}

func (h *HealthChecker) checkAPIHealth(api models.API) {
	client := &http.Client{}
	req, err := http.NewRequest(api.RequestMethod, api.RequestURL, bytes.NewBuffer([]byte(*api.RequestBody)))
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range helpers.ParseHeaders(*api.RequestHeaders) {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to send request: %v", err)
		err := h.healthCheckRepo.LogHealthCheck(api.ID, api.RequestURL, strconv.Itoa(http.StatusServiceUnavailable))
		if err != nil {
			log.Printf("%s", "could not update the health check on the database")
		}
		go h.notifySlack(api, http.StatusServiceUnavailable)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("%s", "could not write the health check result")
		}
	}(resp.Body)

	err = h.healthCheckRepo.LogHealthCheck(api.ID, api.RequestURL, strconv.Itoa(resp.StatusCode))
	if err != nil {
		log.Printf("%s", "could not write the health check result")
	}
	if *api.LastStatus != resp.StatusCode {
		*api.LastStatus = resp.StatusCode
		err := h.apiRepo.Update(&api)

		if err != nil {
			log.Printf("%s", "could not update the health status in the database")
		}
		go h.notifySlack(api, resp.StatusCode)
	}
}

func (h *HealthChecker) notifySlack(api models.API, status int) {
	payload := map[string]interface{}{
		"text": "API Health Check Alert",
		"attachments": []map[string]interface{}{
			{
				"title": "API Health Status Change",
				"color": "#ff0000",
				"fields": []map[string]interface{}{
					{
						"title": "API ID",
						"value": api.ID,
						"short": true,
					},
					{
						"title": "New Status",
						"value": status,
						"short": true,
					},
				},
			},
		},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal payload: %v", err)
		return
	}

	req, err := http.NewRequest("POST", h.slackWebhookURL, bytes.NewBuffer(data))
	if err != nil {
		log.Printf("Failed to create Slack webhook request: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to send Slack webhook request: %v", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("%s", "")
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected status code from Slack webhook: %v", resp.StatusCode)
	}
}
