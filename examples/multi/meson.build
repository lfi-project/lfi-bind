subdir('libfoo')
subdir('libbar')

libfoo_wrap = custom_target('libfoo-wrapper',
  output: ['foo_trampolines.S', 'foo_init.c'],
  input: libfoo,
  command: [lfibind, '-gen-trampolines', '@OUTPUT0@', '-gen-init', '@OUTPUT1@', '-lib', 'foo', '-symbols', 'foo', '@INPUT@'])

libbar_wrap = custom_target('libbar-wrapper',
  output: ['bar_trampolines.S', 'bar_init.c'],
  input: libbar,
  command: [lfibind, '-gen-trampolines', '@OUTPUT0@', '-gen-init', '@OUTPUT1@', '-lib', 'bar', '-symbols', 'bar', '@INPUT@'])

multi = executable(
  'multi',
  [libfoo_wrap[0], libfoo_wrap[1], libbar_wrap[0], libbar_wrap[1], files('main.c')],
  dependencies: lfi_linux.as_link_whole())

test('multi', multi)
