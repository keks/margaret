package margaret // import "cryptoscope.co/go/margaret"

import (
	"cryptoscope.co/go/luigi"
)

type Seq int64

const (
	SeqEmpty Seq = -1
)

type Log interface {
	Seq() luigi.Observable
	Get(Seq) (interface{}, error)
	Query(...QuerySpec) (luigi.Source, error)
	Append(interface{}) error
}

type oob struct{}

var OOB oob

func (_ oob) Error() string {
	return "out of bounds"
}

// IsOutOfBounds returns whether a particular error is an out-of-bounds error
func IsOutOfBounds(err error) bool {
	_, ok := err.(oob)
	return ok
}
