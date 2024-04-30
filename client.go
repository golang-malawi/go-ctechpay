package ctechpay

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"net/http"
	"net/url"
	"time"
)

const ProductionURL = "https://api-sandbox.ctechpay.com/"
const SandboxURL = "https://api-sandbox.ctechpay.com/"
const defaultCancelText = "Cancel Payment"

// Client is the client type for interacting with the CTechPay API
type Client struct {
	APIToken    string
	BaseURL     string
	httpClient  *http.Client
	redirectURL string
	cancelURL   string
	cancelText  string
	Logger      *slog.Logger
}

// OrderRequest initiates an order request
type OrderRequest struct {
	Token              string    `json:"token" form:"token"`
	Amount             big.Float `json:"amount" form:"amount"`
	MerchantAttributes bool      `json:"merchantAttributes,omitempty" form:"merchantAttributes"` //	Merchant attributes enable you to include your custom options when the customer is interacting with the payment page. Set to true to enable the following compulsory and optional parameters: (redirectUrl, cancelUrl, cancelText) 	true
	RedirectUrl        string    `json:"redirectUrl,omitempty" form:"redirectUrl"`               //	A compulsory URL to redirect the user after a payment has been made. 	https://example.com/ Note: Only https requests are supported
	CancelUrl          string    `json:"cancelUrl,omitempty" form:"cancelUrl"`                   // A compulsory URL to redirect the user when they choose to cancel the payment. 	https://example.com/ Note: Only https requests are supported
	CancelText         string    `json:"cancelText,omitempty" form:"cancelText"`                 //	An optional custom text on the cancelUrl link 	“Go back to shop,” “Cancel Payment”
}

// OrderResponse contains information about order response from CTechPay
type OrderResponse struct {
	TxnID          string `json:"-"`
	OrderReference string `json:"order_reference"`
	PaymentPageURL string `json:"payment_page_URL"`
}

// NewClient creates a new Production client for CTechPay
func NewClient(token string, timeout time.Duration) *Client {
	return &Client{
		APIToken: token,
		BaseURL:  ProductionURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		Logger: slog.Default(),
	}
}

// NewSandboxClient creates a new Sandbox/Testing client for CTechPay
func NewSandboxClient(token string, timeout time.Duration) *Client {
	return &Client{
		APIToken: token,
		BaseURL:  SandboxURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		Logger: slog.Default(),
	}
}

// SetRedirectURL sets the redirect url for subsequent order requests
func (c *Client) SetRedirectURL(redirectURL string) {
	c.redirectURL = redirectURL
}

// SetCancelURL sets the cancel url and cancel text for subsequent order requests. cancel text defaults to "Cancel Payment"
func (c *Client) SetCancelURL(cancelURL, cancelText string) {
	c.cancelURL = cancelURL
	if cancelText != "" {
		c.cancelText = cancelText
	} else {
		c.cancelText = defaultCancelText
	}
}

// InitiateCardOrder initiates a payment transaction that will be fulfilled on the CTechPay platform
func (c *Client) InitiateCardOrder(txnID string, amount big.Float, merchant bool) (*OrderResponse, error) {
	form := url.Values{}
	form.Set("token", c.APIToken)
	form.Set("amount", amount.String())

	if merchant {
		if c.redirectURL == "" {
			return nil, fmt.Errorf("redirectURL must be set via *Client.SetRedirectURL(string)")
		}

		if c.cancelURL == "" {
			return nil, fmt.Errorf("cancelURL must be set via *Client.SetCancelURL(string, string)")
		}

		if c.cancelText == "" {
			return nil, fmt.Errorf("cancelText must be set via *Client.SetCancelURL(string, string)")
		}

		form.Set("merchantAttributes", "true")
		form.Set("redirectUrl", c.redirectURL)
		form.Set("cancelUrl", c.cancelURL)
		form.Set("cancelText", c.cancelText)
	}

	url := fmt.Sprintf("%s/?endpoint=order", c.BaseURL)
	c.Logger.Debug("Sending request to CTechPay at", "url", url)
	resp, err := c.httpClient.PostForm(url, form)
	if err != nil {
		return nil, fmt.Errorf("failed to create order, got: %w", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body, got: %w", err)
	}
	var orderResponse OrderResponse
	err = json.Unmarshal(data, &orderResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to parse order response as JSON %w", err)
	}

	orderResponse.TxnID = txnID
	return &orderResponse, nil
}
