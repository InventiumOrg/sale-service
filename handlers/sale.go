package handlers

import (
	"log/slog"
	"net/http"
	models "sale-service/models/sqlc"
	"sale-service/observability"
	"sale-service/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Handlers struct {
	queries           *models.Queries
	tracer            trace.Tracer
	db                *pgx.Conn
	prometheusMetrics *observability.PrometheusMetrics
}

func NewHandlers(db *pgx.Conn, prometheusMetrics *observability.PrometheusMetrics) *Handlers {
	return &Handlers{
		db:                db,
		queries:           models.New(db),
		tracer:            otel.Tracer("sale-service/handlers"),
		prometheusMetrics: prometheusMetrics,
	}
}

func (h *Handlers) GetSaleUnit(ctx *gin.Context) {
	// Start a new span for this operation
	_, span := h.tracer.Start(ctx.Request.Context(), "GetSaleUnit")
	defer span.End()

	// Get sale unit ID from URL parameter
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid sale unit ID",
		})
		return
	}

	// Add attributes to the span
	span.SetAttributes(attribute.Int64("saleUnit.id", id))

	dbStart := time.Now()
	saleUnit, err := h.queries.GetSaleUnit(ctx, int64(id))
	dbDuration := time.Since(dbStart)

	// Record database operation duration (Prometheus)
	if h.prometheusMetrics != nil {
		h.prometheusMetrics.RecordDBOperation("get", "sale_units", dbDuration, err)
	}

	if err != nil {
		slog.Error("Got an error while getting sale unit: ", slog.Any("err", err.Error()))
		span.RecordError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get sale unit",
		})
		return
	}

	// Record successful retrieval (Prometheus)
	if h.prometheusMetrics != nil {
		h.prometheusMetrics.RecordSaleUnitOperation("get", saleUnit.ID)
	}

	// Record successful operation
	span.SetAttributes(
		attribute.String("operation.status", "success"),
	)

	ctx.JSON(200, gin.H{
		"message": "Get Sale Unit Successfully",
		"data":    saleUnit,
	})
}

func (h *Handlers) ListSaleUnit(ctx *gin.Context) {
	// Start a new span for this operation
	_, span := h.tracer.Start(ctx.Request.Context(), "ListSaleUnits")
	defer span.End()

	// Add attributes to the span
	span.SetAttributes(
		attribute.Int("saleUnit.limit", 10),
		attribute.Int("saleUnit.offset", 0),
	)

	dbStart := time.Now()
	saleUnits, err := h.queries.ListSaleUnit(ctx, models.ListSaleUnitParams{
		Limit:  10,
		Offset: 0,
	})
	dbDuration := time.Since(dbStart)

	// Record database operation duration (Prometheus)
	if h.prometheusMetrics != nil {
		h.prometheusMetrics.RecordDBOperation("list", "sale_units", dbDuration, err)
	}

	if err != nil {
		slog.Error("Got an error while listing sale units: ", slog.Any("err", err.Error()))
		span.RecordError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list sale units",
		})
		return
	}

	// Record successful list operation (Prometheus)
	if h.prometheusMetrics != nil {
		h.prometheusMetrics.RecordSaleUnitOperation("list", 0)
	}

	span.SetAttributes(
		attribute.Int("saleUnit.count", len(saleUnits)),
		attribute.String("operation.status", "success"),
	)

	ctx.JSON(200, gin.H{
		"message": "List Sale Units Successfully",
		"data":    saleUnits,
		"count":   len(saleUnits),
	})
}

func (h *Handlers) CreateSaleUnit(ctx *gin.Context) {
	// Start a new span for this operation
	_, span := h.tracer.Start(ctx.Request.Context(), "CreateSaleUnit")
	defer span.End()

	id, ok := utils.ParseSaleUnitID(ctx, "create sale unit rejected")
	if !ok {
		return
	}

	posID, price, recipeID, orderID, ok := utils.SaleFormFields(ctx, "create sale unit rejected", &id)
	if !ok {
		return
	}

	params := models.CreateSaleUnitParams{
		PosID:    posID,
		Price:    price,
		RecipeID: recipeID,
		OrderID:  orderID,
	}

	// Add attributes to the span
	span.SetAttributes(
		attribute.Int64("saleUnit.pos_id", int64(posID)),
		attribute.Int64("saleUnit.price", int64(price)),
		attribute.Int64("saleUnit.recipe_id", int64(recipeID)),
		attribute.Int64("saleUnit.order_id", int64(orderID)),
	)

	dbStart := time.Now()
	saleUnit, err := h.queries.CreateSaleUnit(ctx, params)
	dbDuration := time.Since(dbStart)

	// Record database operation duration (Prometheus)
	if h.prometheusMetrics != nil {
		h.prometheusMetrics.RecordDBOperation("create", "sale_units", dbDuration, err)
	}

	if err != nil {
		slog.Error("Could not create sale unit: ", slog.Any("err", err.Error()))
		span.RecordError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create sale unit",
		})
		return
	}

	// Record successful creation (Prometheus)
	if h.prometheusMetrics != nil {
		h.prometheusMetrics.RecordSaleUnitOperation("create", saleUnit.ID)
		h.prometheusMetrics.UpdateSaleUnitsCount(1)
	}

	// Record successful operation
	span.SetAttributes(
		attribute.Int64("saleUnit.id", int64(saleUnit.ID)),
		attribute.String("operation.status", "success"),
	)

	ctx.JSON(200, gin.H{
		"message": "Create Sale Unit Successfully",
		"data":    saleUnit,
	})
}

func (h *Handlers) UpdateSaleUnit(ctx *gin.Context) {
	// Start a new span for this operation
	_, span := h.tracer.Start(ctx.Request.Context(), "UpdateSaleUnit")
	defer span.End()

	id, ok := utils.ParseSaleUnitID(ctx, "update sale unit: ")
	if !ok {
		return
	}

	posID, price, recipeID, orderID, ok := utils.SaleFormFields(ctx, "update sale unit: ", &id)
	if !ok {
		return
	}

	params := models.UpdateSaleUnitParams{
		ID:       id,
		PosID:    posID,
		Price:    price,
		RecipeID: recipeID,
		OrderID:  orderID,
	}

	span.SetAttributes(
		attribute.Int64("saleUnit.pos_id", int64(posID)),
		attribute.Int64("saleUnit.price", int64(price)),
		attribute.Int64("saleUnit.recipe_id", int64(recipeID)),
		attribute.Int64("saleUnit.order_id", int64(orderID)),
	)

	dbStart := time.Now()
	saleUnit, err := h.queries.UpdateSaleUnit(ctx, params)
	dbDuration := time.Since(dbStart)

	// Record database operation duration (Prometheus)
	if h.prometheusMetrics != nil {
		h.prometheusMetrics.RecordDBOperation("update", "sale_units", dbDuration, err)
	}

	if err != nil {
		slog.Error("failed to update sale unit", slog.Int64("sale.id", id), slog.Any("err", err))
		span.RecordError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update sale unit",
		})
		return
	}

	// Record successful update (Prometheus)
	if h.prometheusMetrics != nil {
		h.prometheusMetrics.RecordSaleUnitOperation("update", saleUnit.ID)
	}

	span.SetAttributes(attribute.String("operation.status", "success"))

	slog.Info("order updated", slog.Int64("sale.id", id), slog.Int64("sale.pos_id", int64(saleUnit.PosID)))
	ctx.JSON(200, gin.H{
		"message": "Update Sale Unit Successfully",
		"data":    saleUnit,
	})
}

func (h *Handlers) DeleteSaleUnit(ctx *gin.Context) {
	// Start a new span for this operation
	_, span := h.tracer.Start(ctx.Request.Context(), "DeleteSaleUnit")
	defer span.End()

	id, ok := utils.ParseSaleUnitID(ctx, "delete sale unit: ")
	if !ok {
		return
	}

	span.SetAttributes(attribute.Int64("saleUnit.id", id))

	dbStart := time.Now()
	err := h.queries.DeleteSaleUnit(ctx, id)
	dbDuration := time.Since(dbStart)

	// Record database operation duration (Prometheus)
	if h.prometheusMetrics != nil {
		h.prometheusMetrics.RecordDBOperation("delete", "sale_units", dbDuration, err)
	}

	if err != nil {
		slog.Error("failed to delete sale unit", slog.Int64("sale.id", id), slog.Any("err", err))
		span.RecordError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete sale unit",
		})
		return
	}

	// Record successful deletion (Prometheus)
	if h.prometheusMetrics != nil {
		h.prometheusMetrics.RecordSaleUnitOperation("delete", id)
	}

	span.SetAttributes(attribute.String("operation.status", "success"))

	slog.Info("sale unit deleted", slog.Int64("sale.id", id))
	ctx.JSON(200, gin.H{
		"message": "Delete Sale Unit Successfully",
	})
}
