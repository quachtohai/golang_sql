package mappings

import (
	"golang_sql/pkg/mapper"
	"golang_sql/services/inventory_service/inventory/consumers/events"
	"golang_sql/services/inventory_service/inventory/dtos"
	"golang_sql/services/inventory_service/inventory/models"
)

func ConfigureMappings() error {
	err := mapper.CreateMap[*models.Inventory, *dtos.InventoryDto]()
	if err != nil {
		return err
	}

	err = mapper.CreateMap[*models.ProductItem, *events.InventoryUpdated]()
	if err != nil {
		return err
	}

	return nil
}
