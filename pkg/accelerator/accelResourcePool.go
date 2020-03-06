// Copyright 2018 Intel Corp. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package accelerator

import (
	"github.com/golang/glog"
	"github.com/intel/sriov-network-device-plugin/pkg/resources"
	"github.com/intel/sriov-network-device-plugin/pkg/types"
	pluginapi "k8s.io/kubernetes/pkg/kubelet/apis/deviceplugin/v1beta1"
)

type accelResourcePool struct {
	*resources.ResourcePoolImpl
}

var _ types.ResourcePool = &accelResourcePool{}

// NewAccelResourcePool returns an instance of resourcePool
func NewAccelResourcePool(rc *types.ResourceConfig, apiDevices map[string]*pluginapi.Device, devicePool map[string]types.PciDevice) types.ResourcePool {
	rp := resources.NewResourcePool(rc, apiDevices, devicePool)
	return &accelResourcePool{
		ResourcePoolImpl: rp,
	}
}

// Overrides GetDeviceSpecs
func (rp *accelResourcePool) GetDeviceSpecs(deviceIDs []string) []*pluginapi.DeviceSpec {
	glog.Infof("GetDeviceSpecs(): for devices: %v", deviceIDs)
	devSpecs := make([]*pluginapi.DeviceSpec, 0)

	devicePool := rp.GetDevicePool()

	// Add vfio group specific devices
	for _, id := range deviceIDs {
		if dev, ok := devicePool[id]; ok {
			netDev := dev.(types.AccelDevice) // convert generic PciDevice to PciNetDevice
			newSpecs := netDev.GetDeviceSpecs()
			for _, ds := range newSpecs {
				if !rp.DeviceSpecExist(devSpecs, ds) {
					devSpecs = append(devSpecs, ds)
				}

			}

		}
	}
	return devSpecs
}
