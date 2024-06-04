package converter

import (
	"database/sql"
	"github.com/KRUL-marketplace/favorite-service/internal/repository/model"
	desc "github.com/KRUL-marketplace/favorite-service/pkg/favorite-service"
	product_service "github.com/KRUL-marketplace/product-catalog-service/pkg/product-catalog-service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToFavoriteListDescFromService(favoriteList *model.FavoriteList) *desc.GetFavoriteListByIdResponse {
	var favoriteListUpdatedAt *timestamppb.Timestamp
	if favoriteList.UpdatedAt.Valid {
		favoriteListUpdatedAt = timestamppb.New(favoriteList.UpdatedAt.Time)
	}

	result := desc.GetFavoriteListByIdResponse{
		FavoriteList: &desc.FavoriteList{
			FavoriteListId: favoriteList.FavoriteListID,
			UserId:         favoriteList.UserID,
			CreatedAt:      timestamppb.New(favoriteList.CreatedAt),
			UpdatedAt:      favoriteListUpdatedAt,
		},
	}

	for _, product := range favoriteList.Items {
		result.FavoriteList.Items = append(result.FavoriteList.Items, ToFavoriteItemDescFromService(&product))
	}

	return &result
}

func ToFavoriteItemDescFromService(favoriteItem *model.FavoriteItem) *desc.FavoriteItem {
	var favoriteItemUpdatedAt *timestamppb.Timestamp
	if favoriteItem.UpdatedAt.Valid {
		favoriteItemUpdatedAt = timestamppb.New(favoriteItem.UpdatedAt.Time)
	}

	return &desc.FavoriteItem{
		ItemId:      favoriteItem.ItemID,
		ProductId:   favoriteItem.ProductID,
		CreatedAt:   timestamppb.New(favoriteItem.CreatedAt),
		UpdatedAt:   favoriteItemUpdatedAt,
		ProductInfo: ToFavoriteInfoDescFromService(&favoriteItem.Info),
	}

}

func ToFavoriteInfoDescFromService(favoriteInfo *model.FavoriteProductInfo) *desc.FavoriteProductInfo {
	return &desc.FavoriteProductInfo{
		Name:  favoriteInfo.Name,
		Slug:  favoriteInfo.Slug,
		Image: favoriteInfo.Image,
		Price: favoriteInfo.Price,
		Brand: &desc.Brand{
			Id: favoriteInfo.Brand.ID,
			Info: &desc.BrandInfo{
				Name:        favoriteInfo.Brand.Info.Name,
				Slug:        favoriteInfo.Brand.Info.Slug,
				Description: favoriteInfo.Brand.Info.Description,
			},
			CreatedAt: timestamppb.New(favoriteInfo.Brand.CreatedAt),
			UpdatedAt: timestamppb.New(favoriteInfo.Brand.UpdatedAt.Time),
		},
	}

}

func ToFavoriteItemModelFromDesc(favoriteItem *desc.FavoriteItem) *model.FavoriteItem {
	var favoriteItemUpdatedAt *timestamppb.Timestamp
	if favoriteItem.GetUpdatedAt().IsValid() {
		favoriteItemUpdatedAt = favoriteItem.GetUpdatedAt()
	}

	return &model.FavoriteItem{
		ItemID:    favoriteItem.GetItemId(),
		ProductID: favoriteItem.GetProductId(),
		CreatedAt: favoriteItem.GetCreatedAt().AsTime(),
		UpdatedAt: sql.NullTime{Time: favoriteItemUpdatedAt.AsTime(), Valid: favoriteItemUpdatedAt.IsValid()},
		Info:      *ToFavoriteProductInfoModelFromDesc(favoriteItem.GetProductInfo()),
	}
}

func ToFavoriteProductInfoModelFromDesc(favoriteProductInfo *desc.FavoriteProductInfo) *model.FavoriteProductInfo {
	return &model.FavoriteProductInfo{
		Name:  favoriteProductInfo.GetName(),
		Slug:  favoriteProductInfo.GetSlug(),
		Image: favoriteProductInfo.GetImage(),
		Price: favoriteProductInfo.GetPrice(),
		Brand: model.Brand{
			ID: favoriteProductInfo.GetBrand().GetId(),
			Info: model.BrandInfo{
				Name:        favoriteProductInfo.GetBrand().GetInfo().GetName(),
				Slug:        favoriteProductInfo.GetBrand().GetInfo().GetSlug(),
				Description: favoriteProductInfo.GetBrand().GetInfo().GetDescription(),
			},
			CreatedAt: favoriteProductInfo.GetBrand().GetCreatedAt().AsTime(),
			UpdatedAt: sql.NullTime{
				Time:  favoriteProductInfo.GetBrand().GetUpdatedAt().AsTime(),
				Valid: favoriteProductInfo.GetBrand().GetUpdatedAt().IsValid(),
			},
		},
	}
}

func ToFavoriteProductInfoModelFromProductInfo(favoriteProductInfo *product_service.ProductInfo) *model.FavoriteProductInfo {
	return &model.FavoriteProductInfo{
		Name:  favoriteProductInfo.GetName(),
		Slug:  favoriteProductInfo.GetSlug(),
		Image: "test",
		Price: favoriteProductInfo.GetPrice(),
		Brand: model.Brand{
			ID: favoriteProductInfo.GetBrand().GetId(),
			Info: model.BrandInfo{
				Name:        favoriteProductInfo.GetBrand().GetInfo().GetName(),
				Slug:        favoriteProductInfo.GetBrand().GetInfo().GetSlug(),
				Description: favoriteProductInfo.GetBrand().GetInfo().GetDescription(),
			},
			CreatedAt: favoriteProductInfo.GetBrand().GetCreatedAt().AsTime(),
			UpdatedAt: sql.NullTime{
				Time:  favoriteProductInfo.GetBrand().GetUpdatedAt().AsTime(),
				Valid: favoriteProductInfo.GetBrand().GetUpdatedAt().IsValid(),
			},
		},
	}
}
