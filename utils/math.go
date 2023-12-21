package utils

func LeastCommonMultiple(nums []int) int {
	dedupFactors := map[int]any{}

	for _, n := range nums {
		factors := PrimeFactorization(n)
		for _, f := range factors {
			dedupFactors[f] = nil
		}
	}

	lcm := 1
	for f := range dedupFactors {
		lcm *= f
	}

	return lcm
}

func PrimeFactorization(num int) []int {
	primeFactors := []int{}

	// get all 2s
	for num%2 == 0 {
		primeFactors = append(primeFactors, 2)
		num /= 2
	}

	// go through odd nums
	for i := 3; i*i <= num; i += 2 {
		for num%i == 0 {
			primeFactors = append(primeFactors, i)
			num /= i
		}
	}

	// whatever is left must still be prime
	if num > 2 {
		primeFactors = append(primeFactors, num)
	}

	return primeFactors
}
