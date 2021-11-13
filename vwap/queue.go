package vwap

// WindowQueue is a fixed size and allocated upfront queue of cached
// data points specialized in the service of the VWAP algorithm needs.
type WindowQueue struct {
	content  []*vwapCache
	readHead uint16
	writeHead uint16
	len uint16
	size uint16
}

func NewWindowQueue(size uint16) *WindowQueue{
	return &WindowQueue{
		// allocate it upfront
		content:   make([]*vwapCache, size),
		size: size,
	}
}

func (q *WindowQueue) Peek(pos uint16) (*vwapCache, bool) {
	if pos >= q.size || q.content[pos] == nil {
		return nil, false
	}

	return q.content[pos], true
}

func (q *WindowQueue) PeekLast()(*vwapCache, bool){
	return q.Peek(q.Last())
}

func (q *WindowQueue) Last()uint16{
	if q.len == 0 && q.readHead == 0 && q.writeHead == 0 {
		return 0
	}
	if q.writeHead == 0 {
		return q.size-1
	}

	return q.writeHead-1
}

func (q *WindowQueue) Pop() (*vwapCache, bool) {
	if q.len <= 0 {
		return nil, false
	}

	result := q.content[q.readHead]
	q.content[q.readHead] = nil
	q.readHead = (q.readHead + 1) % q.size
	q.len--

	return result, true
}

func (q *WindowQueue) Push(v *vwapCache) bool {
	if q.len >= q.size {
		return false
	}

	q.content[q.writeHead] = v
	q.writeHead = (q.writeHead + 1) % q.size
	q.len++

	return true
}