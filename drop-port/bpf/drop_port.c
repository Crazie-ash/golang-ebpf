#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/ip.h>
#include <linux/tcp.h>
#include <linux/in.h>
#include <bpf/bpf_helpers.h>
#define ntohs __builtin_bswap16

// Define your eBPF map here correctly
struct bpf_map_def SEC("maps") port_map = {
    .type = BPF_MAP_TYPE_HASH,
    .key_size = sizeof(__u16),
    .value_size = sizeof(__u16),
    .max_entries = 1,
};

SEC("xdp")
int drop_port_bpf(struct xdp_md *ctx) {
    char msg[] = "TCP Packet received\n";
    bpf_trace_printk(msg, sizeof(msg));
    void *data = (void *)(long)ctx->data;
    void *data_end = (void *)(long)ctx->data_end;

    // Ensure Ethernet header is within packet bounds
    if (data + sizeof(struct ethhdr) > data_end)
        return XDP_PASS;

    struct ethhdr *eth = data;
    struct iphdr *ip = data + sizeof(struct ethhdr);

    // Ensure IP header is within packet bounds
    if ((void *)ip + sizeof(struct iphdr) > data_end)
        return XDP_PASS;

    if (ip->protocol == IPPROTO_TCP) {
        struct tcphdr *tcp = (void *)ip + sizeof(struct iphdr);

        // Ensure TCP header is within packet bounds
        if ((void *)tcp + sizeof(struct tcphdr) > data_end)
            return XDP_PASS;
            
        int dest_port = ntohs(tcp->dest);
        // Correctly use the map here
        int *port = bpf_map_lookup_elem(&port_map, &tcp->dest);
        if (port && *port != 0) {
           bpf_trace_printk("Dropped packet to port %d\n", *port);
           return XDP_DROP;
        }

        // Temporary line for testing
        char msg[] = "TCP Packet Processed\n";
        bpf_trace_printk(msg, sizeof(msg));
    }

    return XDP_PASS;
}

char _license[] SEC("license") = "GPL";
