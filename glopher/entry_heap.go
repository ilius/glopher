package glopher

type EntryHeap []*Entry

func (h EntryHeap) Len() int           { return len(h) }
func (h EntryHeap) Less(i, j int) bool { return h[i].Word < h[j].Word }
func (h EntryHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *EntryHeap) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(*Entry))
}

func (h *EntryHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
