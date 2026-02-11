package response

import (
	"github.com/diogenes/costforensics/backend/internal/domain/importjob"
)

type ImportJob struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Source        string   `json:"source"`
	EntityType    string   `json:"entity_type"`
	Status        string   `json:"status"`
	TotalRows     int      `json:"total_rows"`
	ProcessedRows int      `json:"processed_rows"`
	FailedRows    int      `json:"failed_rows"`
	ErrorLog      []string `json:"error_log"`
	SourceURI     string   `json:"source_uri"`
	StartedAt     *string  `json:"started_at"`
	CompletedAt   *string  `json:"completed_at"`
	CreatedAt     string   `json:"created_at"`
}

func ImportJobFromDomain(j *importjob.ImportJob) ImportJob {
	var startedAt *string
	if j.StartedAt() != nil {
		s := j.StartedAt().Format("2006-01-02T15:04:05Z07:00")
		startedAt = &s
	}
	var completedAt *string
	if j.CompletedAt() != nil {
		s := j.CompletedAt().Format("2006-01-02T15:04:05Z07:00")
		completedAt = &s
	}
	return ImportJob{
		ID: j.ID().String(), Name: j.Name(),
		Source: string(j.Source()), EntityType: j.EntityType(),
		Status: string(j.Status()), TotalRows: j.TotalRows(),
		ProcessedRows: j.ProcessedRows(), FailedRows: j.FailedRows(),
		ErrorLog: j.ErrorLog(), SourceURI: j.SourceURI(),
		StartedAt: startedAt, CompletedAt: completedAt,
		CreatedAt: j.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}
}
