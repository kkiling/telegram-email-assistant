package storage_impl

import (
	"github.com/kiling91/telegram-email-assistant/internal/factory"
	"github.com/kiling91/telegram-email-assistant/internal/storage"
	"github.com/kiling91/telegram-email-assistant/internal/types"
)

type service struct {
	fact factory.Factory
}

func NewStorage(fact factory.Factory) storage.Storage {
	return &service{
		fact: fact,
	}
}

func (s *service) AddUser(user *types.EmailUser) (types.UID, error) {
	return 0, nil
}

func (s *service) GetUser(uid types.UID) (*types.EmailUser, error) {
	return nil, nil
}
