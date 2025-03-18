package database

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/zeebo/errs"

	"one-help/app/users"
)

// ErrUsers indicates that there was an error in the database.
var ErrUsers = errs.Class("users repository")

// usersDB provides access to users db.
//
// architecture: Database
type usersDB struct {
	conn *sql.DB
}

// newUsersDB is a constructor for base usersDB.
func newUsersDB(baseConn *sql.DB) users.DB {
	return &usersDB{
		conn: baseConn,
	}
}

// Create inserts user into the database.
func (db *usersDB) Create(ctx context.Context, user users.User) error {
	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		return ErrUsers.Wrap(err)
	}

	defer DeferCommitRollback(tx, &err)

	query := `INSERT INTO users(user_id, first_name, last_name, website, file_name)
              VALUES ($1, $2, $3, $4, $5)`
	_, err = tx.ExecContext(ctx, query, user.ID, user.FirstName, user.LastName, user.Website, user.FileName)
	if err != nil {
		return ErrUsers.Wrap(err)
	}

	if !user.IsDeliveryAddressEmpty() {
		query = `INSERT INTO delivery_addresses(user_id, city, post, post_department)
                 VALUES ($1, $2, $3, $4)`
		_, err = tx.ExecContext(ctx, query, user.ID, user.City, user.Post, user.PostDepartment)
		if err != nil {
			return ErrUsers.Wrap(err)
		}
	}

	return ErrUsers.Wrap(err)
}

// Get returns user from the database by ID.
func (db *usersDB) Get(ctx context.Context, id uuid.UUID) (users.User, error) {
	var (
		user           users.User
		city           sql.NullString
		post           sql.NullString
		postDepartment sql.NullString
	)

	query := `SELECT u.user_id, first_name, last_name, website, file_name, city, post, post_department
              FROM users u LEFT JOIN delivery_addresses d ON u.user_id = d.user_id
              WHERE u.user_id = $1`

	row := db.conn.QueryRowContext(ctx, query, id)
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Website, &user.FileName, &city, &post, &postDepartment)
	if err != nil {
		if errs.Is(err, sql.ErrNoRows) {
			return users.User{}, ErrUsers.Wrap(users.ErrNoUser)
		}

		return user, ErrUsers.Wrap(err)
	}

	user.City = city.String
	user.Post = post.String
	user.PostDepartment = postDepartment.String

	return user, nil
}

// Update updates user in database by id.
func (db *usersDB) Update(ctx context.Context, user users.User) error {
	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		return ErrUsers.Wrap(err)
	}

	defer DeferCommitRollback(tx, &err)

	query := `UPDATE users
              SET first_name = $2, last_name = $3, website = $4, file_name = $5
              WHERE user_id = $1`
	_, err = tx.ExecContext(ctx, query, user.ID, user.FirstName, user.LastName, user.Website, user.FileName)
	if err != nil {
		return ErrUsers.Wrap(err)
	}

	if !user.IsDeliveryAddressEmpty() {
		query = `INSERT INTO delivery_addresses(user_id, city, post, post_department)
                 VALUES ($1, $2, $3, $4)
                 ON CONFLICT DO UPDATE
                 SET city = $2, post = $3, post_department = $4`
		_, err = tx.ExecContext(ctx, query, user.ID, user.City, user.Post, user.PostDepartment)
	} else {
		query = `DELETE FROM delivery_addresses WHERE user_id = $1`
		_, err = tx.ExecContext(ctx, query, user.ID)
	}
	if err != nil {
		return ErrUsers.Wrap(err)
	}

	return nil
}

// Delete removes user from the database.
func (db *usersDB) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE user_id = $1`
	_, err := db.conn.ExecContext(ctx, query, id)
	return ErrUsers.Wrap(err)
}
