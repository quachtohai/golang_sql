package mappings

import (
	"golang_sql/pkg/mapper"
	"golang_sql/services/product_service/product/dtos"
	"golang_sql/services/product_service/product/features/creating_product/v1/events"
	events2 "golang_sql/services/product_service/product/features/updating_product/v1/events"
	"golang_sql/services/product_service/product/models"
)

func ConfigureMappings() error {
	err := mapper.CreateMap[*models.Product, *dtos.ProductDto]()
	if err != nil {
		return err
	}

	err = mapper.CreateMap[*models.Product, *events.ProductCreated]()
	if err != nil {
		return err
	}

	err = mapper.CreateMap[*models.Product, *events2.ProductUpdated]()
	if err != nil {
		return err
	}
	return nil
}
