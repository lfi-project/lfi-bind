libfoo = custom_target('libfoo.lfi',
  output: 'libfoo.lfi',
  input: files('foo.c'),
  command: [lficc, '-static-pie', '@INPUT@', '-o', '@OUTPUT@', '-lboxrt', '-Wl,--export-dynamic', '-O2'])
