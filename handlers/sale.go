package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	models "sale-service/models/sqlc"
	"sale-service/observability"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type Handlers struct {
	queries         *models.Queries
	tracer          trace.Tracer
	db              *pgx.Conn
	businessMetrics *observability.BusinessMetrics
}

func NewHandlers(db *pgx.Conn, businessMetrics *observability.BusinessMetrics) *Handlers {
	return &Handlers{
		db:              db,
		queries:         models.New(db),
		tracer:          otel.Tracer("sale-service/handlers"),
		businessMetrics: businessMetrics,
	}
}

func (h *Handlers) GetSaleUnit(ctx *gin.Context) {
	// Start a new span for this operation
	spanCtx, span := h.tracer.Start(ctx.Request.Context(), "Get Sale Unit")
	defer span.End()

	// Record authentication attempt
	if h.businessMetrics != nil {
		h.businessMetrics.AuthenticationAttempts.Add(spanCtx, 1,
			metric.WithAttributes(attribute.String("operation", "get_sale_unit")))
	}

	_, existed := ctx.Get("claims")
	if !existed {
		span.RecordError(fmt.Errorf("claims not found in context"))
		span.SetAttributes(attribute.String("error", "claims_not_found"))

		// Record authentication failure
		if h.businessMetrics != nil {
			h.businessMetrics.AuthenticationAttempts.Add(spanCtx, 1,
				metric.WithAttributes(
					attribute.String("operation", "get_sale_unit"),
					attribute.String("status", "failed"),
				))
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Claims not found in context",
		})
		return
	}

	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Sale Unit ID format",
		})
		return
	}

	// Add attributes to the span
	span.SetAttributes(attribute.Int64("saleUnit.id", id))

	// Record business operation
	if h.businessMetrics != nil {
		h.businessMetrics.SaleUnitOperations.Add(spanCtx, 1,
			metric.WithAttributes(
				attribute.String("operation", "get"),
				attribute.Int64("sale_unit_id", id),
			))
	}

	// Measure database operation duration
	dbStart := time.Now()
	saleUnit, err := h.queries.GetSaleUnit(spanCtx, int32(id))
	dbDuration := time.Since(dbStart).Seconds()

	// Record database operation metrics
	if h.businessMetrics != nil {
		h.businessMetrics.DBOperationDuration.Record(spanCtx, dbDuration,
			metric.WithAttributes(
				attribute.String("operation", "get_sale_unit"),
				attribute.String("table", "sale_units"),
			))
	}

	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String("error", "database_query_failed"))

		// Record database error
		if h.businessMetrics != nil {
			h.businessMetrics.DBOperationErrors.Add(spanCtx, 1,
				metric.WithAttributes(
					attribute.String("operation", "get_sale_unit"),
					attribute.String("error_type", "query_failed"),
				))
		}

		slog.Error("Got an error while getting sale units: ", slog.Any(err.Error(), "err"))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get sale unit",
		})
		return
	}

	// Record successful retrieval
	if h.businessMetrics != nil {
		h.businessMetrics.SaleUnitRetrievals.Add(spanCtx, 1,
			metric.WithAttributes(
				attribute.String("sale_unit_name", saleUnit.Name),
				attribute.String("status", "success"),
			))
	}

	// Record successful operation
	span.SetAttributes(
		attribute.String("saleUnit.name", saleUnit.Name),
		attribute.String("operation.status", "success"),
	)

	ctx.JSON(200, gin.H{
		"message": "Get Sale Unit Successfully",
		"data":    saleUnit,
	})
}

func (h *Handlers) ListSaleUnit(ctx *gin.Context) {
	// Start a new span for this operation
	spanCtx, span := h.tracer.Start(ctx.Request.Context(), "List Sale Units")
	defer span.End()

	// Record authentication attempt
	if h.businessMetrics != nil {
		h.businessMetrics.AuthenticationAttempts.Add(spanCtx, 1,
			metric.WithAttributes(attribute.String("operation", "list_sale_units")))
	}

	_, existed := ctx.Get("claims")
	if !existed {
		span.RecordError(fmt.Errorf("claims not found in context"))
		span.SetAttributes(attribute.String("error", "claims_not_found"))

		// Record authentication failure
		if h.businessMetrics != nil {
			h.businessMetrics.AuthenticationAttempts.Add(spanCtx, 1,
				metric.WithAttributes(
					attribute.String("operation", "list_sale_units"),
					attribute.String("status", "failed"),
				))
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Claims not found in context",
		})
		return
	}

	// Add attributes to the span
	span.SetAttributes(
		attribute.Int("saleUnit.limit", 10),
		attribute.Int("saleUnit.offset", 0),
	)

	// Record list request
	if h.businessMetrics != nil {
		h.businessMetrics.SaleUnitListRequests.Add(spanCtx, 1,
			metric.WithAttributes(
				attribute.Int("limit", 10),
				attribute.Int("offset", 0),
			))
	}

	// Measure database operation duration
	dbStart := time.Now()
	saleUnits, err := h.queries.ListSaleUnit(spanCtx, models.ListSaleUnitParams{
		Limit:  10,
		Offset: 0,
	})
	dbDuration := time.Since(dbStart).Seconds()

	// Record database operation metrics
	if h.businessMetrics != nil {
		h.businessMetrics.DBOperationDuration.Record(spanCtx, dbDuration,
			metric.WithAttributes(
				attribute.String("operation", "list_sale_units"),
				attribute.String("table", "sale_units"),
			))
	}

	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String("error", "database_query_failed"))

		// Record database error
		if h.businessMetrics != nil {
			h.businessMetrics.DBOperationErrors.Add(spanCtx, 1,
				metric.WithAttributes(
					attribute.String("operation", "list_sale_units"),
					attribute.String("error_type", "query_failed"),
				))
		}

		slog.Error("Got an error while listing sale units: ", slog.Any("err", err.Error()))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list sale units",
		})
		return
	}

	// Record successful operation
	span.SetAttributes(
		attribute.Int("saleUnit.count", len(saleUnits)),
		attribute.String("operation.status", "success"),
	)

	// Record business metrics for successful list operation
	if h.businessMetrics != nil {
		h.businessMetrics.SaleUnitOperations.Add(spanCtx, 1,
			metric.WithAttributes(
				attribute.String("operation", "list"),
				attribute.Int("result_count", len(saleUnits)),
			))
	}

	ctx.JSON(200, gin.H{
		"message": "List Sale Units Successfully",
		"data":    saleUnits,
	})
}

// func (h *Handlers) UpdateStorageRoom(ctx *gin.Context) {
// 	_, existed := ctx.Get("claims")
// 	if !existed {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Claims not found in context",
// 		})
// 		return
// 	}

// 	// Get inventory ID from URL parameter
// 	idStr := ctx.Param("id")
// 	id, err := strconv.ParseInt(idStr, 10, 64)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Invalid storage room ID",
// 		})
// 		return
// 	}

// 	// Start database transaction
// 	tx, err := h.db.Begin(ctx)
// 	if err != nil {
// 		slog.Error("Failed to start transaction", slog.Any("err", err.Error()))
// 		ctx.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Failed to start transaction",
// 		})
// 		return
// 	}
// 	defer tx.Rollback(ctx) // This will be ignored if tx.Commit() succeeds

// 	// Create queries with transaction
// 	qtx := h.queries.WithTx(tx)

// 	// Check if storage room exists before updating
// 	_, err = qtx.GetStorageRoom(ctx, int32(id))
// 	if err != nil {
// 		slog.Error("Storage room not found", slog.Any("err", err.Error()))
// 		ctx.JSON(http.StatusNotFound, gin.H{
// 			"error": "Storage room not found",
// 		})
// 		return
// 	}

// 	// Parse WarehouseID from string to int32
// 	warehouseIDStr := ctx.PostForm("WarehouseId")
// 	warehouseID, err := strconv.ParseInt(warehouseIDStr, 10, 32)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Invalid warehouse ID format",
// 		})
// 		return
// 	}

// 	// Update storage room within transaction
// 	param := models.UpdateStorageRoomParams{
// 		ID:          int32(id),
// 		Name:        ctx.PostForm("Name"),
// 		Number:      ctx.PostForm("Number"),
// 		WarehouseID: int32(warehouseID),
// 	}

// 	storageRoom, err := qtx.UpdateStorageRoom(ctx, param)
// 	if err != nil {
// 		slog.Error("Could not update storage room", slog.Any("err", err.Error()))
// 		ctx.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Failed to update storage room",
// 		})
// 		return
// 	}

// 	// Commit transaction
// 	if err := tx.Commit(ctx); err != nil {
// 		slog.Error("Failed to commit transaction", slog.Any("err", err.Error()))
// 		ctx.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Failed to commit transaction",
// 		})
// 		return
// 	}

// 	ctx.JSON(200, gin.H{
// 		"message": "Update Storage Room Successfully",
// 		"data":    storageRoom,
// 	})
// }

func (h *Handlers) CreateSaleUnit(ctx *gin.Context) {
	// Start a new span for this operation
	spanCtx, span := h.tracer.Start(ctx.Request.Context(), "Create Sale Unit")
	defer span.End()

	// Record authentication attempt
	if h.businessMetrics != nil {
		h.businessMetrics.AuthenticationAttempts.Add(spanCtx, 1,
			metric.WithAttributes(attribute.String("operation", "create_sale_unit")))
	}

	_, existed := ctx.Get("claims")
	if !existed {
		span.RecordError(fmt.Errorf("claims not found in context"))
		span.SetAttributes(attribute.String("error", "claims_not_found"))

		// Record authentication failure
		if h.businessMetrics != nil {
			h.businessMetrics.AuthenticationAttempts.Add(spanCtx, 1,
				metric.WithAttributes(
					attribute.String("operation", "create_sale_unit"),
					attribute.String("status", "failed"),
				))
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Claims not found in context",
		})
		return
	}

	unitName := ctx.PostForm("Name")
	param := models.CreateSaleUnitParams{
		Name: unitName,
	}

	// Add attributes to the span
	span.SetAttributes(attribute.String("saleUnit.name", unitName))

	// Record business operation
	if h.businessMetrics != nil {
		h.businessMetrics.SaleUnitOperations.Add(spanCtx, 1,
			metric.WithAttributes(
				attribute.String("operation", "create"),
				attribute.String("sale_unit_name", unitName),
			))
	}

	// Measure database operation duration
	dbStart := time.Now()
	saleUnit, err := h.queries.CreateSaleUnit(spanCtx, param)
	dbDuration := time.Since(dbStart).Seconds()

	// Record database operation metrics
	if h.businessMetrics != nil {
		h.businessMetrics.DBOperationDuration.Record(spanCtx, dbDuration,
			metric.WithAttributes(
				attribute.String("operation", "create_sale_unit"),
				attribute.String("table", "sale_units"),
			))
	}

	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String("error", "database_insert_failed"))

		// Record database error
		if h.businessMetrics != nil {
			h.businessMetrics.DBOperationErrors.Add(spanCtx, 1,
				metric.WithAttributes(
					attribute.String("operation", "create_sale_unit"),
					attribute.String("error_type", "insert_failed"),
				))
		}

		slog.Error("Could not create sale unit: ", slog.Any("err", err.Error()))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create sale unit",
		})
		return
	}

	// Record successful creation
	if h.businessMetrics != nil {
		h.businessMetrics.SaleUnitCreated.Add(spanCtx, 1,
			metric.WithAttributes(
				attribute.String("sale_unit_name", saleUnit.Name),
				attribute.Int64("sale_unit_id", int64(saleUnit.ID)),
			))

		// Increment active sale units counter
		h.businessMetrics.ActiveSaleUnits.Add(spanCtx, 1,
			metric.WithAttributes(attribute.String("operation", "created")))
	}

	// Record successful operation
	span.SetAttributes(
		attribute.String("saleUnit.created_name", saleUnit.Name),
		attribute.Int64("saleUnit.created_id", int64(saleUnit.ID)),
		attribute.String("operation.status", "success"),
	)

	ctx.JSON(200, gin.H{
		"message": "Create Sale Unit Successfully",
		"data":    saleUnit,
	})
}

// func (h *Handlers) DeleteStorageRoom(ctx *gin.Context) {
// 	_, existed := ctx.Get("claims")
// 	if !existed {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Claims not found in context",
// 		})
// 		return
// 	}

// 	idStr := ctx.Param("id")
// 	id, err := strconv.ParseInt(idStr, 10, 32)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Invalid storage room ID format",
// 		})
// 		return
// 	}

// 	err = h.queries.DeleteStorageRoom(ctx, int32(id))
// 	if err != nil {
// 		slog.Error("Failed to delete storage room: ", slog.Any("err", err.Error()))
// 		ctx.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Failed to delete storage room",
// 		})
// 		return
// 	} else {
// 		ctx.JSON(200, gin.H{"message": "Delete Storage Room Successfully"})
// 	}

// }
