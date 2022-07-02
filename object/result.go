package object

const (
	SuccessMessage = "success"
	FailMessage    = "fail"
	SuccessCode    = 200
	FailCode       = 400
)

type Result struct {
	isSuccess bool
	code      int
	message   string
	data      interface{}
}

func NewResult(isSuccess bool, code int, msg string, data interface{}) *Result {
	return &Result{
		isSuccess: isSuccess,
		code:      code,
		message:   msg,
		data:      data,
	}
}

func SuccessResult() *Result {
	return SuccessData(nil)
}

func SuccessData(data interface{}) *Result {
	return NewResult(true, SuccessCode, SuccessMessage, data)
}

func SuccessMsg(msg string) *Result {
	return NewResult(true, SuccessCode, msg, nil)
}

func FailResult() *Result {
	return SuccessData(nil)
}

func FailData(data interface{}) *Result {
	return NewResult(false, FailCode, FailMessage, data)
}

func FailMsg(msg string) *Result {
	return NewResult(false, FailCode, msg, nil)
}

func (r *Result) SetData(data interface{}) {
	r.data = data
}

func (r *Result) Data() interface{} {
	return r.data
}

func (r *Result) Success() bool {
	return r.isSuccess
}

func (r *Result) Msg() string {
	return r.message
}
