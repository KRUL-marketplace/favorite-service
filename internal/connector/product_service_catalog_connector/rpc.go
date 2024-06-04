package product_service_catalog_connector

import (
	"context"
	product_service "github.com/KRUL-marketplace/product-catalog-service/pkg/product-catalog-service"
)

type ProductCatalogServiceClient interface {
	GetById(ctx context.Context, id []string) (*product_service.GetProductByIdResponse, error)
}
