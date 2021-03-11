package main

import (
	"flag"
	"io"
	"log"
	"os"
	"strings"
	"text/template"
)

var (
	st     = flag.String("s", "", "set type, default: ${tp}s")
	pkg    = flag.String("p", "", "element package, default: don't import package")
	tp     = flag.String("t", "", "element type")
	output = flag.String("o", "", "Output file; defaults to current ${st}.go")
	light  = flag.Bool("l", false, "set need ErrBreakEach error, light predicates whether generate code imports github.com/SeananXu/go-set")
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
	if *pkg != "" {
		if i := strings.LastIndex(*pkg, "/"); i != -1 {
			pkgWithStruct = (*pkg)[i+1:] + "." + structName
		} else {
			pkgWithStruct = *pkg + "." + structName
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
	t, err := template.New("setgen").Parse(tmp)
	if err != nil {
		log.Fatalf("parse template file error: %v", err)
	}
	var writer io.Writer
	switch *output {
	case "":
		writer, err = os.OpenFile(strings.ToLower(*st)+".go", os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Fatalf("create writer steam error: %v", err)
		}
	case "-":
		writer = os.Stdout
	default:
		writer, err = os.OpenFile(*output, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Fatalf("create writer steam error: %v", err)
		}
	}
	if err = t.Execute(writer, map[string]interface{}{
		"st":    st,
		"tp":    tp,
		"obj":   obj,
		"light": *light,
		"pkg":   *pkg,
	}); err != nil {
		log.Fatalf("excute output file error: %v", err)
	}
}
