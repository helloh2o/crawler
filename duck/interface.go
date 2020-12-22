package duck

import (
	"io"
	"net/url"
)

type Result interface {
	Value() Result
}

type Parser interface {
	Parse(*url.URL, io.Reader, []string, func(string)) Result
}
