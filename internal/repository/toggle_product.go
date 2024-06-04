package repository

import (
	"context"
	"github.com/KRUL-marketplace/favorite-service/client/db"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"log"
)

func (r *repo) ToggleProduct(ctx context.Context, userID, productID string) error {
	builderSelectFavoriteListID := sq.Select("favorite_list_id").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"user_id": userID}).
		From("favorite_lists")

	query, args, err := builderSelectFavoriteListID.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "favorite_repository.FindFavoriteList",
		QueryRaw: query,
	}

	var favoriteListID string
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&favoriteListID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			favoriteListID, err = r.createFavoriteList(ctx, userID)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	_, err = r.productCatalogServiceClient.GetById(ctx, []string{productID})

	if err != nil {
		return err
	}

	builderSelectFavoriteItem := sq.Select("1").
		PlaceholderFormat(sq.Dollar).
		From("favorite_items").
		Where(sq.Eq{"favorite_list_id": favoriteListID, "product_id": productID})

	query, args, err = builderSelectFavoriteItem.ToSql()
	if err != nil {
		return err
	}

	q = db.Query{
		Name:     "favorite_repository.FindFavoriteItem",
		QueryRaw: query,
	}

	var exists int
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&exists)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	if exists == 1 {
		deleteQuery := sq.Delete("favorite_items").
			PlaceholderFormat(sq.Dollar).
			Where(sq.Eq{"favorite_list_id": favoriteListID, "product_id": productID})

		query, args, err = deleteQuery.ToSql()
		if err != nil {
			return err
		}

		q = db.Query{
			Name:     "favorite_repository.DeleteFavoriteItem",
			QueryRaw: query,
		}

		_, err = r.db.DB().ExecContext(ctx, q, args...)
		if err != nil {
			return err
		}
	} else {
		insertQuery := sq.Insert("favorite_items").
			PlaceholderFormat(sq.Dollar).
			Columns("item_id", "favorite_list_id", "product_id").
			Values(uuid.New(), favoriteListID, productID).
			Suffix("RETURNING item_id")

		query, args, err = insertQuery.ToSql()
		if err != nil {
			return err
		}

		q = db.Query{
			Name:     "favorite_repository.InsertFavoriteItem",
			QueryRaw: query,
		}

		_, err = r.db.DB().ExecContext(ctx, q, args...)
		if err != nil {
			return err
		}
	}

	log.Printf("Toggled product %s in favorite list %s for user %s\n", productID, favoriteListID, userID)

	return nil
}

func (r *repo) createFavoriteList(ctx context.Context, userID string) (string, error) {
	builder := sq.Insert("favorite_lists").
		PlaceholderFormat(sq.Dollar).
		Columns("user_id").
		Values(userID).
		Suffix("RETURNING favorite_list_id")

	query, args, err := builder.ToSql()
	if err != nil {
		return "", err
	}

	var favoriteListID string
	q := db.Query{
		Name:     "favorite_repository.createFavoriteList",
		QueryRaw: query,
	}

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&favoriteListID)
	if err != nil {
		return "", err
	}

	return favoriteListID, nil
}
