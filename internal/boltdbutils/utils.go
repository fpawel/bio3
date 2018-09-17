package boltdbutils

import (
	"fmt"
	"github.com/boltdb/bolt"
)



type BucketSearcher struct {
	tx *bolt.Tx
	buck *bolt.Bucket
	err error
	canCreate bool
}

type Bucket struct {
	b *bolt.Bucket
}

var ErrorBucketNotExist = fmt.Errorf("%s", "Bucket not exists")

func (x Bucket) SetByte(key string, value byte) {
	if err := x.b.Put([]byte(key), []byte{value}); err != nil {
		panic(err)
	}
}

func (x Bucket) Byte(key string) byte {
	v := x.b.Get([]byte(key))
	if len(v) == 0{
		return 0
	}
	return v[0]
}

func NewBucketSearcher(tx *bolt.Tx, canCreate bool) *BucketSearcher{
	return &BucketSearcher{
		tx:tx,
		canCreate:canCreate,
	}
}

func (x *BucketSearcher) findNext(name []byte) {
	if x.canCreate {
		if x.buck == nil {
			x.buck, x.err = x.tx.CreateBucketIfNotExists(name)
		} else {
			x.buck, x.err = x.buck.CreateBucketIfNotExists(name)
		}
	} else {
		if x.buck == nil {
			x.buck = x.tx.Bucket(name)
		} else {
			x.buck = x.buck.Bucket(name)
		}
		if x.buck == nil {
			x.err = ErrorBucketNotExist
		}
	}

}
func (x *BucketSearcher) Find(path [][]byte)  {
	for _, k := range path  {
		if x.err != nil {
			return
		}
		x.findNext(k)
	}
}
func (x *BucketSearcher) Bucket()*bolt.Bucket{
	return x.buck
}

func (x *BucketSearcher) Error() error{
	return x.err
}
