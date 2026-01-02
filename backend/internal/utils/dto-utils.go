package utils

import "github.com/jinzhu/copier"

func EntityToDto[T any](entity interface{}) (*T, error) {
	var dto T
	if err := copier.Copy(&dto, entity); err != nil {
		return nil, err
	}
	return &dto, nil
}

func EntitiesToDto[T any](entities interface{}) ([]*T, error) {
	var dtos []*T
	if err := copier.Copy(&dtos, entities); err != nil {
		return nil, err
	}
	return dtos, nil
}
