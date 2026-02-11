package tenant

import (
	"context"
	"errors"

	"github.com/diogenes/costforensics/backend/internal/adapter/postgres/mapper"
	"github.com/diogenes/costforensics/backend/internal/adapter/postgres/model"
	"github.com/diogenes/costforensics/backend/internal/domain/anomaly"
	"github.com/diogenes/costforensics/backend/internal/domain/budget"
	"github.com/diogenes/costforensics/backend/internal/domain/cloudaccount"
	"github.com/diogenes/costforensics/backend/internal/domain/costreport"
	"github.com/diogenes/costforensics/backend/internal/domain/costrecord"
	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CloudAccountRepo implements cloudaccount.Repository against a tenant DB.
type CloudAccountRepo struct {
	db *gorm.DB
}

func NewCloudAccountRepo(db *gorm.DB) *CloudAccountRepo {
	return &CloudAccountRepo{db: db}
}

func (r *CloudAccountRepo) Create(ctx context.Context, a *cloudaccount.CloudAccount) error {
	return r.db.WithContext(ctx).Create(mapper.CloudAccountToModel(a)).Error
}

func (r *CloudAccountRepo) GetByID(ctx context.Context, id uuid.UUID) (*cloudaccount.CloudAccount, error) {
	var m model.CloudAccount
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, cloudaccount.ErrAccountNotFound
		}
		return nil, err
	}
	return mapper.CloudAccountToDomain(&m), nil
}

func (r *CloudAccountRepo) List(ctx context.Context) ([]*cloudaccount.CloudAccount, error) {
	var models []model.CloudAccount
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]*cloudaccount.CloudAccount, len(models))
	for i := range models {
		result[i] = mapper.CloudAccountToDomain(&models[i])
	}
	return result, nil
}

func (r *CloudAccountRepo) Update(ctx context.Context, a *cloudaccount.CloudAccount) error {
	return r.db.WithContext(ctx).Save(mapper.CloudAccountToModel(a)).Error
}

func (r *CloudAccountRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.CloudAccount{}).Error
}

// CostRecordRepo implements costrecord.Repository.
type CostRecordRepo struct {
	db *gorm.DB
}

func NewCostRecordRepo(db *gorm.DB) *CostRecordRepo {
	return &CostRecordRepo{db: db}
}

func (r *CostRecordRepo) Create(ctx context.Context, cr *costrecord.CostRecord) error {
	return r.db.WithContext(ctx).Create(mapper.CostRecordToModel(cr)).Error
}

func (r *CostRecordRepo) GetByID(ctx context.Context, id uuid.UUID) (*costrecord.CostRecord, error) {
	var m model.CostRecord
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, costrecord.ErrRecordNotFound
		}
		return nil, err
	}
	return mapper.CostRecordToDomain(&m), nil
}

func (r *CostRecordRepo) List(ctx context.Context, filter costrecord.ListFilter) (*shared.PagedResult[*costrecord.CostRecord], error) {
	q := r.db.WithContext(ctx).Model(&model.CostRecord{})

	if filter.CloudAccountID != nil {
		q = q.Where("cloud_account_id = ?", *filter.CloudAccountID)
	}
	if filter.Service != "" {
		q = q.Where("service = ?", filter.Service)
	}
	if filter.Category != "" {
		q = q.Where("category = ?", string(filter.Category))
	}
	if filter.StartDate != nil {
		q = q.Where("usage_date >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		q = q.Where("usage_date <= ?", *filter.EndDate)
	}

	var total int64
	q.Count(&total)

	var models []model.CostRecord
	if err := q.Offset(filter.Pagination.Offset()).Limit(filter.Pagination.PageSize()).
		Order("usage_date DESC").Find(&models).Error; err != nil {
		return nil, err
	}

	items := make([]*costrecord.CostRecord, len(models))
	for i := range models {
		items[i] = mapper.CostRecordToDomain(&models[i])
	}

	return &shared.PagedResult[*costrecord.CostRecord]{
		Items:      items,
		TotalCount: total,
		Page:       filter.Pagination.Page(),
		PageSize:   filter.Pagination.PageSize(),
	}, nil
}

// BudgetRepo implements budget.Repository.
type BudgetRepo struct {
	db *gorm.DB
}

func NewBudgetRepo(db *gorm.DB) *BudgetRepo {
	return &BudgetRepo{db: db}
}

func (r *BudgetRepo) Create(ctx context.Context, b *budget.Budget) error {
	return r.db.WithContext(ctx).Create(mapper.BudgetToModel(b)).Error
}

func (r *BudgetRepo) GetByID(ctx context.Context, id uuid.UUID) (*budget.Budget, error) {
	var m model.Budget
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, budget.ErrBudgetNotFound
		}
		return nil, err
	}
	return mapper.BudgetToDomain(&m), nil
}

func (r *BudgetRepo) List(ctx context.Context) ([]*budget.Budget, error) {
	var models []model.Budget
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]*budget.Budget, len(models))
	for i := range models {
		result[i] = mapper.BudgetToDomain(&models[i])
	}
	return result, nil
}

func (r *BudgetRepo) Update(ctx context.Context, b *budget.Budget) error {
	return r.db.WithContext(ctx).Save(mapper.BudgetToModel(b)).Error
}

func (r *BudgetRepo) CreateAlert(ctx context.Context, a *budget.Alert) error {
	return r.db.WithContext(ctx).Create(mapper.BudgetAlertToModel(a)).Error
}

// AnomalyRepo implements anomaly.Repository.
type AnomalyRepo struct {
	db *gorm.DB
}

func NewAnomalyRepo(db *gorm.DB) *AnomalyRepo {
	return &AnomalyRepo{db: db}
}

func (r *AnomalyRepo) Create(ctx context.Context, a *anomaly.CostAnomaly) error {
	return r.db.WithContext(ctx).Create(mapper.AnomalyToModel(a)).Error
}

func (r *AnomalyRepo) GetByID(ctx context.Context, id uuid.UUID) (*anomaly.CostAnomaly, error) {
	var m model.Anomaly
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, anomaly.ErrAnomalyNotFound
		}
		return nil, err
	}
	return mapper.AnomalyToDomain(&m), nil
}

func (r *AnomalyRepo) List(ctx context.Context) ([]*anomaly.CostAnomaly, error) {
	var models []model.Anomaly
	if err := r.db.WithContext(ctx).Order("detected_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]*anomaly.CostAnomaly, len(models))
	for i := range models {
		result[i] = mapper.AnomalyToDomain(&models[i])
	}
	return result, nil
}

func (r *AnomalyRepo) Update(ctx context.Context, a *anomaly.CostAnomaly) error {
	return r.db.WithContext(ctx).Save(mapper.AnomalyToModel(a)).Error
}

// CostReportRepo implements costreport.Repository.
type CostReportRepo struct {
	db *gorm.DB
}

func NewCostReportRepo(db *gorm.DB) *CostReportRepo {
	return &CostReportRepo{db: db}
}

func (r *CostReportRepo) Create(ctx context.Context, cr *costreport.CostReport) error {
	return r.db.WithContext(ctx).Create(mapper.CostReportToModel(cr)).Error
}

func (r *CostReportRepo) GetByID(ctx context.Context, id uuid.UUID) (*costreport.CostReport, error) {
	var m model.CostReport
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, costreport.ErrReportNotFound
		}
		return nil, err
	}
	return mapper.CostReportToDomain(&m), nil
}

func (r *CostReportRepo) List(ctx context.Context) ([]*costreport.CostReport, error) {
	var models []model.CostReport
	if err := r.db.WithContext(ctx).Order("generated_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]*costreport.CostReport, len(models))
	for i := range models {
		result[i] = mapper.CostReportToDomain(&models[i])
	}
	return result, nil
}
