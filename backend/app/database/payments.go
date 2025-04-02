package database

import (
	"context"
	"database/sql"
	"errors"

	"one-help/app/payments"

	"github.com/google/uuid"
	"github.com/zeebo/errs"
)

// ErrPayments indicates that there was an error in the database.
var ErrPayments = errs.Class("payments repository")

// paymentsDB provides access to payments db.
//
// architecture: Database
type paymentsDB struct {
	conn *sql.DB
}

// newPaymentsDB is a constructor for base paymentsDB.
func newPaymentsDB(baseConn *sql.DB) payments.DB {
	return &paymentsDB{
		conn: baseConn,
	}
}

// Create inserts payment into the database.
func (db *paymentsDB) Create(ctx context.Context, payment payments.Payment) error {
	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		return ErrPayments.Wrap(err)
	}

	defer DeferCommitRollback(tx, &err)

	query := `INSERT INTO payments(donation_id, payment_type, transaction_id, confirmed)
              VALUES ($1, $2, $3, $4)`
	_, err = tx.ExecContext(ctx, query, payment.DonationId, payment.PaymentType, payment.TransactionId, payment.Confirmed)
	return ErrPayments.Wrap(err)
}

// Get returns payment from the database by donation ID.
func (db *paymentsDB) Get(ctx context.Context, id uuid.UUID) (payments.Payment, error) {
	var (
		payment payments.Payment
	)

	query := `SELECT donation_id, payment_type, transaction_id, confirmed
	          FROM payments
              WHERE donation_id = $1`

	row := db.conn.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&payment.DonationId,
		&payment.PaymentType,
		&payment.TransactionId,
		&payment.Confirmed,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return payments.Payment{}, ErrPayments.Wrap(payments.ErrNoPayment)
		}
		return payment, ErrPayments.Wrap(err)
	}

	return payment, nil
}

// List returns all the payments.
func (db *paymentsDB) List(ctx context.Context) ([]payments.Payment, error) {
	query := `SELECT donation_id, payment_type, transaction_id, confirmed
              FROM payments`
	rows, err := db.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, ErrPayments.Wrap(err)
	}
	defer func() { err = errs.Combine(err, rows.Close()) }()

	var paymentsList []payments.Payment
	for rows.Next() {
		var (
			payment payments.Payment
		)
		err := rows.Scan(
			&payment.DonationId,
			&payment.PaymentType,
			&payment.TransactionId,
			&payment.Confirmed,
		)
		if err != nil {
			return nil, ErrPayments.Wrap(err)
		}
		paymentsList = append(paymentsList, payment)
	}

	if err = rows.Err(); err != nil {
		return nil, ErrPayments.Wrap(err)
	}

	return paymentsList, nil
}

// Update updates payment in database by donation ID.
func (db *paymentsDB) Update(ctx context.Context, payment payments.Payment) error {
	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		return ErrPayments.Wrap(err)
	}
	defer DeferCommitRollback(tx, &err)

	query := `UPDATE payments
	          SET payment_type = $2, transaction_id = $3, confirmed = $4
	          WHERE donation_id = $1`

	_, err = tx.ExecContext(ctx, query,
		payment.DonationId,
		payment.PaymentType,
		payment.TransactionId,
		payment.Confirmed,
	)
	if err != nil {
		return ErrPayments.Wrap(err)
	}

	return nil
}

// Delete removes payment from the database.
func (db *paymentsDB) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM payments WHERE donation_id = $1`
	_, err := db.conn.ExecContext(ctx, query, id)
	return ErrPayments.Wrap(err)
}
