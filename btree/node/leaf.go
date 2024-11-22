package node

import (
	"fmt"
	"strings"
	"btree/utils"
)



func NewLeafNode(order uint32) *LeafNode {
	return &LeafNode{
		Order:  order,
		Parent: nil,
		Next:   nil,
		Prev:   nil,
		Keys:   make([]int32, 0, order),
		Vals:   make([][]byte, 0, order),
	}
}


func (n *LeafNode) getParent() *InternalNode {
	return n.Parent
}

func (n *LeafNode) setParent(parent *InternalNode) error {
	n.Parent = parent
	return nil
}

// returns the index of where the KVP would go
func (n *LeafNode) find(key int32) (int, bool) {
	for i, curKey := range n.Keys {
		if key == curKey {
			return i, true
		} else if key < curKey {
			return i, false
		}
	}
	return len(n.Keys), false
}

func (n *LeafNode) get(key int32) ([]byte, bool) {
	idx, wasFound := n.find(key)
	if wasFound {
		return n.Vals[idx], true
	}
	var zero []byte
	return zero, false
}

func (n *LeafNode) getRange(start, end int32, res [][]byte) [][]byte {
	startIdx, _ := n.find(start)
	endIdx, _ := n.find(end)
	res = append(res, n.Vals[startIdx:endIdx]...)
	if n.Next != nil && endIdx == len(n.Keys) {
		return n.Next.getRange(start, end, res)
	}
	return res
}

func (n *LeafNode) traverse(res [][]byte) [][]byte {
	res = append(res, n.Vals...)
	if n.Next == nil {
		return res
	}
	return n.Next.traverse(res)
}

func (n *LeafNode) set(key int32, val []byte) (bool, error) {
	idx, wasFound := n.find(key)
	if wasFound {
		n.Vals[idx] = val
	} else {
		n.Keys = utils.Insert(n.Keys, idx, key)
		n.Vals = utils.Insert(n.Vals, idx, val)
	}
	return !wasFound, n.split()
}

func (n *LeafNode) split() error {
	if len(n.Keys) < int(n.Order) {
		return nil
	}

	siblingNode := NewLeafNode(n.Order)
	if n.Parent == nil {
		n.Parent = newInternalNode(n.Order)
		n.Parent.insertChild(0, n)
	}
	siblingNode.Parent = n.Parent

	mid := len(n.Keys) / 2

	idx, _ := n.Parent.find(n.Keys[mid])
	n.Parent.Keys = utils.Insert(n.Parent.Keys, idx, n.Keys[mid])
	n.Parent.insertChild(idx+1, siblingNode)

	siblingNode.Keys = make([]int32, len(n.Keys[mid:]), n.Order)
	copy(siblingNode.Keys, n.Keys[mid:])
	n.Keys = n.Keys[:mid]

	siblingNode.Vals = make([][]byte, len(n.Vals[mid:]), n.Order)
	copy(siblingNode.Vals, n.Vals[mid:])
	n.Vals = n.Vals[:mid]

	siblingNode.Next = n.Next
	if n.Next != nil {
		n.Next.Prev = siblingNode
	}
	n.Next = siblingNode
	siblingNode.Prev = n
	return n.Parent.split()
}


func (n *LeafNode) delete(key int32) (bool, error) {
	idx, wasFound := n.find(key)
	if !wasFound {
		return false, nil
	}
	old := n.Keys[0]
	n.Keys = utils.Delete(n.Keys, idx)
	n.Vals = utils.Delete(n.Vals, idx)
	
	err := n.merge()
	if n.Parent != nil && len(n.Keys) > 0 {
		n.Parent.replaceKey(old, n.Keys[0])
	}
	return true, err
}

func (n *LeafNode) merge() error {
	if n.Parent == nil || len(n.Keys) >= int(n.Order-1)/2 {
		return nil
	}
	if n.Next != nil && n.Next.Parent == n.Parent &&
	len(n.Next.Keys) > int(n.Order-1)/2 {
		// steal from next
		n.Parent.replaceKey(n.Next.Keys[0], n.Next.Keys[1])
		n.Keys = append(n.Keys, n.Next.Keys[0])
		n.Vals = append(n.Vals, n.Next.Vals[0])
		n.Next.Keys = n.Next.Keys[1:]
		n.Next.Vals = n.Next.Vals[1:]
	} else if n.Prev != nil && n.Prev.Parent == n.Parent &&
	len(n.Prev.Keys) > int(n.Order-1)/2 {
		// steal from prev
		n.Keys = utils.Insert(n.Keys, 0, n.Prev.Keys[len(n.Prev.Keys)-1])
		n.Vals = utils.Insert(n.Vals, 0, n.Prev.Vals[len(n.Prev.Vals)-1])
		n.Prev.Keys = n.Prev.Keys[:len(n.Prev.Keys)-1]
		n.Prev.Vals = n.Prev.Vals[:len(n.Prev.Vals)-1]
	} else if n.Prev != nil && n.Prev.Parent == n.Parent {
		// merge with prev
		n.Prev.Keys = append(n.Prev.Keys, n.Keys...)
		n.Prev.Vals = append(n.Prev.Vals, n.Vals...)
		ourIdx, _ := n.Parent.findChildIdx(n)
		n.Parent.Keys = utils.Delete(n.Parent.Keys, ourIdx-1)
		n.Parent.Children = utils.Delete(n.Parent.Children, ourIdx)
		n.Prev.Next = n.Next
		if n.Next != nil {
			n.Next.Prev = n.Prev
		}
	} else {
		// merge with next
		n.Keys = append(n.Keys, n.Next.Keys...)
		n.Vals = append(n.Vals, n.Next.Vals...)
		ourIdx, _ := n.Parent.findChildIdx(n)
		n.Parent.Keys = utils.Delete(n.Parent.Keys, ourIdx)
		n.Parent.Children = utils.Delete(n.Parent.Children, ourIdx+1)
		if n.Next.Next != nil {
			n.Next.Next.Prev = n
		}
		n.Next = n.Next.Next
	}
	return n.Parent.merge()
}

func (n *LeafNode) getNewRoot() BTreeNode {
	if len(n.Keys) == 0 {
		return nil
	}
	return n
}

func (n *LeafNode) print(level int) {
    indent := strings.Repeat("    ", level)
	next := n.Next
	prev := n.Prev
	nextKeys := []int32{}
	if next != nil {
		nextKeys = next.Keys
	}
	prevKeys := []int32{}
	if prev != nil {
		prevKeys = prev.Keys
	}
    fmt.Printf("%sLeaf Node: keys=%v, next=%v, prev=%v\n", indent, n.Keys, nextKeys, prevKeys)
}

func (n *LeafNode) verify() (int32, int32) {
	return n.Keys[0], n.Keys[len(n.Keys)-1]
}

