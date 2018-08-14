package main

import "github.com/JoeSelvik/hdm-service/models"

type ResourceController interface {
	Path() string
	DBTableName() string

	Create(m []models.Resource) ([]int, *ApplicationError)
	Read(id int) (models.Resource, *ApplicationError)
	Update(m []models.Resource) *ApplicationError
	Destroy(ids []int) *ApplicationError // todo: add DestroyCollection()?
	ReadCollection() ([]models.Resource, *ApplicationError)
}
