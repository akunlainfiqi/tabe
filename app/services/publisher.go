package services

import "time"

type TenantPaidPayload struct {
	TenantID  string    `json:"tenant_id"`
	ProductID string    `json:"product_id"`
	Timestamp time.Time `json:"timestamp"`
}

type Publisher interface {
	Publish(topicID string, data []byte) error
}
