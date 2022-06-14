package gpa

import (
	"reflect"

	"github.com/blugelabs/bluge"
	"github.com/xuender/oils/times"
)

type Field struct {
	Type string
	Name string
}

func NewField(fieldType, fieldName string) *Field {
	return &Field{Type: fieldType, Name: fieldName}
}

func (p *Field) Field(value reflect.Value) bluge.Field {
	switch p.Type {
	case "number", "int", "float", "int32", "int64", "float64", "float32":
		return bluge.NewNumericField(p.Name, float64(value.Int()))
	case "time", "date", "datetime", "timestamp", "stamp", "ts":
		if i64 := value.Int(); i64 != 0 {
			return bluge.NewDateTimeField(p.Name, times.ParseNumber(i64))
		}
	default:
		return bluge.NewTextField(p.Name, value.String())
	}

	return nil
}
