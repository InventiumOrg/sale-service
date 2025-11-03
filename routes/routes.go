package routes

import (
	handlers "sale-service/handlers"
	"sale-service/middlewares"
	"sale-service/observability"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type Route struct {
	db       *pgx.Conn
	handlers *handlers.Handlers
}

func NewRoute(db *pgx.Conn, businessMetrics *observability.BusinessMetrics) *Route {
	return &Route{
		db:       db,
		handlers: handlers.NewHandlers(db, businessMetrics),
	}
}

func (r *Route) AddSaleRoutes(router *gin.Engine) {
	v1 := router.Group("/v1")
	{
		sale := v1.Group("/sale")
		sale.Use(middlewares.ClerkAuth(r.db))
		{
			sale.GET("/:id", r.handlers.GetSaleUnit)
			sale.GET("/list", r.handlers.ListSaleUnit)
			sale.POST("/create", r.handlers.CreateSaleUnit)
			// TODO: Add update and delete handlers when implemented
			// sale.PUT("/:id", r.handlers.UpdateSaleUnit)
			// sale.DELETE("/:id", r.handlers.DeleteSaleUnit)
		}
	}
}

func (r *Route) AddHealthRoutes(router *gin.Engine) {
	// Health check endpoints (no authentication required)
	router.GET("/healthz", r.handlers.HealthzHandler)
	router.GET("/readyz", r.handlers.ReadyzHandler)
}
