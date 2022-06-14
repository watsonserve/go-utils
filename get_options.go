package goutils

import (
	"fmt"
	"os"
	"strings"
)

type Option struct {
	Opt       byte
	Option    string
	HasParams bool
	Desc      string
}

func parseOption(argp string, length int) (string, string) {
	kv := strings.Split(argp[2:length], "=")
	key := kv[0]
	value := ""
	if 1 < len(kv) {
		value = kv[1]
	}
	return key, value
}

// parse commandline params
func GetOptions(options []Option) (map[string]string, []string) {
	table := make(map[string]string)
	params := make([]string, 0)

	argv := os.Args
	argc := len(argv)

	// no any params
	if 2 > argc {
		return table, params
	}

	optionsMap := make(map[byte]*Option)
	for i := 0; i < len(options); i++ {
		opt := options[i].Opt
		optionsMap[opt] = &options[i]
	}

	// get option
	for i := 1; i < argc; i++ {
		argp := argv[i]

		// value
		if '-' != argp[0] {
			params = append(params, argp)
			continue
		}

		argpLen := len(argp)
		// 只有一个横线
		if argpLen < 2 {
			// TODO
			continue
		}

		// 长选项 option
		if '-' == argp[1] {
			// 只有两个横线
			if 2 < argpLen {
				key, value := parseOption(argp, argpLen)
				table[key] = value
			}
			continue
		}

		// 短选项 opt
		argp = argp[1:argpLen]
		argpLen -= 1
		lst := argpLen - 1
		for j := 0; j < argpLen; j++ {
			opt := optionsMap[argp[j]]
			if nil == opt {
				continue
			}
			key := opt.Option
			table[key] = ""
		}
		// 最后一个选项opt需要判断是否需要参数
		opt := optionsMap[argp[lst]]
		if opt.HasParams && i+1 < argc {
			payload := argv[i+1]
			if '-' == payload[0] {
				continue
			}
			table[opt.Option] = payload
			i++
		}
	}

	return table, params
}

func GenHelp(options []Option) string {
	ret := os.Args[0] + " [opt argument]\n"
	for i := 0; i < len(options); i++ {
		item := options[i]
		opt := "   "
		if '\x00' != item.Opt {
			opt = fmt.Sprintf("-%c,", item.Opt)
		}
		option := "  \t"
		if "" != item.Option {
			option = fmt.Sprintf("--%s", item.Option)
		}
		ret += fmt.Sprintf("  %s  %s %s\n", opt, option, item.Desc)
	}
	return ret + "\n"
}
