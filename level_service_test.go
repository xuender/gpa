package gpa_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/lithammer/shortuuid"
	"github.com/xuender/gpa"
	"github.com/xuender/gpa/pb_test"
	"github.com/xuender/oils/assert"
)

func TestLevelService_Save(t *testing.T) {
	t.Parallel()

	level := gpa.NewLevelService[*pb_test.Book](filepath.Join(os.TempDir(), shortuuid.New()))
	assert.Nil(t, level.Save([]*pb_test.Book{{Id: 1}, {Id: 2}}))
}

func TestLevelService_Load(t *testing.T) {
	t.Parallel()

	level := gpa.NewLevelService[*pb_test.Book](filepath.Join(os.TempDir(), shortuuid.New()))
	assert.Nil(t, level.Save([]*pb_test.Book{{Id: 1, Name: "t1"}, {Id: 2, Name: "t2"}}))

	data := []*pb_test.Book{{Id: 1}, {Id: 2}}
	assert.Nil(t, level.Load(data))
	assert.Equal(t, "t1", data[0].Name)
	assert.Equal(t, "t2", data[1].Name)
}
