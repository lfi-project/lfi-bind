#include <assert.h>
#include <stdio.h>
#include <stdlib.h>

int add(int, int);

int main() {
    int r = add(10, 32);
    assert(r == 42);
    printf("add(10, 32) = %d\n", r);
}
