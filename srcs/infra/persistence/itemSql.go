package persistence

import (
	"database/sql"
	"todoapi/domain/model"
	"todoapi/domain/repository"
)

type itemSqlRepository struct {
	db *sql.DB
}

func NewItemSqlRepository(db *sql.DB) repository.ItemRepository {
	return &itemSqlRepository{
		db: db,
	}
}

func (repo itemSqlRepository) Save(item model.Item) (model.Item, error) {
	row := repo.db.QueryRow("INSERT INTO items (name, done) values ($1, $2) RETURNING id", item.Name, item.Done)
	err := row.Scan(&item.Id)
	if err != nil {
		return model.Item{}, err
	}
	return item, err
}

func (repo itemSqlRepository) GetAll() ([]model.Item, error) {
	rows, err := repo.db.Query(("SELECT id, name, done FROM items"))
	if err != nil {
		return nil, err
	}
	var items []model.Item
	for rows.Next() {
		var item model.Item
		err := rows.Scan(&item.Id, &item.Name, &item.Done)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (repo itemSqlRepository) GetById(id model.ID) (model.Item, error) {
	var item model.Item
	row := repo.db.QueryRow("SELECT id, name, done FROM items WHERE id = ?", id)
	err := row.Scan(item.Id, item.Name, item.Done)
	if err != nil {
		return item, err
	}
	return item, nil
}

func (repo itemSqlRepository) DeleteById(id model.ID) error {
	_, err := repo.db.Exec("DELETE FROM items WHERE id = ?", id)
	return err
}

func (repo itemSqlRepository) UpdateDone(id model.ID, done bool) error {
	_, err := repo.db.Exec("UPDATE items SET done=? where id=?", done, id)
	return err
}

func (repo itemSqlRepository) DeleteDone() error {
	_, err := repo.db.Exec("DELETE FROM items WHERE done=true")
	return err
}
