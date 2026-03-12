libthread = custom_target('libthread.lfi',
  output: 'libthread.lfi',
  input: files('lib.c'),
  command: [lficc, '-static-pie', '@INPUT@', '-o', '@OUTPUT@', '-lboxrt', '-Wl,--export-dynamic', '-O2'])

thread_gen = custom_target('libthread-wrapper',
  output: ['lib_trampolines.S', 'lib_init.c', 'threadbox.h'],
  input: libthread,
  command: [lfibind, '-embed', '-gen-trampolines', '@OUTPUT0@', '-gen-init', '@OUTPUT1@', '-lib', 'threadbox', '-symbols', 'foo', '@INPUT@'])

thread = executable(
  'thread',
  [thread_gen[0], thread_gen[1], thread_gen[2], files('main.c')],
  dependencies: lfi_linux.as_link_whole())

test('thread', thread)
