package mysql

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
func NewMySQLUserRepository(conn *sql.DB) *UserRepository {
	return &UserRepository{conn}
}

func (m *UserRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.User, err error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
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

func (m *UserRepository) Fetch(ctx context.Context) (result []domain.User, err error) {
	query := "SELECT id, name, email, role FROM users"
	res, err := m.fetch(ctx, query)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (m *UserRepository) Create(ctx context.Context, user domain.User) error {
	query := `
		INSERT INTO users (name, email, password, role)
		VALUES (?, ?, ?, ?)
	`
	var createdUser domain.User
	if err := user.HashPassword(); err != nil {
		return err
	}

	err := m.Conn.QueryRowContext(ctx, query, user.Name, user.Email, user.Password, user.Role).
		Scan(&createdUser.ID, &createdUser.Name, &createdUser.Email, &createdUser.Role)
	if err != nil {
		logrus.Error(err)
	}
	return nil
}

func (m *UserRepository) GetByEmail(ctx context.Context, email string) (result domain.User, err error) {
	query := `SELECT id, name, email, role FROM users WHERE email = ?`
	res, err := m.fetch(ctx, query, email)
	if err != nil {
		return domain.User{}, err
	}
	if len(res) == 0 {
		return domain.User{}, domain.ErrNotFound
	}
	return res[0], nil
}
