package initialize

import (
	"github.com/ppthinh/ChatApp/services/user-service/internal/handler"
	model "github.com/ppthinh/ChatApp/services/user-service/internal/models"
	"github.com/ppthinh/ChatApp/services/user-service/internal/repository"
	"github.com/ppthinh/ChatApp/services/user-service/internal/service"
	"github.com/ppthinh/ChatApp/services/user-service/pkg/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func Run() {
	cfg := config.LoadConfig()

	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Can't conn to database: %v", err)
	}
	// gen db
	db.AutoMigrate(&model.User{}, &model.FriendShip{})

	repo := repository.NewUserRepository(db)
	svc := service.NewUserService(repo)

	handler.RunUserGRPCServer(svc, cfg.GRPCPort)

}
