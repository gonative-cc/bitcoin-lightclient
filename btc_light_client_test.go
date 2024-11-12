package main

import (
	"testing"
)

func TestCleanup(t *testing.T) {
	// test-case1
	// b1 <-b2 <- b3  .... b9
	// run cleanup
	// b2 <- b3  .... b9

	// test-case2
	// b3  ... b8 <- b9 <- b10
	//                  \- c10
	// run cleanup
	// b3  ... b8 <- b9 <- b10
	//                  \- c10

	// test-case3
	// b1 <- b2 <- b3  .... b9
	//          \- c3
	//    \- d2 <- d3 <- d4
	// run cleanup
	// b2 <- b3  .... b9
	//    \- c3

}
