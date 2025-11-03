#include <assert.h>
#include <setjmp.h>
#include <stdio.h>
#include <stdlib.h>
#include <divbox.h>

void *div_info_init(void);
void *div_jmpbuf(void *div_info);
int divide(void *div_info_p, int a, int b);

int main() {
    void *div_info = div_info_init();

    jmp_buf host_buf;
    if (setjmp(host_buf)) {
        printf("error in div()\n");
        exit(1);
    }
    divbox_setjmp(host_buf, div_jmpbuf(div_info));

    int r = divide(div_info, 16, 4);
    printf("16 / 4 = %d\n", r);
    assert(r == 4);

    r = divide(div_info, 16, 0);
    printf("16 / 0 = %d\n", r);
    assert(!"unreachable");
}
