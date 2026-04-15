package routes

import (
	handlers "sale-service/handlers"
	"sale-service/observability"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type Route struct {
	db       *pgx.Conn
	handlers *handlers.Handlers
}

func NewRoute(db *pgx.Conn, prometheusMetrics *observability.PrometheusMetrics) *Route {
	return &Route{
		db:       db,
		handlers: handlers.NewHandlers(db, prometheusMetrics),
	}
}

func (r *Route) AddSaleRoutes(router *gin.Engine) {
	v1 := router.Group("/v1")
	{
		sale := v1.Group("/sale")
		{
			sale.GET("/:id", r.handlers.GetSaleUnit)
			sale.GET("/list", r.handlers.ListSaleUnit)
			sale.POST("/create", r.handlers.CreateSaleUnit)
			sale.PUT("/:id", r.handlers.UpdateSaleUnit)
			sale.DELETE("/:id", r.handlers.DeleteSaleUnit)
		}
	}
}

func (r *Route) AddHealthRoutes(router *gin.Engine) {
	// Health check endpoints (no authentication required)
	router.GET("/healthz", r.handlers.HealthzHandler)
	router.GET("/readyz", r.handlers.ReadyzHandler)
}
