package api

import (
	"context"
	"github.com/KRUL-marketplace/favorite-service/internal/converter"
	desc "github.com/KRUL-marketplace/favorite-service/pkg/favorite-service"
)

func (i *Implementation) GetFavoriteListById(ctx context.Context, req *desc.GetFavoriteListByIdRequest) (*desc.GetFavoriteListByIdResponse, error) {
	favoriteList, err := i.favoriteService.GetFavoriteListById(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}

	return converter.ToFavoriteListDescFromService(favoriteList), nil
}
