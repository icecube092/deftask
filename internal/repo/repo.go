package repo

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Repo interface {
	IsExistsSameAddrForUsers(ctx context.Context, userID1, userID2 int64) (bool, error)
}

type repo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) Repo {
	return &repo{
		db: db,
	}
}

func (r *repo) IsExistsSameAddrForUsers(
	ctx context.Context,
	userID1, userID2 int64,
) (bool, error) {
	const query = `
SELECT exists (
	SELECT addr FROM (
		SELECT addr, count(*)
		FROM conn_log
		WHERE user_id = $1
		GROUP BY addr
		HAVING count(*) > 1
	) AS q1
	
	INTERSECT
	
	SELECT addr FROM (
		SELECT addr, count(*)
		FROM conn_log
		WHERE user_id = $2
		GROUP BY addr
		HAVING count(*) > 1
	) AS q2
)
`

	var exists bool
	if err := r.db.GetContext(ctx, &exists, query, userID1, userID2); err != nil {
		return false, fmt.Errorf("db.GetContext: %w", err)
	}

	return exists, nil
}
