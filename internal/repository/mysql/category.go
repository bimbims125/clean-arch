package mysql

import (
	"context"
	"database/sql"

	"github.com/bimbims125/clean-arch/domain"
	"github.com/sirupsen/logrus"
)

type CategoryRepository struct {
	Conn *sql.DB
}

func NewMySQLCategoryRepository(conn *sql.DB) *CategoryRepository {
	return &CategoryRepository{conn}
}

func (m *CategoryRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Category, err error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			logrus.Error(errRow)
		}
	}()

	result = make([]domain.Category, 0)
	for rows.Next() {
		c := domain.Category{}
		err := rows.Scan(&c.ID, &c.Name)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, c)
	}
	return result, nil
}

func (m *CategoryRepository) Fetch(ctx context.Context) (result []domain.Category, err error) {
	query := "SELECT id, name FROM categories"
	res, err := m.fetch(ctx, query)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return res, nil
}

func (m *CategoryRepository) GetByID(ctx context.Context, id string) (result domain.Category, err error) {
	query := "SELECT id, name FROM categories WHERE id = ?"
	res, err := m.fetch(ctx, query, id)
	if err != nil {
		logrus.Error(err)
		return domain.Category{}, err
	}
	if len(res) == 0 {
		return domain.Category{}, domain.ErrNotFound
	}

	return res[0], nil
}

func (m *CategoryRepository) Create(ctx context.Context, category domain.Category) error {
	query := `INSERT INTO categories (id, name) VALUES (?, ?)`

	var createdCategory domain.Category
	err := m.Conn.QueryRowContext(ctx, query, category.ID, category.Name).Scan(&createdCategory.ID, &createdCategory.Name)
	if err != nil {
		logrus.Error(err)
	}
	return nil
}
