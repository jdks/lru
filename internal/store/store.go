package store

import (
	"sync"
)

type Store struct {
	capacity uint
	size     uint
	data     sync.Map
	txLog    txLog
	top      *node
	bottom   *node
}

type txLog map[string]bool

const txLogSize = 1000

func New(capacity uint) *Store {
	return &Store{
		capacity: capacity,
		txLog:    make(txLog, txLogSize),
	}
}

func (s *Store) Put(key, value string) {
	if stored, ok := s.data.Load(key); ok {
		if node, ok := stored.(*node); ok {
			node.update(value)
			s.reorder(node)
			return
		}
	}
	if s.size == s.capacity {
		s.evictLeastRecentlyUsed()
	}
	node := newNode(key, value)
	s.data.Store(key, node)
	s.reorder(node)
	s.size++
}

func (s *Store) Get(key string) string {
	if stored, ok := s.data.Load(key); ok {
		if node, ok := stored.(*node); ok {
			s.reorder(node)
			return node.value
		}
	}
	return ""
}

func (s *Store) Head(n int) EntryList {
	if n > int(s.capacity) {
		n = int(s.capacity)
	}
	res := make(EntryList, 0, n)
	current := s.top
	for i := 0; i < n; i++ {
		if current == nil {
			return res
		}
		res = append(res, NewEntry(current.key, current.value))
		current = current.prev
	}
	return res
}

func (s *Store) CommitTx(txID string) {
	// TODO: implement a transaction log
	s.txLog[txID] = true
}

func (s *Store) isEmpty() bool {
	return s.top == nil
}

func (s *Store) isNodeAtTop(node *node) bool {
	return s.top != nil && s.top.key == node.key
}

func (s *Store) isNodeAtBottom(node *node) bool {
	return s.bottom != nil && s.bottom.key == node.key
}

// initializeCache sets up the cache with the first element
//
// Before:  top: nil, bottom: nil
// After:   top: [A] <-> bottom: [A]
func (s *Store) initializeCache(node *node) {
	s.top = node
	s.bottom = node
}

// reorder adjusts the order of nodes to maintain LRU property
func (s *Store) reorder(node *node) {
	switch {
	case s.isEmpty():
		s.initializeCache(node)
	case !s.isNodeAtTop(node):
		s.moveNodeToTop(node)
	}
}

// moveNodeToTop moves the given node to the top of the cache
//
// Before:  top: [B] <-> [C] <-> [A] <-> bottom: [D]
// After:   top: [A] <-> [B] <-> [C] <-> bottom: [D]
func (s *Store) moveNodeToTop(node *node) {
	node.prev = s.top
	s.top.next = node
	s.top = node

	if s.isNodeAtBottom(node) {
		s.removeNodeFromBottom()
	}
}

// removeNodeFromBottom removes the node from the bottom of the cache
//
// Before:  top: [A] <-> [B] <-> [C] <-> bottom: [D]
// After:   top: [A] <-> [B] <-> bottom: [C]
func (s *Store) removeNodeFromBottom() {
	s.bottom = s.bottom.next
	if s.bottom != nil {
		s.bottom.prev = nil
	}
}

// evictLeastRecentlyUsed removes the least recently used item from the cache
//
// Before:  top: [A] <-> [B] <-> [C] <-> bottom: [D]
// After:   top: [A] <-> [B] <-> bottom: [C]
func (s *Store) evictLeastRecentlyUsed() {
	s.data.Delete(s.bottom.key)
	s.removeNodeFromBottom()
	s.size--
}
