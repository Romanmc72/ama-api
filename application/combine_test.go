package application

import (
	"ama/api/constants"
	"reflect"
	"testing"
)

// Ensure it works with zero inputs
func TestCombineZero(t *testing.T) {
	starting := []string{}
	want := []string{}
	actual := Combine(starting, constants.SearchTagDelimiter)
	if !reflect.DeepEqual(actual, want) {
		t.Fatalf("Wanted %s but received %s for the Combine method input of %s", want, actual, starting)
	}
}

// Ensure it works with one input
func TestCombineOne(t *testing.T) {
	starting := []string{"a"}
	want := []string{"a"}
	actual := Combine(starting, constants.SearchTagDelimiter)
	if !reflect.DeepEqual(actual, want) {
		t.Fatalf("Wanted %s but received %s for the Combine method input of %s", want, actual, starting)
	}
}

// Ensure it works with two inputs
func TestCombineTwo(t *testing.T) {
	starting := []string{"a", "b"}
	want := []string{"a", "b", "a|b"}
	actual := Combine(starting, constants.SearchTagDelimiter)
	if !reflect.DeepEqual(actual, want) {
		t.Fatalf("Wanted %s but received %s for the Combine method input of %s", want, actual, starting)
	}
}

// Test that the slice combine method does what is expected with 4 elements to combine.
// There is nothing special about 4, it is just higher than 2 by a decent enough amount
// that it proves that this works. This also confirms the compact and sort works.
func TestCombineFour(t *testing.T) {
	starting := []string{"d", "a", "a", "b", "c", "d", constants.SearchTagDelimiter}
	want := []string{
		"a", "b", "c", "d", // iteration = 0; i = [0, 1, 2, 3]
		"a|b", "a|c", "a|d", // iteration = 1; i = [4, 5, 6]
		"b|c", "b|d", // iteration = 1; i = [7, 8]
		"c|d",            // iteration = 1; i = [9]
		"a|b|c", "a|b|d", // iteration = 2; i = [10, 11]
		"a|c|d",   // iteration = 2; i = [12]
		"b|c|d",   // iteration = 2; i = [13]
		"a|b|c|d", // iteration = 3; i = [14]
	}
	actual := Combine(starting, constants.SearchTagDelimiter)
	if !reflect.DeepEqual(actual, want) {
		t.Fatalf("Wanted %s but received %s for the Combine method input of %s", want, actual, starting)
	}
}
