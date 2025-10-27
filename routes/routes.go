package routes

import (
	handlers "sale-service/handlers"
	"sale-service/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type Route struct {
	db       *pgx.Conn
	handlers *handlers.Handlers
}

func NewRoute(db *pgx.Conn) *Route {
	return &Route{
		db:       db,
		handlers: handlers.NewHandlers(db),
	}
}

func (r *Route) AddSaleRoutes(router *gin.Engine) {
	v1 := router.Group("/v1")
	{
		inventory := v1.Group("/sale")
		inventory.Use(middlewares.ClerkAuth(r.db))
		{
			inventory.GET("/:id", r.handlers.CreateSaleUnit)
			inventory.GET("/list", r.handlers.CreateSaleUnit)
			inventory.POST("/create", r.handlers.CreateSaleUnit)
			inventory.PUT("/:id", r.handlers.CreateSaleUnit)
			inventory.DELETE("/:id", r.handlers.CreateSaleUnit)
		}
	}
}

func (r *Route) AddHealthRoutes(router *gin.Engine) {
	// Health check endpoints (no authentication required)
	router.GET("/healthz", r.handlers.HealthzHandler)
	router.GET("/readyz", r.handlers.ReadyzHandler)
}
