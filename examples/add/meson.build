subdir('libadd')

gen = custom_target('libadd-wrapper',
  output: ['lib_trampolines.S', 'lib_init.c', 'addbox.h'],
  input: libadd,
  command: [lfibind, '-embed', '-gen-trampolines', '@OUTPUT0@', '-gen-init', '@OUTPUT1@', '-lib', 'addbox', '-symbols', 'add', '@INPUT@'])

add = executable(
  'add',
  [gen[0], gen[1], gen[2], files('main.c')],
  dependencies: lfi_linux.as_link_whole())

test('add', add)
