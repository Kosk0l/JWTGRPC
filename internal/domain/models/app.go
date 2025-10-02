package models

type App struct {
	ID 		int
	Name 	string
	Secret 	string // подписывать токены для дальнейшей валидации
}