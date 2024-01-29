// Package cerrors is package for custom errors
package cerrors

// custom dynamodb exceptions
const (
	ConditionalCheckFailedException = "ccf"
	ResourceNotFoundException       = "rnf"
)

type CustomError struct {
	Message  string
	TypesMap map[string]bool
}

func NewCustomError(message string, errorType string) error {
	typesMap := map[string]bool{
		ConditionalCheckFailedException: false,
		ResourceNotFoundException:       false,
	}

	if _, ok := typesMap[errorType]; ok {
		typesMap[errorType] = true
	}

	return &CustomError{
		Message:  message,
		TypesMap: typesMap,
	}
}

func (ce *CustomError) Error() string {
	return ce.Message
}
