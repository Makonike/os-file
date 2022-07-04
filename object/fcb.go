package object

import (
	"fmt"
	"strconv"
	"time"
)

type Fcb struct {
	isDir        bool
	Name         string         // 文件名
	StartBlock   int            // 起始盘块号
	BlockNum     int            // 占用盘块数
	ProtectTypes []*ProtectType // 访问控制列表
	CreateTime   time.Time
	UpdateTime   time.Time
}

func NewFcb(isDir bool, name string, startBlock, blockNum int, pt []*ProtectType, ct, ut time.Time) *Fcb {
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

func (f *Fcb) String() string {
	return fmt.Sprintf("%20s %10s %20s %20s", f.Name, strconv.FormatInt(int64(f.BlockNum), 10), f.CreateTime.Format(time.RFC822), f.UpdateTime.Format(time.RFC822))
}
