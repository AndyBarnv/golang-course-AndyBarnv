package main

import (
	_ "gateway/docs"
	"gateway/internal/adapter/grpc"
	"gateway/internal/delivery/http"
	"gateway/internal/usecase"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title GitHub Repository API
// @version 1.0
// @description A microservice API to fetch GitHub repository info.
// @host localhost:8080
// @BasePath /
func main() {
	collectorAddr := os.Getenv("COLLECTOR_ADDR")
	if collectorAddr == "" {
		collectorAddr = "localhost:50051"
	}

	grpcAdapter, err := grpc.NewAdapter(collectorAddr)
	if err != nil {
		log.Fatalf("Failed to connect to Collector: %v", err)
	}
	defer grpcAdapter.Close()

	service := usecase.NewService(grpcAdapter)
	handler := http.NewHandler(service)

	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/repo/:owner/:repo", handler.GetRepoInfo)

	log.Println("API Gateway listening on port 8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
