package dependency

import (
	usersvc "be/internal/services/user"
)

func NewUserService(infra *Infra, outbox usersvc.OutboxEnqueuer) *usersvc.Service {
	userRepo := newUserRepository(infra.DB)
	return usersvc.NewService(userRepo, outbox)
}
