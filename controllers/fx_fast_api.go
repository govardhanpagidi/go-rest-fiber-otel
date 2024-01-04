package controllers

import (
	"github.com/PeerIslands/aci-fx-go/model/dto/request"
	"github.com/PeerIslands/aci-fx-go/model/dto/response"
	"github.com/PeerIslands/aci-fx-go/model/validation"
	"github.com/PeerIslands/aci-fx-go/service/common"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
	"os"

	"strconv"
)

var tracerName = "OTEL_SERVICE_NAME"

func FhAddRoutes(e *fiber.App) {
	// Convert currency
	e.Get("/api/convert", common.ParamValidationMiddlewareFiber[response.ForexDataResponse]([]validation.ValidationRule{
		{ParamName: "tenantId", Required: true, ParamType: "int"},
		{ParamName: "bankId", Required: true, ParamType: "int"},
		{ParamName: "baseCurrency", Required: true, ParamType: "string"},
		{ParamName: "targetCurrency", Required: true, ParamType: "string"},
	}), FhGetConvertedRate)

	// POST /api/forexrates
	e.Post("/api/forexrates", InsertForexRate)

	// POST /api/forexrates
	e.Post("/api/forexrates/batch", BulkInsertForexRate)

	// DELETE /api/forexrates?id=1
	e.Delete("/api/forexrates/:id", DeleteForexById)

	// GET /api/forexrate?id=1
	e.Get("/api/forexrates", common.ParamValidationMiddlewareFiber[response.ForexDataResponse]([]validation.ValidationRule{
		{ParamName: "id", Required: true, ParamType: "int"},
	}), FhGetForexRateById)

	// PUT /api/forexrate?tenantId=1&bankId=1&baseCurrency=USD&targetCurrency=INR&tier=1
	e.Put("/api/forexrates", common.ParamValidationMiddlewareFiber[response.ForexDataResponse]([]validation.ValidationRule{
		{ParamName: "tenantId", Required: true, ParamType: "int"},
		{ParamName: "bankId", Required: true, ParamType: "int"},
		{ParamName: "baseCurrency", Required: true, ParamType: "string"},
		{ParamName: "targetCurrency", Required: true, ParamType: "string"},
	}), UpdateForexRate)

	// not in use
	e.Put("/api/forexrate", common.ParamValidationMiddlewareFiber[response.ForexDataResponse]([]validation.ValidationRule{
		{ParamName: "id", Required: true, ParamType: "int"},
	}), UpdateForexById)
}

func InsertForexRate(c *fiber.Ctx) error {
	tracer := otel.Tracer(os.Getenv(tracerName))
	ctx, span := tracer.Start(c.Context(), c.Path())
	defer span.End()
	// convert body to forex_data_request
	var forexRateReq request.CreateForexDataRequest
	if err := c.BodyParser(&forexRateReq); err != nil {
		return err
	}
	err := c.Status(fiber.StatusOK).JSON(fxService.CreateForexData(&ctx, forexRateReq))
	if err != nil {
		return err
		return err
	}
	return nil
}

func BulkInsertForexRate(c *fiber.Ctx) error {
	tracer := otel.Tracer(os.Getenv(tracerName))
	ctx, span := tracer.Start(c.Context(), c.Path())
	defer span.End()
	// convert body to forex_data_request
	var forexRatesReq []request.CreateForexDataRequest
	if err := c.BodyParser(&forexRatesReq); err != nil {
		return err
	}

	err := c.Status(fiber.StatusOK).JSON(fxService.BulkInsertForexData(&ctx, forexRatesReq))
	if err != nil {
		return err
	}
	return nil
}

func FhGetConvertedRate(c *fiber.Ctx) error {
	tracer := otel.Tracer(os.Getenv(tracerName))
	ctx, span := tracer.Start(c.Context(), c.Path())
	defer span.End()
	tenantId, _ := strconv.Atoi(c.Query("tenantId"))
	bankId, _ := strconv.Atoi(c.Query("bankId"))
	baseCurrency := c.Query("baseCurrency")
	targetCurrency := c.Query("targetCurrency")
	amount, _ := strconv.ParseFloat(c.Query("amount"), 64)
	tier := c.Query("tier")

	err := c.Status(fiber.StatusOK).JSON(fxService.GetConvertedRate(&ctx, tenantId, bankId, amount, baseCurrency, targetCurrency, tier))
	if err != nil {
		return err
	}
	return nil
}

func FhGetForexRateById(c *fiber.Ctx) error {
	tracer := otel.Tracer(os.Getenv(tracerName))
	ctx, span := tracer.Start(c.Context(), c.Path())
	defer span.End()
	id, _ := strconv.Atoi(c.Query("id"))
	err := c.Status(fiber.StatusOK).JSON(fxService.GetConvertedRateById(&ctx, id))
	if err != nil {
		return err
	}
	return nil
}

func UpdateForexRate(c *fiber.Ctx) error {
	tracer := otel.Tracer(os.Getenv(tracerName))
	ctx, span := tracer.Start(c.Context(), c.Path())
	defer span.End()
	tenantId, _ := strconv.Atoi(c.Query("tenantId"))
	bankId, _ := strconv.Atoi(c.Query("bankId"))
	baseCurrency := c.Query("baseCurrency")
	targetCurrency := c.Query("targetCurrency")
	tier := c.Query("tier")

	err := c.Status(fiber.StatusOK).JSON(fxService.UpdateForexRate(&ctx, tenantId, bankId, baseCurrency, targetCurrency, tier))
	if err != nil {
		return err
	}
	return nil
}

func UpdateForexById(c *fiber.Ctx) error {
	tracer := otel.Tracer(os.Getenv(tracerName))
	ctx, span := tracer.Start(c.Context(), c.Path())
	defer span.End()
	id, _ := strconv.Atoi(c.Query("id"))

	err := c.Status(fiber.StatusOK).JSON(fxService.UpdateForexById(&ctx, id))
	if err != nil {
		return err
	}
	return nil
}

func DeleteForexById(c *fiber.Ctx) error {
	tracer := otel.Tracer(os.Getenv(tracerName))
	ctx, span := tracer.Start(c.Context(), c.Path())
	defer span.End()
	err := c.Status(fiber.StatusOK).JSON(fxService.DeleteForexRateById(&ctx, c.Params("id")))
	if err != nil {
		return err
	}
	return nil
}
