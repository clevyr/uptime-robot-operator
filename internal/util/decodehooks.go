package util

import (
	"reflect"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DecodeHookMetav1Duration(from, to reflect.Type, data any) (any, error) {
	if from.Kind() != reflect.String {
		return data, nil
	}

	var result metav1.Duration
	if to != reflect.TypeOf(result) {
		return data, nil
	}

	if from.Kind() == reflect.String {
		d, err := time.ParseDuration(data.(string))
		result.Duration = d
		return result, err
	}
	return data, nil
}
