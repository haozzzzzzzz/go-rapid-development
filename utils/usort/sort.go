package usort

import (
	"sort"
)

func SortAsc(data sort.Interface) {
	sort.Sort(data)
}

func SortDesc(data sort.Interface) {
	sort.Sort(sort.Reverse(data))
}

type Int64Slice []int64

func (p Int64Slice) Len() int           { return len(p) }
func (p Int64Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Int64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Sort is a convenience method.
func (p Int64Slice) SortAsc() { SortAsc(p) }

func (p Int64Slice) SortDesc() { SortDesc(p) }
