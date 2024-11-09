package main

import (
	"flag"
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

func main() {
	processName := flag.String("pc", "", "")
	flag.Parse()
	// Get all process alive
	processes, err := getProcesses()
	if err != nil {
		fmt.Printf("Error getting processes: %v\n", err)
		return
	}
	for _, p := range processes {
		if p.name == *processName {
			fmt.Println("Found process", p.pid)
			dlls, err := getLoadedDLLs(p.pid)
			if err != nil {
				fmt.Printf("Error getting DLLs for process %d: %v\n", p.pid, err)
				continue
			}
			for _, dll := range dlls {
				fmt.Println("  ", dll)
			}
		} else {
			continue
		}
	}
}

type process struct {
	pid  uint32
	name string
}

func getProcesses() ([]process, error) {
	var processes []process
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPALL, 0)
	if err != nil {
		return nil, err
	}
	defer windows.CloseHandle(snapshot)

	var pe32 windows.ProcessEntry32
	pe32.Size = uint32(unsafe.Sizeof(pe32))

	if err := windows.Process32First(snapshot, &pe32); err != nil {
		return nil, err
	}
	for {
		processes = append(processes, process{pid: pe32.ProcessID, name: syscall.UTF16ToString(pe32.ExeFile[:])})
		if err := windows.Process32Next(snapshot, &pe32); err != nil {
			break
		}
	}
	return processes, nil
}

func getLoadedDLLs(pid uint32) ([]string, error) {
	var dlls []string
	handle, err := windows.OpenProcess(windows.PROCESS_QUERY_INFORMATION|windows.PROCESS_VM_READ, false, pid)
	if err != nil {
		return nil, err
	}
	defer windows.CloseHandle(handle)

	hModuleSnapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPMODULE, pid)
	if err != nil {
		return nil, err
	}
	defer windows.CloseHandle(hModuleSnapshot)

	var moduleEntry windows.ModuleEntry32
	moduleEntry.Size = uint32(unsafe.Sizeof(moduleEntry))

	if err := windows.Module32First(hModuleSnapshot, &moduleEntry); err != nil {
		return nil, err
	}
	for {
		dlls = append(dlls, syscall.UTF16ToString(moduleEntry.Module[:]))
		if err := windows.Module32Next(hModuleSnapshot, &moduleEntry); err != nil {
			break
		}
	}
	return dlls, nil
}
