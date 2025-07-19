#include <assert.h>
#include <stdio.h>
#include <time.h>
#include <stdlib.h>

#include "cbbox.h"

int cb(int (*fn)(void));

static int fn(void) {
    printf("hello from fn!\n");
    return 32;
}

int main() {
    void *box_fn = cbbox_register_cb(&fn);
    int r = cb(box_fn);
    assert(r == 42);
    printf("cb(fn) = %d\n", r);
}
