package biological_driven_architecture

type Runtime interface {
	Init() Error
	Run() Error
	Stop() Error
	HandleError(Error) Error
}

type Builder[T any] interface {
	Spawn() (T, Error)
}

type ErrorType string
type ErrorLevel int

const (
	Warning ErrorLevel = iota
	Fatal
)

type ErrorStruct struct {
	Type      ErrorType
	Level     ErrorLevel
	Message   string
	SubErrors []Error
}

type Error *ErrorStruct

func NewError(errorType ErrorType, msg string, subErrors []Error) Error {
	return &ErrorStruct{
		Type:      errorType,
		Message:   msg,
		SubErrors: subErrors,
	}
}
