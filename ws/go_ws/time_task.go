package go_ws

import (
	"github.com/robfig/cron/v3"
)

func CleanOfflineConn() {

	c := cron.New()

	// 每天定时执行的条件
	spec := `* * * * *`

	c.AddFunc(spec, func() {
		// fmt.Println("CleanOfflineConn")
		HandelOfflineCoon()
	})

	go c.Start()
}
