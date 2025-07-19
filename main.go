package main

import (
	"bytes"
	"debug/elf"
	"embed"
	"encoding/binary"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//go:embed embed
var files embed.FS

func fatal(v ...interface{}) {
	fmt.Fprintln(os.Stderr, v...)
	os.Exit(1)
}

type StackArgInfo struct {
	Sret uint32
	Args []StackArg
}

type StackArg struct {
	Offset uint32
	Size   uint32
}

func ObjGetStackArgs(file *elf.File) (map[string]StackArgInfo, bool) {
	sec := file.Section(".stack_args")
	if sec == nil {
		return nil, false
	}

	syms, err := file.Symbols()
	if err != nil {
		log.Fatal(err)
	}
	symtab := make(map[uint64]string)
	for _, sym := range syms {
		symtab[sym.Value] = sym.Name
	}

	info := make(map[string]StackArgInfo)

	b64 := make([]byte, 8)
	b32 := make([]byte, 4)
	idx := uint64(0)
	for idx < sec.Size {
		sec.ReadAt(b64, int64(idx))
		idx += 8
		fn := binary.LittleEndian.Uint64(b64)

		sec.ReadAt(b32, int64(idx))
		idx += 4
		sret := binary.LittleEndian.Uint32(b32)

		sec.ReadAt(b32, int64(idx))
		idx += 4
		entries := binary.LittleEndian.Uint32(b32)

		var args []StackArg
		for i := uint32(0); i < entries; i++ {
			// stack offset
			sec.ReadAt(b32, int64(idx))
			idx += 4
			offset := binary.LittleEndian.Uint32(b32)
			// size
			sec.ReadAt(b32, int64(idx))
			idx += 4
			size := binary.LittleEndian.Uint32(b32)

			args = append(args, StackArg{
				Offset: offset,
				Size:   size,
			})
		}

		sym := symtab[fn]
		info[sym] = StackArgInfo{
			Sret: sret,
			Args: args,
		}
	}

	return info, true
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

type Options struct {
	Input     string
	Syms      []string
	Lib       string
	LibPrefix string
	LibPath   string
	Dynamic   bool
	Embed     bool
	NoVerify  bool
	StackArgs map[string]StackArgInfo
}

// Returns sret, nstack.
func GetStackInfo(stackArgs map[string]StackArgInfo, s string, warn bool) (int, int) {
	if stackArgs == nil {
		return 0, 0
	}

	info, ok := stackArgs[s]
	if !ok {
		return 0, 0
	}

	if info.Sret != 0 && warn {
		fmt.Fprintf(os.Stderr, "warning: %s has struct return (unsupported)\n", s)
	}

	args := info.Args
	n := 0
	for _, a := range args {
		n += int(a.Size)
	}
	if n != 0 && warn {
		fmt.Fprintf(os.Stderr, "warning: %s has %d bytes of stack arguments (unsupported)\n", s, n)
	}
	return int(info.Sret), n
}

func GenTrampolines(file string, opts Options) {
	w, err := os.Create(file)
	if err != nil {
		fatal(err)
	}

	for _, s := range opts.Syms {
		GetStackInfo(opts.StackArgs, s, true)
	}

	ExecTemplate(w, file, ReadEmbed("embed/lib_trampolines.S.in"), map[string]any{
		"lib":        opts.Lib,
		"lib_prefix": opts.LibPrefix,
		"syms":       opts.Syms,
	}, map[string]any{
		"n_stack_args": func(s string) int {
			_, n := GetStackInfo(opts.StackArgs, s, false)
			return n
		},
	})

	w.Close()
}

func GenInit(file string, opts Options) {
	w, err := os.Create(file)
	if err != nil {
		fatal(err)
	}

	embedData := ""
	if opts.Embed {
		data, err := os.ReadFile(opts.Input)
		if err != nil {
			fatal(err)
		}
		buf := &bytes.Buffer{}
		for _, b := range data {
			fmt.Fprintf(buf, "%d,", b)
		}
		embedData = buf.String()
	}

	ExecTemplate(w, file, ReadEmbed("embed/lib_init.c.in"), map[string]any{
		"lib":        opts.Lib,
		"lib_path":   opts.LibPath,
		"syms":       opts.Syms,
		"dynamic":    opts.Dynamic,
		"no_verify":  opts.NoVerify,
		"embed":      opts.Embed,
		"embed_data": embedData,
	}, nil)

	w.Close()
}

func GenInitHeader(file string, lib string) {
	w, err := os.Create(file)
	if err != nil {
		fatal(err)
	}

	ExecTemplate(w, file, ReadEmbed("embed/lib.h.in"), map[string]any{
		"lib": lib,
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
	embedF := flag.Bool("embed", false, "fully embed the input library into the data segment")
	noVerify := flag.Bool("no-verify", false, "disable verification")

	flag.Parse()

	args := flag.Args()

	if len(args) <= 0 {
		fatal("no input")
	}

	input := args[0]

	dynamic := false
	if strings.HasSuffix(input, ".so") {
		dynamic = true

		if *embedF {
			fatal("-embed is not supported with shared libraries")
		}
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
	}
	if *symbols != "" {
		syms = append(syms, strings.Split(*symbols, ",")...)
	}

	if *symbols == "" && *symbolsFile == "" {
		syms = FindDynamicSymbols(input, *symPrefix)
	}

	if *libPath == "" {
		*libPath = input
	}

	f, err := elf.Open(input)
	if err != nil {
		fatal("failed to open ELF file: %w", err)
	}
	stackArgs, ok := ObjGetStackArgs(f)
	if !ok {
		fmt.Fprintln(os.Stderr, "warning: no .stack_args section found")
	}
	f.Close()

	opts := Options{
		Input:     input,
		Syms:      syms,
		Lib:       *lib,
		LibPrefix: *libPrefix,
		LibPath:   *libPath,
		Dynamic:   dynamic,
		Embed:     *embedF,
		NoVerify:  *noVerify,
		StackArgs: stackArgs,
	}

	if *genTrampolines != "" {
		GenTrampolines(*genTrampolines, opts)
	}

	if *genInit != "" {
		GenInit(*genInit, opts)
		GenInitHeader(filepath.Join(filepath.Dir(*genInit), *lib+".h"), *lib)
	}
}
