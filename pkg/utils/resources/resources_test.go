package resources

import (
	"testing"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestFits(t *testing.T) {
	tests := []struct {
		name       string
		candidate  v1.ResourceList
		total      v1.ResourceList
		ignored    v1.ResourceList
		wantResult bool
	}{
		{
			name: "Custom resource requests added to the ignore list makes the candidate fit",
			candidate: v1.ResourceList{
				v1.ResourceCPU:                  resource.MustParse("200m"),
				v1.ResourceMemory:               resource.MustParse("100Mi"),
				"devices.kubevirt.io/kvm":       resource.MustParse("1"),
				"devices.kubevirt.io/tun":       resource.MustParse("1"),
				"devices.kubevirt.io/vhost-net": resource.MustParse("1"),
			},
			total: v1.ResourceList{
				v1.ResourceCPU:    resource.MustParse("200m"),
				v1.ResourceMemory: resource.MustParse("100Mi"),
			},
			ignored: v1.ResourceList{
				"devices.kubevirt.io/kvm":       resource.MustParse("1"),
				"devices.kubevirt.io/tun":       resource.MustParse("1"),
				"devices.kubevirt.io/vhost-net": resource.MustParse("1"),
			},
			wantResult: true,
		},
		{
			name: "Custom resource request pattern on the ignore list makes the candidate fit",
			candidate: v1.ResourceList{
				v1.ResourceCPU:    resource.MustParse("200m"),
				v1.ResourceMemory: resource.MustParse("100Mi"),
				"scheduling.node.kubevirt.io/tsc-frequency-2999987000": resource.MustParse("1"),
			},
			total: v1.ResourceList{
				v1.ResourceCPU:    resource.MustParse("200m"),
				v1.ResourceMemory: resource.MustParse("100Mi"),
			},
			ignored: v1.ResourceList{
				"scheduling.node.kubevirt.io/tsc-frequency-*": resource.MustParse("1"),
			},
			wantResult: true,
		},
		{
			name: "Custom resource request not on the ignore list makes the candidate unfit",
			candidate: v1.ResourceList{
				v1.ResourceCPU:                    resource.MustParse("200m"),
				v1.ResourceMemory:                 resource.MustParse("100Mi"),
				"devices.kubevirt.io/not-ignored": resource.MustParse("1"),
			},
			total: v1.ResourceList{
				v1.ResourceCPU:    resource.MustParse("200m"),
				v1.ResourceMemory: resource.MustParse("100Mi"),
			},
			ignored: v1.ResourceList{
				"devices.kubevirt.io/kvm": resource.MustParse("1"),
			},
			wantResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Fits(tt.candidate, tt.total, tt.ignored)
			if got != tt.wantResult {
				t.Errorf("Fits() = %v, want = %v", got, tt.wantResult)
			}
		})
	}
}
