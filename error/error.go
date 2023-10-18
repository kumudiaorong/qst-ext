package error

type Kind int

const (
	// 未知错误
	Unknown Kind = iota
	// 无效的参数
	InvalidArgument
	MaxKind
)

type Error struct {
	// 错误类型
	kind Kind
	// 错误信息
	info string
}

func (e *Error) Kind() Kind {
	return e.kind
}
func (e *Error) Info() string {
	return e.info
}
