package usecase

import "fmt"

type Server struct {
}

func NewServerUC() *Server {
	return &Server{}
}

func (s *Server) RegisterUser(name string) {
	fmt.Println("register success: ", name)
}
