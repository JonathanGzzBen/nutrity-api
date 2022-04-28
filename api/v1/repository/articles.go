package repository

import (
	"github.com/JonathanGzzBen/ingenialists/api/v1/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ArticlesRepository interface {
	GetAllArticles() ([]models.Article, error)
	GetArticle(uint) (*models.Article, error)
	CreateArticle(*models.Article) (*models.Article, error)
	UpdateArticle(*models.Article) (*models.Article, error)
	DeleteArticle(uint) error
}

type ArticlesGormRepository struct {
	db *gorm.DB
}

func NewArticlesGormRepository(db *gorm.DB) *ArticlesGormRepository {
	db.AutoMigrate(&models.Article{})
	return &ArticlesGormRepository{
		db: db,
	}
}

func (r ArticlesGormRepository) GetAllArticles() ([]models.Article, error) {
	var articles []models.Article
	res := r.db.Preload(clause.Associations).Find(&articles)
	if res.Error != nil {
		return nil, ErrCouldNotRetrieve
	}
	return articles, nil
}

func (r ArticlesGormRepository) GetArticle(id uint) (*models.Article, error) {
	var article *models.Article
	res := r.db.Preload(clause.Associations).Find(&article, id)
	if res.Error == gorm.ErrRecordNotFound || res.RowsAffected != 1 {
		return nil, ErrNotFound
	}
	if res.Error != nil {
		return nil, ErrCouldNotRetrieve
	}
	return article, nil
}

func (r ArticlesGormRepository) CreateArticle(a *models.Article) (*models.Article, error) {
	res := r.db.Create(&a)
	if res.Error != nil {
		return nil, ErrCouldNotCreate
	}
	a, err := r.GetArticle(a.ID)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (r ArticlesGormRepository) UpdateArticle(a *models.Article) (*models.Article, error) {
	res := r.db.Save(&a)
	if res.Error != nil {
		return nil, ErrCouldNotUpdate
	}
	a, err := r.GetArticle(a.ID)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (r ArticlesGormRepository) DeleteArticle(id uint) error {
	a, err := r.GetArticle(id)
	if err != nil {
		return err
	}
	res := r.db.Delete(&a)
	if res.Error != nil {
		return ErrCouldNotDelete
	}
	return nil
}
