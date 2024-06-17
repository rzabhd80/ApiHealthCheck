package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/rzabhd80/healthCheck/api/healthCheckApi/repository"
	"github.com/rzabhd80/healthCheck/models"
	"net/http"
	"strconv"
)

type APIHandler struct {
	apiRepo repository.APIRepository
}

func NewAPIHandler(apiRepo repository.APIRepository) *APIHandler {
	return &APIHandler{
		apiRepo: apiRepo,
	}
}

func (h *APIHandler) CreateAPI(c *gin.Context) {
	var api models.API
	if err := c.ShouldBindJSON(&api); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.apiRepo.Create(&api); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, api)

}

func (h *APIHandler) GetAPIs(c *gin.Context) {
	apis, err := h.apiRepo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, apis)
}

func (h *APIHandler) DeleteAPI(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.apiRepo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
