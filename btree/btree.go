package btree

import (
	"btree/node"
	"btree/serialization"
	"bytes"
	"encoding/gob"
	"fmt"
	"errors"
)


type BTree interface {
	Set(key int32, val []byte) error
	Delete(key int32) error
	Get(key int32) ([]byte, bool)
	GetRange(start, end int32) [][]byte
	Traverse() [][]byte
	Print()
	Verify()
}

type btree struct {
	Order uint32 // branching factor (max number of children per node)
	Root  node.BTreeNode
	Size  int
}

func NewBTree(order uint32) BTree {
	return &btree{
		Order: order,
		Root:  nil,
		Size:  0,
	}
}


func Deserialize(data []byte) (BTree, error) {
	b := new(btree)
	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(b)
	return b, err
}


func (b *btree) Get(key int32) ([]byte, bool) {
	return b.Root.get(key)
}

func (b *btree) Traverse() [][]byte {
	if b.Root == nil {
		return nil
	}
	res := b.Root.traverse(make([][]byte, 0, b.Size))
	return res
}

func (b *btree) Set(key int32, val []byte) error {
	if int(b.Order) * (serialization.KeySize + len(val)) +
	serialization.LeafHeaderSize > serialization.PageSize {
		return errors.New("value too large")
	}

	if b.Root == nil {
		b.Root = NewLeafNode(b.Order)
	}

	wasAdded, err := b.Root.set(key, val)
	if b.Root.getParent() != nil {
		// fmt.Println("hello")
		b.Root = b.Root.getParent()
	}
	if wasAdded {
		b.Size++
	}
	//fmt.Println(a)
	return err
}

func (b *btree) Delete(key int32) error {
	if b.Root == nil {
		return nil
	}
	wasDeleted, err := b.Root.delete(key)
	if wasDeleted {
		b.Size--
	}
	b.Root = b.Root.getNewRoot()
	if b.Root != nil {
		b.Root.setParent(nil)
	}
	return err
}

func (b *btree) Print() {
    if b.Root == nil {
        fmt.Println("Empty tree")
        return
    }
    b.Root.print(0)
}

func (b *btree) Verify() {
	if b.Root == nil {
		return
	}
	b.Root.verify()
}

func (b *btree) GetRange(start, end int32) [][]byte {
	res := make([][]byte, 0)
	res = b.Root.getRange(start, end, res)
	return res
}
