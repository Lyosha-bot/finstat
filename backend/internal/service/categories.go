package service

import (
	ewrap "finstat/internal/lib"
	"finstat/internal/repository"
)

type Category = repository.Category

type CategoryRepo interface {
	SystemCategories() ([]Category, error)
	UserCategories(userID uint) ([]Category, error)
	Categories(userID uint) ([]Category, error)
}

type CategoryService struct {
	repo             CategoryRepo
	systemCategories []Category
}

func NewCategoryService(repo CategoryRepo) (*CategoryService, error) {
	systemCategories, err := repo.SystemCategories()
	if err != nil {
		return nil, ewrap.Wrap("Couldn't create category service", err)
	}

	return &CategoryService{
		repo:             repo,
		systemCategories: systemCategories,
	}, nil
}

func (s *CategoryService) AddCategory(userID uint, category string) {

}

func (s *CategoryService) SystemCategories() []Category {
	return s.systemCategories
}

func (s *CategoryService) UserCategories(userID uint) ([]Category, error) {
	return s.repo.UserCategories(userID)
}

func (s *CategoryService) Categories(userID uint) ([]Category, error) {
	return s.repo.Categories(userID)
}
