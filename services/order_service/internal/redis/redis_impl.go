package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/randnull/Lessons/internal/custom_errors"

	//"github.com/lib/pq"
	"log"

	//"github.com/gofiber/fiber/v2/log"
	"github.com/randnull/Lessons/internal/config"
	"github.com/randnull/Lessons/internal/models"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

type Redis struct {
	client *redis.Client
}

func NewRedis(cfg config.RedisConfig) *Redis {
	addr := fmt.Sprintf("%v:%v", cfg.Host, cfg.Port)
	fmt.Println(addr)

	redisClient := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  500 * time.Millisecond,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := redisClient.Ping(ctx).Result()

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Redis: Ready!")

	return &Redis{
		client: redisClient,
	}
}

func (red *Redis) SaveNewOrder(order *models.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	OrderKey := fmt.Sprintf("order:%v", order.ID)
	allOrdersKey := "orders:all"
	userOrdersKey := fmt.Sprintf("orders:student:%v", order.StudentID)

	tagsJson, err := json.Marshal(order.Tags)

	if err != nil {
		log.Println(err)
		return err
	}

	_, err = red.client.HSet(ctx, OrderKey, map[string]interface{}{
		"id":             order.ID,
		"student_id":     order.StudentID,
		"title":          order.Title,
		"description":    order.Description,
		"grade":          order.Grade,
		"min_price":      order.MinPrice,
		"max_price":      order.MaxPrice,
		"tags":           tagsJson,
		"status":         order.Status,
		"response_count": order.ResponseCount,
		"created_at":     order.CreatedAt.Format(time.RFC3339),
		"updated_at":     order.UpdatedAt.Format(time.RFC3339),
	}).Result()

	err = red.client.LPush(ctx, allOrdersKey, order.ID).Err()

	if err != nil {
		log.Println(err)
		return err
	}

	err = red.client.LPush(ctx, userOrdersKey, order.ID).Err()

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (red *Redis) AddOrder(order *models.OrderDetails) error {
	// объединить с верхним
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	OrderKey := fmt.Sprintf("order:%v", order.ID)
	allOrdersKey := "orders:all"
	userOrdersKey := fmt.Sprintf("orders:student:%v", order.StudentID)

	tagsJson, err := json.Marshal(order.Tags)

	if err != nil {
		log.Println(err)
		return err
	}

	_, err = red.client.HSet(ctx, OrderKey, map[string]interface{}{
		"id":             order.ID,
		"student_id":     order.StudentID,
		"title":          order.Title,
		"description":    order.Description,
		"grade":          order.Grade,
		"min_price":      order.MinPrice,
		"max_price":      order.MaxPrice,
		"tags":           tagsJson,
		"status":         order.Status,
		"response_count": order.ResponseCount,
		"created_at":     order.CreatedAt.Format(time.RFC3339),
		"updated_at":     order.UpdatedAt.Format(time.RFC3339),
	}).Result()

	err = red.client.LPush(ctx, allOrdersKey, order.ID).Err()

	if err != nil {
		log.Println(err)
		return err
	}

	err = red.client.LPush(ctx, userOrdersKey, order.ID).Err()

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (red *Redis) GetOrderHSET(ctx context.Context, orderID string) (*models.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	OrderKey := fmt.Sprintf("order:%v", orderID)

	orderRaw, err := red.client.HGetAll(ctx, OrderKey).Result()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if len(orderRaw) == 0 {
		log.Println("no data")
		return nil, nil // log.Println(" заказ не найден")
	}

	responseCount, _ := strconv.Atoi(orderRaw["response_count"])
	minPrice, _ := strconv.Atoi(orderRaw["min_price"])
	maxPrice, _ := strconv.Atoi(orderRaw["max_price"])

	createdAt, _ := time.Parse(time.RFC3339, orderRaw["created_at"])
	updatedAt, _ := time.Parse(time.RFC3339, orderRaw["updated_at"])

	//tags_raw := data["tags"].([]string)
	//tags := pq.StringArray{data["tags"].([]string)}
	//fmt.Println(data["tags"], tags)

	order := models.Order{
		ID:            orderRaw["id"],
		StudentID:     orderRaw["student_id"],
		Title:         orderRaw["title"],
		Description:   orderRaw["description"],
		Grade:         orderRaw["grade"],
		MinPrice:      minPrice,
		MaxPrice:      maxPrice,
		Tags:          nil,
		Status:        orderRaw["status"],
		ResponseCount: responseCount,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}

	log.Println("no data", order)

	return &order, nil
}

func (red *Redis) DeleteOrderHSET(orderID string) error {
	OrderKey := fmt.Sprintf("order:%v", orderID)

	_, err := red.client.Del(context.Background(), OrderKey).Result()
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (red *Redis) UpdateOrderHSET(orderID string, updateOrder *models.UpdateOrder) error {
	OrderKey := fmt.Sprintf("order:%v", orderID)

	updater := map[string]interface{}{}

	if updateOrder.Title != "" {
		updater["title"] = updateOrder.Title
	}

	if updateOrder.Description != "" {
		updater["description"] = updateOrder.Description
	}

	if updateOrder.Grade != "" {
		updater["grade"] = updateOrder.Grade
	}

	_, err := red.client.HSet(context.Background(), OrderKey, updater).Result()

	if err != nil {
		return err
	}

	return nil
}

func (red *Redis) GetAllOrders(ctx context.Context) ([]*models.Order, error) {
	allOrderKey := "orders:all"

	orderIds, err := red.client.LRange(ctx, allOrderKey, 0, -1).Result()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	if len(orderIds) == 0 {
		log.Println("no data")
		return nil, custom_errors.ErrNoOrders
	}

	var orders []*models.Order

	for _, orderID := range orderIds {
		order, err := red.GetOrderHSET(ctx, orderID)
		if err == nil {
			orders = append(orders, order)
		}
	}

	return orders, nil
}

func (red *Redis) GetAllUsersOrders(ctx context.Context, userID string) ([]*models.Order, error) {
	userOrdersKey := fmt.Sprintf("orders:student:%v", userID)

	orderIds, err := red.client.LRange(ctx, userOrdersKey, 0, -1).Result()

	if err != nil {
		return nil, err
	}

	var orders []*models.Order

	for _, orderID := range orderIds {
		order, err := red.GetOrderHSET(ctx, orderID)
		if err == nil {
			orders = append(orders, order)
		}
	}

	return orders, nil
}

func (red *Redis) SaveNewResponse(response *models.ResponseDB) error {
	responseKey := fmt.Sprintf("response:%v", response.ID)

	_, err := red.client.HSet(context.Background(), responseKey, map[string]interface{}{
		"id":             response.ID,
		"order_id":       response.OrderID,
		"tutor_id":       response.TutorID,
		"tutor_username": response.TutorUsername,
		"name":           response.Name,
		"created_at":     response.CreatedAt.Format(time.RFC3339),
	}).Result()

	if err != nil {
		return err
	}

	orderResponsesKey := fmt.Sprintf("order:%v:responses", response.OrderID)

	err = red.client.LPush(context.Background(), orderResponsesKey, response.ID).Err()

	if err != nil {
		return err
	}

	_, err = red.client.HIncrBy(context.Background(), "order:"+response.OrderID, "response_count", 1).Result()

	if err != nil {
		return err
	}

	return nil
}

func (red *Redis) AddResponse(response *models.Response) error {
	responseKey := fmt.Sprintf("response:%v", response.ID)

	_, err := red.client.HSet(context.Background(), responseKey, map[string]interface{}{
		"id":             response.ID,
		"order_id":       response.OrderID,
		"tutor_id":       response.TutorID,
		"tutor_username": response.TutorUsername,
		"name":           response.Name,
		"created_at":     response.CreatedAt.Format(time.RFC3339),
	}).Result()

	if err != nil {
		return err
	}

	orderResponsesKey := fmt.Sprintf("order:%v:responses", response.OrderID)

	err = red.client.LPush(context.Background(), orderResponsesKey, response.ID).Err()

	if err != nil {
		return err
	}

	_, err = red.client.HIncrBy(context.Background(), "order:"+response.OrderID, "response_count", 1).Result()

	if err != nil {
		return err
	}

	return nil
}

func (red *Redis) AddResponses(response []*models.Response) error {
	responseKey := fmt.Sprintf("response:%v", response.ID)

	_, err := red.client.HSet(context.Background(), responseKey, map[string]interface{}{
		"id":             response.ID,
		"order_id":       response.OrderID,
		"tutor_id":       response.TutorID,
		"tutor_username": response.TutorUsername,
		"name":           response.Name,
		"created_at":     response.CreatedAt.Format(time.RFC3339),
	}).Result()

	if err != nil {
		return err
	}

	orderResponsesKey := fmt.Sprintf("order:%v:responses", response.OrderID)

	err = red.client.LPush(context.Background(), orderResponsesKey, response.ID).Err()

	if err != nil {
		return err
	}

	_, err = red.client.HIncrBy(context.Background(), "order:"+response.OrderID, "response_count", 1).Result()

	if err != nil {
		return err
	}

	return nil
}

func (red *Redis) GetAllResponses(ctx context.Context, orderID string) ([]models.Response, error) {
	responseKey := fmt.Sprintf("order:%v:responses", orderID)

	responsesIds, err := red.client.LRange(ctx, responseKey, 0, -1).Result()

	if err != nil {
		return nil, err
	}

	var responses []models.Response

	for _, responseID := range responsesIds {
		response, err := red.GetShortResponseById(ctx, responseID)
		if err == nil {
			responses = append(responses, *response)
		}
	}

	return responses, nil
}

func (red *Redis) GetResponseById(ctx context.Context, responseID string) (*models.ResponseDB, error) {
	responseKey := fmt.Sprintf("response:%v", responseID)

	responseRaw, err := red.client.HGetAll(ctx, responseKey).Result()

	if err != nil || len(responseRaw) == 0 {
		return nil, err
	}

	createdAt, _ := time.Parse(time.RFC3339, responseRaw["created_at"])

	response := &models.ResponseDB{
		ID:            responseRaw["id"],
		OrderID:       responseRaw["order_id"],
		TutorID:       responseRaw["tutor_id"],
		TutorUsername: responseRaw["tutor_username"],
		Name:          responseRaw["name"],
		CreatedAt:     createdAt,
	}

	return response, nil
}

func (red *Redis) GetShortResponseById(ctx context.Context, responseID string) (*models.Response, error) {
	responseKey := fmt.Sprintf("response:%v", responseID)

	responseRaw, err := red.client.HGetAll(ctx, responseKey).Result()
	if err != nil || len(responseRaw) == 0 {
		return nil, err
	}

	createdAt, _ := time.Parse(time.RFC3339, responseRaw["created_at"])

	response := &models.Response{
		ID:        responseRaw["id"],
		TutorID:   responseRaw["tutor_id"],
		Name:      responseRaw["name"],
		CreatedAt: createdAt,
	}

	return response, nil
}
