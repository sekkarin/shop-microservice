package itemUsecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/sekkarin/shop-microservice/modules/item"
	itemPb "github.com/sekkarin/shop-microservice/modules/item/itemPb"
	"github.com/sekkarin/shop-microservice/modules/item/itemRepository"
	"github.com/sekkarin/shop-microservice/modules/models"
	"github.com/sekkarin/shop-microservice/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	ItemUsecaseService interface {
		CreateItem(pctx context.Context, req *item.CreateItemReq) (*item.ItemShowCase, error)
		FindOneItem(pctx context.Context, itemId string) (*item.ItemShowCase, error)
		FindManyItems(pctx context.Context, basePaginateUrl string, req *item.ItemSearchReq) (*models.PaginateRes, error)
		EditItem(pctx context.Context, itemId string, req *item.ItemUpdateReq) (*item.ItemShowCase, error)
		EnableOrDisableItem(pctx context.Context, itemId string) (bool, error)
		FindItemInIds(pctx context.Context, req *itemPb.FindItemsInIdsReq) (*itemPb.FindItemsInIdsRes, error)
	}

	itemUsecase struct {
		itemRepository itemRepository.ItemRepositoryService
	}
)

func NewItemUsecase(itemRepository itemRepository.ItemRepositoryService) ItemUsecaseService {
	return &itemUsecase{itemRepository}
}

func (u *itemUsecase) CreateItem(pctx context.Context, req *item.CreateItemReq) (*item.ItemShowCase, error) {
	if !u.itemRepository.IsUniqueItem(pctx, req.Title) {
		return nil, errors.New("error: this title is already exist")
	}

	itemId, err := u.itemRepository.InsertOneItem(pctx, &item.Item{
		Title:       req.Title,
		Price:       req.Price,
		Damage:      req.Damage,
		UsageStatus: true,
		ImageUrl:    req.ImageUrl,
		CreatedAt:   utils.LocalTime(),
		UpdatedAt:   utils.LocalTime(),
	})
	if err != nil {
		return nil, err
	}

	return u.FindOneItem(pctx, itemId.Hex())
}

func (u *itemUsecase) FindOneItem(pctx context.Context, itemId string) (*item.ItemShowCase, error) {
	result, err := u.itemRepository.FindOneItem(pctx, itemId)
	if err != nil {
		return nil, err
	}

	return &item.ItemShowCase{
		ItemId:   "item:" + result.Id.Hex(),
		Title:    result.Title,
		Price:    result.Price,
		Damage:   result.Damage,
		ImageUrl: result.ImageUrl,
	}, nil
}

func (u *itemUsecase) FindManyItems(pctx context.Context, basePaginateUrl string, req *item.ItemSearchReq) (*models.PaginateRes, error) {
	findItemsFilter := bson.D{}
	findItemsOpts := make([]*options.FindOptions, 0)

	countItemsFilter := bson.D{}

	// Filter
	if req.Start != "" {
		req.Start = strings.TrimPrefix(req.Start, "item:")
		findItemsFilter = append(findItemsFilter, bson.E{"_id", bson.D{{"$gt", utils.ConvertToObjectId(req.Start)}}})
	}

	if req.Title != "" {
		findItemsFilter = append(findItemsFilter, bson.E{"title", primitive.Regex{Pattern: req.Title, Options: "i"}})
		countItemsFilter = append(countItemsFilter, bson.E{"title", primitive.Regex{Pattern: req.Title, Options: "i"}})
	}

	findItemsFilter = append(findItemsFilter, bson.E{"usage_status", true})
	countItemsFilter = append(countItemsFilter, bson.E{"usage_status", true})

	// Options
	findItemsOpts = append(findItemsOpts, options.Find().SetSort(bson.D{{"_id", 1}}))
	findItemsOpts = append(findItemsOpts, options.Find().SetLimit(int64(req.Limit)))

	// Find
	results, err := u.itemRepository.FindManyItems(pctx, findItemsFilter, findItemsOpts)
	if err != nil {
		return nil, err
	}

	// Count
	total, err := u.itemRepository.CountItems(pctx, countItemsFilter)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return &models.PaginateRes{
			Data:  make([]*item.ItemShowCase, 0),
			Total: total,
			Limit: req.Limit,
			First: models.FirstPaginate{
				Href: fmt.Sprintf("%s?limit=%d&title=%s", basePaginateUrl, req.Limit, req.Title),
			},
			Next: models.NextPaginate{
				Start: "",
				Href:  "",
			},
		}, nil
	}

	return &models.PaginateRes{
		Data:  results,
		Total: total,
		Limit: req.Limit,
		First: models.FirstPaginate{
			Href: fmt.Sprintf("%s?limit=%d&title=%s", basePaginateUrl, req.Limit, req.Title),
		},
		Next: models.NextPaginate{
			Start: results[len(results)-1].ItemId,
			Href:  fmt.Sprintf("%s?limit=%d&title=%s&start=%s", basePaginateUrl, req.Limit, req.Title, results[len(results)-1].ItemId),
		},
	}, nil
}

func (u *itemUsecase) EditItem(pctx context.Context, itemId string, req *item.ItemUpdateReq) (*item.ItemShowCase, error) {
	// Update logical
	updateReq := bson.M{}
	if req.Title != "" {
		if !u.itemRepository.IsUniqueItem(pctx, req.Title) {
			log.Println("Error: EditItem failed: this title is already exist")
			return nil, errors.New("error: this title is already exist")
		}

		updateReq["title"] = req.Title
	}
	if req.ImageUrl != "" {
		updateReq["image_url"] = req.ImageUrl
	}
	if req.Damage > 0 {
		updateReq["damage"] = req.Damage
	}
	if req.Price >= 0 {
		updateReq["price"] = req.Price
	}
	updateReq["updated_at"] = utils.LocalTime()

	if err := u.itemRepository.UpdateOneItem(pctx, itemId, updateReq); err != nil {
		return nil, err
	}

	return u.FindOneItem(pctx, itemId)
}

func (u *itemUsecase) EnableOrDisableItem(pctx context.Context, itemId string) (bool, error) {
	result, err := u.itemRepository.FindOneItem(pctx, itemId)
	if err != nil {
		return false, err
	}

	if err := u.itemRepository.EnableOrDisableItem(pctx, itemId, !result.UsageStatus); err != nil {
		return false, err
	}

	return !result.UsageStatus, nil
}

func (u *itemUsecase) FindItemInIds(pctx context.Context, req *itemPb.FindItemsInIdsReq) (*itemPb.FindItemsInIdsRes, error) {
	filter := bson.D{}

	objectIds := make([]primitive.ObjectID, 0)
	for _, itemId := range req.Ids {
		objectIds = append(objectIds, utils.ConvertToObjectId(strings.TrimPrefix(itemId, "item:")))
	}

	filter = append(filter, bson.E{"_id", bson.D{{"$in", objectIds}}})
	filter = append(filter, bson.E{"usage_status", true})

	results, err := u.itemRepository.FindManyItems(pctx, filter, nil)
	if err != nil {
		return nil, err
	}

	resultsToRes := make([]*itemPb.Item, 0)
	for _, result := range results {
		resultsToRes = append(resultsToRes, &itemPb.Item{
			Id:       result.ItemId,
			Title:    result.Title,
			Price:    result.Price,
			Damage:   int32(result.Damage),
			ImageUrl: result.ImageUrl,
		})
	}

	return &itemPb.FindItemsInIdsRes{
		Items: resultsToRes,
	}, nil
}
