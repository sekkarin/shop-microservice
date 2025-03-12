package server

import (
	"log"

	"github.com/sekkarin/shop-microservice/modules/item/itemHandler"
	itemPb "github.com/sekkarin/shop-microservice/modules/item/itemPb"
	"github.com/sekkarin/shop-microservice/modules/item/itemRepository"
	"github.com/sekkarin/shop-microservice/modules/item/itemUsecase"
	"github.com/sekkarin/shop-microservice/pkg/grpccon"
)

func (s *server) itemService() {
	repo := itemRepository.NewItemRepository(s.db)
	usecase := itemUsecase.NewItemUsecase(repo)
	httpHandler := itemHandler.NewItemHttpHandler(s.cfg, usecase)
	grpcHandler := itemHandler.NewItemGrpcHandler(usecase)

	// gRPC test3
	go func() {
		grpcServer, lis := grpccon.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.ItemUrl)

		itemPb.RegisterItemGrpcServiceServer(grpcServer, grpcHandler)

		log.Printf("Item gRPC server listening on %s", s.cfg.Grpc.ItemUrl)
		grpcServer.Serve(lis)
	}()

	_ = grpcHandler

	item := s.app.Group("/item_v1")

	// Health Check
	item.GET("", s.healthCheckService)

	item.POST("/item", s.middleware.JwtAuthorization(s.middleware.RbacAuthorization(httpHandler.CreateItem, []int{1, 0})))
	item.GET("/item/:item_id", httpHandler.FindOneItem)
	item.GET("/item", httpHandler.FindManyItems)
	item.PATCH("/item/:item_id", s.middleware.JwtAuthorization(s.middleware.RbacAuthorization(httpHandler.EditItem, []int{1, 0})))
	item.PATCH("/item/:item_id/is-activated", s.middleware.JwtAuthorization(s.middleware.RbacAuthorization(httpHandler.EnableOrDisableItem, []int{1, 0})))
}
