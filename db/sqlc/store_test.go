package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestTransferTx tests the TransferTx function, which handles transferring money between accounts.
// It performs multiple concurrent transfer transactions and checks the results for consistency.
func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB) // Initialize a new store with the test database connection.

	// Create two random accounts for testing the transfer between them.
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before:", account1.Balance, account2.Balance)
	n := 10             // Number of concurrent transactions to run.
	amount := int64(10) // Amount of money to transfer in each transaction.

	// Channels to receive errors and results from goroutines.
	errs := make(chan error)
	// results := make(chan TransferTxResult)

	// Start n concurrent transfer transactions.
	for i := 0; i < n; i++ {
		// txName := fmt.Sprintf("tx %d", i+1)
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		go func() {
			// Execute the transfer transaction from account1 to account2.
			// ctx := context.WithValue(context.Background(), txkey, txName)

			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountId: fromAccountID, // Source account ID.
				ToAccountID:   toAccountID,   // Destination account ID.
				Amount:        amount,        // Amount to transfer.
			})
			errs <- err // Send any error encountered to the errs channel.
			// results <- result // Send the result to the results channel.
		}()
	}

	// existed := make(map[int]bool)
	// Process the results and errors for each transaction.
	for i := 0; i < n; i++ {
		err := <-errs           // Receive the error from the errs channel.
		require.NoError(t, err) // Ensure there was no error in the transaction.

		/*
			result := <-results         // Receive the result from the results channel.
			require.NotEmpty(t, result) // Ensure the result is not empty.

			// Check the transfer details.
			transfer := result.Transfer
			require.NotEmpty(t, transfer)                         // Ensure the transfer record is not empty.
			require.Equal(t, account1.ID, transfer.FromAccountID) // Verify the source account ID.
			require.Equal(t, account2.ID, transfer.ToAccountID)   // Verify the destination account ID.
			require.Equal(t, amount, transfer.Amount)             // Verify the transfer amount.
			require.NotZero(t, transfer.ID)                       // Ensure the transfer has a valid ID.
			require.NotZero(t, transfer.CreatedAt)                // Ensure the transfer has a creation timestamp.

			// Verify that the transfer record can be retrieved by its ID.
			_, err = store.GetTranfer(context.Background(), transfer.ID)
			require.NoError(t, err)

			// Check the fromEntry details.
			fromEntry := result.FromEntry
			require.NotEmpty(t, fromEntry)                     // Ensure the fromEntry record is not empty.
			require.Equal(t, account1.ID, fromEntry.AccountID) // Verify the account ID for the fromEntry.
			require.Equal(t, -amount, fromEntry.Amount)        // Verify the entry amount is negative (withdrawal).
			require.NotZero(t, fromEntry.ID)                   // Ensure the fromEntry has a valid ID.
			require.NotZero(t, fromEntry.CreatedAt)            // Ensure the fromEntry has a creation timestamp.

			// Verify that the fromEntry record can be retrieved by its ID.
			_, err = store.GetEntry(context.Background(), fromEntry.ID)
			require.NoError(t, err)

			// Check the toEntry details.
			toEntry := result.ToEntry
			require.NotEmpty(t, toEntry)                     // Ensure the toEntry record is not empty.
			require.Equal(t, account2.ID, toEntry.AccountID) // Verify the account ID for the toEntry.
			require.Equal(t, amount, toEntry.Amount)         // Verify the entry amount is positive (deposit).
			require.NotZero(t, toEntry.ID)                   // Ensure the toEntry has a valid ID.
			require.NotZero(t, toEntry.CreatedAt)            // Ensure the toEntry has a creation timestamp.

			// Verify that the toEntry record can be retrieved by its ID.
			_, err = store.GetEntry(context.Background(), toEntry.ID)
			require.NoError(t, err)
			/////////////////
			// check account
			fromAccount := result.FromAccount
			require.NotEmpty(t, fromAccount)
			require.Equal(t, account1.ID, fromAccount.ID)

			toAccount := result.ToAccount
			require.NotEmpty(t, toAccount)
			require.Equal(t, account2.ID, toAccount.ID)

			//check account balance
			fmt.Println(">> tx:", fromAccount.Balance, toAccount.Balance)
			diff1 := account1.Balance - fromAccount.Balance
			diff2 := toAccount.Balance - account2.Balance

			require.Equal(t, diff1, diff2)
			require.True(t, diff1 > 0)
			require.True(t, diff1%amount == 0)

			k := int(diff1 / amount)
			require.True(t, k >= 1 && k <= n)
			require.NotContains(t, existed, k)
			existed[k] = true
		*/
	}
	// check the final updated balance
	updateAccount1, err := testQuries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updateAccount2, err := testQuries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updateAccount1.Balance, updateAccount2.Balance)

	require.Equal(t, account1.Balance /* -int64(n)*amount */, updateAccount1.Balance)
	require.Equal(t, account2.Balance /* +int64(n)*amount */, updateAccount2.Balance)

	// +int64(n)*amount && +int64(n)*amount WE moved because we have total 10 transactions in which we need to 5 goes account1 and 5 goes account2 or these need to be in the same transaction - mean no change in account balance
}
