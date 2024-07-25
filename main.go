package main

import (
	"fmt"

	"github.com/calderwd/jobframework/test"
)

type TestJob struct {
}

func (j TestJob) Process() (bool, error) {
	fmt.Println("Test job running")
	return true, nil
}

type TestJobProfile struct {
}

func main() {
	fmt.Println("Start")

	test.RunAddTest()

}

// import (
// 	"container/list"
// 	"fmt"
// )

// // Find the maximums in each sliding window of size k
// func slidingWindowMax(nums []int, k int) []int {
// 	if len(nums) == 0 || k == 0 {
// 		return []int{}
// 	}

// 	result := make([]int, 0, len(nums)-k+1)
// 	deque := list.New()

// 	for i := 0; i < len(nums); i++ {
// 		// Remove elements not within the sliding window
// 		if deque.Len() > 0 && deque.Front().Value.(int) < i-k+1 {
// 			deque.Remove(deque.Front())
// 		}

// 		// Remove elements smaller than the current element
// 		for deque.Len() > 0 && nums[deque.Back().Value.(int)] < nums[i] {
// 			deque.Remove(deque.Back())
// 		}

// 		// Add current element at the back of the deque
// 		deque.PushBack(i)

// 		// The front of the deque is the largest element in the current window
// 		if i >= k-1 {
// 			result = append(result, nums[deque.Front().Value.(int)])
// 		}
// 	}

// 	return result
// }

// func main() {
// 	nums := []int{1, 3, -1, -3, 5, 3, 6, 7}
// 	k := 3
// 	fmt.Println(slidingWindowMax(nums, k)) // Output: [3 3 5 5 6 7]
// }
