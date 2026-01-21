package utils

import (
	"reflect"

	"github.com/jinzhu/copier"
)

func EntityToDto[T any](entity interface{}, options ...copier.Option) (T, error) {
	var dto T

	if reflect.TypeOf(dto) != nil && reflect.TypeOf(dto).Kind() == reflect.Ptr {
		val := reflect.New(reflect.TypeOf(dto).Elem())
		dto = val.Interface().(T)
	}

	finalOpt := copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
	}

	if len(options) > 0 {
		finalOpt = options[0]
	}

	if err := copier.CopyWithOption(&dto, entity, finalOpt); err != nil {
		return dto, err
	}
	return dto, nil
}

func EntitiesToDto[T any](entities interface{}, options ...copier.Option) ([]T, error) {
	dtos := []T{}

	finalOpt := copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
	}

	if len(options) > 0 {
		finalOpt = options[0]
	}

	if err := copier.CopyWithOption(&dtos, entities, finalOpt); err != nil {
		return dtos, err
	}
	return dtos, nil
}
