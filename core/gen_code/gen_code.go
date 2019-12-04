package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"

	. "github.com/candid82/joker/core"
	_ "github.com/candid82/joker/std/string"
)

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
	// {
	// 	name:     "<joker.repl>",
	// 	filename: "repl.joke",
	// },
	// {
	// 	name:     "<joker.walk>",
	// 	filename: "walk.joke",
	// },
	// {
	// 	name:     "<joker.template>",
	// 	filename: "template.joke",
	// },
	// {
	// 	name:     "<joker.test>",
	// 	filename: "test.joke",
	// },
	// {
	// 	name:     "<joker.set>",
	// 	filename: "set.joke",
	// },
	// {
	// 	name:     "<joker.tools.cli>",
	// 	filename: "tools_cli.joke",
	// },
	// {
	// 	name:     "<joker.core>",
	// 	filename: "linter_all.joke",
	// },
	// {
	// 	name:     "<joker.core>",
	// 	filename: "linter_joker.joke",
	// },
	// {
	// 	name:     "<joker.core>",
	// 	filename: "linter_cljx.joke",
	// },
	// {
	// 	name:     "<joker.core>",
	// 	filename: "linter_clj.joke",
	// },
	// {
	// 	name:     "<joker.core>",
	// 	filename: "linter_cljs.joke",
	// },
}

const hextable = "0123456789abcdef"
const masterFile = "a_code.go"

func main() {
	codeWriterEnv := &CodeWriterEnv{
		Need:      map[string]Finisher{},
		Generated: map[interface{}]interface{}{},
	}

	GLOBAL_ENV.FindNamespace(MakeSymbol("user")).ReferAll(GLOBAL_ENV.CoreNamespace)
	for _, f := range files {
		fileTemplate := `// Generated by gen_code. Don't modify manually!

package core

func init() {
	{name}NamespaceInfo = internalNamespaceInfo{init: {name}Init, generated: {name}NamespaceInfo.generated, available: true}
}

{statics}
func {name}Init() {
{interns}
}
`

		GLOBAL_ENV.SetCurrentNamespace(GLOBAL_ENV.CoreNamespace)
		content, err := ioutil.ReadFile("data/" + f.filename)
		if err != nil {
			panic(err)
		}

		var statics, interns string
		statics, interns, err = CodeWriter(NewReader(bytes.NewReader(content), f.name), codeWriterEnv)
		PanicOnErr(err)

		name := f.filename[0 : len(f.filename)-5] // assumes .joke extension
		newFile := "a_" + name + "_code.go"
		if newFile <= masterFile {
			panic(fmt.Sprintf("I think Go initializes file-scopes vars alphabetically by filename, so %s must come after %s due to dependencies; rename accordingly",
				newFile, masterFile))
		}
		fileContent := strings.Replace(strings.Replace(strings.ReplaceAll(fileTemplate, "{name}", name), "{statics}", statics, 1), "{interns}", interns, 1)
		ioutil.WriteFile(newFile, []byte(fileContent), 0666)
	}

	statics := []string{}
	runtime := []func() string{}

	oldWriterEnv := codeWriterEnv

	for {
		newWriterEnv := &CodeWriterEnv{
			Need:      map[string]Finisher{},
			Generated: oldWriterEnv.Generated,
		}

		env := &CodeEnv{
			CodeWriterEnv: newWriterEnv,
			Namespace:     nil,
		}

		for name, obj := range oldWriterEnv.Need {
			if _, ok := newWriterEnv.Generated[name]; ok {
				continue
			}
			s := obj.Finish(name, env)
			newWriterEnv.Generated[name] = struct{}{}
			if env.Interns != "" {
				panic("non-null interns for a_code.go")
			}
			if s != "" {
				statics = append(statics, s)
			}
		}

		runtime = append(runtime, env.Runtime...)
		if len(env.Runtime) == 0 {
			break
		}

		oldWriterEnv = newWriterEnv
		fmt.Printf("ONE!! MORE!! TIME!!\n")
	}

	// 		bindingDefs = append(bindingDefs, fmt.Sprintf(`
	// var binding_%s = Binding{
	// 	name: sym_%s,
	// 	index: %d,
	// 	frame: %d,
	// 	isUsed: %v,
	// }`[1:],
	// 			id, symName, b.Index(), b.Frame(), b.IsUsed()))

	// 		codeWriterEnv.NeedSyms[b.SymName()] = b.Symbol()
	// 	}
	// 	sort.Strings(bindingDefs)

	// 	symDefs := []string{}
	// 	symInterns := []string{}
	// 	for s, sym := range codeWriterEnv.NeedSyms {
	// 		name := NameAsGo(*s)

	// 		fields := []string{}
	// 		fields = InfoHolderField(name, sym.InfoHolder, fields, codeEnv)
	// 		fields = MetaHolderField(name, sym.MetaHolder, fields, codeEnv)
	// 		meta := strings.Join(fields, "\n")
	// 		if !IsGoExprEmpty(meta) {
	// 			meta = "\n" + meta + "\n"
	// 		}

	// 		symDefs = append(symDefs, fmt.Sprintf(`
	// var sym_%s = Symbol{%s}`[1:],
	// 			name, meta))

	// 		codeWriterEnv.NeedStrs[*s] = struct{}{}
	// 		symInterns = append(symInterns, fmt.Sprintf(`
	// 	sym_%s.name = s_%s`[1:],
	// 			name, name))
	// 	}
	// 	sort.Strings(symDefs)
	// 	sort.Strings(symInterns)

	// 	kwDefs := []string{}
	// 	kwHashes := []string{}
	// 	for _, k := range codeWriterEnv.NeedKeywords {
	// 		strName := "s_" + NameAsGo(*k.NameField())

	// 		strNs := "nil"
	// 		if k.NsField() != nil {
	// 			ns := *k.NsField()
	// 			nsName := NameAsGo(ns)
	// 			strNs = "s_" + nsName
	// 		}

	// 		name := "kw_" + k.UniqueId()

	// 		initNs := ""
	// 		if strNs != "nil" {
	// 			initNs = fmt.Sprintf(`
	// 	%s.ns = %s
	// `[1:],
	// 				name, strNs)
	// 		}

	// 		fields := []string{}
	// 		fields = InfoHolderField(name, k.InfoHolder, fields, codeEnv)
	// 		meta := strings.Join(fields, "\n")
	// 		if !IsGoExprEmpty(meta) {
	// 			meta = "\n" + meta + "\n"
	// 		}

	// 		kwDefs = append(kwDefs, fmt.Sprintf(`
	// var %s = Keyword{%s}`[1:],
	// 			name, meta))

	// 		kwHashes = append(kwHashes, fmt.Sprintf(`
	// %s	%s.name = %s
	// 	%s.hash = hashSymbol(%s, %s)`[1:],
	// 			initNs, name, strName, name, strNs, strName))
	// 	}
	// 	sort.Strings(kwDefs)
	// 	sort.Strings(kwHashes)

	// 	strDefs := []string{}
	// 	strInterns := []string{}
	// 	for s, _ := range codeWriterEnv.NeedStrs {
	// 		name := NameAsGo(s)
	// 		strDefs = append(strDefs, fmt.Sprintf(`
	// var s_%s *string`[1:],
	// 			name))

	// 		strInterns = append(strInterns, fmt.Sprintf(`
	// 	s_%s = STRINGS.Intern("%s")`[1:],
	// 			name, s))
	// 	}

	sort.Strings(statics)
	r := JoinStringFns(runtime)

	var tr = [][2]string{
		{"{statics}", strings.Join(statics, "")},
		{"{runtime}", r},
	}

	fileContent := `// Generated by gen_code. Don't modify manually!

package core

{statics}

func init() {
{runtime}
}
`

	for _, t := range tr {
		fileContent = strings.Replace(fileContent, t[0], t[1], 1)
	}

	ioutil.WriteFile(masterFile, []byte(fileContent), 0666)
}
