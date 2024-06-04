package repository

import (
	"context"
	"github.com/KRUL-marketplace/favorite-service/client/db"
	productCatalogServiceClient "github.com/KRUL-marketplace/favorite-service/internal/connector/product_service_catalog_connector"
	"github.com/KRUL-marketplace/favorite-service/internal/repository/model"
)

type Repository interface {
	ToggleProduct(ctx context.Context, userID, productID string) error
	GetFavoriteListById(ctx context.Context, userId string) (*model.FavoriteList, error)
}

type repo struct {
	db                          db.Client
	productCatalogServiceClient productCatalogServiceClient.ProductCatalogServiceClient
}

func NewRepository(db db.Client, productCatalogServiceClient productCatalogServiceClient.ProductCatalogServiceClient) Repository {
	return &repo{
		db:                          db,
		productCatalogServiceClient: productCatalogServiceClient,
	}
}
