package app

import "github.com/phillipmugisa/go_resume_generator/storage"

type AppServer struct {
	port    string
	storage storage.Storage
}

func NewAppServer(p string, s storage.Storage) *AppServer {
	return &AppServer{
		port:    p,
		storage: s,
	}
}

type HandlerError struct {
	message string
	code    int
}

func (e HandlerError) Error() string {
	return e.message
}

type ApiError struct{}
