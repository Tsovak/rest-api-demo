package api

import (
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/tsovak/rest-api-demo/api/model"
	"github.com/tsovak/rest-api-demo/service"
	"net/http"
)

type ApiServer struct {
	accountManager service.AccountManager
	paymentManager service.PaymentManager
	logger         *logrus.Logger
}

func NewApiServer(accountManager service.AccountManager, paymentManager service.PaymentManager, logger *logrus.Logger) *ApiServer {
	return &ApiServer{
		accountManager: accountManager,
		paymentManager: paymentManager,
		logger:         logger,
	}
}

func (s *ApiServer) CreateAccount(ctx echo.Context) error {
	context := ctx.Request().Context()
	var account = new(model.Account)

	body := ctx.Request().Body
	if body == nil {
		s.logger.Warn("Body is nil")
		return ctx.JSON(http.StatusBadRequest, errors.New("Body is nil").Error())
	}

	decoder := json.NewDecoder(body)
	err := decoder.Decode(&account)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errors.New("cannot decode").Error())
	}

	err = s.accountManager.Save(context, account)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, account)
}

func (s *ApiServer) GetAllAccounts(ctx echo.Context) error {
	context := ctx.Request().Context()
	accounts, err := s.accountManager.GetAllAccounts(context)
	if err != nil {
		s.logger.WithContext(context).Error(err)
		return ctx.JSON(http.StatusInternalServerError, struct{}{})
	}

	return ctx.JSON(http.StatusOK, accounts)
}
