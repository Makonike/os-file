package object

const (
	PathSeparator = "/"
	BackToPre     = ".."
)

type Dir struct {
	Fcb      *Fcb
	Index    int    // 当前索引
	Children []*Dir // 子目录项集合
	ParIndex int    // 父目录项位置
}

func NewDir(fcb *Fcb, index, parIndex int) *Dir {
	return &Dir{
		Fcb:      fcb,
		Index:    index,
		Children: make([]*Dir, 0),
		ParIndex: parIndex,
	}
}

func (d *Dir) IIndex() int {
	for i, v := range DCache.Disk.Dirs {
		if v == d {
			return i
		}
	}
	return 0
}
