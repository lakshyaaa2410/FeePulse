package repository

import (
	"context"
	"fee-reminder/model"
	"fmt"
)

func (repo *Repository) AddMembers(members []model.Members) error {
	if len(members) == 0 {
		return nil
	}

	sql := `INSERT INTO members (name, phone, joining_date, duration, expiry_date) VALUES ($1, $2, $3, $4, $5)`
	pool := repo.db.GetPool()

	for _, member := range members {
		_, err := pool.Exec(context.Background(), sql,
			member.Name,
			member.Phone,
			member.JoiningDate,
			member.Duration,
			member.ExpiryDate,
		)
		if err != nil {
			return fmt.Errorf("in repository.AddMembers(): failed to insert member %s: %w", member.Name, err)
		}
	}

	return nil
}
