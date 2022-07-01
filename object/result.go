package object

const (
	OperateSuccessMessage = "success"
	OperateFailMessage    = "fail"
	SuccessCode           = 200
	FailCode              = 400
)

type result struct {
	isSuccess bool
	code      int
	message   string
	data      interface{}
}
