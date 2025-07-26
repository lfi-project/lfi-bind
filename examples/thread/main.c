#include <stddef.h>
#include <assert.h>
#include <unistd.h>
#include <pthread.h>

#include "threadbox.h"

void foo(void);

static void *threadfn(void *arg) {
    foo();

    return NULL;
}

int main() {
    pthread_t t1, t2;
    pthread_create(&t1, NULL, threadfn, NULL);
    pthread_create(&t2, NULL, threadfn, NULL);
    pthread_join(t1, NULL);
    pthread_join(t2, NULL);
}
