package repository

import (
	"context"
	"fee-reminder/model"
	"fmt"
	"time"
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

func (repository *Repository) AddMember(member model.MembersDB) error {

	sql := `INSERT INTO members (name, phone, joining_date, duration, expiry_date) VALUES ($1, $2, $3, $4, $5)`
	pool := repository.db.GetPool()

	_, err := pool.Exec(context.Background(), sql,
		member.Name,
		member.Phone,
		member.JoiningDate,
		member.Duration,
		member.ExpiryDate,
	)

	if err != nil {
		return fmt.Errorf("in repository.AddMember(): failed to insert member %s: %w", member.Name, err)
	}

	return nil
}

func (repository *Repository) GetAllExpiringMemberships() ([]model.Members, error) {
	// Load IST timezone for accurate date calculations
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		return nil, fmt.Errorf("in repository.GetAllExpiringMemberships(): error loading location: %w", err)
	}

	// Get current time in IST
	now := time.Now().In(loc)

	// Generate a 5-day window: 2 days before and 2 days after today
	// dateRange[0] = -2 days, [1] = -1 day, [2] = today, [3] = +1 day, [4] = +2 days
	dateRange := make([]string, 5)
	for i := -2; i <= 2; i++ {
		dateRange[i+2] = now.AddDate(0, 0, i).Format("02-01-2006")
	}

	// Use ANY($1) to match expiry_date against all 5 dates in a single query
	sql := `SELECT name, phone, joining_date, duration FROM members WHERE expiry_date = ANY($1)`

	pool := repository.db.GetPool()
	rows, err := pool.Query(context.Background(), sql, dateRange)
	if err != nil {
		return nil, fmt.Errorf("in repository.GetAllExpiringMemberships(): failed to query expiring memberships: %w", err)
	}
	defer rows.Close()

	var expiringMemberships []model.Members
	for rows.Next() {
		var member model.Members
		err := rows.Scan(&member.Name, &member.Phone, &member.JoiningDate, &member.Duration)
		if err != nil {
			return nil, fmt.Errorf("in repository.GetAllExpiringMemberships(): failed to scan expiring membership row: %w", err)
		}
		expiringMemberships = append(expiringMemberships, member)
	}

	// Check for any errors encountered during row iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("in repository.GetAllExpiringMemberships(): error iterating expiring membership rows: %w", err)
	}

	return expiringMemberships, nil
}
