#include <linux/tcp.h>
#include <linux/ip.h>
#include <linux/udp.h>
#include <linux/inet.h>
#include <linux/pid.h>
#include <linux/sched.h>

#define MY_PORT 4040
#define PROC_NAME "myprocess"

int filter_port(struct __sk_buff *skb) {
    u8 *cursor = 0;

    struct ethernet_t *ethernet = cursor_advance(cursor, sizeof(*ethernet));
    if (!(ethernet->type == 0x0800)) { // IP
        return 0; // Pass non-IP packets
    }

    struct ip_t *ip = cursor_advance(cursor, sizeof(*ip));
    if (ip->nextp != IPPROTO_TCP) {
        return 0; // Pass non-TCP packets
    }

    struct tcp_t *tcp = cursor_advance(cursor, sizeof(*tcp));
    u16 dport = ntohs(tcp->dst_port);
    u16 sport = ntohs(tcp->src_port);

    // Check if the packet's process is 'myprocess'
    u32 pid = bpf_get_current_pid_tgid() >> 32;
    struct task_struct *task = (struct task_struct *)bpf_get_current_task();
    char comm[TASK_COMM_LEN];
    bpf_get_current_comm(&comm, sizeof(comm));

    if (strcmp(comm, PROC_NAME) == 0 && dport != MY_PORT && sport != MY_PORT) {
        return -1; // Drop packet
    }

    return 0; // Pass the packet
}
