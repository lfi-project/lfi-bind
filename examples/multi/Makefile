include ../common.mk

main: main.c libfoo.a libbar.a
	$(CC) -O2 $^ $(LIBLFI) -o $@

libfoo.a: libfoo/libfoo.lfi
	mkdir -p gen
	$(LFIBIND) -gen-trampolines gen/foo_trampolines.S -gen-init gen/foo_init.c -lib-path $(PWD)/$< -lib foo -symbols foo $<
	$(CC) -c -O2 gen/foo_init.c -I$(INCLUDELFI) -fPIC -o gen/foo_init.o
	$(CC) -c -O2 gen/foo_trampolines.S -I$(INCLUDELFI) -fPIC -o gen/foo_trampolines.o
	$(AR) rcs $@ gen/foo_trampolines.o gen/foo_init.o

libbar.a: libbar/libbar.lfi
	mkdir -p gen
	$(LFIBIND) -gen-trampolines gen/bar_trampolines.S -gen-init gen/bar_init.c -lib-path $(PWD)/$< -lib bar -symbols bar $<
	$(CC) -c -O2 gen/bar_init.c -I$(INCLUDELFI) -fPIC -o gen/bar_init.o
	$(CC) -c -O2 gen/bar_trampolines.S -I$(INCLUDELFI) -fPIC -o gen/bar_trampolines.o
	$(AR) rcs $@ gen/bar_trampolines.o gen/bar_init.o

libfoo/libfoo.lfi: libfoo/foo.c
	$(LFICC) -static-pie $< -o $@ -O2 -lboxrt -Wl,--export-dynamic

libbar/libbar.lfi: libbar/bar.c
	$(LFICC) -static-pie $< -o $@ -O2 -lboxrt -Wl,--export-dynamic

run: main
	./main

clean:
	rm -rf gen
	rm -f libfoo/libfoo.lfi libbar/libbar.lfi libfoo.a libbar.a main

.PHONY: run
