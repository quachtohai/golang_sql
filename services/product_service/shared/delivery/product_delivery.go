package delivery

import (
	"context"
	"golang_sql/pkg/logger"
	"golang_sql/pkg/rabbitmq"
	"golang_sql/services/product_service/config"
	"golang_sql/services/product_service/product/data/contracts"

	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

type ProductDeliveryBase struct {
	Log               logger.ILogger
	Cfg               *config.Config
	RabbitmqPublisher rabbitmq.IPublisher
	ConnRabbitmq      *amqp.Connection
	HttpClient        *resty.Client
	JaegerTracer      trace.Tracer
	Gorm              *gorm.DB
	Echo              *echo.Echo
	ProductRepository contracts.ProductRepository
	Ctx               context.Context
}
