/*
Copyright 2016 The Kubernetes Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package localkube

import (
	kubelet "k8s.io/kubernetes/cmd/kubelet/app"
	"k8s.io/kubernetes/cmd/kubelet/app/options"
)

const (
	HostnameOverride = "127.0.0.1"
)

func (lk LocalkubeServer) NewKubeletServer() Server {
	return NewSimpleServer("kubelet", serverInterval, StartKubeletServer(lk))
}

func changeToRkt(cfg string) {
	//cfg.ContainerRuntime = "rkt"
	//cfg.NetworkPluginName = "kubenet"
	//RktPath:                          "",
	//RktStage1Image:                   ""
}

func StartKubeletServer(lk LocalkubeServer) func() error {
	config := options.NewKubeletServer()
	config.ContainerRuntime = "rkt"
	config.NetworkPluginName = "kubenet"

	// Master details
	config.APIServerList = []string{lk.GetAPIServerInsecureURL()}

	// Set containerized based on the flag
	config.Containerized = lk.Containerized

	config.AllowPrivileged = true
	config.Config = "/etc/kubernetes/manifests"

	// Networking
	config.ClusterDomain = lk.DNSDomain
	config.ClusterDNS = lk.DNSIP.String()
	config.HostnameOverride = HostnameOverride

	// Use the host's resolver config
	if lk.Containerized {
		config.ResolverConfig = "/rootfs/etc/resolv.conf"
	} else {
		config.ResolverConfig = "/etc/resolv.conf"
	}

	return func() error {
		return kubelet.Run(config, nil)
	}
}
