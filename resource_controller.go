package main

type ResourceController interface {
	Path() string
	DBTableName() string

	//Create(m Resource) (Resource, error)
	//Read(id int, m Resource) (Resource, error)
	//Update(id int, m Resource) (Resource, error)
	//Destroy(id int) error
	ReadCollection(m Resource) (*[]Resource, error)
	//DestroyCollection() error
}
