package store

type Entry struct {
	Key   string
	Value string
}

type EntryList []Entry

func NewEntry(key, value string) Entry {
	return Entry{Key: key, Value: value}
}
