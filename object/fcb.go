package object

import "time"

type Fcb struct {
	isDir        bool
	name         string         // 文件名
	startBlock   int            // 起始盘块号
	blockNum     int            // 占用盘块数
	protectTypes *[]ProtectType // 访问控制列表
	createTime   time.Time
	updateTime   time.Time
}
