// Copyright (c) 2017 Intel Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import "github.com/golang/glog"

func (c *controller) EvacuateNode(nodeID string) error {
	// should I bother to see if nodeID is valid?
	go func() {
		if err := c.client.EvacuateNode(nodeID); err != nil {
			glog.Warningf("Error evacuating node")
		}
	}()
	return nil
}

func (c *controller) RestoreNode(nodeID string) error {
	go func() {
		if err := c.client.RestoreNode(nodeID); err != nil {
			glog.Warning("Error restoring node")
		}
	}()
	return nil
}
