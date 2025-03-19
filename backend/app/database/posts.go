package database

import (
	"context"
	"database/sql"
	"one-help/app/posts"

	"github.com/zeebo/errs"
)

// ErrPosts indicates that there was an error in the database.
var ErrPosts = errs.Class("post services repository")

// postsDB provides access to post services db.
//
// architecture: Database
type postsDB struct {
	conn *sql.DB
}

// newPostsDB is a constructor for base postsDB.
func newPostsDB(baseConn *sql.DB) posts.DB {
	return &postsDB{
		conn: baseConn,
	}
}

// Create inserts post service into the database.
func (db *postsDB) Create(ctx context.Context, post string) error {
	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		return ErrPosts.Wrap(err)
	}

	defer DeferCommitRollback(tx, &err)

	query := `INSERT INTO posts(post)
              VALUES ($1)`
	_, err = tx.ExecContext(ctx, query, post)
	if err != nil {
		return ErrPosts.Wrap(err)
	}

	return ErrPosts.Wrap(err)
}

// Delete removes user from the database.
func (db *postsDB) Delete(ctx context.Context, post string) error {
	query := `DELETE FROM posts WHERE post = $1`
	_, err := db.conn.ExecContext(ctx, query, post)
	return ErrPosts.Wrap(err)
}

// List returns all available post services.
func (db *postsDB) List(ctx context.Context) ([]string, error) {
	query := `SELECT post FROM posts`
	rows, err := db.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, ErrPosts.Wrap(err)
	}

	defer rows.Close()

	var posts []string
	for rows.Next() {
		var post string
		if err := rows.Scan(&post); err != nil {
			return nil, ErrPosts.Wrap(err)
		}

		posts = append(posts, post)
	}

	return posts, nil
}
