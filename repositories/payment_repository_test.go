// +build integration

package repositories

import (
	"context"
	"github.com/go-pg/pg/v9"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/tsovak/rest-api-demo/config"
	"github.com/tsovak/rest-api-demo/service/db"
	"github.com/tsovak/rest-api-demo/testutils"
	"strconv"
	"testing"
)

type paymentRepoSuiteTest struct {
	suite.Suite
	db                *pg.DB
	paymentRepository PaymentRepository
}

func (p *paymentRepoSuiteTest) SetupSuite() {
	config, err := config.GetTestConfig()
	require.Nil(p.T(), err, "Config is nil")

	setup, err := testutils.SetupTestDB(db.GetPgConnectionOptions(config), "../scripts/migrations/")
	require.NoError(p.T(), err)

	client := db.NewPostgresClient(config.Logger, setup.Db)
	p.db = client.GetConnection()
	p.paymentRepository = NewPaymentRepository(p.db)
}

func (p *paymentRepoSuiteTest) TearDownSuite() {
	p.db.Close()
}

func (p *paymentRepoSuiteTest) Test_paymentRepository_FindById() {
	payment := testutils.GetTestPayment()
	err := p.paymentRepository.Save(context.Background(), payment)
	require.NoError(p.T(), err, "Cannot save payment")

	dbPayment, err := p.paymentRepository.FindById(context.Background(), strconv.FormatInt(payment.ID, 10))
	require.NoError(p.T(), err, "Cannot find user")
	require.Equal(p.T(), payment, &dbPayment)
}

func (p *paymentRepoSuiteTest) Test_paymentRepository_GetAll() {
	t := p.T()
	payments, err := p.paymentRepository.GetAll(context.Background())
	require.NoError(t, err)
	require.NotNil(t, payments)
}

func (p *paymentRepoSuiteTest) Test_paymentRepository_Save() {
	payment := testutils.GetTestPayment()
	err := p.paymentRepository.Save(context.Background(), payment)
	require.NoError(p.T(), err, "Cannot save payment")
}

func TestAllPaymentsTests(t *testing.T) {
	suite.Run(t, new(paymentRepoSuiteTest))
}
