package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"repo-stat/api/internal/adapter/processor"
	"repo-stat/api/internal/adapter/subscriber"
	"repo-stat/api/internal/dto"
	"sync"
)

// @Summary      Check services health
// @Description  Pings Processor and Subscriber services
// @Tags         health
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.PingResponse
// @Failure      503  {object}  dto.PingResponse
// @Router       /api/ping [get]
func NewPingHandler(log *slog.Logger, subClient *subscriber.Client, procClient *processor.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var wg sync.WaitGroup
		var procStatus, subStatus string = "up", "up"

		wg.Add(2)

		go func() {
			defer wg.Done()
			if err := procClient.Ping(r.Context()); err != nil {
				procStatus = "down"
			}
		}()

		go func() {
			defer wg.Done()
			if subClient.Ping(r.Context()) != "up" {
				subStatus = "down"
			}
		}()

		wg.Wait()

		code := http.StatusOK
		statusStr := "ok"
		if procStatus == "down" || subStatus == "down" {
			code = http.StatusServiceUnavailable
			statusStr = "degraded"
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(dto.PingResponse{
			Status: statusStr,
			Services: []dto.PingService{
				{Name: "processor", Status: procStatus},
				{Name: "subscriber", Status: subStatus},
			},
		})
	}
}
