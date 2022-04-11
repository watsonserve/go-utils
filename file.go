package goutils

import (
    "io/ioutil"
    "strings"
)

// 从一个文件中读取前n行和后面所有内容
func ReadLineN(filename string, n int) ([]string, error) {
    content, err := ioutil.ReadFile(filename)
    if nil != err {
        return make([]string, 0), err
    }

    lines := strings.SplitN(string(content), "\n", n)
    return lines, nil
}

