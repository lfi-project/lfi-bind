subdir('libadd')

add_so_gen = custom_target('libadd-so-wrapper',
  output: ['lib_trampolines.S', 'lib_init.c'],
  input: libadd_so,
  command: [lfibind, '-gen-trampolines', '@OUTPUT0@', '-gen-init', '@OUTPUT1@', '-lib', 'addbox', '-symbols', 'add', '@INPUT@'])

add_so = executable(
  'add-so',
  [add_so_gen[0], add_so_gen[1], files('main.c')],
  dependencies: lfi_linux.as_link_whole())

test('add-so', add_so)
