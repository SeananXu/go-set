/*
MIT License

Copyright (c) 2021 Seanan Xu

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package main

import (
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var (
	st     = flag.String("s", "", "Set name, default: element type add 's'.")
	ipt    = flag.String("i", "", "Import element package, default: don't import package.")
	pkg    = flag.String("p", "", "Generated go file package, default: directory name.")
	tp     = flag.String("t", "", "Set storage element type, this options must be set.")
	output = flag.String("o", "", "Output file name, default: set name add '.go'.")
	light  = flag.Bool("l", false, "Whether go file imports 'ErrBreakEach' of 'github.com/SeananXu/go-set', default: import.")
)

func main() {
	flag.Parse()
	if *tp == "" {
		log.Fatalf("empty element type, please use -t set up the type")
	}
	if strings.HasPrefix(*tp, "**") {
		log.Fatalf("element type invalid, the prefix only allows one '*'")
	}
	structName := *tp
	if strings.HasPrefix(*tp, "*") {
		structName = (*tp)[1:]
	}
	if *st == "" {
		*st = structName + "s"
	}
	pkgWithStruct := structName
	if *ipt != "" {
		if i := strings.LastIndex(*ipt, "/"); i != -1 {
			pkgWithStruct = (*ipt)[i+1:] + "." + structName
		} else {
			pkgWithStruct = *ipt + "." + structName
		}
	}
	var obj string
	if strings.HasPrefix(*tp, "*") {
		obj += "&" + pkgWithStruct + "{}"
		*tp = "*" + pkgWithStruct
	} else {
		*tp = pkgWithStruct
		obj = pkgWithStruct + "{}"
	}
	if *pkg == "" {
		pwd, _ := os.Getwd()
		*pkg = filepath.Base(pwd)
	}
	t, err := template.New("setgen").Parse(tmp)
	if err != nil {
		log.Fatalf("parse template file error: %v", err)
	}
	var writer io.Writer
	switch *output {
	case "":
		f, err := os.OpenFile(strings.ToLower(*st)+".go", os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0644)
		if err != nil {
			log.Fatalf("create writer steam error: %v", err)
		}
		defer f.Close()
		writer = f
	case "-":
		writer = os.Stdout
	default:
		f, err := os.OpenFile(*output, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0644)
		if err != nil {
			log.Fatalf("create writer steam error: %v", err)
		}
		defer f.Close()
		writer = f
	}
	if err = t.Execute(writer, map[string]interface{}{
		"st":    st,
		"tp":    tp,
		"obj":   obj,
		"light": *light,
		"ipt":   *ipt,
		"pkg":   *pkg,
	}); err != nil {
		log.Fatalf("excute output file error: %v", err)
	}
}
