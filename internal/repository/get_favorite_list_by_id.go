package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/KRUL-marketplace/common-libs/pkg/client/db"
	"github.com/KRUL-marketplace/favorite-service/internal/converter"
	"github.com/KRUL-marketplace/favorite-service/internal/repository/model"
	sq "github.com/Masterminds/squirrel"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"time"
)

const (
	REDIS_EXPIRATION = 5 * time.Minute
)

func (r *repo) GetFavoriteListById(ctx context.Context, userId string) (*model.FavoriteList, error) {
	cachedFavoriteList, err := r.getFavoriteListFromRedis(ctx, userId)
	if err != nil {
		return nil, err
	}

	if cachedFavoriteList != nil {
		return cachedFavoriteList, nil
	}

	builder := sq.Select(
		"f.favorite_list_id",
		"f.user_id",
		"f.created_at",
		"f.updated_at",
		"fi.item_id",
		"fi.product_id",
		"fi.created_at",
		"fi.updated_at",
	).
		From("favorite_lists f").
		LeftJoin("favorite_items fi ON f.favorite_list_id = fi.favorite_list_id").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"f.user_id": userId})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "favorite_repository.GetByUserId " + userId,
		QueryRaw: query,
	}

	var favoriteList model.FavoriteList
	rows, err := r.db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	favoriteList.Items = []model.FavoriteItem{}
	var productIDS []string

	for rows.Next() {
		var favoriteItem model.FavoriteItem
		var itemID, productID sql.NullString
		var createdAt, updatedAt sql.NullTime

		err := rows.Scan(
			&favoriteList.FavoriteListID,
			&favoriteList.UserID,
			&favoriteList.CreatedAt,
			&favoriteList.UpdatedAt,
			&itemID,
			&productID,
			&createdAt,
			&updatedAt,
		)

		if err != nil {
			return nil, err
		}

		if itemID.Valid {
			favoriteItem.ItemID = itemID.String
		} else {
			favoriteItem.ItemID = ""
		}
		if productID.Valid {
			favoriteItem.ProductID = productID.String
		} else {
			favoriteItem.ProductID = ""
		}
		if createdAt.Valid {
			favoriteItem.CreatedAt = createdAt.Time
		} else {
			favoriteItem.CreatedAt = time.Time{}
		}
		if updatedAt.Valid {
			favoriteItem.UpdatedAt = updatedAt
		} else {
			favoriteItem.UpdatedAt = sql.NullTime{}
		}

		if productID.Valid {
			productIDS = append(productIDS, productID.String)
			favoriteList.Items = append(favoriteList.Items, favoriteItem)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	result, err := r.productCatalogServiceClient.GetById(ctx, productIDS)
	if err != nil {
		return nil, err
	}

	for i, product := range result.GetProduct() {
		favoriteList.Items[i] = model.FavoriteItem{
			ItemID:    favoriteList.Items[i].ItemID,
			ProductID: favoriteList.Items[i].ProductID,
			CreatedAt: favoriteList.Items[i].CreatedAt,
			UpdatedAt: favoriteList.Items[i].UpdatedAt,
			Info:      *converter.ToFavoriteProductInfoModelFromProductInfo(product.GetInfo()),
		}
	}

	if err := r.setFavoriteListFromRedis(ctx, userId, &favoriteList); err != nil {
		return nil, err
	}

	return &favoriteList, nil
}

func (r *repo) getFavoriteListFromRedis(ctx context.Context, userId string) (*model.FavoriteList, error) {
	data, err := r.redisClient.Get(ctx, fmt.Sprintf("favorite_list:%s", userId)).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var favoriteList model.FavoriteList
	if err := json.Unmarshal([]byte(data), &favoriteList); err != nil {
		return nil, err
	}

	return &favoriteList, nil
}

func (r *repo) setFavoriteListFromRedis(ctx context.Context, userId string, favoriteList *model.FavoriteList) error {
	data, err := json.Marshal(favoriteList)
	if err != nil {
		return err
	}

	return r.redisClient.Set(ctx, fmt.Sprintf("favorite_list:%s", userId), data, REDIS_EXPIRATION).Err()
}
