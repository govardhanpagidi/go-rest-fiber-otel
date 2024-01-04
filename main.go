package main

import (
	"context"
	"os"
	"sync"

	"github.com/PeerIslands/aci-fx-go/controllers"
	"github.com/PeerIslands/aci-fx-go/model/dto/response"
	"github.com/PeerIslands/aci-fx-go/service/common"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	//"github.com/gofiber/fiber/v2"
	//"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
)

var (
	resource          *sdkresource.Resource
	initResourcesOnce sync.Once
)

func otelMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tracer := otel.Tracer(os.Getenv("OTEL_SERVICE_NAME"))
		ctx, span := tracer.Start(c.Context(), c.Path())
		defer span.End()

		c.SetUserContext(ctx)

		return c.Next()
	}
}

func main() {
	common.InitLog()

	initTracerProvider()

	/*router := gin.Default()
	router.Use(GlobalErrorHandler)
	controllers.AddRoutes(router)
	log.Fatal(router.Run("0.0.0.0:8080"))*/

	fiberApp := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusOK
			err = ctx.Status(code).JSON(common.GetSimpleResponse[response.ForexDataResponse](nil, response.InternalError, &[]response.Error{
				{Code: "INTERNAL_ERROR", Message: "Please try again later", Details: "Something went wrong. Please try again later"},
			}))
			if err != nil {
				return err
			}
			return nil
		},
	})

	fiberApp.Use(logger.New())
	fiberApp.Use(otelMiddleware())

	controllers.FhAddRoutes(fiberApp)
	err := fiberApp.Listen("0.0.0.0:8080")
	if err != nil {
		log.Fatal("Error:", err)
	}
	//stream.Connect()

}

func initTracerProvider() *sdktrace.TracerProvider {
	ctx := context.Background()

	exporter, err := otlptracegrpc.New(ctx)
	if err != nil {
		log.Fatalf("OTLP Trace gRPC Creation: %v", err)
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(initResource()),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp
}

func initResource() *sdkresource.Resource {
	initResourcesOnce.Do(func() {
		extraResources, _ := sdkresource.New(
			context.Background(),
			sdkresource.WithOS(),
			sdkresource.WithProcess(),
			sdkresource.WithContainer(),
			sdkresource.WithHost(),
		)
		resource, _ = sdkresource.Merge(
			sdkresource.Default(),
			extraResources,
		)
	})
	return resource
}

/*func GlobalErrorHandler(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			// Handle the error here, you can log it or send an error response
			c.IndentedJSON(http.StatusOK, common.GetSimpleResponse[response.ForexDataResponse](nil, response.InternalError, &[]response.Error{
				{Code: "INTERNAL_ERROR", Message: "Please try again later", Details: "Something went wrong. Please try again later"},
			}))
			common.Logger.Errorf("Unhandled Error in %s %s. Exception:%v", c.Request.Method, c.Request.URL, err)
		}
	}()
	c.Next()
}*/
