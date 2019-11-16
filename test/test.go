package main

import (
    "fmt"
    "github.com/watsonserve/goutils"
)

func main() {
    options := []goutils.Option { 
        {
            Opt: 'h',
            Option: "help",
            HasParams: false,
        },
        {
            Opt: 'z',
            Option: "gz",
            HasParams: false,
        },
        {
            Opt: 'C',
            Option: "cc",
            HasParams: true,
        },
        {
            Opt: 'f',
            Option: "file",
            HasParams: true,
        },
    }
    argMap, params := goutils.GetOptions(options)

    fmt.Println(argMap)
    fmt.Println(params)
}
