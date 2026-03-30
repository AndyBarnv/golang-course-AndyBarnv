package http

import (
	"gateway/internal/usecase"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/status"
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
		log.Printf("Ошибка при вызове Collector: %v", err)

		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case 404:
				c.JSON(http.StatusNotFound, gin.H{"error": st.Message()})
				return
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				return
			}
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, repoInfo)
}
