package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/diofanto33/ms-go/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	randAccount := createRandomAccount(t)
	getAccount, err := testQueries.GetAccount(context.Background(), randAccount.ID)
	require.NoError(t, err)
	require.NotEmpty(t, getAccount)

	require.Equal(t, randAccount.ID, getAccount.ID)
	require.Equal(t, randAccount.Owner, getAccount.Owner)
	require.Equal(t, randAccount.Balance, getAccount.Balance)
	require.Equal(t, randAccount.Currency, getAccount.Currency)
	require.WithinDuration(t, randAccount.CreatedAt, getAccount.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	randAccount := createRandomAccount(t) // create random account
	/* we need to update the balance of the random account */
	arg := UpdateAccountParams{
		ID:      randAccount.ID,
		Balance: util.RandomMoney(), // update balance to random value
	}

	updateAccount, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updateAccount)

	require.Equal(t, randAccount.ID, updateAccount.ID)
	require.Equal(t, randAccount.Owner, updateAccount.Owner)
	require.Equal(t, arg.Balance, updateAccount.Balance)
	require.Equal(t, randAccount.Currency, updateAccount.Currency)
	require.WithinDuration(t, randAccount.CreatedAt, updateAccount.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	randAccount := createRandomAccount(t) // create random account
	err := testQueries.DeleteAccount(context.Background(), randAccount.ID)
	require.NoError(t, err)

	getAccount, err := testQueries.GetAccount(context.Background(), randAccount.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, getAccount)
}

func TestListAccounts(t *testing.T) {
	/* create 5 random accounts */
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}
	/* we need to limit the number of accounts to 5 and offset to 5 */
	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}
	/* get the list of accounts */
	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	/* check if there is no error and the length of the accounts is 5 */
	require.NoError(t, err)
	require.Len(t, accounts, 5)
	/* check if the accounts are not empty */
	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
