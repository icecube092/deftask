package service

import (
	"context"
	"fmt"

	"deftask/internal/repo"
)

type Service interface {
	IsUserDuplicate(ctx context.Context, userID1, userID2 int64) (bool, error)
}

type service struct {
	db repo.Repo
}

func New(db repo.Repo) Service {
	return &service{
		db: db,
	}
}

func (s *service) IsUserDuplicate(ctx context.Context, userID1, userID2 int64) (
	bool,
	error,
) {
	if userID1 == userID2 {
		return true, nil
	}

	dupl, err := s.db.IsExistsSameAddrForUsers(ctx, userID1, userID2)
	if err != nil {
		return false, fmt.Errorf("db.IsExistsSameAddrForUsers: %w", err)
	}

	return dupl, nil
}
