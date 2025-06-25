package service

import (
	"errors"
	"gou-pc/internal/agent"
	"gou-pc/internal/api/repository"
)

// ClientService interface định nghĩa các hàm thao tác với client
//go:generate mockgen -source=client_service.go -destination=mock_client_service.go -package=service

type ClientService interface {
	AssignUserToClient(clientID, userID string) error
	GetAllClients() ([]agent.ManagedClient, error)
	GetClientByAgentID(agentID string) (*agent.ManagedClient, error)
	GetClientsByUserID(userID string) ([]agent.ManagedClient, error)
}

type clientServiceImpl struct {
	repo repository.ClientRepository
}

func NewClientService(repo repository.ClientRepository) ClientService {
	return &clientServiceImpl{repo: repo}
}

func (s *clientServiceImpl) AssignUserToClient(clientID, userID string) error {
	clients, err := s.repo.GetAll()
	if err != nil {
		return err
	}
	found := false
	for i, c := range clients {
		if c.ClientID == clientID {
			clients[i].UserID = userID
			found = true
			break
		}
	}
	if !found {
		return errors.New("client not found")
	}
	return s.repo.SaveAll(clients)
}

func (s *clientServiceImpl) GetAllClients() ([]agent.ManagedClient, error) {
	return s.repo.GetAll()
}

func (s *clientServiceImpl) GetClientByAgentID(agentID string) (*agent.ManagedClient, error) {
	return s.repo.FindByAgentID(agentID)
}

func (s *clientServiceImpl) GetClientsByUserID(userID string) ([]agent.ManagedClient, error) {
	return s.repo.FindByUserID(userID)
}
