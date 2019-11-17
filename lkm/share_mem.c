#include "share_mem.h"

static DEFINE_MUTEX(mchar_mutex);
static struct class *class;
static struct device *device;
static int major;
static char *sh_mem = NULL;
static rwlock_t _write_index_lock;

static void lock_init(void);
static inline void write_index_lock(void);
static inline void write_index_unlock(void);
static int device_mmap(struct file *filp, struct vm_area_struct *vma);
static ssize_t device_write(struct file *file,const char __user * buffer, size_t length, loff_t * offsetf);

static const struct file_operations mchar_fops = {
    .owner = THIS_MODULE,
    .mmap = device_mmap,
    .write = device_write,
};

typedef struct protect {
        char *path;
        struct list_head list;
} node;

struct protect *pprotect;
struct list_head protect_list;

#if DELAY_TEST == 1
static char *get_timespec(void)
{
    char *res = NULL;
    struct timespec tmp_time;
    res = kzalloc(64, GFP_ATOMIC);
    tmp_time = current_kernel_time();
    snprintf(res, 64, "%lu.%lu", tmp_time.tv_sec, tmp_time.tv_nsec);
    return res;
}
#endif

static void lock_init(void)
{
    rwlock_init(&_write_index_lock);
}

static inline void write_index_lock(void)
{
    write_lock(&_write_index_lock);
}

static inline void write_index_unlock(void)
{
    write_unlock(&_write_index_lock);
}

static inline void read_index_lock(void)
{
    read_lock(&_write_index_lock);
}

static inline void read_index_unlock(void)
{
    read_unlock(&_write_index_lock);
}

void del_protect_list(char *path)
{
    node *s;
    struct list_head *p;
    strim(path);
    write_index_lock();
    list_for_each(p, &protect_list) {
        s = list_entry(p, node, list);
        if (strcmp(s->path, path) == 0) {
            list_del(p);
            write_index_unlock();
            return ;
        }
    }
    write_index_unlock();
    printk("[*] Del Protect List: %s\n", path);
}

int checkpath(char *path) {
    node *s;
    struct list_head *p;
    strim(path);
    list_for_each(p, &protect_list) {
        s = list_entry(p, node, list);
        strim(s->path);
        if(strlen(s->path) > 4) {
            if(strstr((const char*)path, (const char*)s->path) != NULL) {
                return 1;
            }
        }
    }
    return 0;
}

static void add_protect_list(char *data)
{
    node *s = NULL;
    write_index_lock();
    strim(data);
    if(strlen(data) > 4) {
        s = (node *)kmalloc(sizeof(node), GFP_ATOMIC);
        s->path = data;
        list_add_tail(&(s->list), &protect_list);
        printk("[*] Add Protect List: %s\n", data);
    }
    write_index_unlock();
}

static ssize_t device_write(struct file *file,const char __user * buffer, size_t length, loff_t * offset)
{
    int i;
    int flag = 0;
    char path[PATH_MAX] = "";

    for (i = 0; i < length && i < PATH_MAX; i++) {
        get_user(path[i], buffer + i);

        if(i == 0 && strcmp(path, "+") == 0) {
            flag = 1;
            continue;
        }


        if(i == 0 && strcmp(path, "-") == 0) {
            flag = 2;
            continue;
        }
    }

    if(flag == 1)
        add_protect_list(path+1);

    if(flag == 2)
        del_protect_list(path+1);

    return i;
}

static int device_mmap(struct file *filp, struct vm_area_struct *vma)
{
    int ret = 0;
    struct page *page = NULL;
    unsigned long size = (unsigned long)(vma->vm_end - vma->vm_start);

    vma->vm_flags |= 0;

    if (size > MAX_SIZE) {
        ret = -EINVAL;
        goto out;
    }

    page = virt_to_page((unsigned long)sh_mem + (vma->vm_pgoff << PAGE_SHIFT));
    ret = remap_pfn_range(vma, vma->vm_start, page_to_pfn(page), size, vma->vm_page_prot);
    if (ret != 0) {
        goto out;
    }

out:
    return ret;
}

int protected_init(void)
{
	INIT_LIST_HEAD(&protect_list);
	return 0;
}

int init_share_mem(void)
{
    int i;
    share_mem_flag = -1;

    protected_init();

    major = register_chrdev(0, DEVICE_NAME, &mchar_fops);

    if (major < 0) {
        pr_err("[SMITH] REGISTER_CHRDEV_ERROR\n");
        return -1;
    }

    class = class_create(THIS_MODULE, CLASS_NAME);
    if (IS_ERR(class)) {
        unregister_chrdev(major, DEVICE_NAME);
        pr_err("[SMITH] CLASS_CREATE_ERROR");
        return -1;
    }

    device = device_create(class, NULL, MKDEV(major, 0), NULL, DEVICE_NAME);
    if (IS_ERR(device)) {
        class_destroy(class);
        unregister_chrdev(major, DEVICE_NAME);
        pr_err("[SMITH] DEVICE_CREATE_ERROR");
        return -1;
    }

    sh_mem = kzalloc(MAX_SIZE, GFP_KERNEL);

    if (sh_mem == NULL) {
        device_destroy(class, MKDEV(major, 0));
        class_destroy(class);
        unregister_chrdev(major, DEVICE_NAME);
        pr_err("[SMITH] SHMEM_INIT_ERROR\n");
        return -ENOMEM;
    }
    else {
        for (i = 0; i < MAX_SIZE; i += PAGE_SIZE)
            SetPageReserved(virt_to_page(((unsigned long)sh_mem) + i));
    }

    mutex_init(&mchar_mutex);
    lock_init();
    return 0;
}

void uninstall_share_mem(void)
{
    device_destroy(class, MKDEV(major, 0));
    class_unregister(class);
    class_destroy(class);
    unregister_chrdev(major, DEVICE_NAME);
}