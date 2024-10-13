package store

type node struct {
	key   string
	value string
	prev  *node
	next  *node
}

func newNode(key, value string) *node {
	return &node{key: key, value: value}
}

func (n *node) update(value string) {
	n.value = value
}
