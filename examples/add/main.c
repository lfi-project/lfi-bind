#include <assert.h>
#include <stdio.h>

#include "addbox.h"

int add(int, int);

int main() {
    int r = LFI_CALL(add, 10, 32);
    printf("add(10, 32) = %d\n", r);
    assert(r == 42);
}
