#ifdef CONFIG_X86
static inline unsigned long p_regs_get_arg1(struct pt_regs *p_regs) {
   return p_regs->di;
}

static inline unsigned long p_regs_get_arg2(struct pt_regs *p_regs) {
   return p_regs->si;
}

static inline unsigned long p_regs_get_arg3(struct pt_regs *p_regs) {
   return p_regs->dx;
}

static inline unsigned long p_regs_get_arg4(struct pt_regs *p_regs) {
   return p_regs->r10;
}

static inline unsigned long p_regs_get_arg5(struct pt_regs *p_regs) {
   return p_regs->r8;
}

static inline unsigned long p_regs_get_arg6(struct pt_regs *p_regs) {
   return p_regs->r9;
}

static inline unsigned long p_get_arg1(struct pt_regs *p_regs) {
#if LINUX_VERSION_CODE >= KERNEL_VERSION(4,17,0) && defined(CONFIG_ARCH_HAS_SYSCALL_WRAPPER)
   return p_regs_get_arg1((struct pt_regs *)p_regs_get_arg1(p_regs));
#else
   return p_regs_get_arg1(p_regs);
#endif
}

static inline unsigned long p_get_arg2(struct pt_regs *p_regs) {
#if LINUX_VERSION_CODE >= KERNEL_VERSION(4,17,0) && defined(CONFIG_ARCH_HAS_SYSCALL_WRAPPER)
   return p_regs_get_arg2((struct pt_regs *)p_regs_get_arg1(p_regs));
#else
   return p_regs_get_arg2(p_regs);
#endif
}

static inline unsigned long p_get_arg3(struct pt_regs *p_regs) {
#if LINUX_VERSION_CODE >= KERNEL_VERSION(4,17,0) && defined(CONFIG_ARCH_HAS_SYSCALL_WRAPPER)
   return p_regs_get_arg3((struct pt_regs *)p_regs_get_arg1(p_regs));
#else
   return p_regs_get_arg3(p_regs);
#endif
}

static inline unsigned long p_get_arg4(struct pt_regs *p_regs) {
#if LINUX_VERSION_CODE >= KERNEL_VERSION(4,17,0) && defined(CONFIG_ARCH_HAS_SYSCALL_WRAPPER)
   return p_regs_get_arg4((struct pt_regs *)p_regs_get_arg1(p_regs));
#else
   return p_regs_get_arg4(p_regs);
#endif
}

static inline unsigned long p_get_arg5(struct pt_regs *p_regs) {
#if LINUX_VERSION_CODE >= KERNEL_VERSION(4,17,0) && defined(CONFIG_ARCH_HAS_SYSCALL_WRAPPER)
   return p_regs_get_arg5((struct pt_regs *)p_regs_get_arg1(p_regs));
#else
   return p_regs_get_arg5(p_regs);
#endif
}

static inline unsigned long p_get_arg6(struct pt_regs *p_regs) {
#if LINUX_VERSION_CODE >= KERNEL_VERSION(4,17,0) && defined(CONFIG_ARCH_HAS_SYSCALL_WRAPPER)
   return p_regs_get_arg6((struct pt_regs *)p_regs_get_arg1(p_regs));
#else
   return p_regs_get_arg6(p_regs);
#endif
}
#endif

static inline int get_current_uid(void) {
#if LINUX_VERSION_CODE >= KERNEL_VERSION(3,5,0)
    return current->real_cred->uid.val;
#else
    return current->real_cred->uid;
#endif
}