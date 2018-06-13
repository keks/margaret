package multilog

import (
	"cryptoscope.co/go/librarian"
	"cryptoscope.co/go/margaret"
)

// MultiLog is a collection of logs, keyed by a librarian.Addr
// TODO maybe only call this log to avoid multilog.MultiLog?
type MultiLog interface {
	Get(librarian.Addr) (margaret.Log, error)
}