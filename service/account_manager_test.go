package service

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/tsovak/rest-api-demo/api/model"
	"github.com/tsovak/rest-api-demo/mock"
	"github.com/tsovak/rest-api-demo/testutils"
	"reflect"
	"testing"
)

var account = func() *model.Account {
	return testutils.GetTestUser()
}

var accounts = func() []model.Account {
	account1 := *testutils.GetTestUser()
	account1.Name = "one"
	account1.Balance = 100

	account2 := *testutils.GetTestUser()
	account2.Name = "two"
	account2.Balance = 200
	return []model.Account{account1, account2}
}

func TestGetAllAccountsOk(t *testing.T) {
	mc := gomock.NewController(t)
	ctx := context.Background()
	defer mc.Finish()

	users := accounts()
	mockAccountRepository := mock.NewMockAccountRepository(mc)
	mockAccountRepository.
		EXPECT().
		GetAll(ctx).
		AnyTimes().
		Return(users, nil)

	manager := NewAccountManager(mockAccountRepository)
	receivedUsers, err := manager.GetAllAccounts(ctx)

	require.Nil(t, err)
	require.NotNil(t, receivedUsers)
	require.Len(t, receivedUsers, len(users))
	require.True(t, reflect.DeepEqual(users, receivedUsers))
}

func TestGetAllAccountsFail(t *testing.T) {
	mc := gomock.NewController(t)
	ctx := context.Background()
	defer mc.Finish()

	mockAccountRepository := mock.NewMockAccountRepository(mc)
	mockAccountRepository.
		EXPECT().
		GetAll(ctx).
		AnyTimes().
		Return(nil, errors.New("fail"))

	manager := NewAccountManager(mockAccountRepository)
	allAccounts, err := manager.GetAllAccounts(ctx)
	require.Error(t, err)
	require.Empty(t, allAccounts)
}

func TestSaveAccountsFail(t *testing.T) {
	mc := gomock.NewController(t)
	ctx := context.Background()
	defer mc.Finish()

	mockAccountRepository := mock.NewMockAccountRepository(mc)
	user := account()
	mockAccountRepository.
		EXPECT().
		Save(ctx, user).
		AnyTimes().
		Return(errors.New("could not save"))

	manager := NewAccountManager(mockAccountRepository)
	err := manager.Save(ctx, user)
	require.Error(t, err)
}

func TestSaveAccountsOk(t *testing.T) {
	mc := gomock.NewController(t)
	ctx := context.Background()
	defer mc.Finish()

	mockAccountRepository := mock.NewMockAccountRepository(mc)
	user := account()
	mockAccountRepository.
		EXPECT().
		Save(ctx, user).
		AnyTimes().
		Return(nil)

	manager := NewAccountManager(mockAccountRepository)
	err := manager.Save(ctx, user)
	require.NoError(t, err)
}

func TestFindByIdAccountsOk(t *testing.T) {
	mc := gomock.NewController(t)
	ctx := context.Background()
	defer mc.Finish()

	mockAccountRepository := mock.NewMockAccountRepository(mc)
	user := *account()
	mockAccountRepository.
		EXPECT().
		FindById(ctx, string(user.ID)).
		AnyTimes().
		Return(user, nil)

	manager := NewAccountManager(mockAccountRepository)
	returnedAccount, err := manager.FindByID(ctx, string(user.ID))
	require.NoError(t, err)
	require.Equal(t, user, returnedAccount)
}

func TestFindByIdAccountsFail(t *testing.T) {
	mc := gomock.NewController(t)
	ctx := context.Background()
	defer mc.Finish()

	mockAccountRepository := mock.NewMockAccountRepository(mc)
	user := *account()
	mockAccountRepository.
		EXPECT().
		FindById(ctx, string(user.ID)).
		AnyTimes().
		Return(model.Account{}, errors.New("cannot find account"))

	manager := NewAccountManager(mockAccountRepository)
	returnedAccount, err := manager.FindByID(ctx, string(user.ID))
	require.Error(t, err)
	require.Equal(t, model.Account{}, returnedAccount)
}

func TestDeleteByIdAccountsFail(t *testing.T) {
	mc := gomock.NewController(t)
	ctx := context.Background()
	defer mc.Finish()

	mockAccountRepository := mock.NewMockAccountRepository(mc)
	user := *account()
	mockAccountRepository.
		EXPECT().
		DeleteById(ctx, string(user.ID)).
		AnyTimes().
		Return(errors.New("cannot delete account"))

	manager := NewAccountManager(mockAccountRepository)
	err := manager.DeleteByID(ctx, string(user.ID))
	require.Error(t, err)
}

func TestDeleteByIdAccountsOk(t *testing.T) {
	mc := gomock.NewController(t)
	ctx := context.Background()
	defer mc.Finish()

	mockAccountRepository := mock.NewMockAccountRepository(mc)
	user := *account()
	mockAccountRepository.
		EXPECT().
		DeleteById(ctx, string(user.ID)).
		AnyTimes().
		Return(nil)

	manager := NewAccountManager(mockAccountRepository)
	err := manager.DeleteByID(ctx, string(user.ID))
	require.NoError(t, err)
}

func TestUpdateAccountsOk(t *testing.T) {
	mc := gomock.NewController(t)
	ctx := context.Background()
	defer mc.Finish()

	mockAccountRepository := mock.NewMockAccountRepository(mc)
	user := account()
	mockAccountRepository.
		EXPECT().
		Update(ctx, user, nil).
		AnyTimes().
		Return(nil)

	manager := NewAccountManager(mockAccountRepository)
	err := manager.Update(ctx, user, nil)
	require.NoError(t, err)
}

func TestUpdateAccountsFail(t *testing.T) {
	mc := gomock.NewController(t)
	ctx := context.Background()
	defer mc.Finish()

	mockAccountRepository := mock.NewMockAccountRepository(mc)
	user := account()
	mockAccountRepository.
		EXPECT().
		Update(ctx, user, nil).
		AnyTimes().
		Return(errors.New("cannot update account"))

	manager := NewAccountManager(mockAccountRepository)
	err := manager.Update(ctx, user, nil)
	require.Error(t, err)
}
