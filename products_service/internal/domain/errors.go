package domain

import "errors"

var (
	ErrProductNotFound      = errors.New("товар не найден")
	ErrInsufficientQuantity = errors.New("недостаточное количество товара")
	ErrReservationNotFound  = errors.New("резервирование не найдено")
	ErrInvalidStatus        = errors.New("недопустимый статус резервирования")
)
