package message

import "fmt"

// Message ...
type Message string

// String ...
func (m Message) String() string {
	return string(m)
}

// Format ...
func (m Message) Format(args ...interface{}) string {
	return fmt.Sprintf(string(m), args...)
}

// Equal ...
func (m Message) Equal(src string) bool {
	return string(m) == src
}
