// Author: Tan Shane Lee
// This project is a part of an assessment test for Golang.
// The goal is to create a servlet that accepts a body of text such from a book and return the top 10
// most frequent occurring words.
// This project is created using various external packages, and references.

package main

import (
	"bufio"
	"bytes"
	"container/heap"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/jdkato/prose/tokenize"
)

// An Item is something we manage in a priority queue.
type Item struct {
	value    string // The value of the item; arbitrary.
	priority int    // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].priority > pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(item *Item, value string, priority int) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}

// WordCount Returns a map with the key being the word and values being the number of occurrences.
// Also removes all non-word characters like hyphens, commas, full-stops, etc.
func WordCount(s string) map[string]int {
	tokenizer := tokenize.NewWordPunctTokenizer()
	r, _ := regexp.Compile(`\W`) // Anything that is not a word character like hyphens, commas etc...

	words := tokenizer.Tokenize(s)

	counts := make(map[string]int, len(words))
	for _, word := range words {
		if !r.MatchString(word) {
			counts[strings.ToLower(word)]++ // Makes all lowercase as it it Case-sensitive
		}
	}
	return counts
}

func GetFrequencyMap(content string) PriorityQueue {

	items := WordCount(content)

	// Create a priority queue, put the items in it, and
	// establish the priority queue (heap) invariants.
	pq := make(PriorityQueue, len(items))
	i := 0
	for value, priority := range items {
		pq[i] = &Item{
			value:    value,
			priority: priority,
			index:    i,
		}
		i++
	}

	heap.Init(&pq)

	return pq
}

func PrintTopNFrequentWords(n uint, pq PriorityQueue) {

	// Takes the top 10 from the list.
	nFromList := 10

	if pq.Len() < nFromList {
		fmt.Println(`There are less words than ` + strconv.Itoa(nFromList) + ` retrieving ` + strconv.Itoa(pq.Len()) + ` instead.`)
		nFromList = pq.Len()
	}

	for i := 0; i < nFromList; i++ {
		item := heap.Pop(&pq).(*Item)
		fmt.Printf("%.2d:%s ", item.priority, item.value)
	}
}

func ReadFile(filePath string) string {
	file, _ := os.Open(filePath)
	reader := bufio.NewReader(file)

	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)

	return buf.String()
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Test")
}

func main() {
	// http.HandleFunc("/", IndexHandler)
	// http.ListenAndServe(":8000", nil)
	PrintTopNFrequentWords(10, GetFrequencyMap(ReadFile("file.txt")))
}
