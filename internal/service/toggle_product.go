package service

import (
	"context"
)

func (s *favoriteService) ToggleProduct(ctx context.Context, userID, productID string) error {
	err := s.favoriteRepository.ToggleProduct(ctx, userID, productID)
	return err
}
