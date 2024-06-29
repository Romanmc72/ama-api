package application

import (
	"errors"
	"slices"
	"strings"
)

// Sorts the input slice, removes any duplicates, and ignores any instances of
// the illegal character or blanks.
func SortDedupeAndIgnore(options []string, illegalCharacter string) []string {
	slices.Sort(options)
	options = slices.Compact(options)
	final := make([]string, 0, len(options))
	for _, o := range options {
		if o != illegalCharacter && strings.TrimSpace(o) != "" {
			final = append(final, o)
		}
	}
	return final
}

// For a list of size X, there will return Y elements that represent all of the
// unique possible combinations of each of the Y elements where there are no
// repeated sets (i.e. order does not matter). Any input list will be compacted
// to remove duplicates as well as sorted to ensure that the results are always
// returned the same for an equivalent input. Duplicates cause issues with the
// algorithm. The resulting size Y can be calculated as follows:
//
//  func tot(X) {
// 		Y := 0
// 		for i := range X {
// 			Y = sum + C(X, i)
// 		}
//  }
//
// where
//
// 	func C(n int, k int) int
//
// is of the mathematical form:
//
// 	C(n,k) = n!/((n-k)!k!)
//
// I am pretty sure the time complexity of this algorithm is O(n) which
// I think can be simplified to O(n*n!).
//
// "tot()"" is defined above
//
// The space complexity of this algorithm is O(2*tot(n)+n)
func Combine(options []string, delimiter string) []string {
	options = SortDedupeAndIgnore(options, delimiter)
	numberOfOptions := uint(len(options))
	if numberOfOptions == 0 {
		return options
	}
	// pre-allocatings the array will save time/space
	numberOfResults, err := totalCombinations(numberOfOptions, numberOfOptions)
	if err != nil {
		return options
	}
	results := make([]string, numberOfResults)
	// Pre-allocate a slice size to use for storing intermediate results
	perIterationCombinations := int(combinationsForOneChoiceSet(numberOfOptions, 1))
	intermediateResults := make([][]string, perIterationCombinations)
	previousSlices := make([][]string, perIterationCombinations)
	position := 0
	for iteration := range numberOfOptions {
		optionsThisIteration := iteration + 1
		// Round 1, each item is added as a standalone
		if iteration == 0 {
			for _, o := range options {
				intermediateResults[position] = []string{o}
				position += 1
			}
		// The final iteration uses every item combined as a single entry
		} else if iteration == numberOfOptions - 1 {
			intermediateResults[position] = options
			position += 1
		// Every other iteration will run through carrying a pointer to the previous
		// run's iteration, tracking its rightmost position along with a pointer to
		// the remaining items to combine ensuring that it grabs and adds on items
		// to the running list in order to create all possible combinations without repeats.
		} else {
			for _, previousSlice := range previousSlices {
				previousSliceRightOption := previousSlice[len(previousSlice) - 1]
				// The point on the previous slice we are currently iterating that
				// represents the right most edge of the option set. If the right edge
				// is the same index as the last element in the choices set, then we
				// must continue to the next element of the previous iteration.
				previousSliceRightIndex := 0
				// The index pointing to the next option after the right most edge of the
				// previous slice. When it hits the end of the available options it will
				// reset and move along to the next grouping of the previous slice.
				nextOptionIndex := 0
				for i, o := range options {
					if previousSliceRightOption == o {
						previousSliceRightIndex = i
						nextOptionIndex = i + 1
						break
					}
				}
				if previousSliceRightIndex == nextOptionIndex || nextOptionIndex >= len(options) {
					continue
				}

				for _, o := range options[nextOptionIndex:] {
					currentSlice := make([]string, optionsThisIteration)
					copy(currentSlice, previousSlice)
					currentSlice[iteration] = o
					intermediateResults[position] = currentSlice
					position += 1
				}
			}
		}
		for i, r := range intermediateResults {
			positionOffset := 0
			if iteration > 0 {
				positionOffsetUint, err := totalCombinations(numberOfOptions, iteration)
				if err != nil {
					return options
				}
				positionOffset = int(positionOffsetUint)
			}
			finalPosition := i + int(positionOffset)
			optionToAdd := strings.Join(r, delimiter)
			results[finalPosition] = optionToAdd
		}
		// Copy over the intermediate slice from the iteration that just finished
		// for use in the next iteration.
		previousSlices = make([][]string, perIterationCombinations)
		copy(previousSlices, intermediateResults)
		if int(optionsThisIteration) != len(options) {
			perIterationCombinations = int(combinationsForOneChoiceSet(
				numberOfOptions, optionsThisIteration + 1))
			intermediateResults = make([][]string, perIterationCombinations)
		}
		position = 0
	}
	return results
}

// Get the number of choices for one single combination and choice set.
func combinationsForOneChoiceSet(numberOfOptions uint, choices uint) uint {
	return (
		factorial(numberOfOptions) / (
			factorial(numberOfOptions - choices) * factorial(choices)))
}

// Determines the number of results that will appear in the final combination
// for all choices from n=choices to n=1.
func totalCombinations(numberOfOptions uint, choices uint) (uint, error) {
	if choices > numberOfOptions {
		return 0, errors.New("cannot select more choices than there are options")
	}
	if choices == 1 {
		return numberOfOptions, nil
	}
	lowerAmount, err := totalCombinations(numberOfOptions, choices - 1)
	if err != nil {
		return 0, err
	}
	result := lowerAmount + combinationsForOneChoiceSet(numberOfOptions, choices)
	return result, nil
}

// Perform the factorial of a number
func factorial(number uint) uint {
	if number <= 0 {
		return 1
	}
	return number * factorial(number - 1)
}
