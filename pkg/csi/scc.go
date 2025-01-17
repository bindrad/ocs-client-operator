/*
Copyright 2022 Red Hat, Inc.

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

package csi

import (
	"fmt"

	secv1 "github.com/openshift/api/security/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	SCCName = "ocs-csi-scc"
)

var (
	// serviceaccount names
	cephFSProvisionerServiceAccountName = "ocs-client-operator-csi-cephfs-provisioner-sa"
	cephFSPluginServiceAccountName      = "ocs-client-operator-csi-cephfs-plugin-sa"
	rbdProvisionerServiceAccountName    = "ocs-client-operator-csi-rbd-provisioner-sa"
	rbdPluginServiceAccountName         = "ocs-client-operator-csi-rbd-plugin-sa"
)

func GetSecurityContextConstraints(namespace string) *secv1.SecurityContextConstraints {
	return &secv1.SecurityContextConstraints{
		ObjectMeta: metav1.ObjectMeta{
			Name: SCCName,
		},
		// CSI daemonset pod needs to run as privileged
		AllowPrivilegedContainer: true,
		// CSI daemonset pod needs hostnetworking
		AllowHostNetwork: true,
		// This need to be set to true as we use HostPath
		AllowHostDirVolumePlugin: true,
		// Required for csi addons
		AllowHostPorts: true,
		// Needed as we are setting this in RBD plugin pod
		AllowHostPID: true,
		// Required for multus and encryption
		AllowHostIPC: true,
		// SYS_ADMIN is needed for rbd to execute rbd map command
		AllowedCapabilities: []corev1.Capability{"SYS_ADMIN"},
		// # Set to false as we write to RootFilesystem inside csi containers
		ReadOnlyRootFilesystem: false,
		RunAsUser: secv1.RunAsUserStrategyOptions{
			Type: secv1.RunAsUserStrategyRunAsAny,
		},
		SELinuxContext: secv1.SELinuxContextStrategyOptions{
			Type: secv1.SELinuxStrategyRunAsAny,
		},
		FSGroup: secv1.FSGroupStrategyOptions{
			Type: secv1.FSGroupStrategyRunAsAny,
		},
		SupplementalGroups: secv1.SupplementalGroupsStrategyOptions{
			Type: secv1.SupplementalGroupsStrategyRunAsAny,
		},
		Volumes: []secv1.FSType{
			secv1.FSTypeHostPath,
			secv1.FSTypeConfigMap,
			secv1.FSTypeEmptyDir,
			secv1.FSProjected,
		},
		Users: []string{
			fmt.Sprintf("system:serviceaccount:%s:%s", namespace, cephFSProvisionerServiceAccountName),
			fmt.Sprintf("system:serviceaccount:%s:%s", namespace, cephFSPluginServiceAccountName),
			fmt.Sprintf("system:serviceaccount:%s:%s", namespace, rbdProvisionerServiceAccountName),
			fmt.Sprintf("system:serviceaccount:%s:%s", namespace, rbdPluginServiceAccountName),
		},
	}
}
