package service

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/tsovak/rest-api-demo/api/model"
	"github.com/tsovak/rest-api-demo/mock"
	"github.com/tsovak/rest-api-demo/testutils"
	"reflect"
	"testing"
)

var payments = func() []model.Payment {
	payment1 := *testutils.GetTestPayment()
	payment2 := *testutils.GetTestPayment()

	return []model.Payment{payment1, payment2}
}

func TestGetAllPaymentsOk(t *testing.T) {
	mc := gomock.NewController(t)
	ctx := context.Background()
	defer mc.Finish()

	payments := payments()
	mockPaymentRepository := mock.NewMockPaymentRepository(mc)
	mockPaymentRepository.
		EXPECT().
		GetAll(ctx).
		AnyTimes().
		Return(payments, nil)

	manager := NewPaymentManager(mockPaymentRepository)
	receivedPayments, err := manager.GetAccountPayments(ctx)

	require.Nil(t, err)
	require.NotNil(t, receivedPayments)
	require.Len(t, receivedPayments, len(payments))
	require.True(t, reflect.DeepEqual(payments, receivedPayments))
}

func TestGetAllPaymentsFail(t *testing.T) {
	mc := gomock.NewController(t)
	ctx := context.Background()
	defer mc.Finish()

	mockPaymentRepository := mock.NewMockPaymentRepository(mc)
	mockPaymentRepository.
		EXPECT().
		GetAll(ctx).
		AnyTimes().
		Return(nil, errors.New("cannot get payments"))

	manager := NewPaymentManager(mockPaymentRepository)
	receivedPayments, err := manager.GetAccountPayments(ctx)

	require.Error(t, err)
	require.Nil(t, receivedPayments)
}

func TestCreatePaymentFail(t *testing.T) {
	mc := gomock.NewController(t)
	ctx := context.Background()
	defer mc.Finish()

	mockPaymentRepository := mock.NewMockPaymentRepository(mc)
	mockPaymentRepository.
		EXPECT().
		Save(ctx, nil).
		AnyTimes().
		Return(errors.New("cannot create payment"))

	manager := NewPaymentManager(mockPaymentRepository)
	err := manager.CreatePayments(ctx, nil)

	require.Error(t, err)
}

func TestCreatePaymentOk(t *testing.T) {
	mc := gomock.NewController(t)
	ctx := context.Background()
	defer mc.Finish()

	payment := testutils.GetTestPayment()
	mockPaymentRepository := mock.NewMockPaymentRepository(mc)
	mockPaymentRepository.
		EXPECT().
		Save(ctx, payment).
		AnyTimes().
		Return(nil)

	manager := NewPaymentManager(mockPaymentRepository)
	err := manager.CreatePayments(ctx, payment)

	require.NoError(t, err)
}

func TestGetPaymentsByAccountIdOk(t *testing.T) {
	mc := gomock.NewController(t)
	ctx := context.Background()
	defer mc.Finish()

	payment := testutils.GetTestPayment()
	payments := []model.Payment{*payment}

	mockPaymentRepository := mock.NewMockPaymentRepository(mc)
	mockPaymentRepository.
		EXPECT().
		GetPaymentsByAccountId(ctx, payment.FromAccountID).
		AnyTimes().
		Return(payments, nil)

	manager := NewPaymentManager(mockPaymentRepository)
	receivedPayments, err := manager.GetPaymentsByAccountId(ctx, payment.FromAccountID)

	require.NoError(t, err)
	require.NotNil(t, receivedPayments)
	require.Len(t, receivedPayments, len(payments))
	require.True(t, reflect.DeepEqual(payments, receivedPayments))
}

func TestGetPaymentsByAccountIdFail(t *testing.T) {
	mc := gomock.NewController(t)
	ctx := context.Background()
	defer mc.Finish()

	payment := testutils.GetTestPayment()

	mockPaymentRepository := mock.NewMockPaymentRepository(mc)
	mockPaymentRepository.
		EXPECT().
		GetPaymentsByAccountId(ctx, payment.FromAccountID).
		AnyTimes().
		Return(nil, errors.New("cannot find the payment"))

	manager := NewPaymentManager(mockPaymentRepository)
	receivedPayments, err := manager.GetPaymentsByAccountId(ctx, payment.FromAccountID)

	require.Error(t, err)
	require.Nil(t, receivedPayments)
}
