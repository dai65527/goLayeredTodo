package usecase

import (
	"todoapi/domain/model"
	"todoapi/domain/repository"
)

type ItemUseCase interface {
	AddItem(name string) error
	GetAll() ([]model.Item, error)
	Done(id model.ID) error
	UnDone(id model.ID) error
	DeleteDone() error
	Delete(id model.ID) error
}

type itemUseCase struct {
	repository repository.ItemRepository
}

func NewItemUseCase(repo repository.ItemRepository) ItemUseCase {
	return &itemUseCase{
		repository: repo,
	}
}

func (usecase itemUseCase) AddItem(name string) error {
	_, err := usecase.repository.Save(model.Item{Name: name, Done: false})
	return err
}

func (usecase itemUseCase) GetAll() ([]model.Item, error) {
	return usecase.repository.GetAll()
}

func (usecase itemUseCase) DeleteDone() error {
	return usecase.repository.DeleteDone()
}

func (usecase itemUseCase) Delete(id model.ID) error {
	return usecase.repository.DeleteById(id)
}

func (usecase itemUseCase) Done(id model.ID) error {
	return usecase.repository.UpdateDone(id, true)
}

func (usecase itemUseCase) UnDone(id model.ID) error {
	return usecase.repository.UpdateDone(id, false)
}
