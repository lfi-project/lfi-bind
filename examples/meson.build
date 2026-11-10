cpu = target_machine.cpu_family()

lficc = find_program(cpu + '_lfi-linux-musl-clang', required: false)
if not lficc.found()
  warning('could not find LFI compiler')
  subdir_done()
endif

lfi_linux = subproject('lfi-runtime').get_variable('lfi_linux')

# The runtime has no mmap by default (only with -Dlinux_full=true), so
# libraries that allocate memory link minialloc, a brk-based allocator, in
# place of the libc allocator.
minialloc_srcs = subproject('minialloc').get_variable('minialloc_srcs')

subdir('add')
subdir('add-vscript')
subdir('callback')
subdir('memory')
subdir('thread')
# Disabled due to lib_engine supporting only 1 sandbox.
# subdir('multi')
# Disabled due to setjmp/longjmp support being experimental.
# subdir('setjmp')
