package gpa

import (
	"context"
	"reflect"

	"github.com/blugelabs/bluge"
	"github.com/xuender/oils/base"
	"github.com/xuender/oils/logs"
)

// IndexService 索引服务.
type IndexService[T Document] struct {
	config bluge.Config
	fields map[int]*Field
}

// NewIndexService 新建索引服务.
func NewIndexService[T Document](dir string) *IndexService[T] { // nolint
	var doc T

	elem := reflect.TypeOf(doc).Elem()
	fields := map[int]*Field{}

	for index := 0; index < elem.NumField(); index++ {
		field := elem.Field(index)
		fieldType := field.Tag.Get("gpa")

		if fieldType != "" {
			fields[index] = NewField(fieldType, field.Name)
		}
	}

	return &IndexService[T]{
		config: bluge.DefaultConfig(dir),
		fields: fields,
	}
}

// Index 索引.
func (p *IndexService[T]) Index(docs []T) (err error) {
	if len(docs) == 0 {
		return nil
	}

	defer base.Recover(func(call error) { err = call })

	writer := base.Must1(bluge.OpenWriter(p.config))
	defer writer.Close()

	batch := bluge.NewBatch()

	for _, doc := range docs {
		bdoc := p.Parse(doc)

		batch.Update(bdoc.ID(), bdoc)
	}

	return writer.Batch(batch)
}

func (p *IndexService[T]) Parse(doc Document) *bluge.Document {
	bdoc := bluge.NewDocument(base.Itoa(doc.GetId()))
	value := reflect.ValueOf(doc).Elem()

	for index, field := range p.fields {
		if bfield := field.Field(value.Field(index)); bfield != nil {
			logs.Debugw("parse", "name", field.Name, "field", field.Type, "value", value.Field(index))
			bdoc.AddField(bfield)
		}
	}

	return bdoc
}

func (p *IndexService[T]) Query(values map[string]string) (ids []uint64, err error) {
	if len(values) == 0 {
		return nil, nil
	}

	defer base.Recover(func(call error) { err = call })

	reader := base.Must1(bluge.OpenReader(p.config))
	defer reader.Close()

	querys := make([]bluge.Query, len(values))
	index := 0

	for key, value := range values {
		querys[index] = bluge.NewMatchQuery(value).SetField(key)
		index++
	}

	query := bluge.NewBooleanQuery().AddMust(querys...)
	request := bluge.NewTopNSearch(base.Ten, query).
		WithStandardAggregations()

	iterator := base.Must1(reader.Search(context.Background(), request))
	ids = []uint64{}

	for {
		match, iterr := iterator.Next()
		if iterr != nil || match == nil {
			err = iterr

			return
		}

		base.Must(match.VisitStoredFields(func(field string, value []byte) bool {
			if field == "_id" {
				ids = append(ids, base.Must1(base.ParseInteger[uint64](string(value))))
			}

			return true
		}))
	}
}
