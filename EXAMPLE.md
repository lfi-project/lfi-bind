# Manually creating an LFI bind example

This guide will explain how to manually recreate the example in `examples/add`.

First make sure that you have an LFI toolchain available. You can get a
prebuilt one from here:
https://github.com/lfi-project/lfi/releases/latest.

You will also need `liblfi.a`. You can get that from the prebuilts listed above
(`aarch64-lfi-tools/lib` should contain it), or you can build it from source
from https://github.com/lfi-project/lfi-runtime.

You also need `lfi-bind`, which you can build by running `go build` in this
repository.

## Steps

Create your sandboxed library function in `add.c`.

```c
int add(int a, int b) {
    return a + b;
}
```

Compile this into a library.

```
$ aarch64-lfi-linux-musl-clang add.c -c -O2 -fPIC
$ llvm-ar rcs libadd.a add.o
```

Create an "LFI library" by compiling `libadd.a` into a static PIE and linking
with `boxrt`.

```
$ aarch64-lfi-linux-musl-clang -Wl,--whole-archive libadd.a -Wl,--no-whole-archive -Wl,--export-dynamic -lboxrt -static-pie -o libadd.lfi
```

Next invoke `lfi-bind` to create trampolines and initialization code.

```
$ lfi-bind -embed -gen-trampolines lib_trampolines.S -gen-init lib_init.c -lib addbox -symbols add libadd.lfi
```

This will create `lib_trampolines.S` and `lib_init.c` and expose `add` from `libadd.lfi`. It will also embed `libadd.lfi` inside the generated code (because we passed `-embed`).

Finally create your host program.

```c
#include <stdio.h>

int add(int, int);

int main() {
    printf("add(10, 32) = %d\n", add(10, 32));
}
```

Compile it with the host compiler and pass in `lib_trampolines.S` and `lib_init.c` and link with liblfi:

```
$ cc main.c lib_trampolines.S lib_init.c -O2 /path/to/liblfi.a -o main
$ ./main
add(10, 32) = 42
```

You can use `LFI_VERBOSE=1` to enable verbose mode:

```
$ LFI_VERBOSE=1 ./main
[lfi-core] initialized LFI engine: 80 GiB
[lfi-linux] initialized LFI Linux engine
[lfi-linux] elf_load [0xffed00015000, 0xffed00025000] (P: 1)
[lfi-linux] elf_load [0xffed00025000, 0xffed00035000] (P: 5)
[lfi-linux] elf_load [0xffed00035000, 0xffed00045000] (P: 3)
[lfi-linux] elf_load [0xffed00045000, 0xffed00055000] (P: 3)
[lfi-linux] sys_mmap(ffed00055000 (ffed00055000), 4096, 0, 50, -1, 0) = ffed00055000
[lfi-linux] sys_mmap(ffed20076000 (0), 4096, 3, 34, -1, 0) = ffed20076000
[lfi-linux] rt_sigprocmask: ignored
[lfi-linux] sys_mmap(ffed20077000 (0), 143360, 0, 34, -1, 0) = ffed20077000
[lfi-linux] rt_sigprocmask: ignored
[lfi-linux] sys_clone(3d0f00, ffed20099ef0, ffed20099f40, ffed00045270, ffed20099fe8) = 10001
[lfi-linux] rt_sigprocmask: ignored
[lfi-linux] rt_sigprocmask: ignored
add(10, 32) = 42
```
