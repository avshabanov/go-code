// Given a pool of hardware nodes, each has a fixed amount of CPU cores, RAM and Disk, write a code
// that allocates a new VM, consuming certain number of cores, ram and disk

package hwalloc

import (
	"fmt"
)

// CapacityUnit represents a hardware unit's capacity description
type CapacityUnit struct {
	Cores int
	RAM   int
	Disk  int
}

func (t *CapacityUnit) String() string {
	return fmt.Sprintf("{Cores=%d, Disk=%d, RAM=%d}", t.Cores, t.Disk, t.RAM)
}

// VM represents a node data
type VM struct {
	CapacityUnit
	ID string
}

// Host models a VM hypervisor
type Host struct {
	CapacityUnit
	Serial string
}

//
// Simplest solution
//

// PooledHost represents a pooled host record
type PooledHost struct {
	host         *Host
	capacity     CapacityUnit
	allocatedVMs []*VM
}

// Pool models a pool of VMs
type Pool struct {
	pooledHosts []*PooledHost
}

// Add adds a host to the host pool
func (t *Pool) Add(newHost *Host) error {
	// check, that host has not yet been pooled
	for _, h := range t.pooledHosts {
		if h.host.Serial == newHost.Serial {
			return fmt.Errorf("host has already been pooled: %s", newHost.Serial)
		}
	}

	newPooledHost := &PooledHost{host: newHost}

	// Initialize available capacity
	newPooledHost.capacity.Cores = newHost.Cores
	newPooledHost.capacity.RAM = newHost.RAM
	newPooledHost.capacity.Disk = newHost.Disk

	t.pooledHosts = append(t.pooledHosts, newPooledHost)
	return nil
}

// AllocateVM creates a VM on one of the pooled hosts
func (t *Pool) AllocateVM(vm *VM) (string, error) {
	// brain-dead solution
	for _, ph := range t.pooledHosts {
		cap := &ph.capacity
		newCores := cap.Cores - vm.Cores
		newRAM := cap.RAM - vm.RAM
		newDisk := cap.Disk - vm.Disk
		if newCores < 0 || newRAM < 0 || newDisk < 0 {
			continue
		}

		// capacity can be reserved, do it now
		cap.Cores = newCores
		cap.RAM = newRAM
		cap.Disk = newDisk
		ph.allocatedVMs = append(ph.allocatedVMs, vm)
		return ph.host.Serial, nil
	}

	return "", fmt.Errorf("can not allocate vm with id=%s", vm.ID)
}

// GetCapacity returns total available capacity inside the host pool
func (t *Pool) GetCapacity() *CapacityUnit {
	c := &CapacityUnit{}
	for _, ph := range t.pooledHosts {
		c.Cores = c.Cores + ph.capacity.Cores
		c.RAM = c.RAM + ph.capacity.RAM
		c.Disk = c.Disk + ph.capacity.Disk
	}
	return c
}
