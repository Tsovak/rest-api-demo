package api

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/tsovak/rest-api-demo/api/model"
	"github.com/tsovak/rest-api-demo/service"
	"net/http"
)

// need to hide server errors.
// internalErrorMessage returns commons error message when something went wrong
var internalErrorMessage = NewErrorMessage("internal server error")

// ApiServer struct saves needed managers
type ApiServer struct {
	accountManager service.AccountManager
	paymentManager service.PaymentManager
	logger         *logrus.Logger
}

// NewApiServer returns an ApiServer structure
func NewApiServer(accountManager service.AccountManager, paymentManager service.PaymentManager, logger *logrus.Logger) *ApiServer {
	return &ApiServer{
		accountManager: accountManager,
		paymentManager: paymentManager,
		logger:         logger,
	}
}

// CreateAccount represents an account creation handler
func (s *ApiServer) CreateAccount(ctx echo.Context) error {
	context := ctx.Request().Context()
	var accountRequest = new(model.AccountRequest)

	body := ctx.Request().Body
	if body == nil {
		s.logger.Warn("Body is nil")
		return ctx.JSON(http.StatusBadRequest, NewErrorMessage("Body is nil"))
	}

	decoder := json.NewDecoder(body)
	err := decoder.Decode(&accountRequest)
	if err != nil {
		s.logger.WithContext(context).Error(err)
		return ctx.JSON(http.StatusBadRequest, NewErrorMessage("cannot decode"))
	}

	// create an account model with zero ID. Id will be stored by DB
	account := model.Account{
		ID:       0,
		Name:     accountRequest.Name,
		Currency: accountRequest.Currency,
		Balance:  accountRequest.Balance,
	}
	err = s.accountManager.Save(context, &account)
	if err != nil {
		s.logger.WithContext(context).Error(err)
		return ctx.JSON(http.StatusInternalServerError, internalErrorMessage)
	}

	return ctx.JSON(http.StatusOK, account)
}

// GetAllAccounts returns all created accounts
func (s *ApiServer) GetAllAccounts(ctx echo.Context) error {
	context := ctx.Request().Context()
	accounts, err := s.accountManager.GetAllAccounts(context)
	if err != nil {
		s.logger.WithContext(context).Error(err)
		return ctx.JSON(http.StatusInternalServerError, internalErrorMessage)
	}

	return ctx.JSON(http.StatusOK, accounts)
}

// GetAccountPayments returns an account payments
func (s *ApiServer) GetAccountPayments(ctx echo.Context) error {
	context := ctx.Request().Context()
	accountId := ctx.Param("id")
	payments, err := s.paymentManager.GetPaymentsByAccountId(context, accountId)
	if err != nil {
		s.logger.WithContext(context).Error(err)
		return ctx.JSON(http.StatusInternalServerError, NewErrorMessage("cannot get an account payments"))
	}

	// preparing the response
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

// CreatePayment represents an account payment creation
func (s *ApiServer) CreatePayment(ctx echo.Context) error {
	context := ctx.Request().Context()

	// here we will decode request
	var localPayment = new(model.PaymentRequest)

	body := ctx.Request().Body
	if body == nil {
		s.logger.Warn("Body is nil")
		return ctx.JSON(http.StatusBadRequest, NewErrorMessage("body is nil"))
	}

	decoder := json.NewDecoder(body)
	err := decoder.Decode(&localPayment)
	if err != nil {
		s.logger.WithContext(context).Error(err)
		return ctx.JSON(http.StatusBadRequest, NewErrorMessage("cannot decode"))
	}

	fromAccount, err := s.accountManager.FindById(context, localPayment.FromAccountID)
	if err != nil {
		s.logger.WithContext(context).Error(err)
		return ctx.JSON(http.StatusInternalServerError, internalErrorMessage)
	}
	// we cannot debit an account which does not exist
	if fromAccount == (model.Account{}) {
		return ctx.JSON(http.StatusNotFound,
			NewErrorMessage(fmt.Sprintf("account id=%v does not exist", localPayment.FromAccountID)))
	}

	if fromAccount.Balance < localPayment.Amount {
		// we cannot credit if you have not enough money
		return ctx.JSON(http.StatusBadRequest,
			NewErrorMessage(fmt.Sprintf("account id=%v has not enough money", localPayment.FromAccountID)))
	}

	// we cannot credit an account which does not exist
	toAccount, err := s.accountManager.FindById(context, localPayment.ToAccountID)
	if err != nil {
		s.logger.WithContext(context).Error(err)
		return ctx.JSON(http.StatusInternalServerError, internalErrorMessage)
	}
	if toAccount == (model.Account{}) {
		return ctx.JSON(http.StatusNotFound,
			NewErrorMessage(fmt.Sprintf("account id=%v does not exist", localPayment.ToAccountID)))
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
	// get the transaction function for execute with account update inside one transaction
	doPaymentTransaction := s.paymentManager.GetSaveTransaction(context, payments...)

	err = s.accountManager.Update(context, &toAccount, doPaymentTransaction)
	if err != nil {
		s.logger.WithContext(context).Error(err)
		return ctx.JSON(http.StatusInternalServerError, internalErrorMessage)
	}

	return ctx.JSON(http.StatusOK, payments)
}
