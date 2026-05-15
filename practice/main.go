package main

import "fmt"

func main() {
	var a []int
	//a = []int{1, 6, 2, 6, 8, 22, 9, 5}
	a = []int{3, 1, 5}
	sort := mergeSort(a)
	fmt.Println("sorted : ", sort)

}

func mergeSort(a []int) []int {
	if len(a) <= 1 {
		return a
	}
	mid := len(a) / 2
	first := mergeSort(a[:mid])
	last := mergeSort(a[mid:])
	// fmt.Println(a)
	fmt.Println("first ", first)
	fmt.Println("last", last)
	return merge(first, last)

}

func merge(first, last []int) []int {
	var result []int
	i, j := 0, 0
	for i < len(first) && j < len(last) {
		fmt.Println("first ", first)
		fmt.Println("last ", last)
		if first[i] <= last[j] {
			fmt.Println("first[i]", first[i])
			result = append(result, first[i])
			fmt.Println("result", result)
			i++
		} else {
			fmt.Println("last[j]", last[j])
			result = append(result, last[j])
			fmt.Println("result", result)
			j++
		}
		//fmt.Println("Appended result : ", result)
	}

	// Append remaining elements
	result = append(result, first[i:]...)
	result = append(result, last[j:]...)

	return result

}
