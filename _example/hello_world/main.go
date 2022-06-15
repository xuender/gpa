package main

import (
	"os"
	"path/filepath"

	"github.com/lithammer/shortuuid"
	"github.com/xuender/gpa"
	"github.com/xuender/gpa/pb_test"
	"github.com/xuender/oils/base"
	"github.com/xuender/oils/logs"
)

func main() {
	db := base.Must1(gpa.NewDBByDir[*pb_test.Book](filepath.Join(os.TempDir(), shortuuid.New())))
	defer db.Close()

	base.Must(db.Save(&pb_test.Book{Id: 1, Name: "test"}, &pb_test.Book{Id: 2, Name: "name"}))
	books := base.Must1(db.Match("test"))

	logs.Info(books[0])
}
