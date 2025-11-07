package model

import (
	"encoding/json"
	"fmt"
)

// Value универсальное значение для метаданных
// Только одно из полей должно быть установлено (не nil)
type Value struct {
	StringValue *string  `json:"string_value,omitempty"`
	Int64Value  *int64   `json:"int64_value,omitempty"`
	DoubleValue *float64 `json:"double_value,omitempty"`
	BoolValue   *bool    `json:"bool_value,omitempty"`
}

// Конструкторы

func NewStringValue(v string) *Value {
	return &Value{StringValue: &v}
}

func NewInt64Value(v int64) *Value {
	return &Value{Int64Value: &v}
}

func NewFloat64Value(v float64) *Value {
	return &Value{DoubleValue: &v}
}

func NewBoolValue(v bool) *Value {
	return &Value{BoolValue: &v}
}

// Геттеры

func (v *Value) GetStringValue() (string, bool) {
	if v.StringValue != nil {
		return *v.StringValue, true
	}
	return "", false
}

func (v *Value) GetInt64Value() (int64, bool) {
	if v.Int64Value != nil {
		return *v.Int64Value, true
	}
	return 0, false
}

func (v *Value) GetFloat64Value() (float64, bool) {
	if v.DoubleValue != nil {
		return *v.DoubleValue, true
	}
	return 0, false
}

func (v *Value) GetBoolValue() (bool, bool) {
	if v.BoolValue != nil {
		return *v.BoolValue, true
	}
	return false, false
}

// String возвращает строковое представление значения
func (v *Value) String() string {
	if v.StringValue != nil {
		return *v.StringValue
	}
	if v.Int64Value != nil {
		return fmt.Sprintf("%d", *v.Int64Value)
	}
	if v.DoubleValue != nil {
		return fmt.Sprintf("%f", *v.DoubleValue)
	}
	if v.BoolValue != nil {
		return fmt.Sprintf("%t", *v.BoolValue)
	}
	return ""
}

// IsEmpty проверяет, установлено ли какое-либо значение
func (v *Value) IsEmpty() bool {
	return v.StringValue == nil && v.Int64Value == nil && v.DoubleValue == nil && v.BoolValue == nil
}

// MarshalJSON кастомная сериализация
// Сериализует только установленное значение
func (v *Value) MarshalJSON() ([]byte, error) {
	if v.StringValue != nil {
		return json.Marshal(*v.StringValue)
	}
	if v.Int64Value != nil {
		return json.Marshal(*v.Int64Value)
	}
	if v.DoubleValue != nil {
		return json.Marshal(*v.DoubleValue)
	}
	if v.BoolValue != nil {
		return json.Marshal(*v.BoolValue)
	}
	return []byte("null"), nil
}

// UnmarshalJSON кастомная десериализация
func (v *Value) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	// bool (проверяем первым)
	if string(data) == "true" || string(data) == "false" {
		var b bool
		if err := json.Unmarshal(data, &b); err == nil {
			v.BoolValue = &b
			return nil
		}
	}

	// number
	var num float64
	if err := json.Unmarshal(data, &num); err == nil {
		// Если целое число, сохраняем как int64
		if num == float64(int64(num)) {
			i := int64(num)
			v.Int64Value = &i
		} else {
			v.DoubleValue = &num
		}
		return nil
	}

	// string
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		v.StringValue = &str
		return nil
	}

	return fmt.Errorf("unable to unmarshal Value from %s", string(data))
}
