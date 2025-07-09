project('lfi-bind',
  ['c'],
  default_options: ['warning_level=1', 'buildtype=release'],
  meson_version: '>=0.60.3')

src = ['main.go']

go = find_program('go')

lfibind = custom_target('lfi-bind',
  output: 'lfi-bind',
  input: src,
  command: [go, 'build', '-o', '@OUTPUT@', '@INPUT@'],
  build_by_default: true,
  install: true,
  install_dir: 'bin')

subdir('examples')
