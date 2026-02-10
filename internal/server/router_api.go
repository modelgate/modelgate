package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"

	apiv1 "github.com/modelgate/modelgate/internal/app/api/v1"
	"github.com/modelgate/modelgate/internal/server/middleware"
)

func registerRouters(container do.Injector, engine *gin.Engine) {
	// 健康检查
	engine.GET("health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	relayService := do.MustInvoke[*apiv1.RelayService](container)

	// v1
	rgv1 := engine.Group("/v1", middleware.RateLimit(container), middleware.CheckApiKey(container))

	rgv1.POST("/completions", relayService.Run)
	rgv1.POST("/chat/completions", relayService.Run)
	rgv1.POST("/embeddings", relayService.Run)
	rgv1.POST("/relay/:provider/*path", relayService.RunWithProvider)
}
