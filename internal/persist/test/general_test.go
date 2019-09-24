package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"go.cryptoscope.co/margaret/internal/persist"
	"go.cryptoscope.co/margaret/internal/persist/fs"
	"go.cryptoscope.co/margaret/internal/persist/mkv"
	"go.cryptoscope.co/margaret/internal/persist/sqlite"

	_ "github.com/mattn/go-sqlite3"
)

func SimpleSaver(mk func(*testing.T) persist.Saver) func(*testing.T) {

	return func(t *testing.T) {
		p := mk(t)

		r := require.New(t)

		l, err := p.List()
		r.NoError(err)
		r.Len(l, 0, "%v", l)

		k := persist.Key{0, 0, 0, 1}
		d, err := p.Get(k)
		r.EqualError(err, persist.ErrNotFound.Error())
		r.Nil(d)

		testData := []byte("fooo")

		err = p.Put(k, testData)
		r.NoError(err)

		l, err = p.List()
		r.NoError(err)
		r.Len(l, 1)
		r.Equal(k, l[0])

		d, err = p.Get(k)
		r.NoError(err)
		r.Equal(d, testData)

		r.NoError(p.Close())
	}
}

func TestSaver(t *testing.T) {
	t.Run("fs", SimpleSaver(makeFS))
	t.Run("sqlite", SimpleSaver(makeSqlite))
	t.Run("kv", SimpleSaver(makeSqlite))

}

func makeFS(t *testing.T) persist.Saver {
	base := filepath.Join("testrun", t.Name())
	os.RemoveAll(base)
	return fs.New(base)
}

func makeSqlite(t *testing.T) persist.Saver {
	base := filepath.Join("testrun", t.Name())
	os.RemoveAll(base)
	s, err := sqlite.New(base)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func makeMKV(t *testing.T) persist.Saver {
	base := filepath.Join("testrun", t.Name())
	os.RemoveAll(base)
	s, err := mkv.New(base)
	if err != nil {
		t.Fatal(err)
	}
	return s
}