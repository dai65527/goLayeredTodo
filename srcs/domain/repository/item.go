package repository

import "todoapi/domain/model"

type ItemRepository interface {
	Save(item model.Item) (model.Item, error)
	GetAll() ([]model.Item, error)
	GetById(id model.ID) (model.Item, error)
	DeleteById(id model.ID) error
	UpdateDone(id model.ID, done bool) error
	DeleteDone() error
}
