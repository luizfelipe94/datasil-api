package auth

import (
	"database/sql"
	"log"

	users "github.com/luizfelipe94/datasil/modules/users/models"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) GetUserByEmail(email string) (*users.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	user := new(users.User)
	for rows.Next() {
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
			&user.CompanyID,
		)
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}
	return user, nil
}
