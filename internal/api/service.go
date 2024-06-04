package api

import (
	"github.com/KRUL-marketplace/favorite-service/internal/service"
	desc "github.com/KRUL-marketplace/favorite-service/pkg/favorite-service"
)

type Implementation struct {
	desc.UnimplementedFavoriteServiceServer
	favoriteService service.FavoriteService
}

func NewImplementation(favoriteService service.FavoriteService) *Implementation {
	return &Implementation{
		favoriteService: favoriteService,
	}
}
