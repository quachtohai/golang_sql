package otel

import (
	"context"
	"fmt"
	"golang_sql/pkg/logger"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

type JaegerConfig struct {
	Server      string `mapstructure:"server"`
	ServiceName string `mapstructure:"serviceName"`
	TracerName  string `mapstructure:"tracerName"`
}

func TracerProvider(ctx context.Context, cfg *JaegerConfig, log logger.ILogger) (trace.Tracer, error) {
	var serverUrl = fmt.Sprintf(cfg.Server+"%s", "/api/traces")

	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(serverUrl)))

	if err != nil {
		return nil, err
	}

	env := os.Getenv("APP_ENV")

	if env != "production" {
		env = "development"
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.ServiceName),
			attribute.String("environment", env),
		)),
	)

	go func() {
		for {
			select {
			case <-ctx.Done():
				err = tp.Shutdown(ctx)
				log.Info("open-telemetry exited properly")
				if err != nil {
					return
				}
			}
		}
	}()
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}))
	t := tp.Tracer(cfg.TracerName)
	return t, nil

}
