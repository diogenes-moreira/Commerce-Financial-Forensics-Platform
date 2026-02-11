package importjob

import (
	"context"
	"io"

	"github.com/diogenes/costforensics/backend/internal/domain/importjob"
	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/diogenes/costforensics/backend/internal/platform/importer"
	"github.com/google/uuid"
)

type Service struct {
	repo      importjob.Repository
	csvParser *importer.CSVParser
}

func NewService(repo importjob.Repository, csvParser *importer.CSVParser) *Service {
	return &Service{repo: repo, csvParser: csvParser}
}

func (s *Service) CreateFromUpload(ctx context.Context, name, entityType string, reader io.Reader) (*importjob.ImportJob, error) {
	job, err := importjob.NewImportJob(name, importjob.SourceCSVUpload, entityType, "upload")
	if err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, job); err != nil {
		return nil, err
	}

	rows, err := s.csvParser.ParseCSV(reader)
	if err != nil {
		job.Fail(err.Error())
		_ = s.repo.Update(ctx, job)
		return job, nil
	}

	if err := job.Begin(); err != nil {
		return nil, err
	}
	job.SetTotalRows(len(rows))
	_ = s.repo.Update(ctx, job)

	s.processRows(ctx, job, rows)

	job.Complete()
	if err := s.repo.Update(ctx, job); err != nil {
		return nil, err
	}

	return job, nil
}

func (s *Service) processRows(ctx context.Context, job *importjob.ImportJob, rows []map[string]string) {
	for i := range rows {
		// Placeholder: actual entity creation logic will be wired per entity type later.
		// For now, just count each row as successfully processed.
		job.IncrementProcessed()

		// Periodically persist progress (every 100 rows)
		if (i+1)%100 == 0 {
			_ = s.repo.Update(ctx, job)
		}
	}
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*importjob.ImportJob, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context, filter importjob.ListFilter) (*shared.PagedResult[*importjob.ImportJob], error) {
	return s.repo.List(ctx, filter)
}
