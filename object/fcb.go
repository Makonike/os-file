package object

import "time"

type Fcb struct {
	isDir        bool
	Name         string         // 文件名
	StartBlock   int            // 起始盘块号
	BlockNum     int            // 占用盘块数
	ProtectTypes *[]ProtectType // 访问控制列表
	CreateTime   time.Time
	UpdateTime   time.Time
}

func NewFcb(isDir bool, name string, startBlock, blockNum int, pt *[]ProtectType, ct, ut time.Time) *Fcb {
	return &Fcb{
		isDir:        isDir,
		Name:         name,
		StartBlock:   startBlock,
		BlockNum:     blockNum,
		ProtectTypes: pt,
		CreateTime:   ct,
		UpdateTime:   ut,
	}
}

func (f *Fcb) IsDir() bool {
	return f.isDir
}
