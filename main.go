package main

import (
	"debug/elf"
	"embed"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"strings"
)

//go:embed embed
var files embed.FS

func fatal(v ...interface{}) {
	fmt.Fprintln(os.Stderr, v...)
	os.Exit(1)
}

func FindDynamicSymbols(input string, symPrefix string) []string {
	var syms []string

	f, err := elf.Open(input)
	if err != nil {
		log.Fatalf("Failed to open ELF file: %v", err)
	}
	defer f.Close()

	symbols, err := f.DynamicSymbols()
	if err != nil {
		log.Fatalf("Failed to read dynamic symbols: %v", err)
	}

	for _, sym := range symbols {
		if elf.ST_TYPE(sym.Info) == elf.STT_FUNC && elf.ST_BIND(sym.Info) == elf.STB_GLOBAL && sym.Section != elf.SHN_UNDEF {
			if strings.HasPrefix(sym.Name, symPrefix) {
				syms = append(syms, sym.Name)
			}
		}
	}
	return syms
}

func ReadEmbed(s string) string {
	data, err := files.ReadFile(s)
	if err != nil {
		fatal(err)
	}
	return string(data)
}

func ExecTemplate(w io.Writer, name string, data string, vars map[string]any, funcs template.FuncMap) {
	tmpl := template.New(name)
	tmpl.Funcs(funcs)
	tmpl, err := tmpl.Parse(data)
	if err != nil {
		fatal(err)
	}
	err = tmpl.Execute(w, vars)
	if err != nil {
		fatal(err)
	}
}

func GenTrampolines(file string, syms []string, lib, libPrefix string) {
	w, err := os.Create(file)
	if err != nil {
		fatal(err)
	}

	ExecTemplate(w, file, ReadEmbed("embed/lib_trampolines.S.in"), map[string]any{
		"lib":        lib,
		"lib_prefix": libPrefix,
		"syms":       syms,
	}, nil)

	w.Close()
}

func GenInit(file string, syms []string, lib, libPath string) {
	w, err := os.Create(file)
	if err != nil {
		fatal(err)
	}

	ExecTemplate(w, file, ReadEmbed("embed/lib_init.c.in"), map[string]any{
		"lib":      lib,
		"lib_path": libPath,
		"syms":     syms,
	}, nil)

	w.Close()
}

func main() {
	lib := flag.String("lib", "lib", "library name for function prefixes")
	libPath := flag.String("lib-path", "", "path to library executable at runtime")
	genTrampolines := flag.String("gen-trampolines", "", "output file for trampolines")
	genInit := flag.String("gen-init", "", "output file for initialization functions")
	symbols := flag.String("symbols", "", "comma-separated list of exported symbols")
  symbolsFile := flag.String("symbols-file", "", "list of symbols in a file, one line per symbol")
	symPrefix := flag.String("sym-prefix", "", "prefix used to match exported symbols")
	libPrefix := flag.String("lib-prefix", "", "prefix to put on library symbols")

	flag.Parse()

	args := flag.Args()

	if len(args) <= 0 {
		fatal("no input")
	}

  var syms []string

  if *symbolsFile != "" {
    data, err := os.ReadFile(*symbolsFile)
    if err != nil {
      log.Fatal(err)
    }
    lines := strings.Split(string(data), "\n")
    for _, l := range lines {
      l = strings.TrimSpace(l)
      if l != "" {
        syms = append(syms, l)
      }
    }
	} else if *symbols != "" {
		syms = strings.Split(*symbols, ",")
	} else {
    syms = FindDynamicSymbols(args[0], *symPrefix)
  }

	if *genTrampolines != "" {
		GenTrampolines(*genTrampolines, syms, *lib, *libPrefix)
	}

	if *genInit != "" {
		if *libPath != "" {
			*libPath = args[0]
		}
		GenInit(*genInit, syms, *lib, *libPath)
	}
}
