libdiv = custom_target('libdiv.lfi',
  output: 'libdiv.lfi',
  input: files('lib.c'),
  command: [lficc, '-static-pie', '@INPUT0@', '-o', '@OUTPUT@', '-lboxrt', '-Wl,--export-dynamic', '-O2'])

div_gen = custom_target('libdiv-wrapper',
  output: ['lib_trampolines.S', 'lib_init.c', 'divbox.h'],
  input: libdiv,
  command: [lfibind, '-embed', '-gen-trampolines', '@OUTPUT0@', '-gen-init', '@OUTPUT1@', '-lib', 'divbox', '-symbols', 'divide,div_info_init,div_jmpbuf', '@INPUT@'])

div = executable(
  'div',
  [div_gen[0], div_gen[1], div_gen[2], files('main.c')],
  dependencies: lfi_linux.as_link_whole())

test('div', div)
