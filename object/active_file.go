package object

// ActiveFile 当前打开的文件
type ActiveFile struct {
	Fcb      *Fcb
	Record   []rune // 文件历史记录
	ReadPtr  int    // 读指针
	WritePtr int    // 写指针
}

func NewActiveFile(fcb *Fcb, record []rune, readPtr, writePtr int) *ActiveFile {
	return &ActiveFile{
		Fcb:      fcb,
		Record:   record,
		ReadPtr:  readPtr,
		WritePtr: writePtr,
	}
}
