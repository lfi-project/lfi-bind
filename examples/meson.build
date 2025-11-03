cpu = target_machine.cpu_family()

lficc = find_program(cpu + '_lfi-linux-musl-clang', required: false)
if not lficc.found()
  warning('could not find LFI compiler')
  subdir_done()
endif

lfi_linux = subproject('lfi-runtime').get_variable('lfi_linux')

subdir('add')
subdir('add-vscript')
subdir('multi')
subdir('callback')
subdir('memory')
subdir('thread')
# subdir('setjmp')

# Disabled because there isn't currently a good mechanism to link boxrt.lfi
# with the correct dynamic linker.
# subdir('so')
