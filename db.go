package gpa

import (
	"os"
	"path/filepath"
	"reflect"

	"github.com/lithammer/shortuuid"
	"github.com/xuender/gpa/pb"
	"github.com/xuender/oils/base"
	"github.com/xuender/oils/oss"
)

// DB 数据对象.
type DB[T Document] struct {
	ds *LevelService[T]
	is *IndexService[T]
}

// NewDB 新建数据对象.
func NewDB[T Document](config *pb.Config) (*DB[T], error) {
	return NewDBByDir[T](config.Dir)
}

func NewDBByDir[T Document](dir string) (db *DB[T], err error) { // nolint
	defer base.Recover(func(call error) { err = call })

	docDir := base.Must1(oss.Abs(filepath.Join(dir, "doc")))
	base.Must(os.MkdirAll(docDir, oss.DefaultDirFileMod))

	indexDir := base.Must1(oss.Abs(filepath.Join(dir, "index")))
	base.Must(os.MkdirAll(indexDir, oss.DefaultDirFileMod))

	return &DB[T]{
		ds: NewLevelService[T](docDir),
		is: NewIndexService[T](indexDir),
	}, nil
}

func NewTest[T Document]() (*DB[T], error) {
	return NewDBByDir[T](filepath.Join(os.TempDir(), shortuuid.New()))
}

func (p *DB[T]) Save(docs ...T) (err error) {
	if err := p.ds.Save(docs); err != nil {
		return err
	}

	return p.is.Index(docs)
}

func (p *DB[T]) Load(docs ...T) (err error) {
	return p.ds.Load(docs)
}

func (p *DB[T]) Query(values map[string]string) (docs []T, err error) {
	ids, err := p.is.Query(values)
	if err != nil {
		return nil, err
	}

	defer base.Recover(func(call error) { err = call })

	var doc T

	docType := reflect.TypeOf(doc)
	slice := reflect.New(reflect.SliceOf(docType)).Elem()
	sliceValue := make([]reflect.Value, len(ids))

	for index, id := range ids {
		value := reflect.New(docType.Elem()).Elem()
		value.FieldByName("Id").Set(reflect.ValueOf(id))
		sliceValue[index] = value.Addr()
	}

	va := reflect.Append(slice, sliceValue...)
	slice.Set(va)
	// nolint
	return slice.Interface().([]T), err
}
