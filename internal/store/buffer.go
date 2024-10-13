package store

import (
	"strconv"
	"time"
)

type WriteBuffer struct {
	ch        chan envelope
	generator func() int
	st        *Store
}

func NewWriterBuffer(bufferSize uint, st *Store) (WriteBuffer, chan envelope) {
	ch := make(chan envelope, bufferSize)
	return WriteBuffer{
		ch:        ch,
		generator: generateSequence(),
		st:        st,
	}, ch
}

type envelope struct {
	txID     string
	received time.Time
	entry    Entry
}

func newEnvelope(txID string, d Entry, received time.Time) envelope {
	return envelope{
		txID:     txID,
		received: received,
		entry:    d,
	}
}

func generateSequence() func() int {
	i := 1
	return func() int {
		i++
		return i
	}
}

func (p WriteBuffer) Write(ds ...Entry) {
	received := time.Now()
	for _, entry := range ds {
		p.ch <- newEnvelope(strconv.Itoa(p.generator()), entry, received)
	}
}

func (p WriteBuffer) Start(done chan<- struct{}) {
	for d := range p.ch {
		p.st.Put(d.entry.Key, d.entry.Value)
		p.st.CommitTx(d.txID)
	}
	close(done)
}
