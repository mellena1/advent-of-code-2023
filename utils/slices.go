package utils

import "strconv"

func StrSliceToIntSlice(strs []string) ([]int, error) {
	nums := make([]int, len(strs))
	for i, s := range strs {
		n, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		nums[i] = n
	}
	return nums, nil
}
