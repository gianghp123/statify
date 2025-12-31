package utils

import "github.com/jinzhu/copier"

func EntityToDto[T any](entity interface{}) (*T, error) {
	var dto T
	if err := copier.Copy(&dto, entity); err != nil {
		return nil, err
	}
	return &dto, nil
}
