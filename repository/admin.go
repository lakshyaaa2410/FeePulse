package repository

import (
	"context"
	"fee-reminder/model"
	"fmt"
)

func (repo *Repository) AddMembers(members []model.MembersDB) error {
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

func (repository *Repository) GetAllMembers() ([]model.Members, error) {

	sql := `SELECT name, phone, joining_date, duration FROM members`
	pool := repository.db.GetPool()

	rows, err := pool.Query(context.Background(), sql)
	if err != nil {
		return nil, fmt.Errorf("in repository.GetAllMembers(): failed to query members: %w", err)
	}
	defer rows.Close()

	var members []model.Members
	for rows.Next() {
		var member model.Members
		err := rows.Scan(&member.Name, &member.Phone, &member.JoiningDate, &member.Duration)
		if err != nil {
			return nil, fmt.Errorf("in repository.GetAllMembers(): failed to scan member row: %w", err)
		}
		members = append(members, member)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("in repository.GetAllMembers(): error iterating member rows: %w", err)
	}

	return members, nil
}
