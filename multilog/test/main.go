// SPDX-License-Identifier: MIT

package test // import "go.cryptoscope.co/margaret/multilog/test"

import (
	"testing"

	"go.cryptoscope.co/margaret/multilog"
)

type NewLogFunc func(name string, tipe interface{}, testdir string) (multilog.MultiLog, string, error)

func SinkTest(f NewLogFunc) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("Simple", SinkTestSimple(f))
	}
}

func MultiLogTest(f NewLogFunc) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("MultiSimple", MultiLogTestSimple(f))
		t.Run("Live", MultilogLiveQueryCheck(f))
	}
}

func SubLogTest(f NewLogFunc) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("Get", SubLogTestGet(f))
	}
}
