package badger

import (
	"encoding/binary"
	"reflect"

	"cryptoscope.co/go/luigi"
	"cryptoscope.co/go/margaret"
	"github.com/dgraph-io/badger"
	"github.com/pkg/errors"
)

type sublog struct {
	mlog   *mlog
	prefix []byte
	seq    luigi.Observable
}

func (log *sublog) Seq() luigi.Observable {
	return log.seq
}

func (log *sublog) Get(seq margaret.Seq) (interface{}, error) {
	t := reflect.TypeOf(log.mlog.tipe)
	v := reflect.New(t).Interface()

	seqBs := make([]byte, 8)
	binary.BigEndian.PutUint64(seqBs, uint64(seq))
	key := append(log.prefix, seqBs...)

	err := log.mlog.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return errors.Wrap(err, "error getting item")
		}

		data, err := item.Value()
		if err != nil {
			return errors.Wrap(err, "error getting value")
		}

		v, err = log.mlog.codec.Unmarshal(data)
		return errors.Wrap(err, "error unmarshaling data")
	})

	if err != nil && errors.Cause(err) != badger.ErrKeyNotFound {
		return nil, errors.Wrap(err, "error in badger transaction (view)")
	}

	return v, nil
}

func (log *sublog) Query(specs ...margaret.QuerySpec) (luigi.Source, error) {
	qry := &query{
		log: log,

		lt:      margaret.SeqEmpty,
		nextSeq: margaret.SeqEmpty,

		limit: -1, //i.e. no limit
	}

	for _, spec := range specs {
		err := spec(qry)
		if err != nil {
			return nil, err
		}
	}

	return qry, nil
}

func (log *sublog) Append(v interface{}) (margaret.Seq, error) {
	var seq margaret.Seq

	data, err := log.mlog.codec.Marshal(v)
	if err != nil {
		return -2, errors.Wrap(err, "error marshaling value")
	}

	err = log.mlog.db.Update(func(txn *badger.Txn) error {
		seqIface, err := log.seq.Value()
		if err != nil {
			return errors.Wrap(err, "error getting value from seq observable")
		}

		seq = seqIface.(margaret.Seq) + 1
		seqBs := make([]byte, 8)
		binary.BigEndian.PutUint64(seqBs, uint64(seq))
		key := append(log.prefix, seqBs...)

		err = txn.Set(key, data)
		if err != nil {
			return errors.Wrap(err, "errors setting value")
		}

		log.seq.Set(seq)
		return nil
	})
	if err != nil {
		return -2, errors.Wrap(err, "error in write transaction")
	}

	return seq, nil
}