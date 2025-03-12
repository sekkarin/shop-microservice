package server

import (
	"log"

	"github.com/sekkarin/shop-microservice/modules/auth/authHandler"
	authPb "github.com/sekkarin/shop-microservice/modules/auth/authPb"
	"github.com/sekkarin/shop-microservice/modules/auth/authRepository"
	"github.com/sekkarin/shop-microservice/modules/auth/authUsecase"
	"github.com/sekkarin/shop-microservice/pkg/grpccon"
)

func (s *server) authService() {
	repo := authRepository.NewAuthRepository(s.db)
	usecase := authUsecase.NewAuthUsecase(repo)
	httpHandler := authHandler.NewAuthHttpHandler(s.cfg, usecase)
	grpcHandler := authHandler.NewAuthGrpcHandler(usecase)

	// gRPC test5
	go func() {
		grpcServer, lis := grpccon.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.AuthUrl)

		authPb.RegisterAuthGrpcServiceServer(grpcServer, grpcHandler)

		log.Printf("Auth gRPC server listening on %s", s.cfg.Grpc.AuthUrl)
		grpcServer.Serve(lis)
	}()
	auth := s.app.Group("/auth_v1")

	// Health Check
	auth.GET("", s.healthCheckService)

	auth.GET("/test/:player_id", s.healthCheckService)
	auth.POST("/auth/login", httpHandler.Login)
	auth.POST("/auth/refresh-token", httpHandler.RefreshToken)
	auth.POST("/auth/logout", httpHandler.Logout)
}
