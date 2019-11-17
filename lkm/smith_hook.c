
#include "share_mem.h"
#include "smith_hook.h"
#include "struct_wrap.h"

#define EXIT_PROTECT 0

int share_mem_flag = -1;
char fsnotify_kprobe_state = 0x0;

static char *str_replace(char *orig, char *rep, char *with)
{
    char *result, *ins, *tmp;
    int len_rep, len_with, len_front, count;

    if (!orig || !rep)
        return NULL;

    len_rep = strlen(rep);
    if (len_rep == 0)
        return NULL;

    if (!with)
        with = "";

    len_with = strlen(with);

    ins = orig;
    for (count = 0; (tmp = strstr(ins, rep)); ++count)
        ins = tmp + len_rep;

    tmp = result = kzalloc(strlen(orig) + (len_with - len_rep) * count + 1, GFP_ATOMIC);

    if (!result)
        return NULL;

    while (count--) {
        ins = strstr(orig, rep);
        len_front = ins - orig;
        tmp = strncpy(tmp, orig, len_front) + len_front;
        tmp = strcpy(tmp, with) + len_with;
        orig += len_front + len_rep;
    }

    strcpy(tmp, orig);
    return result;
}

static inline int _kill_task_by_task(struct task_struct *p_task, char *path) {
   int uid = get_current_uid();
   if(current->pid > 1000 && current->real_parent->pid > 1000) {
      printk("[!!!] Don't Touch Me(%s) %d|%d!\n",path, uid, current->pid);
      return send_sig_info(SIGKILL, SEND_SIG_PRIV, p_task);
   }
   return 0;
}

static inline int _kill_task_by_pid(pid_t p_pid) {
   return _kill_task_by_task(pid_task(find_vpid(p_pid), PIDTYPE_PID), "");
}

static void fsnotify_post_handler(struct kprobe *p, struct pt_regs *regs, unsigned long flags)
{
    struct path *path;
    __u32 flag = (__u32)p_get_arg2(regs);
    if (flag == FS_OPEN || flag == FS_ACCESS) {
        char buffer[PATH_MAX];
        memset(buffer, 0, sizeof(PATH_MAX));
        path = (struct path *)p_get_arg3(regs);
        char *pathstr = dentry_path_raw(path->dentry, buffer, PATH_MAX);

        if(strlen(pathstr) > 5 && checkpath(pathstr) == 1)
            _kill_task_by_task(current, pathstr);
    }
}

static struct kprobe fsnotify_kprobe = {
    .symbol_name = "fsnotify",
	.post_handler = fsnotify_post_handler,
};

static int fsnotify_register_kprobe(void)
{
	int ret;
	ret = register_kprobe(&fsnotify_kprobe);

	if (ret == 0)
        fsnotify_kprobe_state = 0x1;

	return ret;
}

static void unregister_kprobe_fsnotify(void)
{
	unregister_kprobe(&fsnotify_kprobe);
}

static void uninstall_kprobe(void)
{
    if (fsnotify_kprobe_state == 0x1)
	    unregister_kprobe_fsnotify();
}

static int __init smith_init(void)
{
	int ret;

    ret = init_share_mem();

    if (ret != 0)
        return ret;
    else
        printk(KERN_INFO "[SMITH] init_share_mem success \n");

    ret = fsnotify_register_kprobe();
	if (ret < 0) {
	    uninstall_kprobe();
	    uninstall_share_mem();
	    printk(KERN_INFO "[SMITH] fsnotify register_kprobe failed, returned %d\n", ret);
	    return -1;
	}

	return 0;
}

static void __exit smith_exit(void)
{
	uninstall_kprobe();
	uninstall_share_mem();
	printk(KERN_INFO "[SMITH] uninstall_kprobe success\n");
}

module_init(smith_init)
module_exit(smith_exit)

MODULE_LICENSE("GPL v2");
MODULE_VERSION("0.0.1");
MODULE_AUTHOR("E_Bwill <cy_sniper@yeah.net>");