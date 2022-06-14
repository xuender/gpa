package gpa

import (
	"reflect"

	"github.com/blugelabs/bluge"
	"github.com/xuender/oils/times"
)

type Field struct {
	Type FieldType
	Name string
}

func NewField(fieldType, fieldName string) *Field {
	return &Field{Type: ToFieldType(fieldType), Name: fieldName}
}

func (p *Field) Field(value reflect.Value) bluge.Field {
	switch p.Type {
	case Number:
		return bluge.NewNumericField(p.Name, float64(value.Int()))
	case Time:
		if i64 := value.Int(); i64 != 0 {
			return bluge.NewDateTimeField(p.Name, times.ParseNumber(i64))
		}
	case Text:
		return bluge.NewTextField(p.Name, value.String())
	default:
		return bluge.NewTextField(p.Name, value.String())
	}

	return nil
}
