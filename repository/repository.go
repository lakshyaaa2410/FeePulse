// repository/repository.go
package repository

import (
	"fee-reminder/db"
	"fee-reminder/model"
)

type RepositoryInterface interface {
	AddMembers(members []model.MembersDB) error
	GetAllMembers() ([]model.Members, error)
	GetAllExpiringMemberships() ([]model.Members, error)
}

type Repository struct {
	db db.DBInterface
}

func InitializeRepository(db db.DBInterface) *Repository {
	return &Repository{db: db}
}
