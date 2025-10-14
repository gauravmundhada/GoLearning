package main

import (
	"fmt"
	"sort"
)

// Slices are the view on the top of the array - like an abstraction level on array
// Slice syntax = sliceName []T -> ex : names []string
// Access elements = names[inclusive_index : exclusive_index]
// Syntax to define slice = mySlice := make([]T, Len, Cap) - this creates an array; it is not necessary to reference an existing array
type Numbers []int

type byInc struct {
	Numbers
}

func (n byInc) Len() int           { return len(n.Numbers) }
func (n byInc) Swap(i, j int)      { n.Numbers[i], n.Numbers[j] = n.Numbers[j], n.Numbers[i] }
func (n byInc) Less(i, j int) bool { return n.Numbers[i] < n.Numbers[j] }

type decOrder struct {
	Numbers
}

func (n decOrder) Len() int           { return len(n.Numbers) }
func (n decOrder) Swap(i, j int)      { n.Numbers[i], n.Numbers[j] = n.Numbers[j], n.Numbers[i] }
func (n decOrder) Less(i, j int) bool { return n.Numbers[i] > n.Numbers[j] }

func main() {
	numbers := Numbers{1, 99, 1000, 3, 5, 6, 1, 69}
	numbers = removeFromSlice(numbers, 3)
	fmt.Printf("After remove: %d \n", numbers)
	numbers = removeFromSliceInOrder(numbers, 2)
	fmt.Printf("After remove and maintaining order: %d \n", numbers)
	sort.Sort(byInc{numbers}) // this is a custom sort (need to implement len, swap, less func for this )
	fmt.Printf("Sort in asc: %d \n", numbers)
	sort.Sort(decOrder{numbers})
	fmt.Printf("Sort in desc: %d \n", numbers)

}

// doesn't maintain the order
func removeFromSlice(slice []int, index int) []int {
	slice[index] = slice[len(slice)-1] //
	return slice[:len(slice)-1]        //
}

// maintain the order
func removeFromSliceInOrder(slice []int, index int) []int {
	slice = append(slice[:index], slice[index+1:]...)
	return slice
}
