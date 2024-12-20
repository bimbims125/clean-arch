package postgresql

import (
	"context"
	"database/sql"
	"log"

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
			log.Println("Error while scanning product: ", err)
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

func (p *ProductRepository) FetchPaginated(ctx context.Context, offset, limit int) (total int, products []domain.Product, err error) {
	err = p.Conn.QueryRowContext(ctx, "SELECT COUNT(*) FROM products").Scan(&total)
	if err != nil {
		return 0, nil, err
	}

	rows, err := p.Conn.QueryContext(ctx,
		`SELECT
			p.id,
			p.name,
			p.price,
			p.category_id,
			p.stock,
			p.sold,
			p.image_url,
			c.name as category_name
			FROM
					products p
			JOIN
					categories c
			ON
					p.category_id = c.id
			ORDER BY
					p.id ASC
			LIMIT
					$1
			OFFSET
					$2;`, limit, offset)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var product domain.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Category.ID, &product.Stock, &product.Sold, &product.ImageURL, &product.Category.Name); err != nil {
			return 0, nil, err
		}
		products = append(products, product)
	}

	return total, products, nil
}

func (p *ProductRepository) GetByID(ctx context.Context, id int) (result domain.Product, err error) {
	query := `SELECT p.id, p.name, p.price, p.image_url, p.stock, p.sold, p.category_id, c.name as category_name
						FROM products p
						JOIN categories c ON p.category_id = c.id
						WHERE p.id = $1`
	res, err := p.fetch(ctx, query, id)
	if err != nil {
		return domain.Product{}, err
	}
	if len(res) == 0 {
		return domain.Product{}, domain.ErrNotFound
	}

	return res[0], nil
}
