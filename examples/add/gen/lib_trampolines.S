.text

#if defined(__aarch64__)

.macro get_ctx var, reg
    mov x11, x0
    mov x12, x30

    adrp x0, :tlsdesc:\var
    ldr \reg, [x0, :tlsdesc_lo12:\var]
    add x0, x0, :tlsdesc_lo12:\var
    .tlsdesccall \var
    blr \reg
    mrs \reg, tpidr_el0
    add \reg, \reg, x0

    mov x30, x12
    mov x0, x11
.endm

.p2align 4
.add_trampoline:
	adrp x12, :got:add_box
	ldr	x12, [x12, :got_lo12:add_box]
	ldr	x12, [x12]
    get_ctx add_ctx, x11
    stp x11, x12, [sp, #-16]!
    stp x10, xzr, [sp, #-16]!
    b lfi_trampoline_direct
.p2align 4
.global add
add:
    adrp x10, :got:add_addr_add
    ldr	x10, [x10, :got_lo12:add_addr_add]
    ldr	x10, [x10]
    b .add_trampoline

#elif defined(__x86_64__)

.macro get_ctx var, reg
#if defined(ENABLE_TLS_DESC)
    pushq %rax
    leaq \var@TLSDESC(%rip), %rax
    callq *\var@TLSCALL(%rax)
    addq %fs:0x0, %rax
    movq %rax, \reg
    popq %rax
#else
    // Save/restore registers used for arguments.
    pushq %rdi
    pushq %rsi
    pushq %rdx
    pushq %rcx
    pushq %r8
    pushq %r9
    subq $136, %rsp
    movdqu %xmm0, 0(%rsp)
    movdqu %xmm1, 16(%rsp)
    movdqu %xmm2, 32(%rsp)
    movdqu %xmm3, 48(%rsp)
    movdqu %xmm4, 64(%rsp)
    movdqu %xmm5, 80(%rsp)
    movdqu %xmm6, 96(%rsp)
    movdqu %xmm7, 112(%rsp)

    .byte 0x66
    leaq \var@TLSGD(%rip), %rdi
    .byte 0x66, 0x66
    rex64 callq __tls_get_addr@PLT
    movq %rax, \reg

    movdqu 0(%rsp), %xmm0
    movdqu 16(%rsp), %xmm1
    movdqu 32(%rsp), %xmm2
    movdqu 48(%rsp), %xmm3
    movdqu 64(%rsp), %xmm4
    movdqu 80(%rsp), %xmm5
    movdqu 96(%rsp), %xmm6
    movdqu 112(%rsp), %xmm7
    addq $136, %rsp
    popq %r9
    popq %r8
    popq %rcx
    popq %rdx
    popq %rsi
    popq %rdi
#endif
.endm

.p2align 4
.add_trampoline:
    movq add_box@GOTPCREL(%rip), %r11
    pushq (%r11)
    get_ctx add_ctx, %r11
    pushq %r11
    jmp lfi_trampoline_direct
.p2align 4
.global add
add:
    movq add_addr_add@GOTPCREL(%rip), %r10
    pushq (%r10)
    jmp .add_trampoline

#else

#error "unsupported architecture"

#endif

.data
.global add_addr_add
.hidden add_addr_add
add_addr_add:
    .quad 0

.section .note.GNU-stack,"",@progbits
