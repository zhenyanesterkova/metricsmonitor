package metric

import (
	"errors"
	"fmt"
	"reflect"
)

type metric struct {
	name       string
	metricType string
	value      any
}

func (m *metric) StringValue() string {
	return fmt.Sprint(m.value)
}

func (m *metric) Type() string {
	return m.metricType
}

func (m *metric) Name() string {
	return m.name
}

func (m *metric) update(val any) {
	m.value = val
}
func (m *metric) updateCounter() error {

	refValue := reflect.ValueOf(m.value)
	refTypeInt64 := reflect.TypeOf(reflect.Int64)

	canConvertToInt64 := refValue.CanConvert(refTypeInt64)
	if !canConvertToInt64 {
		return errors.New("can not convert value of metric PollCount")
	}

	value := refValue.Int()
	m.value = value + 1

	return nil
}
