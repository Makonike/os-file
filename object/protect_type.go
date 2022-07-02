package object

var (
	CanRead = &ProtectType{
		value: 1,
		desc:  "可读",
	}
	CanWrite = &ProtectType{
		value: 2,
		desc:  "可读",
	}
	CanExecute = &ProtectType{
		value: 3,
		desc:  "可读",
	}
)

type ProtectType struct {
	value int    // 状态码
	desc  string // 描述
}

func All() []*ProtectType {
	return []*ProtectType{
		CanRead,
		CanWrite,
		CanExecute,
	}
}
