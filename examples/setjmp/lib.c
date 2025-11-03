#include <setjmp.h>
#include <stdlib.h>
#include <stdio.h>

struct DivInfo {
    jmp_buf div_buf;
};

void *div_info_init(void) {
    return malloc(sizeof(struct DivInfo));
}

void *div_jmpbuf(void *div_info_p) {
    struct DivInfo *div_info = (struct DivInfo *) div_info_p;
    return div_info->div_buf;
}

int divide(void *div_info_p, int a, int b) {
    struct DivInfo *div_info = (struct DivInfo *) div_info_p;
    if (b == 0) {
        longjmp(div_info->div_buf, 1);
    }
    return a / b;
}
