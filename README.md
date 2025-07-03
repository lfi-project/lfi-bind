# LFI Bind

This program makes it easy to do library sandboxing with LFI. It generates
routines to initialize the library sandbox, and trampolines for calling
functions from the library.

# Installation

```
go build
```

# Usage

The process for creating a sandboxed library is the following:

1. Compile your static library using the LFI compiler. This produces
   `libfoo.a`.
2. Create a "LFI library" by compiling `libfoo.a` to a static PIE and linking
   with `boxrt` (startup code that runs in the sandbox). Usually this command
   looks like `$(LFICC) -Wl,--whole-archive libfoo.a -Wl,--no-whole-archive
   -Wl,--export-dynamic -lboxrt -static-pie -o libfoo.lfi`.
3. Invoke `lfi-bind`, passing it `libfoo.lfi` and instructing it to generate an
   initialization file `lib_init.c` and trampoline file `lib_trampolines.S`.
4. Compile `lib_init.c` and `lib_trampolines.S` into your host application,
   allowing you to use the sandboxed library.

See the `examples` directory for real examples that you can try out.

```
Usage of ./lfi-bind:
  -gen-init string
    	output file for initialization functions
  -gen-trampolines string
    	output file for trampolines
  -lib string
    	library name for function prefixes (default "lib")
  -lib-path string
    	path to library executable at runtime
  -lib-prefix string
    	prefix to put on library symbols
  -sym-prefix string
    	prefix used to match exported symbols
  -symbols string
    	comma-separated list of exported symbols
  -symbols-file string
    	list of symbols in a file, one line per symbol
```
