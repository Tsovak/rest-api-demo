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
