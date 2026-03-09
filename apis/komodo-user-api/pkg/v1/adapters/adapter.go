package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"komodo-user-api/pkg/v1/models"
)

// Client calls the user-api internal server on behalf of sibling services.
//
// Pass the base URL from the calling service's own config:
//
//	local:  http://user-api-public:7052   (shared network namespace via docker-compose)
//	aws:    http://user-api.komodo.internal:7052  (Cloud Map service discovery)
//
// All methods require a service-scoped JWT bearer token, which the calling
// service obtains from komodo-auth-api before making requests.
type Client struct {
	baseURL string
	http    *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		http:    &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) GetUser(ctx context.Context, userID, token string) (*models.User, error) {
	var out models.User
	if err := c.get(ctx, fmt.Sprintf("/users/%s", userID), token, &out); err != nil {
		return nil, fmt.Errorf("adapters.GetUser: %w", err)
	}
	return &out, nil
}

func (c *Client) GetAddresses(ctx context.Context, userID, token string) ([]models.Address, error) {
	var out []models.Address
	if err := c.get(ctx, fmt.Sprintf("/users/%s/addresses", userID), token, &out); err != nil {
		return nil, fmt.Errorf("adapters.GetAddresses: %w", err)
	}
	return out, nil
}

func (c *Client) GetPreferences(ctx context.Context, userID, token string) (*models.Preferences, error) {
	var out models.Preferences
	if err := c.get(ctx, fmt.Sprintf("/users/%s/preferences", userID), token, &out); err != nil {
		return nil, fmt.Errorf("adapters.GetPreferences: %w", err)
	}
	return &out, nil
}

func (c *Client) GetPayments(ctx context.Context, userID, token string) ([]models.PaymentMethod, error) {
	var out []models.PaymentMethod
	if err := c.get(ctx, fmt.Sprintf("/users/%s/payments", userID), token, &out); err != nil {
		return nil, fmt.Errorf("adapters.GetPayments: %w", err)
	}
	return out, nil
}

func (c *Client) get(ctx context.Context, path, token string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(out)
}
