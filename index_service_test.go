package gpa_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	segment "github.com/blugelabs/bluge_segment_api"
	"github.com/lithammer/shortuuid"
	"github.com/xuender/gpa"
	"github.com/xuender/gpa/pb_test"
	"github.com/xuender/oils/assert"
)

func TestNewIndexService(t *testing.T) {
	t.Parallel()

	index := gpa.NewIndexService[*pb_test.Book](filepath.Join(os.TempDir(), shortuuid.New()))
	assert.NotNil(t, index)
}

func TestIndexService_Index(t *testing.T) {
	t.Parallel()

	index := gpa.NewIndexService[*pb_test.Book](filepath.Join(os.TempDir(), shortuuid.New()))
	assert.Nil(t, index.Index([]*pb_test.Book{{Id: 1, Name: "t1"}, {Id: 2}}))
}

func TestIndexService_Query(t *testing.T) {
	t.Parallel()

	index := gpa.NewIndexService[*pb_test.Book](filepath.Join(os.TempDir(), shortuuid.New()))
	assert.Nil(t, index.Index([]*pb_test.Book{{Id: 1, Name: "t1"}, {Id: 2}, {Id: 3, Name: "t1"}}))

	value := map[string]string{"Name": "T1"}
	ids, err := index.Query(value)

	assert.Nil(t, err)
	assert.Equal(t, 2, len(ids))
}

func TestIndexService_Parse(t *testing.T) {
	t.Parallel()

	index := gpa.NewIndexService[*pb_test.Book](filepath.Join(os.TempDir(), shortuuid.New()))
	doc := index.Parse(&pb_test.Book{Id: 3, Name: "t3", Size: 3, Created: time.Now().Unix()})
	count := 0

	doc.EachField(func(f segment.Field) {
		count++
	})

	assert.Equal(t, 4, count)
}
