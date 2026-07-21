package service

import (
	"finstat/internal/lib"
	"finstat/internal/models"
	"finstat/internal/repository"
)

type CategoryService struct {
	repo             repository.Category
	systemCategories []models.Category
}

func NewCategoryService(repo repository.Category) *CategoryService {
	return &CategoryService{
		repo:             repo,
		systemCategories: nil,
	}
}

func (s *CategoryService) InsertCategory(userID uint, categoryName string) (uint, error) {
	formattedName, err := lib.FormatName(categoryName, 3)
	if err != nil {
		return 0, err
	}

	return s.repo.InsertCategory(userID, formattedName)
}

func (s *CategoryService) UpdateCategory(userID, categoryID uint, newCategoryName string) (bool, error) {
	formattedName, err := lib.FormatName(newCategoryName, 3)
	if err != nil {
		return false, err
	}

	return s.repo.UpdateCategory(userID, categoryID, formattedName)
}

func (s *CategoryService) DeleteCategory(userID, categoryID uint) (bool, error) {
	return s.repo.DeleteCategory(userID, categoryID)
}

func (s *CategoryService) SystemCategories() ([]models.Category, error) {
	if s.systemCategories == nil {
		categories, err := s.repo.SystemCategories()
		if err != nil {
			return nil, err
		}
		s.systemCategories = categories
	}

	return s.systemCategories, nil
}

func (s *CategoryService) UserCategories(userID uint) ([]models.Category, error) {
	return s.repo.UserCategories(userID)
}

func (s *CategoryService) Categories(userID uint) ([]models.Category, error) {
	systemCategories, err := s.SystemCategories()
	if err != nil {
		return nil, lib.Ewrap("Couldn't get categories", err)
	}

	userCategories, err := s.UserCategories(userID)
	if err != nil {
		return nil, lib.Ewrap("Couldn't get categories", err)
	}

	result := append(systemCategories, userCategories...)

	return result, nil
}
