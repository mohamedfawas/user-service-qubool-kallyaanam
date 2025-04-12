package di

import (
	"github.com/mohamedfawas/user-service-qubool-kallyaanam/internal/config"
	"github.com/mohamedfawas/user-service-qubool-kallyaanam/internal/repository"
	postgresRepo "github.com/mohamedfawas/user-service-qubool-kallyaanam/internal/repository/postgres"
	"github.com/mohamedfawas/user-service-qubool-kallyaanam/internal/service"
	"github.com/mohamedfawas/user-service-qubool-kallyaanam/internal/utils/logger"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Container holds application dependencies
type Container struct {
	Config             *config.Config
	DB                 *gorm.DB
	Logger             *logger.Logger
	UserProfileRepo    repository.UserProfileRepository
	UserProfileService service.UserProfileService
}

// NewContainer initializes the dependency container
func NewContainer(cfg *config.Config) (*Container, error) {
	// Initialize logger
	log, err := logger.NewLogger(cfg.Logging.IsDevelopment)
	if err != nil {
		return nil, err
	}

	// Initialize database
	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), &gorm.Config{})
	if err != nil {
		log.Error("Failed to connect to database", zap.Error(err))
		return nil, err
	}

	// Initialize repositories
	userProfileRepo := postgresRepo.NewUserProfileRepository(db)

	// Initialize services
	userProfileService := service.NewUserProfileService(userProfileRepo, log)

	return &Container{
		Config:             cfg,
		DB:                 db,
		Logger:             log,
		UserProfileRepo:    userProfileRepo,
		UserProfileService: userProfileService,
	}, nil
}
