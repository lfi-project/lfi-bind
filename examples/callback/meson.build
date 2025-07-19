subdir('libcb')

cb_gen = custom_target('libcb-wrapper',
  output: ['lib_trampolines.S', 'lib_init.c'],
  input: libcb,
  command: [lfibind, '-embed', '-gen-trampolines', '@OUTPUT0@', '-gen-init', '@OUTPUT1@', '-lib', 'cbbox', '-symbols', 'cb', '@INPUT@'])

cb = executable(
  'cb',
  [cb_gen[0], cb_gen[1], files('main.c')],
  dependencies: lfi_linux.as_link_whole())

test('cb', cb)
