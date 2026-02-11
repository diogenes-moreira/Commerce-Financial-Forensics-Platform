package connector

import (
	"context"
	"fmt"
	"time"
)

type MedusaJSConnector struct {
	apiURL   string
	apiToken string
}

func NewMedusaJSConnector(config map[string]string) *MedusaJSConnector {
	return &MedusaJSConnector{
		apiURL:   config["api_url"],
		apiToken: config["api_token"],
	}
}

func (c *MedusaJSConnector) FetchProducts(_ context.Context, _ time.Time) ([]RawProduct, error) {
	return nil, fmt.Errorf("medusajs: FetchProducts not yet implemented")
}

func (c *MedusaJSConnector) FetchOrders(_ context.Context, _ time.Time) ([]RawOrder, error) {
	return nil, fmt.Errorf("medusajs: FetchOrders not yet implemented")
}

func (c *MedusaJSConnector) FetchCustomers(_ context.Context, _ time.Time) ([]RawCustomer, error) {
	return nil, fmt.Errorf("medusajs: FetchCustomers not yet implemented")
}

func (c *MedusaJSConnector) HandleWebhook(_ context.Context, _ string, _ []byte) error {
	return fmt.Errorf("medusajs: HandleWebhook not yet implemented")
}
