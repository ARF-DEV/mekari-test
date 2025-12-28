package paymentob

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/arf-dev/mekari-test/config"
	"github.com/arf-dev/mekari-test/model"
	"github.com/rs/zerolog/log"
)

type Outbound struct {
	config *config.Config
}

func New(config *config.Config) *Outbound {
	return &Outbound{
		config: config,
	}
}
func (service *Outbound) DoPayment(ctx context.Context, paymentRequest model.PaymentRequest) error {
	const (
		max_retry int           = 2
		timeoff   time.Duration = time.Second * 5
	)
	retryCount := max_retry
	retryTimeoff := timeoff
	for ; retryCount > 0; retryCount-- {
		paymentEndpoint := service.config.PAYMENT_GATEWAY_URL + "/v1/payments"
		requestBody, _ := json.Marshal(paymentRequest)
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			paymentEndpoint,
			bytes.NewBuffer(requestBody),
		)
		if err != nil {
			return err
		}

		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{Timeout: 30 * time.Second}

		log.Log().Msgf("hit endpoint %s", paymentEndpoint)
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			if resp.StatusCode == 429 || resp.StatusCode == 0 || (resp.StatusCode >= 500 && resp.StatusCode != 501) {
				// retry if status code is 429, 0, or >= 500 expect 501
				time.Sleep(retryTimeoff)
				retryTimeoff *= 2
				continue
			}
			return fmt.Errorf("unexpected status: %s", resp.Status)
		}

		// if successful then break out of retry
		break
	}
	return nil
}
