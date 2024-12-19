package postgresql

import (
	"context"
	"database/sql"

	"github.com/bimbims125/clean-arch/domain"
	"github.com/sirupsen/logrus"
)

type UserRepository struct {
	Conn *sql.DB
}

// NewUserRepository creates an object representing a user.Repository interface
func NewUserRepository(conn *sql.DB) *UserRepository {
	return &UserRepository{conn}
}

func (p *UserRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.User, err error) {
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

	result = make([]domain.User, 0)
	for rows.Next() {
		u := domain.User{}
		err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Role)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, u)
	}
	return result, nil
}

func (p *UserRepository) Fetch(ctx context.Context) (result []domain.User, err error) {
	query := "SELECT id, name, email,role FROM users"
	res, err := p.fetch(ctx, query)
	if err != nil {
		return nil, err
	}
	return res, nil
}
