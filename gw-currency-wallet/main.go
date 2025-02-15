package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/peesqq/gw-currency-wallet/internal/config"
	"github.com/peesqq/gw-currency-wallet/internal/db"
	"github.com/peesqq/gw-currency-wallet/internal/grpcclient"
	"github.com/peesqq/gw-currency-wallet/internal/storages"

	"github.com/peesqq/proto-exchange/proto"
)

func main() {
	// Инициализация конфига
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Подключение к базе данных
	conn, err := db.ConnectDB(cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBHost, cfg.DBPort)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer conn.Close()

	// Инициализация базы данных
	db.InitDB(conn)

	// Инициализация хранилищ
	userStorage := storages.NewUserStorage(conn)
	walletStorage := storages.NewWalletStorage(conn)

	// Подключение к gRPC-серверу
	grpcClient := grpcclient.NewGRPCClient("localhost:50051")

	// Создание роутера
	r := gin.Default()

	// Обработчики API
	r.POST("/api/v1/wallet/exchange", func(c *gin.Context) {
		var req struct {
			FromCurrency string  `json:"from_currency" binding:"required"`
			ToCurrency   string  `json:"to_currency" binding:"required"`
			Amount       float64 `json:"amount" binding:"required"`
			UserID       int     `json:"user_id" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// Получаем курс обмена через gRPC
		grpcResponse, err := grpcClient.ExchangeService.ConvertCurrency(c, &proto.ConvertCurrencyRequest{
			FromCurrency: req.FromCurrency,
			ToCurrency:   req.ToCurrency,
			Amount:       float32(req.Amount),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert currency", "details": err.Error()})
			return
		}

		convertedAmount := grpcResponse.ConvertedAmount

		err = walletStorage.Withdraw(context.Background(), req.UserID, req.FromCurrency, req.Amount)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to withdraw amount"})
			return
		}

		err = walletStorage.Deposit(context.Background(), req.UserID, req.ToCurrency, float64(convertedAmount))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to deposit converted amount"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":          "Exchange successful",
			"converted_amount": convertedAmount,
		})
	})

	r.POST("/api/v1/register", func(c *gin.Context) {
		var req struct {
			Username string `json:"username" binding:"required"`
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		err := userStorage.CreateUser(context.Background(), req.Username, req.Email, req.Password)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to register user"})
			return
		}

		var userID int
		query := `SELECT id FROM users WHERE username = $1`
		err = conn.QueryRow(context.Background(), query, req.Username).Scan(&userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user ID"})
			return
		}

		err = walletStorage.CreateWallet(context.Background(), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create wallet"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
	})

	r.POST("/api/v1/login", func(c *gin.Context) {
		var req struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		authenticated, err := userStorage.AuthenticateUser(context.Background(), req.Username, req.Password)
		if err != nil || !authenticated {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}

		token := cfg.JWTSecret // Временный токен, заменить на реальный JWT
		c.JSON(http.StatusOK, gin.H{"token": token})
	})

	r.GET("/api/v1/balance", func(c *gin.Context) {
		userID := 1 // Заменить на ID из токена пользователя
		balance, err := walletStorage.GetBalance(context.Background(), userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"balance": balance})
	})

	r.POST("/api/v1/wallet/deposit", func(c *gin.Context) {
		userID := 1 // Заменить на ID из токена пользователя
		var req struct {
			Currency string  `json:"currency" binding:"required"`
			Amount   float64 `json:"amount" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		err := walletStorage.Deposit(context.Background(), userID, req.Currency, req.Amount)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Deposit successful"})
	})

	r.POST("/api/v1/wallet/withdraw", func(c *gin.Context) {
		userID := 1 // Заменить на ID из токена пользователя
		var req struct {
			Currency string  `json:"currency" binding:"required"`
			Amount   float64 `json:"amount" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		err := walletStorage.Withdraw(context.Background(), userID, req.Currency, req.Amount)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Withdrawal successful"})
	})

	log.Println("Currency wallet service running on port 8080")
	r.Run(":8080")
}
