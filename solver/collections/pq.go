package collections

import (
	"container/heap"
)

type OrderedItem interface {
	Priority() int
	UpdatePriority(value int)
	GetIndex() int
	SetIndex(value int)
}

type PriorityQueue []OrderedItem

func (q PriorityQueue) Len() int { return len(q) }

func (q PriorityQueue) Empty() bool { return q.Len() == 0 }

func (q PriorityQueue) Less(i, j int) bool { return q[i].Priority() > q[j].Priority() }

func (q PriorityQueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].SetIndex(i)
	q[j].SetIndex(j)
}

func (q *PriorityQueue) Push(x interface{}) {
	n := len(*q)
	v := x.(OrderedItem)
	v.SetIndex(n)
	*q = append(*q, v)
}

func (q *PriorityQueue) Pop() interface{} {
	old := *q
	n := len(old)
	v := old[n-1]
	v.SetIndex(-1)
	*q = old[0 : n-1]
	return v
}

// update modifies the priority and value of an Item in the queue.
func (q *PriorityQueue) Update(v OrderedItem, newPriority int) {
	v.UpdatePriority(newPriority)
	heap.Fix(q, v.GetIndex())
}
