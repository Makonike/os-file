package object

const (
	BlockSize        = 1024
	BlockNum         = 1024
	DiskSize         = BlockSize * BlockNum
	UserStartBlock   = 0
	UserBlockNum     = 2
	FCBStartBlock    = 2
	FCBBlockNum      = 2
	DirStartBlock    = 4
	DirBlockNum      = 2
	BitMapStartBlock = 6
	BitMapBlockNum   = 1
	RecordStartBlock = 7
	BitMapRowLength  = 32
	BitMapLineLength = 32
	BitMapFree       = 0
	BitMapBusy       = 1
)

type disk struct {
	disk    [][]rune         // 磁盘，包含多个盘块
	uMap    map[string]*User // 系统用户集合
	fcbList []*Fcb           // contain all fcb
	dirs    []*dir           // contain all tree struct dir
	bitMap  [][]int
}
