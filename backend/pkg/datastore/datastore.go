package datastore

type DataStore interface {
	CreateStore(path string)
	CreateSession() string
	Add()
}
