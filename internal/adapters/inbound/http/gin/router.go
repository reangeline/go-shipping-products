package ginadapter

import (
	"net/http"

	"github.com/gin-gonic/gin"
	ctr "github.com/reangeline/go-shipping-products/internal/adapters/inbound/http/order"
	"github.com/reangeline/go-shipping-products/internal/adapters/inbound/http/presenter"
)

func BuildHandler(ctrl *ctr.Controller) http.Handler {
	r := gin.New()
	gin.SetMode(gin.ReleaseMode)

	r.Use(gin.Recovery())
	r.Use(LoggerMiddleware())

	// CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	v1 := r.Group("/v1")
	{
		v1.GET("/packsizes", func(c *gin.Context) {
			res, err := ctrl.HandleGetPackSizes(c.Request.Context())
			if err != nil {
				status, body := presenter.MapError(err)
				c.JSON(status, body)
				return
			}
			c.JSON(http.StatusOK, res)
		})

		v1.POST("/calculate", func(c *gin.Context) {
			var req ctr.CalculateRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, presenter.ErrorBody{
					Code: "invalid_request", Message: "invalid JSON payload",
				})
				return
			}
			res, err := ctrl.HandleCalculate(c.Request.Context(), req)
			if err != nil {
				status, body := presenter.MapError(err)
				c.JSON(status, body)
				return
			}
			c.JSON(http.StatusOK, res)
		})
	}

	r.GET("/healthz", func(c *gin.Context) { c.String(http.StatusOK, "ok") })
	r.GET("/readyz", func(c *gin.Context) { c.String(http.StatusOK, "ready") })

	RegisterDocs(r, "docs/api/v1/openapi.yaml", "/docs")

	return r
}
