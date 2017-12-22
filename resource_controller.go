package main

type ResourceController interface {
	Path() string
	DBTableName() string

	Create(m []Resource) ([]int, error)
	Read(id int) (Resource, error) // todo: or Read(id int, m Resource) (Resource, error)?
	Update(m []Resource) error     // todo: return anything?
	Destroy(ids []int) error       // todo: or destroy collection?
	ReadCollection() ([]Resource, *ApplicationError)
}
