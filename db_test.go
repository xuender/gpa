package gpa_test

import (
	"reflect"
	"testing"

	"github.com/xuender/gpa"
	"github.com/xuender/gpa/pb_test"
	"github.com/xuender/oils/assert"
	"github.com/xuender/oils/base"
)

func TestNewDB(t *testing.T) {
	t.Parallel()

	db, err := gpa.NewMemory[*pb_test.Book]()
	assert.NotNil(t, db)
	assert.Nil(t, err)
}

func TestDB_Load(t *testing.T) {
	t.Parallel()

	db := base.Must1(gpa.NewMemory[*pb_test.Book]())
	assert.Nil(t, db.Save([]*pb_test.Book{{Id: 1, Name: "t1"}, {Id: 2, Name: "t2"}}...))

	data := []*pb_test.Book{{Id: 1}, {Id: 2}}
	assert.Nil(t, db.Load(data...))
	assert.Equal(t, "t1", data[0].Name)
	assert.Equal(t, "t2", data[1].Name)
}

func TestDB_Query(t *testing.T) {
	t.Parallel()

	db := base.Must1(gpa.NewMemory[*pb_test.Book]())
	assert.Nil(t, db.Save([]*pb_test.Book{{Id: 1, Name: "t1"}, {Id: 2}, {Id: 3, Name: "t1"}}...))

	value := map[string]string{"Name": "t1"}
	docs, err := db.Query(value)

	assert.Nil(t, err)
	assert.Equal(t, 2, len(docs))
	assert.Equal(t, 1, docs[0].Id)
	assert.Equal(t, 3, docs[1].Id)
}

func TestDB_Match(t *testing.T) {
	t.Parallel()

	db := base.Must1(gpa.NewMemory[*pb_test.Book]())
	assert.Nil(t, db.Save([]*pb_test.Book{{Id: 1, Name: "t1"}, {Id: 2, Name: "t2"}, {Id: 3, Name: "t1"}}...))

	docs, err := db.Match("t1")
	t.Log(docs)

	assert.Equal(t, 2, len(docs))
	assert.Equal(t, 1, docs[0].Id)
	assert.Equal(t, 3, docs[1].Id)
	assert.Nil(t, err)
}

func TestNew(t *testing.T) {
	t.Parallel()

	book := &pb_test.Book{Id: 3}
	assert.Equal(t, book.Id, newBook(book).Id)
}

func TestSlice(t *testing.T) {
	t.Parallel()

	book := &pb_test.Book{Id: 3}
	bookType := reflect.TypeOf(book)
	t.Log(bookType)
	ret := reflect.New(reflect.SliceOf(bookType)).Elem()
	rvalues := make([]reflect.Value, 2)

	value := reflect.New(bookType.Elem()).Elem()
	t.Log(value)
	value.FieldByName("Id").Set(reflect.ValueOf(book.Id))
	rvalues[0] = value.Addr()
	rvalues[1] = value.Addr()

	t.Log("ret", ret)
	t.Log(rvalues)

	va := reflect.Append(ret, rvalues...)
	ret.Set(va)
}

// nolint
func newBook(book *pb_test.Book) *pb_test.Book {
	bookType := reflect.TypeOf(book).Elem()
	newBook := reflect.New(bookType).Elem()
	field := newBook.FieldByName("Id")
	field.Set(reflect.ValueOf(book.Id))

	return newBook.Addr().Interface().(*pb_test.Book)
}
