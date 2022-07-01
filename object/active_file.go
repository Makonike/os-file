package object

// ActiveFile 当前打开的文件
type ActiveFile struct {
	fcb      *Fcb
	record   []rune // 文件历史记录
	readPtr  int    // 读指针
	writePtr int    // 写指针
}
