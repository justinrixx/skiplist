package skiplist

import (
	"bytes"
	"fmt"
	"math/rand"
)

// invariant assumptions:
// - if a node exists in a higher level, it also exists in every lower level
// - insert/update only (no delete)
// - a level's list is never empty if a front exists
type (
	node struct {
		key   []byte
		value []byte

		next  *node
		below *node
	}

	list struct {
		fronts []*node
	}
)

func (l *list) find(k []byte) ([]byte, bool) {
	if len(l.fronts) < 1 {
		return nil, false
	}

	n, _, found := l.locate(k)
	if found && n != nil {
		return n.value, true
	}
	return nil, false
}

func (l *list) insert(k []byte, v []byte) {
	// list started empty
	if len(l.fronts) < 1 {
		l.fronts = []*node{
			{
				key:   k,
				value: v,
			},
		}
		return
	}

	location, levels, found := l.locate(k)
	if found {
		location.value = v
		return
	}

	var newNode *node
	// nil location means we need to insert before the current head
	if location == nil {
		currFront := l.fronts[0]

		currFront.next = &node{
			key:   currFront.key,
			value: currFront.value,
			next:  currFront.next,
		}
		currFront.key = k
		currFront.value = v

		newNode = currFront
	} else {
		newNode = &node{
			key:   k,
			value: v,
			next:  location.next,
		}
		location.next = newNode
	}

	// add a new level as long as we get heads
	i := len(levels) - 1
	belowNode := newNode
	flip := rand.Intn(2)
	for flip == 0 {
		n := &node{
			key:   k,
			below: belowNode,
		}

		// link into the existing list if appropriate
		if i >= 0 {
			prev := levels[i]
			if prev != nil {
				n.next = prev.next
				prev.next = n
			} else {
				// levels is in reversed order from fronts. also levels does not include the
				// 0th level.
				//
				// fronts: [0, 1, 2, 3, 4]
				// levels: [4, 3, 2, 1]
				//
				// ex: level 4 is at levels[0] but sholud map to fronts[4]
				// len(fronts)-i-1 = 5-0-1 = 4
				//
				// ex level 2 is at levels[2] but should map to fronts[2]
				// len(fronts)-i-1 = 5-2-1 = 2
				n.next = l.fronts[len(l.fronts)-i-1]
			}
		} else {
			// level doesn't exist, create one
			l.fronts = append(l.fronts, n)
		}

		belowNode = n
		flip = rand.Intn(2)
		i--
	}
}

func (l *list) locate(k []byte) (*node, []*node, bool) {
	var n *node
	var levels []*node

	// find where the element belongs, starting with the topmost
	// (sparsest) list
	for i := len(l.fronts) - 1; i > 0; i-- {
		if n == nil {
			// are we able to travel right from this level?
			if bytes.Compare(k, l.fronts[i].key) > 0 {
				n = l.fronts[i]
			} else {
				levels = append(levels, nil)
				continue
			}
		}

		// travel as far right as possible
		for n.next != nil && bytes.Compare(k, n.next.key) > 0 {
			n = n.next
		}

		levels = append(levels, n)
		n = n.below
	}

	// may not have been found in higher levels
	if n == nil {
		n = l.fronts[0]
	}

	// root list
	var prev *node
Loop:
	for n != nil {
		switch bytes.Compare(k, n.key) {
		case 0: // equal
			return n, nil, true
		case -1: // we've passed where it should be, doesn't exist
			break Loop
		case 1:
			prev = n
			n = n.next
		}
	}

	return prev, levels, false
}

func (l *list) print() {
	for i, front := range l.fronts {
		fmt.Printf("\nlevel %d\n", i)
		for front != nil {
			fmt.Printf("%s ", string(front.key))

			front = front.next
		}
	}
	fmt.Println()
}
