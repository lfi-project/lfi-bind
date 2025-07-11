libbar = custom_target('libbar.lfi',
  output: 'libbar.lfi',
  input: files('bar.c'),
  command: [lficc, '-static-pie', '@INPUT@', '-o', '@OUTPUT@', '-lboxrt', '-Wl,--export-dynamic', '-O2'])
