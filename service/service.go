// service/service.go
package service

import (
	"fee-reminder/model"
	"fee-reminder/repository"
)

type ServiceInterface interface {
	AddMembersFromCSV(csvData []byte) error
	GetAllMembers() ([]model.Members, error)
	AddNewMember(member model.Members) error
	GetAllExpiringMemberships() ([]model.Members, error)
}

type Service struct {
	repository repository.RepositoryInterface
}

func InitializeService(repository repository.RepositoryInterface) *Service {
	return &Service{repository: repository}
}
