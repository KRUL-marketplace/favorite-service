package service

import (
	"context"
	"github.com/KRUL-marketplace/favorite-service/internal/repository/model"
)

func (s *favoriteService) GetFavoriteListById(ctx context.Context, userId string) (*model.FavoriteList, error) {
	favoriteList, err := s.favoriteRepository.GetFavoriteListById(ctx, userId)
	if err != nil {
		return nil, err
	}

	return favoriteList, nil
}
