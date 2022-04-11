package goutils_test

import (
    "github.com/watsonserve/goutils"
)

func ExamplePrintf() {
    goutils.Printf("- %d -\n", 404)
}
