include ../common.mk

main: main.c libadd.so
	$(CC) -O2 $^ $(LIBLFI) -o $@

libadd.so: libadd/libadd.lfi
	mkdir -p gen
	$(LFIBIND) -gen-trampolines gen/lib_trampolines.S -gen-init gen/lib_init.c -lib-path $(PWD)/$< -lib add -symbols add $<
	$(CC) -shared -O2 gen/lib_trampolines.S gen/lib_init.c -I$(INCLUDELFI) -fPIC -o $@

libadd/libadd.lfi: libadd/add.c
	$(LFICC) -static-pie $< -o $@ -O2 -lboxrt -Wl,--export-dynamic

clean:
	rm -rf gen
	rm -f libadd/libadd.lfi libadd.so main

run: main
	@LD_LIBRARY_PATH=.:$(ROOT)/examples/lfi-runtime/install/lib ./main

.PHONY: run clean
