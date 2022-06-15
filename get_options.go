package goutils

import (
	"fmt"
	"os"
	"strings"
)

type Option struct {
	Name      string
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
		if 0 == opt {
			continue
		}
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

		if argpLen < 2 {
			// TODO 只有一个横线
			continue
		}

		// 长选项 option
		if '-' == argp[1] {
			if 2 < argpLen {
				key, value := parseOption(argp, argpLen)
				table[key] = value
			} else {
				// 只有两个横线
				params = append(params, strings.Join(argv[i+1:], " "))
				return table, params
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
			table[opt.Name] = ""
		}
		// 最后一个选项opt需要判断是否需要参数
		opt := optionsMap[argp[lst]]
		if opt.HasParams && i+1 < argc {
			payload := argv[i+1]
			if '-' == payload[0] {
				continue
			}
			table[opt.Name] = payload
			i++
		}
	}

	return table, params
}

func GenHelp(options []Option, desc string) string {
	ret := os.Args[0] + " [opt argument]" + desc + "\n"
	for i := 0; i < len(options); i++ {
		item := options[i]
		opt := "   "
		if 0 != item.Opt {
			opt = fmt.Sprintf("-%c,", item.Opt)
		}
		option := "  \t"
		if "" != item.Option {
			option = fmt.Sprintf("--%s", item.Option)
		}
		ret += fmt.Sprintf("  %s  %s\t%s\n", opt, option, item.Desc)
	}
	return ret + "\n"
}
