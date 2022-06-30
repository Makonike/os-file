package object

type disk struct {
	disk    [][]rune        // 磁盘，包含多个盘块
	uMap    map[string]User // 系统用户集合
	fcbList []Fcb           // contain all fcb
	dirs    []dir           // contain all tree struct dir
	bitMap  [][]int
}
