package db

import (
	"context"
	"fmt"
	"testing"

	_ "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDb)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	fmt.Println(">> Before:", account1.Balance, account2.Balance)

	n := 2
	amount := int64(10)

	errorChannel := make(chan error)
	resultChannel := make(chan TransferTxResult)

	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		go func() {
			res, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errorChannel <- err
			resultChannel <- res
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errorChannel
		require.NoError(t, err)

		res := <-resultChannel
		require.NotEmpty(t, res)

		// Test transfer result
		transfer := res.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, transfer.FromAccountID, account1.ID)
		require.Equal(t, transfer.ToAccountID, account2.ID)
		require.Equal(t, transfer.Amount, amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// Test entry result
		fromEntry := res.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := res.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// Check account
		fromAccount := res.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := res.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// Check balances
		fmt.Println(">> Tx:", fromAccount.Balance, toAccount.Balance)

		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) // 1st Tx: 1 * amount; 2nd Tx: 2 * amount; ....... nth Tx: n * amount

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check the final updated balance
	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> After:", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDb)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	fmt.Println(">> Before:", account1.Balance, account2.Balance)

	n := 10
	amount := int64(10)

	errorChannel := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountId := account1.ID
		toAccountId := account2.ID

		if i%2 == 1 {
			fromAccountId = account2.ID
			toAccountId = account1.ID
		}

		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountId,
				ToAccountID:   toAccountId,
				Amount:        amount,
			})
			errorChannel <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errorChannel
		require.NoError(t, err)
	}

	// check the final updated balance
	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> After:", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}
