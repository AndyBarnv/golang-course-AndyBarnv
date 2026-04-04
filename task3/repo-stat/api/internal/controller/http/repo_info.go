package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"repo-stat/api/internal/adapter/processor"
	"repo-stat/api/internal/domain"
	"repo-stat/api/internal/dto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// @Summary      Get repository info
// @Description  Returns basic information about a GitHub repository by URL
// @Tags         repository
// @Accept       json
// @Produce      json
// @Param        url   query      string  true  "GitHub repository URL"
// @Success      200  {object}  dto.RepositoryInfoResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /api/repositories/info [get]
func NewRepoInfoHandler(log *slog.Logger, procClient *processor.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rawURL := r.URL.Query().Get("url")

		if rawURL == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(dto.ErrorResponse{Error: "url query parameter is required"})
			return
		}

		owner, repo, err := domain.ParseGitHubURL(rawURL)

		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
			return
		}

		resp, err := procClient.GetRepoInfo(r.Context(), owner, repo)
		if err != nil {
			st, ok := status.FromError(err)
			if ok && st.Code() == codes.NotFound {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(dto.ErrorResponse{Error: st.Message()})
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(dto.ErrorResponse{Error: "internal server error"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(dto.RepositoryInfoResponse{
			FullName:    resp.FullName,
			Description: resp.Description,
			Stars:       resp.Stars,
			Forks:       resp.Forks,
			CreatedAt:   resp.CreatedAt,
		})
	}
}
