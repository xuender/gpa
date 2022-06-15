package main

import (
	"os"
	"path/filepath"

	"github.com/lithammer/shortuuid"
	"github.com/xuender/gpa"
	"github.com/xuender/gpa/pb_test"
	"github.com/xuender/oils/base"
	"github.com/xuender/oils/dbs"
	"github.com/xuender/oils/logs"
	"github.com/xuender/oils/times"
)

func main() {
	db := base.Must1(gpa.NewDBByDir[*pb_test.Book](filepath.Join(os.TempDir(), shortuuid.New())))
	defer db.Close()

	clockAll := times.ClockStart()
	clockSave := times.ClockStart()

	for f := 0; f < 1000; f++ {
		list := make([]*pb_test.Book, 1000)

		for i := 0; i < 1000; i++ {
			list[i] = &pb_test.Book{
				Id:   dbs.ID(),
				Name: shortuuid.New(),
			}
		}

		clockSave1 := times.ClockStart()
		base.Must(db.Save(list...))
		logs.Debugw("save", "index", f*1000, "time", times.ClockStop(clockSave1))
	}

	logs.Infow("save", "time", times.ClockStop(clockSave), "size", 1000*1000)

	name := shortuuid.New()
	db.Save(&pb_test.Book{Id: dbs.ID(), Name: name})

	clockMatch := times.ClockStart()

	db.Match(name)

	logs.Infow("match", "time", times.ClockStop(clockMatch))
	logs.Infow("all", "time", times.ClockStop(clockAll))
}
