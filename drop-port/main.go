package main

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/rlimit"
)

func main() {
	err := rlimit.RemoveMemlock()
	handleErr(fatal, "could not remove memlock", err)

	log.Info("loading bpf module from file...")
	spec, err := ebpf.LoadCollectionSpec("./drop_port.o")
	handleErr(fatal, "could not load specs from bpf file :", err)

	log.Infof("Number of bpf program found: %v", len(spec.Programs))

	for i, p := range spec.Programs {
		log.Infof("Index: %s, Name: %s, Type: %s, AttachType: %s, Section: %s", i, p.Name, p.Type, p.AttachType, p.SectionName)
	}

	// Load the entire collection based on the spec
	coll, err := ebpf.NewCollection(spec)
	if err != nil {
		log.Fatalf("could not load collection: %v", err)
	}
	defer coll.Close()

	// Access the port_map from the collection
	portMap := spec.Maps["port_map"]
	if err != nil {
		log.Fatalf("map 'port_map' not found in collection")
	}
	fmt.Println(portMap)

	// Insert data into the port_map
	// portNumber := 4040 // example port number
	// mapValue := 1      // example value
	// if err := portMap.Update(portNumber, mapValue, ebpf.UpdateAny); err != nil {
	// 	log.Fatalf("Error updating port_map: %v", err)
	// }

	prog, err := ebpf.NewProgram(spec.Programs["drop_port_bpf"])
	handleErr(fatal, "could not load program from specs :", err)

	fmt.Println(prog)

	// _, err = link.Tracepoint("syscalls", "sys_enter_execve", prog, nil)
	// handleErr(er, "can't link tracepoint", err)

	// err = syscall.Exec("/usr/bin/ls", nil, nil)
	// handleErr(er, "can't trigger syscall :", err)
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
