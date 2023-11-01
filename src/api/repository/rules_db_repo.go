package repository

import (
	"fmt"
	"sync"

	"github.com/martinseco/modak-rate-limiter/src/api/errors"
	"github.com/martinseco/modak-rate-limiter/src/api/models"
)

type (
	rulesDBRepo struct {
		db   map[models.NotificationType]models.NotificationRule
		lock sync.Mutex
	}
)

func NewRulesDBRepo(db map[models.NotificationType]models.NotificationRule) *rulesDBRepo {
	return &rulesDBRepo{
		db:   db,
		lock: sync.Mutex{},
	}
}

func (r *rulesDBRepo) Insert(rule models.NotificationRule) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.db[rule.Type] = rule

	return nil
}

func (r *rulesDBRepo) Get(notificationType models.NotificationType) (models.NotificationRule, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	var rule models.NotificationRule
	rule, found := r.db[notificationType]

	if !found {
		return rule, errors.NotFoundError(fmt.Sprintf("Rule not found for type %s", notificationType))
	}

	return rule, nil
}
