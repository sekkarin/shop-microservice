package itemHandler

import (
	"context"

	itemPb "github.com/sekkarin/shop-microservice/modules/item/itemPb"
	"github.com/sekkarin/shop-microservice/modules/item/itemUsecase"
)

type (
	itemGrpcHandler struct {
		itemUsecase itemUsecase.ItemUsecaseService
		itemPb.UnimplementedItemGrpcServiceServer
	}
)

// test3
func NewItemGrpcHandler(itemUsecase itemUsecase.ItemUsecaseService) *itemGrpcHandler {
	return &itemGrpcHandler{
		itemUsecase: itemUsecase,
	}
}

func (g *itemGrpcHandler) FindItemsInIds(ctx context.Context, req *itemPb.FindItemsInIdsReq) (*itemPb.FindItemsInIdsRes, error) {
	return g.itemUsecase.FindItemInIds(ctx, req)
}
