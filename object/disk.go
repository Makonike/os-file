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
	Disk    [][]rune         // 磁盘，包含多个盘块
	UMap    map[string]*User // 系统用户集合
	FcbList []*Fcb           // contain all fcb
	Dirs    []*Dir           // contain all tree struct dir
	BitMap  [][]int
}
