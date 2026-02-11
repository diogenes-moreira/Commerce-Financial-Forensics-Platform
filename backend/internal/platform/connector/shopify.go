package connector

import (
	"context"
	"fmt"
	"time"
)

type ShopifyConnector struct {
	shopDomain string
	shopToken  string
}

func NewShopifyConnector(config map[string]string) *ShopifyConnector {
	return &ShopifyConnector{
		shopDomain: config["shop_domain"],
		shopToken:  config["shop_token"],
	}
}

func (c *ShopifyConnector) FetchProducts(_ context.Context, _ time.Time) ([]RawProduct, error) {
	return nil, fmt.Errorf("shopify: FetchProducts not yet implemented")
}

func (c *ShopifyConnector) FetchOrders(_ context.Context, _ time.Time) ([]RawOrder, error) {
	return nil, fmt.Errorf("shopify: FetchOrders not yet implemented")
}

func (c *ShopifyConnector) FetchCustomers(_ context.Context, _ time.Time) ([]RawCustomer, error) {
	return nil, fmt.Errorf("shopify: FetchCustomers not yet implemented")
}

func (c *ShopifyConnector) HandleWebhook(_ context.Context, _ string, _ []byte) error {
	return fmt.Errorf("shopify: HandleWebhook not yet implemented")
}
