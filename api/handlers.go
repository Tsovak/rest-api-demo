package api

import (
	"encoding/json"
	"errors"
	"fmt"
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
	var accountRequest = new(model.AccountRequest)

	body := ctx.Request().Body
	if body == nil {
		s.logger.Warn("Body is nil")
		return ctx.JSON(http.StatusBadRequest, errors.New("Body is nil").Error())
	}

	decoder := json.NewDecoder(body)
	err := decoder.Decode(&accountRequest)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errors.New("cannot decode").Error())
	}

	account := model.Account{
		ID:       0,
		Name:     accountRequest.Name,
		Currency: accountRequest.Currency,
		Balance:  accountRequest.Balance,
	}
	err = s.accountManager.Save(context, &account)
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

func (s *ApiServer) GetAccountPayments(ctx echo.Context) error {
	context := ctx.Request().Context()
	accountId := ctx.Param("id")
	payments, err := s.paymentManager.GetPaymentsByAccountId(context, accountId)
	if err != nil {
		s.logger.WithContext(context).Error(err)
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	var localPaymentResponse = make([]model.PaymentResponse, len(payments))
	for i, p := range payments {
		pr := model.PaymentResponse{
			ID:            p.ID,
			Amount:        p.Amount,
			ToAccountID:   p.ToAccountID,
			FromAccountID: p.FromAccountID,
			Direction:     p.Direction,
		}
		localPaymentResponse[i] = pr
	}

	return ctx.JSON(http.StatusOK, localPaymentResponse)
}

func (s *ApiServer) CreatePayment(ctx echo.Context) error {
	context := ctx.Request().Context()

	var localPayment = new(model.PaymentRequest)

	body := ctx.Request().Body
	if body == nil {
		s.logger.Warn("Body is nil")
		return ctx.JSON(http.StatusBadRequest, errors.New("body is nil").Error())
	}

	decoder := json.NewDecoder(body)
	err := decoder.Decode(&localPayment)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errors.New("cannot decode").Error())
	}

	fromAccount, err := s.accountManager.FindById(context, localPayment.FromAccountID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	if fromAccount == (model.Account{}) {
		return ctx.JSON(http.StatusNotFound,
			fmt.Errorf(fmt.Sprintf("account id=%v does not exist", localPayment.FromAccountID)).Error())
	}

	if fromAccount.Balance < localPayment.Amount {
		// we cannot credit if you have not enough money
		return ctx.JSON(http.StatusBadRequest,
			fmt.Errorf(fmt.Sprintf("account id=%v has not enough money", localPayment.FromAccountID)).Error())
	}

	toAccount, err := s.accountManager.FindById(context, localPayment.ToAccountID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	if toAccount == (model.Account{}) {
		return ctx.JSON(http.StatusNotFound,
			fmt.Errorf(fmt.Sprintf("account id=%v does not exist", localPayment.ToAccountID)).Error())
	}

	// now we need to change the balances
	fromAccount.Balance -= localPayment.Amount
	toAccount.Balance += localPayment.Amount

	// preparing the payments to save the DB
	var paymentFrom = model.Payment{
		ID:            0,
		Amount:        -localPayment.Amount,
		ToAccountID:   localPayment.ToAccountID,
		FromAccountID: localPayment.FromAccountID,
		Direction:     model.Outgoing,
	}

	var paymentTo = model.Payment{
		ID:            0,
		Amount:        localPayment.Amount,
		ToAccountID:   localPayment.ToAccountID,
		FromAccountID: localPayment.FromAccountID,
		Direction:     model.Incoming,
	}

	var payments = []*model.Payment{&paymentFrom, &paymentTo}
	doPaymentTransaction := s.paymentManager.GetSaveTransaction(context, payments...)

	err = s.accountManager.Update(context, &toAccount, doPaymentTransaction)
	if err != nil {
		s.logger.WithContext(context).Error(err)
		return ctx.JSON(http.StatusInternalServerError, struct{}{})
	}

	return ctx.JSON(http.StatusOK, payments)
}
