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

type accountRepoSuiteTest struct {
	suite.Suite
	db                *pg.DB
	accountRepository AccountRepository
}

func (a *accountRepoSuiteTest) SetupSuite() {
	config, err := config.GetTestConfig()
	require.Nil(a.T(), err, "Config is nil")

	setup, err := testutils.SetupTestDB(db.GetPgConnectionOptions(config), "../scripts/migrations/")
	require.NoError(a.T(), err)

	client := db.NewPostgresClient(setup.Db)
	a.db = client.GetConnection()
	a.accountRepository = NewAccountRepository(a.db)
}

func (a *accountRepoSuiteTest) TearDownSuite() {
	a.db.Close()
}

func (a *accountRepoSuiteTest) TestAccountRepository_GetAll() {
	t := a.T()
	accounts, err := a.accountRepository.GetAll(context.Background())
	require.NoError(t, err)
	require.NotNil(t, accounts)
}

func (a *accountRepoSuiteTest) TestAccountRepository_Save() {
	user := testutils.GetTestUser()
	err := a.accountRepository.Save(context.Background(), user)
	require.NoError(a.T(), err, "Cannot save user")
}

func (a *accountRepoSuiteTest) TestAccountRepository_FindById() {
	user := testutils.GetTestUser()
	err := a.accountRepository.Save(context.Background(), user)
	require.NoError(a.T(), err, "Cannot save user")
	dbUser, err := a.accountRepository.FindById(context.Background(), strconv.FormatInt(user.ID, 10))
	require.NoError(a.T(), err, "Cannot find user")
	require.Equal(a.T(), user, &dbUser)
}

func (a *accountRepoSuiteTest) TestAccountRepository_DeleteById() {
	user := testutils.GetTestUser()
	userIdString := strconv.FormatInt(user.ID, 10)
	err := a.accountRepository.Save(context.Background(), user)
	require.NoError(a.T(), err, "Cannot save user")

	err = a.accountRepository.DeleteById(context.Background(), userIdString)
	require.NoError(a.T(), err, "Cannot delete user")

	dbUser, err := a.accountRepository.FindById(context.Background(), userIdString)
	require.NoError(a.T(), err, "Cannot find user")
	require.Empty(a.T(), dbUser, "User not deleted")
}

func (a *accountRepoSuiteTest) TestAccountRepository_UpdateBalance() {
	user := testutils.GetTestUser()
	userIdString := strconv.FormatInt(user.ID, 10)
	amount := int64(100)

	err := a.accountRepository.Save(context.Background(), user)
	require.NoError(a.T(), err, "Cannot save user")

	user.Balance += amount
	err = a.accountRepository.Update(context.Background(), user)
	require.NoError(a.T(), err, "Cannot update balance")

	dbUser, err := a.accountRepository.FindById(context.Background(), userIdString)
	require.NoError(a.T(), err, "Cannot find user")

	require.Equal(a.T(), user.Balance, dbUser.Balance, "User balance is different")
}

func TestAllAccountTests(t *testing.T) {
	suite.Run(t, new(accountRepoSuiteTest))
}
