#include <unistd.h>
#include <stdio.h>

void foo(void) {
    // Note: SYS_MINIMAL does not support TLS, and printf uses TLS.
    /* printf("foo\n"); */
}
