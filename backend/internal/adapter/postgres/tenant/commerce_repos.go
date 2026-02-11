package tenant

import (
	"context"
	"errors"

	"github.com/diogenes/costforensics/backend/internal/adapter/postgres/mapper"
	"github.com/diogenes/costforensics/backend/internal/adapter/postgres/model"
	"github.com/diogenes/costforensics/backend/internal/domain/customer"
	"github.com/diogenes/costforensics/backend/internal/domain/order"
	"github.com/diogenes/costforensics/backend/internal/domain/product"
	"github.com/diogenes/costforensics/backend/internal/domain/seller"
	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProductRepo

type ProductRepo struct {
	db *gorm.DB
}

func NewProductRepo(db *gorm.DB) *ProductRepo {
	return &ProductRepo{db: db}
}

func (r *ProductRepo) Create(ctx context.Context, p *product.Product) error {
	return r.db.WithContext(ctx).Create(mapper.ProductToModel(p)).Error
}

func (r *ProductRepo) GetByID(ctx context.Context, id uuid.UUID) (*product.Product, error) {
	var m model.Product
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, product.ErrProductNotFound
		}
		return nil, err
	}
	return mapper.ProductToDomain(&m), nil
}

func (r *ProductRepo) GetBySKU(ctx context.Context, sku string) (*product.Product, error) {
	var m model.Product
	if err := r.db.WithContext(ctx).Where("sku = ?", sku).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, product.ErrProductNotFound
		}
		return nil, err
	}
	return mapper.ProductToDomain(&m), nil
}

func (r *ProductRepo) List(ctx context.Context, filter product.ListFilter) (*shared.PagedResult[*product.Product], error) {
	q := r.db.WithContext(ctx).Model(&model.Product{})

	if filter.Category != "" {
		q = q.Where("category = ?", filter.Category)
	}
	if filter.Status != "" {
		q = q.Where("status = ?", string(filter.Status))
	}

	var total int64
	q.Count(&total)

	var models []model.Product
	if err := q.Offset(filter.Pagination.Offset()).Limit(filter.Pagination.PageSize()).
		Order("name ASC").Find(&models).Error; err != nil {
		return nil, err
	}

	items := make([]*product.Product, len(models))
	for i := range models {
		items[i] = mapper.ProductToDomain(&models[i])
	}

	return &shared.PagedResult[*product.Product]{
		Items:      items,
		TotalCount: total,
		Page:       filter.Pagination.Page(),
		PageSize:   filter.Pagination.PageSize(),
	}, nil
}

func (r *ProductRepo) Update(ctx context.Context, p *product.Product) error {
	return r.db.WithContext(ctx).Save(mapper.ProductToModel(p)).Error
}

// SellerRepo

type SellerRepo struct {
	db *gorm.DB
}

func NewSellerRepo(db *gorm.DB) *SellerRepo {
	return &SellerRepo{db: db}
}

func (r *SellerRepo) Create(ctx context.Context, s *seller.Seller) error {
	return r.db.WithContext(ctx).Create(mapper.SellerToModel(s)).Error
}

func (r *SellerRepo) GetByID(ctx context.Context, id uuid.UUID) (*seller.Seller, error) {
	var m model.Seller
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, seller.ErrSellerNotFound
		}
		return nil, err
	}
	return mapper.SellerToDomain(&m), nil
}

func (r *SellerRepo) GetByExternalID(ctx context.Context, externalID string) (*seller.Seller, error) {
	var m model.Seller
	if err := r.db.WithContext(ctx).Where("external_id = ?", externalID).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, seller.ErrSellerNotFound
		}
		return nil, err
	}
	return mapper.SellerToDomain(&m), nil
}

func (r *SellerRepo) List(ctx context.Context, filter seller.ListFilter) (*shared.PagedResult[*seller.Seller], error) {
	q := r.db.WithContext(ctx).Model(&model.Seller{})

	if filter.Status != "" {
		q = q.Where("status = ?", string(filter.Status))
	}

	var total int64
	q.Count(&total)

	var models []model.Seller
	if err := q.Offset(filter.Pagination.Offset()).Limit(filter.Pagination.PageSize()).
		Order("name ASC").Find(&models).Error; err != nil {
		return nil, err
	}

	items := make([]*seller.Seller, len(models))
	for i := range models {
		items[i] = mapper.SellerToDomain(&models[i])
	}

	return &shared.PagedResult[*seller.Seller]{
		Items:      items,
		TotalCount: total,
		Page:       filter.Pagination.Page(),
		PageSize:   filter.Pagination.PageSize(),
	}, nil
}

func (r *SellerRepo) Update(ctx context.Context, s *seller.Seller) error {
	return r.db.WithContext(ctx).Save(mapper.SellerToModel(s)).Error
}

// CustomerRepo

type CustomerRepo struct {
	db *gorm.DB
}

func NewCustomerRepo(db *gorm.DB) *CustomerRepo {
	return &CustomerRepo{db: db}
}

func (r *CustomerRepo) Create(ctx context.Context, c *customer.Customer) error {
	return r.db.WithContext(ctx).Create(mapper.CustomerToModel(c)).Error
}

func (r *CustomerRepo) GetByID(ctx context.Context, id uuid.UUID) (*customer.Customer, error) {
	var m model.Customer
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, customer.ErrCustomerNotFound
		}
		return nil, err
	}
	return mapper.CustomerToDomain(&m), nil
}

func (r *CustomerRepo) GetByExternalID(ctx context.Context, externalID string) (*customer.Customer, error) {
	var m model.Customer
	if err := r.db.WithContext(ctx).Where("external_id = ?", externalID).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, customer.ErrCustomerNotFound
		}
		return nil, err
	}
	return mapper.CustomerToDomain(&m), nil
}

func (r *CustomerRepo) List(ctx context.Context, filter customer.ListFilter) (*shared.PagedResult[*customer.Customer], error) {
	q := r.db.WithContext(ctx).Model(&model.Customer{})

	if filter.Segment != "" {
		q = q.Where("segment = ?", filter.Segment)
	}

	var total int64
	q.Count(&total)

	var models []model.Customer
	if err := q.Offset(filter.Pagination.Offset()).Limit(filter.Pagination.PageSize()).
		Order("name ASC").Find(&models).Error; err != nil {
		return nil, err
	}

	items := make([]*customer.Customer, len(models))
	for i := range models {
		items[i] = mapper.CustomerToDomain(&models[i])
	}

	return &shared.PagedResult[*customer.Customer]{
		Items:      items,
		TotalCount: total,
		Page:       filter.Pagination.Page(),
		PageSize:   filter.Pagination.PageSize(),
	}, nil
}

func (r *CustomerRepo) Update(ctx context.Context, c *customer.Customer) error {
	return r.db.WithContext(ctx).Save(mapper.CustomerToModel(c)).Error
}

// OrderRepo

type OrderRepo struct {
	db *gorm.DB
}

func NewOrderRepo(db *gorm.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

func (r *OrderRepo) Create(ctx context.Context, o *order.Order) error {
	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Create(mapper.OrderToModel(o)).Error; err != nil {
		tx.Rollback()
		return err
	}
	for _, item := range o.Items() {
		if err := tx.Create(mapper.OrderItemToModel(item)).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func (r *OrderRepo) GetByID(ctx context.Context, id uuid.UUID) (*order.Order, error) {
	var m model.Order
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, order.ErrOrderNotFound
		}
		return nil, err
	}
	items, err := r.loadItems(ctx, m.ID)
	if err != nil {
		return nil, err
	}
	return mapper.OrderToDomain(&m, items), nil
}

func (r *OrderRepo) GetByExternalID(ctx context.Context, externalID string) (*order.Order, error) {
	var m model.Order
	if err := r.db.WithContext(ctx).Where("external_id = ?", externalID).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, order.ErrOrderNotFound
		}
		return nil, err
	}
	items, err := r.loadItems(ctx, m.ID)
	if err != nil {
		return nil, err
	}
	return mapper.OrderToDomain(&m, items), nil
}

func (r *OrderRepo) List(ctx context.Context, filter order.ListFilter) (*shared.PagedResult[*order.Order], error) {
	q := r.db.WithContext(ctx).Model(&model.Order{})

	if filter.SellerID != nil {
		q = q.Where("seller_id = ?", *filter.SellerID)
	}
	if filter.CustomerID != nil {
		q = q.Where("customer_id = ?", *filter.CustomerID)
	}
	if filter.Status != "" {
		q = q.Where("status = ?", string(filter.Status))
	}
	if filter.StartDate != nil {
		q = q.Where("order_date >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		q = q.Where("order_date <= ?", *filter.EndDate)
	}

	var total int64
	q.Count(&total)

	var models []model.Order
	if err := q.Offset(filter.Pagination.Offset()).Limit(filter.Pagination.PageSize()).
		Order("order_date DESC").Find(&models).Error; err != nil {
		return nil, err
	}

	items := make([]*order.Order, len(models))
	for i := range models {
		orderItems, err := r.loadItems(ctx, models[i].ID)
		if err != nil {
			return nil, err
		}
		items[i] = mapper.OrderToDomain(&models[i], orderItems)
	}

	return &shared.PagedResult[*order.Order]{
		Items:      items,
		TotalCount: total,
		Page:       filter.Pagination.Page(),
		PageSize:   filter.Pagination.PageSize(),
	}, nil
}

func (r *OrderRepo) Update(ctx context.Context, o *order.Order) error {
	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Save(mapper.OrderToModel(o)).Error; err != nil {
		tx.Rollback()
		return err
	}
	// Delete existing items and re-create
	if err := tx.Where("order_id = ?", o.ID()).Delete(&model.OrderItem{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	for _, item := range o.Items() {
		if err := tx.Create(mapper.OrderItemToModel(item)).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func (r *OrderRepo) loadItems(ctx context.Context, orderID uuid.UUID) ([]*order.OrderItem, error) {
	var models []model.OrderItem
	if err := r.db.WithContext(ctx).Where("order_id = ?", orderID).Find(&models).Error; err != nil {
		return nil, err
	}
	items := make([]*order.OrderItem, len(models))
	for i := range models {
		items[i] = mapper.OrderItemToDomain(&models[i])
	}
	return items, nil
}
