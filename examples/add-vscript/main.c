#include <assert.h>
#include <stdio.h>

int add(int a, int b);

int main() {
    int r = add(32, 10);
    printf("add(32, 10) = %d\n", r);
    assert(r == 42);

    return 0;
}
