package utils

type NeighbourHeap []NeighbourInfo

func (h NeighbourHeap) Len() int           { return len(h) }
func (h NeighbourHeap) Less(i, j int) bool { return h[i].Distance > h[j].Distance } // max heap
func (h NeighbourHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *NeighbourHeap) Push(x interface{}) {
	*h = append(*h, x.(NeighbourInfo))
}

func (h *NeighbourHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
