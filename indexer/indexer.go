package indexer

import (
	"fmt"
	"github.com/midoks/hammer/configure"
	"github.com/midoks/hammer/cron"
	"github.com/midoks/hammer/ds"
)

func Run(cf *configure.Args) {

	ods := ds.OpenDS(cf)

	/* 导入全量数据 */
	go ods.Import()

	/* 执行任务 */
	go ods.Task()

	if cf.Interval != "" {
		cronExport := fmt.Sprintf("@every %s", cf.Interval)

		cron.Add(cronExport, func() {
			/* 倒入增量数据 */
			ods.DeltaData()
			/* 删除失效数据 */
			ods.DeleteData()
		})
	}
}
