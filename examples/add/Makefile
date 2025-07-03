include ../common.mk

main: main.c libadd/libadd.lfi
	mkdir -p gen
	$(LFIBIND) -gen-trampolines gen/lib_trampolines.S -gen-init gen/lib_init.c -lib-path $(PWD)/libadd/libadd.lfi -lib addbox -symbols add libadd/libadd.lfi
	$(CC) $< -O2 gen/lib_trampolines.S gen/lib_init.c -I$(INCLUDELFI) -o $@ $(LIBLFI)

libadd/libadd.lfi: libadd/add.c
	$(LFICC) -static-pie $< -o $@ -O2 -lboxrt -Wl,--export-dynamic

clean:
	rm -rf gen
	rm -f libadd/libadd.lfi libadd.so main

run: main
	./main

.PHONY: run clean
