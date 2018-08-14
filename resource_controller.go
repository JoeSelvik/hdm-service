package main

type ResourceController interface {
	Path() string
	DBTableName() string

	Create(m []Resource) ([]int, *ApplicationError)
	Read(id int) (Resource, *ApplicationError)
	Update(m []Resource) *ApplicationError
	Destroy(ids []int) *ApplicationError // todo: add DestroyCollection()?
	ReadCollection() ([]Resource, *ApplicationError)
}
