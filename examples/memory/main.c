#include <stddef.h>
#include <assert.h>
#include <unistd.h>

#include "membox.h"

void lib_memset(int *p, int v, size_t n);
int lib_sum(int *p, size_t n);

int main() {
    size_t n = 3;
    int *p = membox_malloc(n * sizeof(int));
    assert(p);

    p[0] = 1;
    p[1] = 2;
    p[2] = 3;

    int s;
    s = lib_sum(p, n);
    assert(s == 6);

    lib_memset(p, 1, n);
    s = lib_sum(p, n);
    assert(s == 3);

    membox_free(p);

    size_t pagesize = getpagesize();
    p = membox_mmap(NULL, pagesize, LFI_PROT_READ | LFI_PROT_WRITE, LFI_MAP_ANONYMOUS | LFI_MAP_PRIVATE, -1, 0);
    assert(p != (int *) -1);

    s = lib_sum(p, n);
    assert(s == 0);
    lib_memset(p, 2, n);
    s = lib_sum(p, n);
    assert(s == 6);

    membox_munmap(p, pagesize);
}
