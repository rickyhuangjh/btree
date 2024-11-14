package btree

import (
    "testing"
)

func TestSet(t *testing.T) {
    btree := NewBTree(4)
	for i:=0; i<1e8; i++ {
		if i % 1e5 == 0 {
			t.Log(float32(i) / 1e8)
		}
        kvp := KVP{int(i)*(1e4+7)%(1e9+7), "hello"}
        // fmt.Printf("Inserting %v\n", i)
        btree.Set(&kvp)
    }
	t.Log(btree.Size())
    res := btree.Traverse()
    for i, kvp := range res {
		if i > 0 && res[i-1].key > res[i].key {
			panic("uh oh")
		}
		continue
		t.Log(kvp)
    }

}


