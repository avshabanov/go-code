// There is a pool of *X* available servers, each hosting a VM. Each VM hosts a single service. Each service is hosted on *K* VMs.
// There are *N* fault domains.
// Options:
// 	* Each server is uniquely assigned to a certain fault domain. The same service is never hosted more than once on a server belonging to the same fault domain.
//	* Each VM is assigned to certain fault domain. The same service is never hosted more than once on a server, where there is a VM belonging to the other fault domain.
//
// Write a code, that:
// 	* Enters new servers, new VMs and new services to the managed environment, that respects fault domains
// 	* Outputs a sequence in which hosts might be rebooted

package faultdomain
