package goutils_test

import (
    "fmt"
    "github.com/watsonserve/goutils"
)

func ExampleGetConf() {
    // # test.conf
    // # this is a config file
    // foo = bar
    // foos=1
    // foos = 2
    // null
    //
    conf, err := goutils.GetConf("test.conf")
    if nil != err {
        return
    }
    fmt.Println(conf)
    // output: { foo: ["bar"], foos: ["1", "2"], null: [] }
}
