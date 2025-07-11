libadd = custom_target('libadd.lfi',
  output: 'libadd.lfi',
  input: files('add.c'),
  command: [lficc, '-static-pie', '@INPUT@', '-o', '@OUTPUT@', '-lboxrt', '-Wl,--export-dynamic', '-O2'])
