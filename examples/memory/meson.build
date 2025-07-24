libmem = custom_target('libmem.lfi',
  output: 'libmem.lfi',
  input: files('lib.c'),
  command: [lficc, '-static-pie', '@INPUT@', '-o', '@OUTPUT@', '-lboxrt', '-Wl,--export-dynamic', '-O2'])

mem_gen = custom_target('libmem-wrapper',
  output: ['lib_trampolines.S', 'lib_init.c', 'membox.h'],
  input: libmem,
  command: [lfibind, '-embed', '-gen-trampolines', '@OUTPUT0@', '-gen-init', '@OUTPUT1@', '-lib', 'membox', '-symbols', 'lib_memset,lib_sum,lib_deref', '@INPUT@'])

mem = executable(
  'mem',
  [mem_gen[0], mem_gen[1], mem_gen[2], files('main.c')],
  dependencies: lfi_linux.as_link_whole())

test('mem', mem)
