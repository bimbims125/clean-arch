package postgresql

import (
	"context"
	"database/sql"

	"github.com/bimbims125/clean-arch/domain"
	"github.com/sirupsen/logrus"
)

type ProductRepository struct {
	Conn *sql.DB
}

func NewProductRepository(conn *sql.DB) *ProductRepository {
	return &ProductRepository{conn}
}

func (p *ProductRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Product, err error) {
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

	result = make([]domain.Product, 0)
	for rows.Next() {
		p := domain.Product{}
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.ImageURL, &p.Stock, &p.Sold, &p.Category.ID, &p.Category.Name)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, p)
	}
	return result, nil
}

func (p *ProductRepository) Fetch(ctx context.Context) (result []domain.Product, err error) {
	query := `SELECT p.id, p.name, p.price, p.image_url, p.stock, p.sold, p.category_id, c.name as category_name
						FROM products p
						JOIN categories c ON p.category_id = c.id
						ORDER BY p.id ASC`

	res, err := p.fetch(ctx, query)
	if err != nil {
		return nil, err
	}
	return res, nil
}
