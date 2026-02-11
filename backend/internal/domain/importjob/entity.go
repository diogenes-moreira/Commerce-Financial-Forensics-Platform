package importjob

import (
	"time"

	"github.com/google/uuid"
)

type ImportSource string

const (
	SourceCSVUpload ImportSource = "csv_upload"
	SourceS3Bucket  ImportSource = "s3_bucket"
	SourceGCSBucket ImportSource = "gcs_bucket"
)

type ImportStatus string

const (
	StatusPending    ImportStatus = "pending"
	StatusProcessing ImportStatus = "processing"
	StatusCompleted  ImportStatus = "completed"
	StatusFailed     ImportStatus = "failed"
)

var allowedEntityTypes = map[string]bool{
	"products":     true,
	"orders":       true,
	"sellers":      true,
	"customers":    true,
	"cost_records": true,
}

type ImportJob struct {
	id            uuid.UUID
	name          string
	source        ImportSource
	entityType    string
	status        ImportStatus
	totalRows     int
	processedRows int
	failedRows    int
	errorLog      []string
	sourceURI     string
	startedAt     *time.Time
	completedAt   *time.Time
	createdAt     time.Time
}

func NewImportJob(name string, source ImportSource, entityType, sourceURI string) (*ImportJob, error) {
	if name == "" {
		return nil, ErrNameEmpty
	}
	if !allowedEntityTypes[entityType] {
		return nil, ErrInvalidEntityType
	}

	now := time.Now().UTC()
	return &ImportJob{
		id:         uuid.New(),
		name:       name,
		source:     source,
		entityType: entityType,
		status:     StatusPending,
		errorLog:   make([]string, 0),
		sourceURI:  sourceURI,
		createdAt:  now,
	}, nil
}

func HydrateImportJob(
	id uuid.UUID, name string, source ImportSource, entityType string,
	status ImportStatus, totalRows, processedRows, failedRows int,
	errorLog []string, sourceURI string,
	startedAt, completedAt *time.Time, createdAt time.Time,
) *ImportJob {
	if errorLog == nil {
		errorLog = make([]string, 0)
	}
	return &ImportJob{
		id: id, name: name, source: source, entityType: entityType,
		status: status, totalRows: totalRows, processedRows: processedRows,
		failedRows: failedRows, errorLog: errorLog, sourceURI: sourceURI,
		startedAt: startedAt, completedAt: completedAt, createdAt: createdAt,
	}
}

func (j *ImportJob) ID() uuid.UUID        { return j.id }
func (j *ImportJob) Name() string          { return j.name }
func (j *ImportJob) Source() ImportSource   { return j.source }
func (j *ImportJob) EntityType() string    { return j.entityType }
func (j *ImportJob) Status() ImportStatus  { return j.status }
func (j *ImportJob) TotalRows() int        { return j.totalRows }
func (j *ImportJob) ProcessedRows() int    { return j.processedRows }
func (j *ImportJob) FailedRows() int       { return j.failedRows }
func (j *ImportJob) ErrorLog() []string    { return j.errorLog }
func (j *ImportJob) SourceURI() string     { return j.sourceURI }
func (j *ImportJob) StartedAt() *time.Time { return j.startedAt }
func (j *ImportJob) CompletedAt() *time.Time { return j.completedAt }
func (j *ImportJob) CreatedAt() time.Time  { return j.createdAt }

func (j *ImportJob) Begin() error {
	if j.status != StatusPending {
		return ErrAlreadyProcessing
	}
	now := time.Now().UTC()
	j.status = StatusProcessing
	j.startedAt = &now
	return nil
}

func (j *ImportJob) SetTotalRows(n int) {
	j.totalRows = n
}

func (j *ImportJob) IncrementProcessed() {
	j.processedRows++
}

func (j *ImportJob) RecordError(row int, msg string) {
	if len(j.errorLog) < 100 {
		j.errorLog = append(j.errorLog, msg)
	}
	j.failedRows++
}

func (j *ImportJob) Complete() {
	now := time.Now().UTC()
	j.status = StatusCompleted
	j.completedAt = &now
}

func (j *ImportJob) Fail(reason string) {
	now := time.Now().UTC()
	j.status = StatusFailed
	j.RecordError(0, reason)
	j.completedAt = &now
}
