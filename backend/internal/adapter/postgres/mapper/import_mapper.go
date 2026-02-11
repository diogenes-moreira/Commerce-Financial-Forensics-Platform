package mapper

import (
	"encoding/json"

	"github.com/diogenes/costforensics/backend/internal/adapter/postgres/model"
	"github.com/diogenes/costforensics/backend/internal/domain/importjob"
)

func ImportJobToModel(j *importjob.ImportJob) *model.ImportJob {
	errorLogJSON, _ := json.Marshal(j.ErrorLog())
	return &model.ImportJob{
		ID:            j.ID(),
		Name:          j.Name(),
		Source:        string(j.Source()),
		EntityType:    j.EntityType(),
		Status:        string(j.Status()),
		TotalRows:     j.TotalRows(),
		ProcessedRows: j.ProcessedRows(),
		FailedRows:    j.FailedRows(),
		ErrorLog:      string(errorLogJSON),
		SourceURI:     j.SourceURI(),
		StartedAt:     j.StartedAt(),
		CompletedAt:   j.CompletedAt(),
		CreatedAt:     j.CreatedAt(),
	}
}

func ImportJobToDomain(m *model.ImportJob) *importjob.ImportJob {
	var errorLog []string
	_ = json.Unmarshal([]byte(m.ErrorLog), &errorLog)
	if errorLog == nil {
		errorLog = make([]string, 0)
	}
	return importjob.HydrateImportJob(
		m.ID, m.Name, importjob.ImportSource(m.Source), m.EntityType,
		importjob.ImportStatus(m.Status), m.TotalRows, m.ProcessedRows, m.FailedRows,
		errorLog, m.SourceURI,
		m.StartedAt, m.CompletedAt, m.CreatedAt,
	)
}
