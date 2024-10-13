package store

import (
	"fmt"
	"math/rand"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkPut_SerialBuffer(b *testing.B) {
	s := New(100)
	wb, ch := NewWriterBuffer(10000, s)
	done := make(chan struct{})
	go wb.Start(done)
	for i := 0; i < b.N; i++ {
		wb.Write(NewEntry(fmt.Sprintf("%d", i), fmt.Sprintf("%d", i)))
	}
	close(ch)
	<-done
}

func TestPut(t *testing.T) {
	s := New(1000)
	done := make(chan struct{})
	wb, ch := NewWriterBuffer(100, s)
	go wb.Start(done)
	testData := EntryList{
		{"a", "1"},
		{"j", "2"},
		{"e", "2"},
		{"b", "3"},
		{"y", "-"},
		{"h", "8"},
	}

	for _, td := range testData {
		wb.Write(td)
	}
	close(ch)
	<-done
	expected := testData[len(testData)-2:]
	slices.Reverse(expected)
	assert.Equal(t, expected, s.Head(2))
}

func TestPut_Eviction(t *testing.T) {
	s := New(3)
	done := make(chan struct{})
	wb, ch := NewWriterBuffer(100, s)
	go wb.Start(done)
	testData := EntryList{
		{"a", "1"},
		{"b", "2"},
		{"c", "3"},
		{"d", "4"},
		{"e", "5"},
	}
	for _, td := range testData {
		wb.Write(td)
	}
	close(ch)
	<-done
	expected := testData[len(testData)-3:]
	slices.Reverse(expected)
	assert.Equal(t, expected, s.Head(5))
}

func TestPut_Update(t *testing.T) {
	s := New(1000)
	done := make(chan struct{})
	wb, ch := NewWriterBuffer(100, s)
	go wb.Start(done)
	testData := EntryList{
		{"a", "1"},
		{"a", "2"},
		{"a", "3"},
		{"a", "4"},
		{"a", "5"},
	}
	for _, td := range testData {
		wb.Write(td)
	}
	close(ch)
	<-done
	expected := testData[len(testData)-1:]
	assert.Equal(t, s.Get("a"), "5")
	assert.Equal(t, expected, s.Head(1))
}

func TestGet(t *testing.T) {
	s := New(1000)
	done := make(chan struct{})
	wb, ch := NewWriterBuffer(100, s)
	go wb.Start(done)
	testData := EntryList{
		{"a", "1"},
		{"b", "2"},
		{"b", "1"},
		{"e", "2"},
		{"b", "3"},
		{"h", "0"},
		{"y", "-"},
		{"h", "8"},
	}
	for _, td := range testData {
		wb.Write(td)
	}
	close(ch)
	<-done
	assert.Equal(t, s.Get("a"), "1")
	assert.Equal(t, s.Get("b"), "3")
}

func TestHead(t *testing.T) {
	s := New(1000)
	done := make(chan struct{})
	wb, ch := NewWriterBuffer(100, s)
	go wb.Start(done)
	testData := EntryList{
		{"a", "1"},
		{"b", "2"},
		{"e", "2"},
		{"b", "3"},
		{"y", "-"},
		{"h", "8"},
	}
	for _, td := range testData {
		wb.Write(td)
	}
	close(ch)
	<-done
	expected := testData[len(testData)-2:]
	slices.Reverse(expected)
	assert.Equal(t, expected, s.Head(2))
}
