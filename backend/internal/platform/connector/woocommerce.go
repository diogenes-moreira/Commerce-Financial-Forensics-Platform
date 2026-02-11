package connector

import (
	"context"
	"fmt"
	"time"
)

type WooCommerceConnector struct {
	apiURL    string
	apiKey    string
	apiSecret string
}

func NewWooCommerceConnector(config map[string]string) *WooCommerceConnector {
	return &WooCommerceConnector{
		apiURL:    config["api_url"],
		apiKey:    config["api_key"],
		apiSecret: config["api_secret"],
	}
}

func (c *WooCommerceConnector) FetchProducts(_ context.Context, _ time.Time) ([]RawProduct, error) {
	return nil, fmt.Errorf("woocommerce: FetchProducts not yet implemented")
}

func (c *WooCommerceConnector) FetchOrders(_ context.Context, _ time.Time) ([]RawOrder, error) {
	return nil, fmt.Errorf("woocommerce: FetchOrders not yet implemented")
}

func (c *WooCommerceConnector) FetchCustomers(_ context.Context, _ time.Time) ([]RawCustomer, error) {
	return nil, fmt.Errorf("woocommerce: FetchCustomers not yet implemented")
}

func (c *WooCommerceConnector) HandleWebhook(_ context.Context, _ string, _ []byte) error {
	return fmt.Errorf("woocommerce: HandleWebhook not yet implemented")
}
