package object

type dir struct {
	fcb      Fcb
	index    int   // 当前索引
	children []dir // 子目录项集合
	parIndex int   // 父目录项位置
}
