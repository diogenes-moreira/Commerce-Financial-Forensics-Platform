package costreport

import (
	"context"
	"time"

	"github.com/diogenes/costforensics/backend/internal/domain/costreport"
	"github.com/diogenes/costforensics/backend/internal/domain/costrecord"
	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/google/uuid"
)

type Service struct {
	reportRepo     costreport.Repository
	costRecordRepo costrecord.Repository
}

func NewService(reportRepo costreport.Repository, costRecordRepo costrecord.Repository) *Service {
	return &Service{reportRepo: reportRepo, costRecordRepo: costRecordRepo}
}

func (s *Service) Generate(ctx context.Context, name string, reportType costreport.ReportType, periodStart, periodEnd time.Time, currency string) (*costreport.CostReport, error) {
	report, err := costreport.NewCostReport(name, reportType, periodStart, periodEnd, currency)
	if err != nil {
		return nil, err
	}

	// Fetch cost records for the period
	filter := costrecord.ListFilter{
		StartDate:  &periodStart,
		EndDate:    &periodEnd,
		Pagination: shared.NewPagination(1, shared.MaxPageSize),
	}
	result, err := s.costRecordRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	for _, cr := range result.Items {
		report.AddServiceCost(cr.Service(), cr.Amount())
	}

	if err := s.reportRepo.Create(ctx, report); err != nil {
		return nil, err
	}
	return report, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*costreport.CostReport, error) {
	return s.reportRepo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]*costreport.CostReport, error) {
	return s.reportRepo.List(ctx)
}
