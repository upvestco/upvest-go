package upvest

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

// Webhook represents an Upvest webhook
type Webhook struct {
	ID           string    `json:"id"`
	URL          string    `json:"url"`
	Name         []Balance `json:"name"`
	HMACSecret   string    `json:"hmac_secret"`
	Headers      string    `json:"headers"`
	Version      string    `json:"version"`
	Status       string    `json:"status"`
	EventFilters []string  `json:"event_filters"`
}

// WebhookService handles operations related to the webhook
type WebhookService struct {
	service
}

// WebhookList is a list object for webhooks.
type WebhookList struct {
	Meta   ListMeta
	Values []Webhook `json:"results"`
}

// Create creates a new webhook
// Unlike other resource enndpoints, we can use the same webhook struct to create a new one
// as the parameters required to create a new one is basically all the fields in the webhook struct.
// Only difference being that it has not yet been saved on Upvest backend
// TODO: validate params
func (s *WebhookService) Create(wh *Webhook) (*Webhook, error) {
	u := "/tenancy/webhooks/"
	webhook := &Webhook{}
	p := NewParams(s.auth)
	err := s.client.Call(http.MethodPost, u, wh, webhook, p)
	return webhook, err
}

// Get retrives and returns a webhook object.
func (s *WebhookService) Get(webhookID string) (*Webhook, error) {
	u := fmt.Sprintf("/tenancy/webhooks/%s", webhookID)
	webhook := &Webhook{}
	p := NewParams(s.auth)
	err := s.client.Call(http.MethodGet, u, nil, webhook, p)
	return webhook, err
}

// List returns list of all webhooks.
func (s *WebhookService) List() (*WebhookList, error) {
	path := "/tenancy/webhooks/"
	u, _ := url.Parse(path)
	p := &Params{}
	p.SetAuthProvider(s.auth)

	var results []Webhook
	webhooks := &WebhookList{}

	for {
		err := s.client.Call(http.MethodGet, u.String(), nil, webhooks, p)
		if err != nil {
			return nil, errors.Wrap(err, "Could not retrieve list of webhooks")
		}
		results = append(results, webhooks.Values...)

		// append page_size param to the returned params
		u1, err := url.Parse(webhooks.Meta.Next)
		q := u1.Query()
		q.Set("page_size", string(MaxPageSize))
		u.RawQuery = q.Encode()
		if webhooks.Meta.Next == "" {
			break
		}
	}

	return &WebhookList{Values: results}, nil
}

// ListN returns a specific number of webhooks
func (s *WebhookService) ListN(count int) (*WebhookList, error) {
	path := "/tenancy/webhooks/"
	u, _ := url.Parse(path)

	p := &Params{}
	p.SetAuthProvider(s.auth)
	var results []Webhook
	webhooks := &WebhookList{}

	total := 0

	for total <= count {
		err := s.client.Call(http.MethodGet, u.String(), nil, webhooks, p)
		if err != nil {
			return nil, errors.Wrap(err, "Could not retrieve list of webhooks")
		}
		results = append(results, webhooks.Values...)
		total += len(webhooks.Values)

		// append page_size param to the returned params
		u1, err := url.Parse(webhooks.Meta.Next)
		q := u1.Query()
		q.Set("page_size", string(MaxPageSize))
		u.RawQuery = q.Encode()
		if webhooks.Meta.Next == "" {
			break
		}
	}

	return &WebhookList{Values: results}, nil
}

// Delete permanently deletes a webhook
func (s *WebhookService) Delete(webhookID string) error {
	u := fmt.Sprintf("/tenancy/webhooks/%s", webhookID)
	resp := &Response{}
	p := NewParams(s.auth)
	err := s.client.Call(http.MethodDelete, u, map[string]string{}, resp, p)
	return err
}

// Verfy a webhook
func (s *WebhookService) Verify(url string) error {
	u := fmt.Sprintf("/tenancy/webhooks-verify/")
	body := map[string]string{"veiry_url": url}
	resp := &Response{}
	p := NewParams(s.auth)
	err := s.client.Call(http.MethodPost, u, body, resp, p)
	return err
}
