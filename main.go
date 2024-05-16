package main

import (
	"fmt"
	"net/http"
	"order-service/domain/order"
	"order-service/domain/order/repository"
	"order-service/grpc/client"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var validate *validator.Validate

func main() {
	// Inisialisasi database
	dsn := "root:@tcp(127.0.0.1:3306)/order-service?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Auto Migrate
	if err := db.AutoMigrate(&order.Order{}); err != nil {
		panic("failed to auto migrate database: " + err.Error())
	}

	// Inisialisasi repository
	orderRepository := repository.NewOrderRepository(db)

	// Inisialisasi server Gin
	r := gin.Default()

	validate = validator.New()

	r.POST("/orders", createOrder(orderRepository))

	if err := r.Run(":8081"); err != nil {
		panic("failed to start server: " + err.Error())
	}
}

func createOrder(orderRepository *repository.OrderRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newOrder order.Order
		if err := c.ShouldBindJSON(&newOrder); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if err := validate.Struct(newOrder); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Call the gRPC client to check if the user exists
		userClient, err := client.NewUserServiceClient("localhost:50051")
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to connect GRPC"})
			return
		}

		_, err = userClient.GetUserByID(uint64(newOrder.UserID))
		if err != nil {

			if st, ok := status.FromError(err); ok && st.Code() == codes.Unavailable {
				c.JSON(500, gin.H{"error": "gRPC service is unavailable"})
				return
			}

			if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
				c.JSON(404, gin.H{"error": "requested resource not found"})
				return
			}

			c.JSON(500, gin.H{"error": fmt.Sprintf("failed to call RPC function: %v", err)})
			return
		}

		if err := orderRepository.Create(&newOrder); err != nil {
			c.JSON(500, gin.H{"error": "failed to create order"})
			return
		}

		c.JSON(http.StatusCreated, newOrder)
	}
}
