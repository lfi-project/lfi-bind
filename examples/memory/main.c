#include <stddef.h>
#include <assert.h>
#include <unistd.h>

#include "membox.h"

int lib_deref(int *p);
void lib_memset(int *p, int v, size_t n);
int lib_sum(int *p, size_t n);

int main() {
    size_t n = 3;
    int *p = membox_malloc(n * sizeof(int));
    assert(p);

    p[0] = 1;
    p[1] = 2;
    p[2] = 3;

    // In this example, we use LFI_CALL to call library functions because it
    // will perform additional debug checks to make sure input pointers given
    // to the sandbox are in-bounds (and therefore won't crash when accessed).
    // Define LFI_NODEBUG before including 'membox.h' to disable the checks.
    int s;
    s = LFI_CALL(lib_sum, p, n);
    assert(s == 6);

    LFI_CALL(lib_memset, p, 1, n);
    s = LFI_CALL(lib_sum, p, n);
    assert(s == 3);

    membox_free(p);

    size_t pagesize = getpagesize();
    p = membox_mmap(NULL, pagesize, LFI_PROT_READ | LFI_PROT_WRITE, LFI_MAP_ANONYMOUS | LFI_MAP_PRIVATE, -1, 0);
    assert(p != (int *) -1);

    s = LFI_CALL(lib_sum, p, n);
    assert(s == 0);
    LFI_CALL(lib_memset, p, 2, n);
    s = LFI_CALL(lib_sum, p, n);
    assert(s == 6);

    membox_munmap(p, pagesize);

    int *x = membox_stack_push(sizeof(int));
    *x = 42;
    s = LFI_CALL(lib_deref, x);
    membox_stack_pop(sizeof(*x));
    assert(s == 42);
}
