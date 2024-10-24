package metric

import (
	"errors"
	"fmt"
	"reflect"
)

type Metric struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
	val   any      `json:"-"`
}

func (m *Metric) StringValue() string {
	if m.Delta == nil && m.Value == nil {
		return fmt.Sprint(m.val)
	}
	if m.Delta == nil {
		return fmt.Sprint(*(m.Value))
	}
	return fmt.Sprint(*(m.Delta))
}

func (m *Metric) updateGauge(val any) {
	m.val = val
}

func (m *Metric) setGaugeValue() error {

	if m.Value == nil {
		temp := float64(0)
		m.Value = &temp
	}

	refValue := reflect.ValueOf(m.val)
	refGaugeValue := reflect.ValueOf(m.Value)
	refGaugeValueElem := refGaugeValue.Elem()
	refGaugeValueType := refGaugeValueElem.Type()

	if refValue.Kind() == reflect.Ptr {
		if refValue.IsNil() {
			return errors.New("the gauge metric value has not been updated: value is zero pointer")
		}
		refValue = refValue.Elem()
	}

	canConvert := refValue.CanConvert(refGaugeValueType)
	if !canConvert {
		return errors.New("can not convert value to gauge; the gauge metric value has not been updated")
	}

	value := refValue.Convert(refGaugeValueType).Float()
	*(m.Value) = value

	return nil
}

func (m *Metric) updateCounter() {
	if m.Delta == nil {
		temp := int64(0)
		m.Delta = &temp
	}
	*(m.Delta) = *(m.Delta) + 1
}
