package http

import (
	"gateway/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	usecase *usecase.Service
}

func NewHandler(uc *usecase.Service) *Handler {
	return &Handler{usecase: uc}
}

// GetRepoInfo godoc
// @Summary Get repository info
// @Description Get information about a GitHub repository
// @Tags repository
// @Accept  json
// @Produce  json
// @Param owner path string true "Repository Owner"
// @Param repo path string true "Repository Name"
// @Success 200 {object} domain.Repository
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /repo/{owner}/{repo} [get]
func (h *Handler) GetRepoInfo(c *gin.Context) {
	owner := c.Param("owner")
	repo := c.Param("repo")

	repoInfo, err := h.usecase.FetchRepository(c.Request.Context(), owner, repo)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, repoInfo)
}
