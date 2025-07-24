#include <stddef.h>

int lib_deref(int *x) {
    return *x;
}

void lib_memset(int *p, int v, size_t n) {
    for (size_t i = 0; i < n; i++) {
        p[i] = v;
    }
}

int lib_sum(int *p, size_t n) {
    int sum = 0;
    for (size_t i = 0; i < n; i++) {
        sum += p[i];
    }
    return sum;
}
