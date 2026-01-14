package utils

import (
	"reflect"

	"github.com/jinzhu/copier"
)

func EntityToDto[T any](entity interface{}) (T, error) {
	var dto T
	// If T is a pointer type -> must initialize it before copying
	if reflect.TypeOf(dto).Kind() == reflect.Ptr {
		val := reflect.New(reflect.TypeOf(dto).Elem())
		// val is Value of type *Dto.
		dto = val.Interface().(T)
	}

	if err := copier.CopyWithOption(&dto, entity, copier.Option{IgnoreEmpty: true}); err != nil {
		return dto, err
	}
	return dto, nil
}

func EntitiesToDto[T any](entities interface{}) ([]T, error) {
	dtos := []T{}
	if err := copier.CopyWithOption(&dtos, entities, copier.Option{IgnoreEmpty: true}); err != nil {
		return dtos, err
	}
	return dtos, nil
}
