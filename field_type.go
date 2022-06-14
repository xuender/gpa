package gpa

// FieldType TODO.
type FieldType int

const (
	Number FieldType = iota
	Time
	Text
)

func ToFieldType(str string) FieldType {
	switch str {
	case "number", "int", "float", "int32", "int64", "float64", "float32":
		return Number
	case "time", "date", "datetime", "timestamp", "stamp", "ts":
		return Time
	default:
		return Text
	}
}
