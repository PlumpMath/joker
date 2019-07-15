package main

import (
	"bytes"
	"io/ioutil"
	"strings"

	. "github.com/candid82/joker/core"
	_ "github.com/candid82/joker/std/csv"
	_ "github.com/candid82/joker/std/string"
)

var template string = `// Generated by gen_data. Don't modify manually!

package core

func init() {
	{name}Data = []byte("{content}")
}
`

type FileInfo struct {
	name     string
	filename string
}

/* The entries must be ordered such that a given namespace depends
/* only upon namespaces loaded above it. E.g. joker.template depends
/* on joker.walk, so is listed afterwards, not in alphabetical
/* order. */
var files []FileInfo = []FileInfo{
	{
		name:     "<joker.core>",
		filename: "core.joke",
	},
	{
		name:     "<joker.csv>",
		filename: "csv.joke",
	},
	{
		name:     "<joker.repl>",
		filename: "repl.joke",
	},
	{
		name:     "<joker.walk>",
		filename: "walk.joke",
	},
	{
		name:     "<joker.template>",
		filename: "template.joke",
	},
	{
		name:     "<joker.test>",
		filename: "test.joke",
	},
	{
		name:     "<joker.set>",
		filename: "set.joke",
	},
	{
		name:     "<joker.tools.cli>",
		filename: "tools_cli.joke",
	},
	{
		name:     "<joker.core>",
		filename: "linter_all.joke",
	},
	{
		name:     "<joker.core>",
		filename: "linter_joker.joke",
	},
	{
		name:     "<joker.core>",
		filename: "linter_cljx.joke",
	},
	{
		name:     "<joker.core>",
		filename: "linter_clj.joke",
	},
	{
		name:     "<joker.core>",
		filename: "linter_cljs.joke",
	},
}

const hextable = "0123456789abcdef"

func main() {
	GLOBAL_ENV.FindNamespace(MakeSymbol("user")).ReferAll(GLOBAL_ENV.CoreNamespace)
	for _, f := range files {
		GLOBAL_ENV.SetCurrentNamespace(GLOBAL_ENV.CoreNamespace)
		content, err := ioutil.ReadFile("data/" + f.filename)
		if err != nil {
			panic(err)
		}
		content, err = PackReader(NewReader(bytes.NewReader(content), f.name), "")
		PanicOnErr(err)

		dst := make([]byte, len(content)*4)
		for i, v := range content {
			dst[i*4] = '\\'
			dst[i*4+1] = 'x'
			dst[i*4+2] = hextable[v>>4]
			dst[i*4+3] = hextable[v&0x0f]
		}
		name := f.filename[0 : len(f.filename)-5] // assumes .joke extension
		fileContent := strings.Replace(template, "{name}", name, 1)
		fileContent = strings.Replace(fileContent, "{content}", string(dst), 1)
		ioutil.WriteFile("a_"+name+"_data.go", []byte(fileContent), 0666)
	}
}
