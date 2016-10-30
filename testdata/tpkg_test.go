package tpkg

import (
	"fmt"
	"testing"
)

func TestFuncOne(t *testing.T) {}

func TestFuncTwo(t *testing.T) {}

func TestFuncThree(t *testing.T) {}

func ExamplePublicStruct() {
	fmt.Println("ExamplePublicStruct")

	// Output:
	// ExamplePublicStruct
}

func ExamplePublicStruct_MethodOne() {
	fmt.Println("ExamplePublicStruct_MethodOne")

	// Output:
	// ExamplePublicStruct_MethodOne
}
