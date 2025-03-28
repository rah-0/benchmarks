package meta

import (
	"runtime"
	"sort"
	"sync"
)

func ParallelIntSortAsc(data []int) {
	parallelIntSort(data, false)
}

func ParallelIntSortDesc(data []int) {
	parallelIntSort(data, true)
}

func parallelIntSort(data []int, reverse bool) {
	n := len(data)
	if n < 10000 {
		sort.Ints(data)
		if reverse {
			reverseInts(data)
		}
		return
	}

	numCPU := runtime.NumCPU()
	chunkSize := (n + numCPU - 1) / numCPU

	chunks := make([][]int, 0, numCPU)
	for i := 0; i < n; i += chunkSize {
		end := i + chunkSize
		if end > n {
			end = n
		}
		chunks = append(chunks, data[i:end])
	}

	var wg sync.WaitGroup
	for _, chunk := range chunks {
		wg.Add(1)
		go func(c []int) {
			defer wg.Done()
			sort.Ints(c)
		}(chunk)
	}
	wg.Wait()

	// Parallel merging loop
	for len(chunks) > 1 {
		mergedCount := (len(chunks) + 1) / 2
		merged := make([][]int, mergedCount)

		var mWg sync.WaitGroup
		for i := 0; i < len(chunks); i += 2 {
			if i+1 == len(chunks) {
				merged[i/2] = chunks[i]
				continue
			}
			mWg.Add(1)
			go func(i int) {
				defer mWg.Done()
				merged[i/2] = mergeSorted(chunks[i], chunks[i+1])
			}(i)
		}
		mWg.Wait()
		chunks = merged
	}

	copy(data, chunks[0])
	if reverse {
		reverseInts(data)
	}
}

func mergeSorted(a, b []int) []int {
	res := make([]int, len(a)+len(b))
	i, j, k := 0, 0, 0
	for i < len(a) && j < len(b) {
		if a[i] <= b[j] {
			res[k] = a[i]
			i++
		} else {
			res[k] = b[j]
			j++
		}
		k++
	}
	for i < len(a) {
		res[k] = a[i]
		i++
		k++
	}
	for j < len(b) {
		res[k] = b[j]
		j++
		k++
	}
	return res
}

func reverseInts(a []int) {
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
}
