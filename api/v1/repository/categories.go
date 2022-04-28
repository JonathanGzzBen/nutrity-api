package repository

import (
	"errors"
	"fmt"

	"github.com/JonathanGzzBen/ingenialists/api/v1/models"
	"gorm.io/gorm"
)

type CategoriesRepository interface {
	GetAllCategories() ([]models.Category, error)
	GetCategory(uint) (*models.Category, error)
	CreateCategory(*models.Category) (*models.Category, error)
	UpdateCategory(*models.Category) (*models.Category, error)
	DeleteCategory(uint) error
}

type CategoriesGormRepository struct {
	db *gorm.DB
}

var (
	ErrCouldNotRetrieve = errors.New("could not retrieve records")
	ErrNotFound         = errors.New("record not found")
	ErrCouldNotCreate   = errors.New("could not insert record")
	ErrCouldNotUpdate   = errors.New("could not update record")
	ErrCouldNotDelete   = errors.New("could not delete record")
)

func NewCategoriesGormRepository(db *gorm.DB) *CategoriesGormRepository {
	db.AutoMigrate(&models.Category{})
	return &CategoriesGormRepository{
		db: db,
	}
}

func (r CategoriesGormRepository) GetAllCategories() ([]models.Category, error) {
	var categories []models.Category
	res := r.db.Find(&categories)
	if res.Error != nil {
		return nil, fmt.Errorf("could not retrieve categories: %s", res.Error.Error())
	}
	return categories, nil
}

func (r *CategoriesGormRepository) GetCategory(id uint) (*models.Category, error) {
	var category *models.Category
	res := r.db.Find(&category, id)
	if res.Error == gorm.ErrRecordNotFound || res.RowsAffected != 1 {
		return nil, ErrNotFound
	}
	if res.Error != nil {
		return nil, ErrCouldNotRetrieve
	}
	return category, nil
}

func (r *CategoriesGormRepository) CreateCategory(c *models.Category) (*models.Category, error) {
	res := r.db.Create(c)
	if res.Error != nil {
		return nil, ErrCouldNotCreate
	}
	return c, nil
}

func (r *CategoriesGormRepository) UpdateCategory(c *models.Category) (*models.Category, error) {
	res := r.db.Save(c)
	if res.Error != nil {
		return nil, ErrCouldNotUpdate
	}
	return c, nil
}

func (r *CategoriesGormRepository) DeleteCategory(id uint) error {
	c, err := r.GetCategory(id)
	if err != nil {
		return err
	}
	res := r.db.Delete(c)
	if res.Error != nil {
		return ErrCouldNotDelete
	}
	return nil
}
