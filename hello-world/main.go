package main

import (
	"fmt"
	"syscall"

	"github.com/charmbracelet/log"
	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/rlimit"
)

func main() {
	err := rlimit.RemoveMemlock()
	handleErr(fatal, "could not remove memlock", err)

	log.Info("loading bpf module from file...")
	spec, err := ebpf.LoadCollectionSpec("./hello.o")
	handleErr(fatal, "could not load specs from bpf file :", err)

	log.Infof("Number of bpf program found: %v", len(spec.Programs))

	for i, p := range spec.Programs {
		log.Infof("Index: %s, Name: %s, Type: %s, AttachType: %s, Section: %s", i, p.Name, p.Type, p.AttachType, p.SectionName)
	}

	prog, err := ebpf.NewProgram(spec.Programs["hello_bpf"])
	handleErr(fatal, "could not load program from specs :", err)

	fmt.Println(prog)

	_, err = link.Tracepoint("syscalls", "sys_enter_execve", prog, nil)
	handleErr(er, "can't link tracepoint", err)

	err = syscall.Exec("/usr/bin/ls", nil, nil)
	handleErr(er, "can't trigger syscall :", err)
}

type level string

const (
	warn  level = "warning"
	er    level = "error"
	fatal level = "fatal"
)

func handleErr(lev level, msg string, err error) {
	if err != nil {
		switch lev {
		case warn:
			log.Warnf(msg+"%s", err)
		case er:
			log.Errorf(msg+"%s", err)
		case fatal:
			log.Fatalf(msg+"%s", err)
		default:
			panic("invalid log level provided")
		}
	}
}
