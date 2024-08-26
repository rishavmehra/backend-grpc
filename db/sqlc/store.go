package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store struct wraps the Queries struct and provides a reference to the database connection.
// It is used to execute database operations.
type Store struct {
	*Queries         // Embedding Queries to access its methods
	db       *sql.DB // Database connection
}

// NewStore initializes and returns a new Store instance.
// It takes a sql.DB pointer as an argument, which is used for database operations.
func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: New(db), // Initialize Queries with the given database connection
		db:      db,
	}
}

// execTx executes a function within a database transaction context.
// It begins a transaction, executes the given function, and then commits or rolls back the transaction.
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	// Begin a new database transaction with the provided context.
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return nil // If transaction start fails, return nil (could return the error instead)
	}

	// Create a new Queries instance using the transaction context.
	q := New(tx)

	// Execute the provided function, passing in the Queries instance.
	err = fn(q)
	if err != nil {
		// If an error occurs, attempt to roll back the transaction.
		if rbErr := tx.Rollback(); rbErr != nil {
			// If rollback also fails, return both the transaction error and rollback error.
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		// Return the transaction error if rollback succeeds.
		return err
	}
	// Commit the transaction if no errors occurred.
	return tx.Commit()
}

// TransferTxParams defines the parameters required for executing a money transfer transaction.
type TransferTxParams struct {
	FromAccountId int64 `json:"from_account_id"` // ID of the account to transfer money from
	ToAccountID   int64 `json:"to_account_id"`   // ID of the account to transfer money to
	Amount        int64 `json:"amount"`          // Amount of money to transfer
}

// TransferTxResult defines the result structure returned after executing a money transfer transaction.
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`     // Details of the transfer
	FromAccount Account  `json:"from_account"` // Account details of the sender
	ToAccount   Account  `json:"to_account"`   // Account details of the receiver
	FromEntry   Entry    `json: "from_entry"`  // Entry record for the sender's account
	ToEntry     Entry    `json: "to_entry"`    // Entry record for the receiver's account
}

var txkey = struct{}{}

// TransferTx executes a money transfer transaction between two accounts.
// It creates a transfer record, creates entry records for both accounts, and updates the account balances.
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	// Execute the transaction logic using execTx.
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// Retrieve the transaction name from the context for logging.
		txName := ctx.Value(txkey)

		fmt.Println(txName, "Create Transfer")
		// Create a transfer record with the given parameters (sender, receiver, amount).
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountId,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err // Return if there's an error during transfer creation.
		}

		fmt.Println(txName, "Create Entry 1")
		// Create an entry record for the sender's account (negative amount).
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountId,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err // Return if there's an error during sender's entry creation.
		}

		fmt.Println(txName, "Create entry 2")
		// Create an entry record for the receiver's account (positive amount).
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err // Return if there's an error during receiver's entry creation.
		}

		if arg.FromAccountId < arg.ToAccountID {
			// fmt.Println(txName, "get account 1")
			// Update the sender's account balance by subtracting the transfer amount.
			// fmt.Println(txName, "update account balance 1")

			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountId, -arg.Amount, arg.ToAccountID, arg.Amount)

			// the below comment code is long hand for above used funnction addMoney
			/* result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID:     arg.FromAccountId,
				Amount: -arg.Amount,
			})
			if err != nil {
				return err // Return if there's an error during the balance update of sender's account.
			}

			// fmt.Println(txName, "get account 2")
			// Update the receiver's account balance by adding the transfer amount.
			result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID:     arg.ToAccountID,
				Amount: arg.Amount,
			})
			if err != nil {
				return err // Return if there's an error during the balance update of receiver's account.
			} */
		} else {

			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountId, -arg.Amount)

			// the below comment code is long hand for above used funnction addMoney
			/* result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID:     arg.ToAccountID,
				Amount: arg.Amount,
			})
			if err != nil {
				return err // Return if there's an error during the balance update of receiver's account.
			}

			result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID:     arg.FromAccountId,
				Amount: -arg.Amount,
			})
			if err != nil {
				return err // Return if there's an error during the balance update of sender's account.
			} */
		}
		// The transaction is successfully executed.
		return nil
	})
	// Return the result of the transaction and any error encountered.
	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	ammount1 int64,
	accountID2 int64,
	ammount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: ammount1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: ammount2,
	})
	if err != nil {
		return
	}
	return
}
