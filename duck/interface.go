package duck

import (
	"io"
	"net/url"
)

type Result interface {
	GetNext() []string
	SetNext([]string)
	Value() interface{}
}

type Parser interface {
	Parse(*url.URL, io.Reader, []string) Result
}
