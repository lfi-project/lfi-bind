libadd_so = custom_target('libadd.lfi.so',
  output: 'libadd.lfi.so',
  input: files('add.c'),
  command: [lficc, '-shared', '@INPUT@', '-o', '@OUTPUT@', '-O2'])
