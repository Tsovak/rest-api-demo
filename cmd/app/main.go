package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tsovak/rest-api-demo/api"
	"github.com/tsovak/rest-api-demo/config"
	"github.com/tsovak/rest-api-demo/repositories"
	"github.com/tsovak/rest-api-demo/service"
	"github.com/tsovak/rest-api-demo/service/db"
	"os"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		println(err.Error())
		os.Exit(-1)
	}
	cfg.Logger.Debugf("DB address is %v", cfg.DbConfig.Address)
	pgClient := db.NewPostgresClientFromConfig(cfg)
	connection := pgClient.GetConnection()
	defer connection.Close()

	e := echo.New()
	e.Use(middleware.Logger())

	accountRepository := repositories.NewAccountRepository(connection)
	paymentRepository := repositories.NewPaymentRepository(connection)

	accountManager := service.NewAccountManager(accountRepository)
	paymentManager := service.NewPaymentManager(paymentRepository)
	apiServer := api.NewApiServer(accountManager, paymentManager, cfg.Logger)

	e.Router().Add("GET", "/accounts", apiServer.GetAllAccounts)
	e.Router().Add("POST", "/accounts", apiServer.CreateAccount)
	e.Router().Add("GET", "/accounts/:id/payments", apiServer.GetAccountPayments)
	e.Router().Add("POST", "/payments", apiServer.CreatePayment)
	err = e.Start(fmt.Sprintf(":%v", cfg.ServerPort))
	if err != nil {
		cfg.Logger.Error("Cannot start the server")
		os.Exit(-1)
	}
}
