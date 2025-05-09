package server

import (
	"log"

	"github.com/sekkarin/shop-microservice/modules/player/playerHandler"
	playerPb "github.com/sekkarin/shop-microservice/modules/player/playerPb"
	"github.com/sekkarin/shop-microservice/modules/player/playerRepository"
	"github.com/sekkarin/shop-microservice/modules/player/playerUsecase"
	"github.com/sekkarin/shop-microservice/pkg/grpccon"
)


func (s *server) playerService() {
	repo := playerRepository.NewPlayerRepository(s.db)
	usecase := playerUsecase.NewPlayerUsecase(repo)
	httpHandler := playerHandler.NewPlayerHttpHandler(s.cfg, usecase)
	grpcHandler := playerHandler.NewPlayerGrpcHandler(usecase)
	queueHandler := playerHandler.NewPlayerQueueHandler(s.cfg, usecase)

	go queueHandler.DockedPlayerMoney()
	go queueHandler.AddPlayerMoney()
	go queueHandler.RollbackPlayerTransaction()

	// gRPC
	go func() {
		grpcServer, lis := grpccon.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.PlayerUrl)

		playerPb.RegisterPlayerGrpcServiceServer(grpcServer, grpcHandler)

		log.Printf("Player gRPC server listening on %s", s.cfg.Grpc.PlayerUrl)
		grpcServer.Serve(lis)
	}()

	player := s.app.Group("/player_v1")

	// Health Check
	player.GET("", s.healthCheckService)

	player.POST("/player/register", httpHandler.CreatePlayer)
	player.POST("/player/add-money", httpHandler.AddPlayerMoney, s.middleware.JwtAuthorization)
	player.GET("/player/:player_id", httpHandler.FindOnePlayerProfile)
	player.GET("/player/saving-account/my-account", httpHandler.GetPlayerSavingAccount, s.middleware.JwtAuthorization)
}
