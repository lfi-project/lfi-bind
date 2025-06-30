#include "lfi_linux.h"

#include <assert.h>
#include <unistd.h>
#include <sys/mman.h>

_Thread_local struct LFIContext *add_ctx;

struct LFIBox *add_box;

const char *library_file = "libadd/libadd.lfi";

static void
resolve(struct LFILinuxProc *proc, lfiptr *loc, const char *name)
{
    lfiptr addr = lfi_proc_sym(proc, name);
    *loc = addr;
}

struct Buf {
    void *data;
    size_t size;
};

static struct Buf
readfile(const char *path)
{
    FILE *f = fopen(path, "r");
    if (!f) {
        fprintf(stderr, "Cannot open %s\n", path);
        return (struct Buf) { 0 };
    }
    fseek(f, 0, SEEK_END);
    size_t sz = ftell(f);
    void *p = mmap(NULL, sz, PROT_READ, MAP_PRIVATE, fileno(f), 0);
    assert(p != (void *) -1);
    fclose(f);
    return (struct Buf) {
        .data = p,
        .size = sz,
    };
}

__attribute__((constructor)) void
add_init(void)
{
    // Create engine if it does not exist.
    bool ok = lfi_linux_lib_init((struct LFIOptions) {
        .boxsize = 4UL * 1024 * 1024 * 1024,
        .pagesize = getpagesize(),
        .verbose = false,
    }, (struct LFILinuxOptions) {
        .stacksize = 2UL * 1024 * 1024,
        .verbose = false,
    });
    assert(ok);

    // Create proc from file.
    struct Buf prog = readfile(library_file);
    assert(prog.data);

    struct LFILinuxProc *proc = lfi_proc_new(lfi_linux_lib_engine());
    assert(proc);

    ok = lfi_proc_load(proc, prog.data, prog.size);
    assert(ok);

    // Initialize return.
    lfi_box_init_ret(lfi_proc_box(proc));

    // Initialize callbacks.
    ok = lfi_box_cbinit(lfi_proc_box(proc));
    assert(ok);

    // Create and run thread.
    const char *argv[] = {
        "add",
        NULL,
    };

    const char *envp[] = {
        NULL,
    };

    struct LFILinuxThread *t = lfi_thread_new(proc,
        sizeof(argv)/sizeof(argv[0]), &argv[0], &envp[0]);
    assert(t);
    int result = lfi_thread_run(t);
    assert(result == 0);

    // Initialize clone.
    lfi_linux_init_clone(t);

    // Initialize all exported symbols.
    extern lfiptr add_addr_add;
    resolve(proc, &add_addr_add, "add");

    add_box = lfi_proc_box(proc);
    add_ctx = *lfi_thread_ctxp(t);
}
