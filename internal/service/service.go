package service

import (
	"context"
	"github.com/KRUL-marketplace/favorite-service/client/db"
	"github.com/KRUL-marketplace/favorite-service/internal/repository"
	"github.com/KRUL-marketplace/favorite-service/internal/repository/model"
)

type favoriteService struct {
	favoriteRepository repository.Repository
	txManager          db.TxManager
}

type FavoriteService interface {
	ToggleProduct(ctx context.Context, userID, productID string) error
	GetFavoriteListById(ctx context.Context, userId string) (*model.FavoriteList, error)
}

func NewService(favoriteRepository repository.Repository, txManager db.TxManager) FavoriteService {
	return &favoriteService{
		favoriteRepository: favoriteRepository,
		txManager:          txManager,
	}
}
