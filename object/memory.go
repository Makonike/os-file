package object

var MMemory = &Memory{}

type Memory struct {
	CurUser *User
	CurDir  *Dir
	BitMap  [][]int
	AcFile  *ActiveFile
}
