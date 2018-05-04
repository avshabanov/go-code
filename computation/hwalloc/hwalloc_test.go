package hwalloc

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHwalloc(t *testing.T) {
	t.Run("fizz buzz 35 1..10", func(t *testing.T) {
		p := &Pool{}
		p.Add(newStdHost("Host1"))
		p.Add(newStdHost("Host2"))

		fmt.Printf(" * current capacity=%s\n", p.GetCapacity())

		vmList := []*VM{
			newStdVM("std-vm-01"),
			newStorageVM("storage-vm-02"),
			newStorageVM("storage-vm-03"),
			newStorageVM("storage-vm-04"),
			newStorageVM("storage-vm-05"),
			newCompVM("comp-vm-06"),
			newStdVM("std-vm-07"),
		}

		for _, vm := range vmList {
			hostSerial, err := p.AllocateVM(newStdVM("std-vm-01"))
			assert.NoError(t, err, "while allocating vm with id=%s", vm.ID)

			fmt.Printf("Allocated vm=%s on host=%s\n", vm.ID, hostSerial)
			fmt.Printf(" * current capacity=%s\n", p.GetCapacity())
		}
	})
}

// creates a record, that designates a standard VM
func newStdVM(id string) *VM {
	vm := &VM{ID: id}
	vm.Cores = 20
	vm.Disk = 2000
	vm.RAM = 200
	return vm
}

// creates a record, that designates a memory-optimized VM
func newStorageVM(id string) *VM {
	vm := &VM{ID: id}
	vm.Cores = 10
	vm.Disk = 3000
	vm.RAM = 300
	return vm
}

// creates a record, that designates a CPU-optimized VM
func newCompVM(id string) *VM {
	vm := &VM{ID: id}
	vm.Cores = 30
	vm.Disk = 1000
	vm.RAM = 100
	return vm
}

func newStdHost(serial string) *Host {
	h := &Host{Serial: serial}
	h.Cores = 100
	h.Disk = 10000 // ~10000 Gb or ~10Tb
	h.RAM = 1000   // ~1000 Gb
	return h
}
