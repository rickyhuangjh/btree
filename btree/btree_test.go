package btree

import (
    "testing"
	"fmt"
	"math/rand"
)

const items = 1e4

func TestSet(t *testing.T) {
	btree := NewBTree(7)
	for i:=0; i<items; i++ {
		num := rand.Intn(items)
		btree.Set(int32(num), []byte(fmt.Sprintf("%d", num)))
		if i % (items/20) == 0 {
			fmt.Printf("%v\n", float64(i)/float64(items))
		}
	}
	btree.Print()
}
