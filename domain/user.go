package domain

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password,omitempty" validate:"required,min=8,password"`
	Role     string `json:"role"`
}

func (u *User) HashPassword() error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashed)
	return nil
}
