// +build go_run_only

package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"zgo.at/utils/mathutil"
	"zgo.at/zpack"
)

func main() {
	dirs, err := lsdomains()
	if err != nil {
		panic(err)
	}

	fp, err := os.Create("pack.go")
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	zpack.Header(fp, "main")
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(fp, "var packmap = map[string]map[string][]byte{\n")
	for _, d := range dirs {
		fmt.Fprintf(fp, "\t\"%s\": %s,\n", d, zpack.Varname(d))
	}
	fmt.Fprintf(fp, "}\n\n")

	l := size(fp)
	for _, d := range dirs {
		err := zpack.Dir(fp, zpack.Varname(d), d)
		if err != nil {
			panic(err)
		}

		fmt.Printf("%s: %s\n", d, mathutil.Byte(size(fp)-l))
		l = size(fp)
	}
}

func size(fp *os.File) int64 {
	fp.Sync()
	st, _ := fp.Stat()
	return st.Size()
}

func lsdomains() ([]string, error) {
	dirs, err := ioutil.ReadDir(".")
	if err != nil {
		return nil, err
	}

	var domains []string
	for _, d := range dirs {
		d, _ = os.Stat(d.Name()) // To follow links.
		if !d.IsDir() || d.Name()[0] == '.' {
			continue
		}
		domains = append(domains, d.Name())
	}

	return domains, nil
}
