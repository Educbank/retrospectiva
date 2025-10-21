package handlers

import (
	"net/http"

	"educ-retro/internal/services"

	"github.com/gin-gonic/gin"
)

type TemplateHandler struct {
	templateService *services.TemplateService
}

func NewTemplateHandler(templateService *services.TemplateService) *TemplateHandler {
	return &TemplateHandler{templateService: templateService}
}

// GetTemplates godoc
// @Summary Get available retrospective templates
// @Description Get all available templates for retrospectives
// @Tags Templates
// @Accept json
// @Produce json
// @Success 200 {array} services.TemplateDefinition "Available templates"
// @Router /templates [get]
func (h *TemplateHandler) GetTemplates(c *gin.Context) {
	templates := h.templateService.GetAvailableTemplates()
	c.JSON(http.StatusOK, templates)
}

// GetTemplate godoc
// @Summary Get specific template
// @Description Get detailed information about a specific template
// @Tags templates
// @Accept json
// @Produce json
// @Param id path string true "Template ID"
// @Success 200 {object} services.TemplateDefinition "Template details"
// @Failure 400 {object} map[string]string "Invalid template ID"
// @Failure 404 {object} map[string]string "Template not found"
// @Router /api/v1/templates/{id} [get]
func (h *TemplateHandler) GetTemplate(c *gin.Context) {
	templateID := c.Param("id")

	template, err := h.templateService.GetTemplate(templateID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, template)
}

// GetTemplateCategories godoc
// @Summary Get template categories
// @Description Get categories for a specific template
// @Tags templates
// @Accept json
// @Produce json
// @Param id path string true "Template ID"
// @Success 200 {array} services.TemplateCategory "Template categories"
// @Failure 400 {object} map[string]string "Invalid template ID"
// @Failure 404 {object} map[string]string "Template not found"
// @Router /api/v1/templates/{id}/categories [get]
func (h *TemplateHandler) GetTemplateCategories(c *gin.Context) {
	templateID := c.Param("id")

	categories, err := h.templateService.GetTemplateCategories(templateID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

func (h *TemplateHandler) SetupRoutes(r *gin.RouterGroup) {
	templates := r.Group("/templates")
	{
		templates.GET("", h.GetTemplates)
		templates.GET("/:id", h.GetTemplate)
		templates.GET("/:id/categories", h.GetTemplateCategories)
	}
}
