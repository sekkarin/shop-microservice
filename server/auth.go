package server

import (
	"log"

	"github.com/labstack/echo"
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

	// gRPC
	go func() {
		grpcServer, lis := grpccon.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.AuthUrl)

		authPb.RegisterAuthGrpcServiceServer(grpcServer, grpcHandler)

		log.Printf("Auth gRPC server listening on %s", s.cfg.Grpc.AuthUrl)
		grpcServer.Serve(lis)
	}()
	s.app.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("X-Content-Type-Options", "nosniff")
			return next(c)
		}
	})
	auth := s.app.Group("")

	// Health Check
	auth.GET("", s.healthCheckService)

	auth.GET("/test/:player_id", s.healthCheckService)
	auth.POST("/auth/login", httpHandler.Login)
	auth.POST("/auth/refresh-token", httpHandler.RefreshToken)
	auth.POST("/auth/logout", httpHandler.Logout)
}
