package model

import (
	"database/sql"
	"time"
)

type FavoriteList struct {
	FavoriteListID string         `db:"favorite_list_id"`
	UserID         string         `db:"user_id"`
	CreatedAt      time.Time      `db:"created_at"`
	UpdatedAt      sql.NullTime   `db:"updated_at"`
	Items          []FavoriteItem `db:"items"`
}

type FavoriteItem struct {
	ItemID    string              `db:"item_id"`
	ProductID string              `db:"product_id"`
	Info      FavoriteProductInfo `db:""`
	CreatedAt time.Time           `db:"created_at"`
	UpdatedAt sql.NullTime        `db:"updated_at"`
}

type FavoriteProductInfo struct {
	Name  string `db:"name" json:"name"`
	Slug  string `db:"slug" json:"slug"`
	Image string `db:"image" json:"image"`
	Price uint32 `db:"price" json:"price"`
	Brand Brand  `db:"brand" json:"brand"`
}

type Brand struct {
	ID        uint32       `db:"id"`
	Info      BrandInfo    `db:""`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}

type BrandInfo struct {
	Name        string `db:"name"`
	Slug        string `db:"slug"`
	Description string `db:"description"`
}
