package object

var memory = Memory{}

type Memory struct {
	curUser User
	curDir  dir
	bitMap  [][]int
	acFile  ActiveFile
}
