package repositories

import (
	"go-boilerplate/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user models.User) (models.User, error)
	GetUserById(id string) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	UpdateUser(user models.User) error
	DeleteUser(id string) error
}

type PostgresUserRepository struct {
	db *gorm.DB
}

func NewPostgresUserRepository(db *gorm.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) CreateUser(user models.User) (models.User, error) {
	err := r.db.Create(&user).Error
	return user, err
}

func (r *PostgresUserRepository) GetUserById(id string) (models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	return user, err
}

func (r *PostgresUserRepository) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return user, err
}

func (r *PostgresUserRepository) UpdateUser(user models.User) error {
	return r.db.Save(&user).Error
}

func (r *PostgresUserRepository) DeleteUser(id string) error {
	var user models.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return err
	}
	return r.db.Delete(&user).Error
}
