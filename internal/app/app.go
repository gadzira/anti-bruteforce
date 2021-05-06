package app

import (
	"fmt"

	"github.com/gadzira/anti-bruteforce/internal/domain"
)

func New(i interface{}) *domain.App {
	str := fmt.Sprintf("%v", i)
	fmt.Println(str)

	a, ok := i.(*domain.App)
	if !ok {
		return nil
	}
	return &domain.App{
		Ctx:     a.Ctx,
		Logger:  a.Logger,
		Storage: a.Storage,
		DB:      a.DB,
	}
}
