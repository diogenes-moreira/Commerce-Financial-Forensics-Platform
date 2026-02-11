package handler

import (
	"net/http"

	"github.com/diogenes/costforensics/backend/internal/adapter/http/dto/request"
	"github.com/diogenes/costforensics/backend/internal/adapter/http/dto/response"
	tenantrepo "github.com/diogenes/costforensics/backend/internal/adapter/postgres/tenant"
	importjobapp "github.com/diogenes/costforensics/backend/internal/application/importjob"
	"github.com/diogenes/costforensics/backend/internal/domain/importjob"
	"github.com/diogenes/costforensics/backend/internal/domain/shared"
	"github.com/diogenes/costforensics/backend/internal/platform/importer"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ImportHandler struct{}

func NewImportHandler() *ImportHandler {
	return &ImportHandler{}
}

func (h *ImportHandler) getService(c *gin.Context) *importjobapp.Service {
	db := c.MustGet("tenant_db").(*gorm.DB)
	repo := tenantrepo.NewImportJobRepo(db)
	csvParser := importer.NewCSVParser()
	return importjobapp.NewService(repo, csvParser)
}

func (h *ImportHandler) Upload(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	defer file.Close()

	entityType := c.PostForm("entity_type")
	name := c.PostForm("name")
	if name == "" {
		name = "CSV Import"
	}

	svc := h.getService(c)
	job, err := svc.CreateFromUpload(c.Request.Context(), name, entityType, file)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, response.ImportJobFromDomain(job))
}

func (h *ImportHandler) List(c *gin.Context) {
	var req request.ListImports
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filter := importjob.ListFilter{
		Status:     importjob.ImportStatus(req.Status),
		EntityType: req.EntityType,
		Pagination: shared.NewPagination(req.Page, req.PageSize),
	}

	svc := h.getService(c)
	result, err := svc.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	items := make([]response.ImportJob, len(result.Items))
	for i, j := range result.Items {
		items[i] = response.ImportJobFromDomain(j)
	}
	c.JSON(http.StatusOK, response.Paginated[response.ImportJob]{
		Items: items, TotalCount: result.TotalCount,
		Page: result.Page, PageSize: result.PageSize,
	})
}

func (h *ImportHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	svc := h.getService(c)
	job, err := svc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.ImportJobFromDomain(job))
}
