libcb = custom_target('libcb.lfi',
  output: 'libcb.lfi',
  input: files('cb.c'),
  command: [lficc, '-static-pie', '@INPUT@', '-o', '@OUTPUT@', '-lboxrt', '-Wl,--export-dynamic', '-O2'])
