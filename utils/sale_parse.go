package utils

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ParseSaleUnitID(ctx *gin.Context, rejectPrefix string) (id int64, ok bool) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		slog.Error(rejectPrefix+"invalid sale unit id", slog.String("id_param", idStr))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid sale unit ID",
		})
		return 0, false
	}
	return id, true
}

func SaleFormFields(ctx *gin.Context, rejectPrefix string, saleID *int64) (orderID, posID, price, recipeID int32, ok bool) {
	posIDStr := ctx.Param("PosID")
	priceStr := ctx.Param("Price")
	recipeIDStr := ctx.Param("RecipeID")
	orderIDStr := ctx.Param("OrderID")
	if posIDStr == "" || priceStr == "" || recipeIDStr == "" || orderIDStr == "" {
		args := []any{
			slog.Bool("has_pos_id", posIDStr != ""),
			slog.Bool("has_price", priceStr != ""),
			slog.Bool("has_recipe_id", recipeIDStr != ""),
			slog.Bool("has_order_id", orderIDStr != ""),
		}
		if saleID != nil {
			args = append([]any{slog.Int64("sale.id", *saleID)}, args...)
		}
		slog.Info(rejectPrefix+": missing parameters", args...)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing required parameters: pos_id, price, recipe_id, order_id",
		})
		return 0, 0, 0, 0, false
	}

	pID, err := strconv.ParseInt(posIDStr, 10, 32)
	if err != nil {
		logInvalidFormField(ctx, rejectPrefix, "invalid pos_id", saleID, slog.String("pos_id", posIDStr), err)
		return 0, 0, 0, 0, false
	}

	p, err := strconv.ParseInt(priceStr, 10, 32)
	if err != nil {
		logInvalidFormField(ctx, rejectPrefix, "invalid price", saleID, slog.String("price", priceStr), err)
		return 0, 0, 0, 0, false
	}

	rID, err := strconv.ParseInt(recipeIDStr, 10, 32)
	if err != nil {
		logInvalidFormField(ctx, rejectPrefix, "invalid recipe_id", saleID, slog.String("recipe_id", recipeIDStr), err)
		return 0, 0, 0, 0, false
	}

	oID, err := strconv.ParseInt(orderIDStr, 10, 32)
	if err != nil {
		logInvalidFormField(ctx, rejectPrefix, "invalid order_id", saleID, slog.String("order_id", orderIDStr), err)
		return 0, 0, 0, 0, false
	}

	return int32(pID), int32(p), int32(rID), int32(oID), true
}

func logInvalidFormField(ctx *gin.Context, rejectPrefix, reason string, saleID *int64, fieldAttr slog.Attr, err error) {
	args := []any{slog.String("reason", reason), fieldAttr, slog.Any("error", err)}
	if saleID != nil {
		args = append([]any{slog.Int64("sale.id", *saleID)}, args...)
	}
	slog.Info(rejectPrefix+": invalid form field", args...)
	ctx.JSON(http.StatusBadRequest, gin.H{
		"error": "Invalid form field: " + reason,
	})
}
