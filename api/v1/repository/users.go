package repository

import (
	"github.com/JonathanGzzBen/ingenialists/api/v1/models"
	"gorm.io/gorm"
)

type UsersRepository interface {
	GetAllUsers() ([]models.User, error)
	GetUser(uint) (*models.User, error)
	GetUserByGoogleSub(string) (*models.User, error)
	CreateUser(*models.User) (*models.User, error)
	UpdateUser(*models.User) (*models.User, error)
}

type UsersGormRepository struct {
	db *gorm.DB
}

func NewUsersGormRepository(db *gorm.DB) *UsersGormRepository {
	db.AutoMigrate(&models.User{})
	return &UsersGormRepository{
		db: db,
	}
}

func (r *UsersGormRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	res := r.db.Find(&users)
	if res.Error != nil {
		return nil, ErrCouldNotRetrieve
	}
	return users, nil
}

func (r *UsersGormRepository) GetUser(id uint) (*models.User, error) {
	var user *models.User
	res := r.db.Find(&user, id)
	if res.Error == gorm.ErrRecordNotFound || res.RowsAffected != 1 {
		return nil, ErrNotFound
	}
	if res.Error != nil {
		return nil, ErrCouldNotRetrieve
	}
	return user, nil
}

func (r *UsersGormRepository) GetUserByGoogleSub(sub string) (*models.User, error) {
	var user *models.User
	res := r.db.Where("google_sub = ?", sub).First(&user)
	if res.Error == gorm.ErrRecordNotFound {
		return nil, ErrCouldNotRetrieve
	}
	if res.Error != nil {
		return nil, ErrCouldNotRetrieve
	}
	return user, nil
}

func (r *UsersGormRepository) CreateUser(u *models.User) (*models.User, error) {
	res := r.db.Create(u)
	if res.Error != nil {
		return nil, ErrCouldNotCreate
	}
	return u, nil
}

func (r *UsersGormRepository) UpdateUser(u *models.User) (*models.User, error) {
	res := r.db.Save(u)
	if res.Error != nil {
		return nil, ErrCouldNotUpdate
	}
	return u, nil
}
