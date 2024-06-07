package app

import (
	"context"
	"github.com/KRUL-marketplace/favorite-service/client/db"
	"github.com/KRUL-marketplace/favorite-service/client/db/pg"
	"github.com/KRUL-marketplace/favorite-service/client/db/transaction"
	"github.com/KRUL-marketplace/favorite-service/internal/api"
	"github.com/KRUL-marketplace/favorite-service/internal/config"
	productCatalogServiceClient "github.com/KRUL-marketplace/favorite-service/internal/connector/product_service_catalog_connector"
	"github.com/KRUL-marketplace/favorite-service/internal/repository"
	"github.com/KRUL-marketplace/favorite-service/internal/service"
	product_service "github.com/KRUL-marketplace/product-catalog-service/pkg/product-catalog-service"
	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type serviceProvider struct {
	favoriteRepository repository.Repository
	favoriteService    service.FavoriteService

	grpcConfig                      config.GRPCConfig
	httpConfig                      config.HTTPConfig
	pgConfig                        config.PGConfig
	swaggerConfig                   config.SwaggerConfig
	productCatalogServiceGRPCConfig config.ProductCatalogServiceGRPCConfig
	redisConfig                     config.RedisConfig

	dbClient  db.Client
	txManager db.TxManager

	favoriteImpl *api.Implementation

	productCatalogServiceClient productCatalogServiceClient.ProductCatalogServiceClient

	redisClient *redis.Client
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := config.NewHTTPConfig()
		if err != nil {
			log.Fatalf("failed to get http config: %s", err.Error())
		}

		s.httpConfig = cfg
	}

	return s.httpConfig
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) SwaggerConfig() config.SwaggerConfig {
	if s.swaggerConfig == nil {
		cfg, err := config.NewSwaggerConfig()
		if err != nil {
			log.Fatalf("failed to get swagger config: %s", err.Error())
		}

		s.swaggerConfig = cfg
	}

	return s.swaggerConfig
}

func (s *serviceProvider) ProductCatalogServiceGRPCConfig() config.ProductCatalogServiceGRPCConfig {
	if s.productCatalogServiceGRPCConfig == nil {
		cfg, err := config.NewProductCatalogServiceGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get product catalog service grpc config: %s", err.Error())
		}

		s.productCatalogServiceGRPCConfig = cfg
	}

	return s.productCatalogServiceGRPCConfig
}

func (s *serviceProvider) RedisConfig() config.RedisConfig {
	if s.redisConfig == nil {
		cfg, err := config.NewRedisConfig()
		if err != nil {
			log.Fatalf("failed to get redis config: %s", err.Error())
		}

		s.redisConfig = cfg
	}

	return s.redisConfig
}

func (s *serviceProvider) FavoriteRepository(ctx context.Context) repository.Repository {
	if s.favoriteRepository == nil {
		s.favoriteRepository = repository.NewRepository(
			s.DBClient(ctx),
			s.RedisClient(ctx),
			s.ProductCatalogServiceClient(ctx),
		)
	}

	return s.favoriteRepository
}

func (s *serviceProvider) FavoriteService(ctx context.Context) service.FavoriteService {
	if s.favoriteService == nil {
		s.favoriteService = service.NewService(
			s.FavoriteRepository(ctx),
			s.TxManager(ctx),
		)
	}

	return s.favoriteService
}

func (s *serviceProvider) FavoriteImpl(ctx context.Context) *api.Implementation {
	if s.favoriteImpl == nil {
		s.favoriteImpl = api.NewImplementation(s.FavoriteService(ctx))
	}

	return s.favoriteImpl
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) RedisClient(ctx context.Context) redis.Client {
	if s.redisClient == nil {
		s.redisClient = redis.NewClient(&redis.Options{
			Addr:     s.RedisConfig().Address(),
			Password: "",
			DB:       0,
		})
	}

	return *s.redisClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) ProductCatalogServiceClient(ctx context.Context) productCatalogServiceClient.ProductCatalogServiceClient {
	if s.productCatalogServiceClient == nil {
		conn, err := grpc.DialContext(ctx,
			s.ProductCatalogServiceGRPCConfig().Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("product catalog service client init error")
		}

		s.productCatalogServiceClient = productCatalogServiceClient.NewProductCatalogServiceClient(
			product_service.NewProductCatalogServiceClient(conn),
		)
	}

	return s.productCatalogServiceClient
}
