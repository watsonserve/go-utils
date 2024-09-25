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

func parseLongOption(argp string) (string, string) {
	kv := strings.SplitN(argp[2:], "=", 2)
	key := kv[0]
	value := ""
	if 1 < len(kv) {
		value = kv[1]
	}
	return key, value
}

/**
 * @returns table and key which want a value
 */
func parseShortOption(table map[string]string, argp string, optionsMap map[byte]*Option) string {
	for i := 1; i < len(argp); i++ {
		info := optionsMap[argp[i]]
		// unknow opt
		if nil == info {
			continue
		}
		// a flag
		if !info.HasParams {
			table[info.Name] = ""
			continue
		}
		// a option
		val := argp[i+1:]
		if "" == val {
			return info.Name
		}

		table[info.Name] = val
		break
	}

	return ""
}

// -o             short option without value
// -oo            short option list without any value
// -oValue        short option and it's value
// -o value       short option and it's value
// -ooValue       short options and a value
// -oo Value      short options and a value
// --option       long option
// --option=value long option
// -- foo bar     payload no split
// value          payload or value of a option

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

	optsMap := make(map[byte]*Option)
	optionsMap := make(map[string]*Option)
	for i := 0; i < len(options); i++ {
		info := &(options[i])
		opt := info.Opt
		option := info.Option

		if 0 != opt {
			optsMap[opt] = info
		}
		optionsMap[option] = info
	}

	waitingKey := ""
	// get option
	for i := 1; i < argc; i++ {
		argp := argv[i]

		// all as one payload after this option signal
		if "--" == argp {
			params = append(params, strings.Join(argv[i+1:], " "))
			return table, params
		}

		// a long option
		if strings.HasPrefix(argp, "--") {
			key, val := parseLongOption(argp)
			info := optionsMap[key]
			if nil != info {
				table[info.Name] = val
				waitingKey = ""
			}
			continue
		}

		// a short option set
		if '-' == argp[0] && 1 < len(argp) {
			waitingKey = parseShortOption(table, argp, optsMap)
			continue
		}

		// a value of short opt
		if "" != waitingKey {
			table[waitingKey] = argp
			waitingKey = ""
			continue
		}

		// a payload
		params = append(params, argp)
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
