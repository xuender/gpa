package gpa

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/xuender/oils/base"
	"github.com/xuender/oils/logs"
	"google.golang.org/protobuf/proto"
)

type LevelService[T Document] struct {
	dir string
}

func NewLevelService[T Document](dir string) *LevelService[T] {
	return &LevelService[T]{dir: dir}
}

func (p *LevelService[T]) Save(docs []T) (err error) {
	if len(docs) == 0 {
		return nil
	}

	defer base.Recover(func(call error) { err = call })

	level := base.Must1(leveldb.OpenFile(p.dir, nil))
	defer level.Close()

	batch := &leveldb.Batch{}

	for _, doc := range docs {
		logs.Debugw("save", "doc", doc)
		batch.Put(base.Number2Bytes(doc.GetId()), base.Must1(proto.Marshal(doc)))
	}

	return level.Write(batch, nil)
}

func (p *LevelService[T]) Load(docs []T) (err error) {
	if len(docs) == 0 {
		return nil
	}

	defer base.Recover(func(call error) { err = call })

	level := base.Must1(leveldb.OpenFile(p.dir, nil))
	defer level.Close()

	for _, doc := range docs {
		data := base.Must1(level.Get(base.Number2Bytes(doc.GetId()), nil))
		base.Must(proto.Unmarshal(data, doc))
	}

	return
}
