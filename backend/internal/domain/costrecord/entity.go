package costrecord

import (
	"time"

	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/google/uuid"
)

type Category string

const (
	CategoryCompute    Category = "compute"
	CategoryStorage    Category = "storage"
	CategoryNetwork    Category = "network"
	CategoryDatabase   Category = "database"
	CategoryAnalytics  Category = "analytics"
	CategoryML         Category = "ml"
	CategoryOther      Category = "other"
)

type CostRecord struct {
	id             uuid.UUID
	cloudAccountID uuid.UUID
	service        string
	category       Category
	amount         shared.Money
	usageDate      time.Time
	tags           map[string]string
	createdAt      time.Time
}

func NewCostRecord(cloudAccountID uuid.UUID, service string, amount shared.Money, usageDate time.Time) (*CostRecord, error) {
	if service == "" {
		return nil, ErrServiceEmpty
	}
	if usageDate.IsZero() {
		return nil, ErrInvalidUsageDate
	}

	return &CostRecord{
		id:             uuid.New(),
		cloudAccountID: cloudAccountID,
		service:        service,
		category:       categorize(service),
		amount:         amount,
		usageDate:      usageDate,
		tags:           make(map[string]string),
		createdAt:      time.Now().UTC(),
	}, nil
}

func HydrateCostRecord(
	id, cloudAccountID uuid.UUID, service string, category Category,
	amount shared.Money, usageDate time.Time, tags map[string]string, createdAt time.Time,
) *CostRecord {
	return &CostRecord{
		id: id, cloudAccountID: cloudAccountID, service: service,
		category: category, amount: amount, usageDate: usageDate,
		tags: tags, createdAt: createdAt,
	}
}

func (r *CostRecord) ID() uuid.UUID             { return r.id }
func (r *CostRecord) CloudAccountID() uuid.UUID  { return r.cloudAccountID }
func (r *CostRecord) Service() string             { return r.service }
func (r *CostRecord) Category() Category          { return r.category }
func (r *CostRecord) Amount() shared.Money        { return r.amount }
func (r *CostRecord) UsageDate() time.Time        { return r.usageDate }
func (r *CostRecord) Tags() map[string]string     { return r.tags }
func (r *CostRecord) CreatedAt() time.Time        { return r.createdAt }

func (r *CostRecord) AddTag(key, value string) {
	r.tags[key] = value
}

// categorize auto-assigns a category based on service name patterns.
func categorize(service string) Category {
	// Simple keyword-based categorization; expand as needed.
	switch {
	case contains(service, "ec2", "lambda", "fargate", "compute", "functions"):
		return CategoryCompute
	case contains(service, "s3", "storage", "ebs", "gcs"):
		return CategoryStorage
	case contains(service, "vpc", "cloudfront", "cdn", "network", "elb"):
		return CategoryNetwork
	case contains(service, "rds", "dynamodb", "aurora", "spanner", "sql"):
		return CategoryDatabase
	case contains(service, "athena", "bigquery", "redshift", "analytics"):
		return CategoryAnalytics
	case contains(service, "sagemaker", "ml", "ai platform"):
		return CategoryML
	default:
		return CategoryOther
	}
}

func contains(s string, substrs ...string) bool {
	lower := toLower(s)
	for _, sub := range substrs {
		if containsStr(lower, sub) {
			return true
		}
	}
	return false
}

func toLower(s string) string {
	b := make([]byte, len(s))
	for i := range s {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		b[i] = c
	}
	return string(b)
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && searchStr(s, sub)
}

func searchStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
