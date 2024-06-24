package skiplist

import "bytes"

// invariant assumptions:
// - if a node exists in a higher level, it also exists in every lower level
// - insert/update only (no delete)
// - a level's list is never empty if a front exists
type (
	node struct {
		key   []byte
		value []byte

		// TODO does this need to be doubly linked?
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

	// TODO use the middle return value
	location, _, found := l.locate(k)
	if found {
		location.value = v
		return
	}

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
	} else {
		location.next = &node{
			key:   k,
			value: v,
			next:  location.next,
		}
	}
	// TODO bubble up
}

func (l *list) locate(k []byte) (*node, []*node, bool) {
	var n *node
	var levels []*node

	// find where the element belongs, starting with the topmost
	// list (greatest i)
	for i := len(l.fronts) - 1; i > 0; i-- {
		if n == nil {
			// are we able to travel right from this level?
			if bytes.Compare(k, l.fronts[i].key) > 0 {
				n = l.fronts[i]
			} else {
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
		case -1: // passed where it should be, doesn't exist
			break Loop
		case 1:
			prev = n
			n = n.next
		}
	}

	return prev, levels, false
}
