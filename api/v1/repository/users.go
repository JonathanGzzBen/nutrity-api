package repository

import (
	"errors"

	"github.com/JonathanGzzBen/nutrity-api/api/v1/models"
	"gorm.io/gorm"
)

var (
	ErrCouldNotRetrieve = errors.New("could not retrieve records")
	ErrNotFound         = errors.New("record not found")
	ErrCouldNotCreate   = errors.New("could not insert record")
	ErrCouldNotUpdate   = errors.New("could not update record")
	ErrCouldNotDelete   = errors.New("could not delete record")
)

type UsersRepository interface {
	GetAllUsers() ([]models.User, error)
	GetUser(uint) (*models.User, error)
	GetUserByGoogleToken(string) (*models.User, error)
	GetUserByAccessToken(string) (*models.User, error)
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

func (r *UsersGormRepository) GetUserByAccessToken(at string) (*models.User, error) {
	var user *models.User
	res := r.db.Where("access_token = ?", at).First(&user)
	if res.Error == gorm.ErrRecordNotFound {
		return nil, ErrNotFound
	}
	if res.Error != nil {
		return nil, ErrCouldNotRetrieve
	}
	return user, nil
}

func (r *UsersGormRepository) GetUserByGoogleToken(at string) (*models.User, error) {
	var user *models.User
	res := r.db.Where("google_token = ?", at).First(&user)
	if res.Error == gorm.ErrRecordNotFound {
		return nil, ErrNotFound
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
