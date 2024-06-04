package api

import (
	"context"
	desc "github.com/KRUL-marketplace/favorite-service/pkg/favorite-service"
)

func (i *Implementation) ToggleProduct(ctx context.Context, r *desc.ToggleProductRequest) (*desc.ToggleProductResponse, error) {
	err := i.favoriteService.ToggleProduct(ctx, r.GetUserId(), r.GetProductId())

	if err != nil {
		return &desc.ToggleProductResponse{
			Success: false,
		}, err
	}

	return &desc.ToggleProductResponse{
		Success: true,
	}, err
}
