package goutils_test

import (
	"fmt"

	"github.com/watsonserve/goutils"
)

func ExampleNewRangeLink() {
	// # test.conf
	// # this is a config file
	// foo = bar
	// foos=1
	// foos = 2
	// null
	//
	link := goutils.NewRangeLink(nil)
	link.Mount(20, 80)
	fmt.Println(link.ToArray())
}
