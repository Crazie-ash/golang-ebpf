// +build ignore
#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>

char __license[] SEC("license") = "Dual MIT/GPL";

SEC("tracepoint/syscalls/sys_enter_execve")
int hello_bpf(void  *ctx)
{
    char msg[] = "hello from bpf module\n";
    bpf_trace_printk(msg, sizeof(msg));
    return 0;
}

// char _license[] SEC("license") = "GPL";
// u32 _version SEC("version") = LINUX_VERSION_CODE;