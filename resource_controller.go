package main

type ResourceController interface {
	Path() string
	DBTableName() string

	Create(m []Resource) ([]int, *ApplicationError)
	Read(id int) (Resource, *ApplicationError) // todo: or Read(id int, m Resource) (Resource, error)?
	Update(m []Resource) *ApplicationError     // todo: return anything?
	Destroy(ids []int) *ApplicationError       // todo: or destroy collection?
	ReadCollection() ([]Resource, *ApplicationError)
}
