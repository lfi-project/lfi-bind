#include <assert.h>
#include <stdio.h>
#include <time.h>
#include <stdlib.h>

int add(int, int);

static inline long long unsigned time_ns() {
    struct timespec ts;
    if (clock_gettime(CLOCK_REALTIME, &ts)) {
        exit(1);
    }
    return ((long long unsigned) ts.tv_sec) * 1000000000LLU +
        (long long unsigned) ts.tv_nsec;
}

int main() {
    int r = add(10, 32);
    assert(r == 42);
    printf("add(10, 32) = %d\n", r);

    size_t iters = 10000000;
    long long unsigned start = time_ns();
    for (size_t i = 0; i < iters; i++) {
        add(10, 32);
    }
    long long unsigned elapsed = time_ns() - start;
    printf("time per add: %.1f ns\n", (float) elapsed / (float) iters);
}
