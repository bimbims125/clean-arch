package postgresql

import (
	"context"
	"database/sql"

	"github.com/bimbims125/clean-arch/domain"
	"github.com/sirupsen/logrus"
)

type CategoryRepository struct {
	Conn *sql.DB
}

func NewCategoryRepository(conn *sql.DB) *CategoryRepository {
	return &CategoryRepository{conn}
}

func (p *CategoryRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Category, err error) {
	rows, err := p.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
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

func (p *CategoryRepository) Fetch(ctx context.Context) (result []domain.Category, err error) {
	query := `SELECT id, name FROM categories ORDER BY id ASC`

	res, err := p.fetch(ctx, query)
	if err != nil {
		return nil, err
	}
	return res, nil
}
