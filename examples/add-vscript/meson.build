libadd_vscript = custom_target('libadd_vscript.lfi',
  output: 'libadd_vscript.lfi',
  input: files('lib.c', 'symbols.map.txt'),
  command: [lficc, '-static-pie', '@INPUT0@', '-o', '@OUTPUT@', '-lboxrt', '-Wl,--export-dynamic', '-O2', '-Wl,--version-script,@INPUT1@'])

add_vscript_gen = custom_target('libadd_vscript-wrapper',
  output: ['lib_trampolines.S', 'lib_init.c', 'add_vscriptbox.h'],
  input: libadd_vscript,
  command: [lfibind, '-embed', '-gen-trampolines', '@OUTPUT0@', '-gen-init', '@OUTPUT1@', '-lib', 'add_vscriptbox', '-symbols-all', '@INPUT@'])

add_vscript = executable(
  'add_vscript',
  [add_vscript_gen[0], add_vscript_gen[1], add_vscript_gen[2], files('main.c')],
  dependencies: lfi_linux.as_link_whole())

test('add_vscript', add_vscript)
