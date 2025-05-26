package service

type MovieService interface {
	GetName() string
}

func GetName() string {
	return "Test"
}
