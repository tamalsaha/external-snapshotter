/*
Copyright 2018 The Kubernetes Authors.

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

package controller

import (
	"errors"
	"testing"
	"time"

	crdv1 "github.com/kubernetes-csi/external-snapshotter/pkg/apis/volumesnapshot/v1alpha1"
	"k8s.io/api/core/v1"
	storage "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var timeNow = time.Now()
var timeNowStamp = timeNow.UnixNano()

var metaTimeNowUnix = &metav1.Time{
	Time: timeNow,
}

var defaultSize int64 = 1000
var deletePolicy = crdv1.VolumeSnapshotContentDelete
var retainPolicy = crdv1.VolumeSnapshotContentRetain
var sameDriverStorageClass = &storage.StorageClass{
	TypeMeta: metav1.TypeMeta{
		Kind: "StorageClass",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name: "sameDriver",
	},
	Provisioner: mockDriverName,
	Parameters:  class1Parameters,
}

var diffDriverStorageClass = &storage.StorageClass{
	TypeMeta: metav1.TypeMeta{
		Kind: "StorageClass",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name: "diffDriver",
	},
	Provisioner: mockDriverName,
	Parameters:  class1Parameters,
}

// Test single call to SyncSnapshot, expecting create snapshot to happen.
// 1. Fill in the controller with initial data
// 2. Call the SyncSnapshot *once*.
// 3. Compare resulting contents with expected contents.
func TestCreateSnapshotSync(t *testing.T) {

	tests := []controllerTest{
		{
			name:              "6-1 - successful create snapshot with snapshot class gold",
			initialContents:   nocontents,
			expectedContents:  newContentArray("snapcontent-snapuid6-1", classGold, "sid6-1", "pv-uid6-1", "volume6-1", "snapuid6-1", "snap6-1", &deletePolicy, &defaultSize, &timeNowStamp, false),
			initialSnapshots:  newSnapshotArray("snap6-1", classGold, "", "snapuid6-1", "claim6-1", false, nil, nil, nil),
			expectedSnapshots: newSnapshotArray("snap6-1", classGold, "snapcontent-snapuid6-1", "snapuid6-1", "claim6-1", false, nil, metaTimeNowUnix, getSize(defaultSize)),
			initialClaims:     newClaimArray("claim6-1", "pvc-uid6-1", "1Gi", "volume6-1", v1.ClaimBound, &classEmpty),
			initialVolumes:    newVolumeArray("volume6-1", "pv-uid6-1", "pv-handle6-1", "1Gi", "pvc-uid6-1", "claim6-1", v1.VolumeBound, v1.PersistentVolumeReclaimDelete, classEmpty, mockDriverName, testNamespace),
			expectedCreateCalls: []createCall{
				{
					snapshotName: "snapshot-snapuid6-1",
					volume:       newVolume("volume6-1", "pv-uid6-1", "pv-handle6-1", "1Gi", "pvc-uid6-1", "claim6-1", v1.VolumeBound, v1.PersistentVolumeReclaimDelete, classEmpty, mockDriverName, testNamespace),
					parameters:   map[string]string{"param1": "value1"},
					// information to return
					driverName:   mockDriverName,
					size:         defaultSize,
					snapshotId:   "sid6-1",
					creationTime: timeNow,
					readyToUse:   true,
				},
			},
			errors: noerrors,
			test:   testSyncSnapshot,
		},
		{
			name:              "6-2 - successful create snapshot with snapshot class silver",
			initialContents:   nocontents,
			expectedContents:  newContentArray("snapcontent-snapuid6-2", classSilver, "sid6-2", "pv-uid6-2", "volume6-2", "snapuid6-2", "snap6-2", &deletePolicy, &defaultSize, &timeNowStamp, false),
			initialSnapshots:  newSnapshotArray("snap6-2", classSilver, "", "snapuid6-2", "claim6-2", false, nil, nil, nil),
			expectedSnapshots: newSnapshotArray("snap6-2", classSilver, "snapcontent-snapuid6-2", "snapuid6-2", "claim6-2", false, nil, metaTimeNowUnix, getSize(defaultSize)),
			initialClaims:     newClaimArray("claim6-2", "pvc-uid6-2", "1Gi", "volume6-2", v1.ClaimBound, &classEmpty),
			initialVolumes:    newVolumeArray("volume6-2", "pv-uid6-2", "pv-handle6-2", "1Gi", "pvc-uid6-2", "claim6-2", v1.VolumeBound, v1.PersistentVolumeReclaimDelete, classEmpty, mockDriverName, testNamespace),
			expectedCreateCalls: []createCall{
				{
					snapshotName: "snapshot-snapuid6-2",
					volume:       newVolume("volume6-2", "pv-uid6-2", "pv-handle6-2", "1Gi", "pvc-uid6-2", "claim6-2", v1.VolumeBound, v1.PersistentVolumeReclaimDelete, classEmpty, mockDriverName, testNamespace),
					parameters:   map[string]string{"param2": "value2"},
					// information to return
					driverName:   mockDriverName,
					size:         defaultSize,
					snapshotId:   "sid6-2",
					creationTime: timeNow,
					readyToUse:   true,
				},
			},
			errors: noerrors,
			test:   testSyncSnapshot,
		},
		{
			name:              "6-3 - successful create snapshot with snapshot class valid-secret-class",
			initialContents:   nocontents,
			expectedContents:  newContentArray("snapcontent-snapuid6-3", validSecretClass, "sid6-3", "pv-uid6-3", "volume6-3", "snapuid6-3", "snap6-3", &deletePolicy, &defaultSize, &timeNowStamp, false),
			initialSnapshots:  newSnapshotArray("snap6-3", validSecretClass, "", "snapuid6-3", "claim6-3", false, nil, nil, nil),
			expectedSnapshots: newSnapshotArray("snap6-3", validSecretClass, "snapcontent-snapuid6-3", "snapuid6-3", "claim6-3", false, nil, metaTimeNowUnix, getSize(defaultSize)),
			initialClaims:     newClaimArray("claim6-3", "pvc-uid6-3", "1Gi", "volume6-3", v1.ClaimBound, &classEmpty),
			initialVolumes:    newVolumeArray("volume6-3", "pv-uid6-3", "pv-handle6-3", "1Gi", "pvc-uid6-3", "claim6-3", v1.VolumeBound, v1.PersistentVolumeReclaimDelete, classEmpty, mockDriverName, testNamespace),
			initialSecrets:    []*v1.Secret{secret()},
			expectedCreateCalls: []createCall{
				{
					snapshotName: "snapshot-snapuid6-3",
					volume:       newVolume("volume6-3", "pv-uid6-3", "pv-handle6-3", "1Gi", "pvc-uid6-3", "claim6-3", v1.VolumeBound, v1.PersistentVolumeReclaimDelete, classEmpty, mockDriverName, testNamespace),
					parameters:   class5Parameters,
					secrets:      map[string]string{"foo": "bar"},
					// information to return
					driverName:   mockDriverName,
					size:         defaultSize,
					snapshotId:   "sid6-3",
					creationTime: timeNow,
					readyToUse:   true,
				},
			},
			errors: noerrors,
			test:   testSyncSnapshot,
		},
		{
			name:              "6-4 - successful create snapshot with snapshot class empty-secret-class",
			initialContents:   nocontents,
			expectedContents:  newContentArray("snapcontent-snapuid6-4", emptySecretClass, "sid6-4", "pv-uid6-4", "volume6-4", "snapuid6-4", "snap6-4", &deletePolicy, &defaultSize, &timeNowStamp, false),
			initialSnapshots:  newSnapshotArray("snap6-4", emptySecretClass, "", "snapuid6-4", "claim6-4", false, nil, nil, nil),
			expectedSnapshots: newSnapshotArray("snap6-4", emptySecretClass, "snapcontent-snapuid6-4", "snapuid6-4", "claim6-4", false, nil, metaTimeNowUnix, getSize(defaultSize)),
			initialClaims:     newClaimArray("claim6-4", "pvc-uid6-4", "1Gi", "volume6-4", v1.ClaimBound, &classEmpty),
			initialVolumes:    newVolumeArray("volume6-4", "pv-uid6-4", "pv-handle6-4", "1Gi", "pvc-uid6-4", "claim6-4", v1.VolumeBound, v1.PersistentVolumeReclaimDelete, classEmpty, mockDriverName, testNamespace),
			initialSecrets:    []*v1.Secret{emptySecret()},
			expectedCreateCalls: []createCall{
				{
					snapshotName: "snapshot-snapuid6-4",
					volume:       newVolume("volume6-4", "pv-uid6-4", "pv-handle6-4", "1Gi", "pvc-uid6-4", "claim6-4", v1.VolumeBound, v1.PersistentVolumeReclaimDelete, classEmpty, mockDriverName, testNamespace),
					parameters:   class4Parameters,
					secrets:      map[string]string{},
					// information to return
					driverName:   mockDriverName,
					size:         defaultSize,
					snapshotId:   "sid6-4",
					creationTime: timeNow,
					readyToUse:   true,
				},
			},
			errors: noerrors,
			test:   testSyncSnapshot,
		},
		{
			name:              "6-5 - successful create snapshot with status uploading",
			initialContents:   nocontents,
			expectedContents:  newContentArray("snapcontent-snapuid6-5", classGold, "sid6-5", "pv-uid6-5", "volume6-5", "snapuid6-5", "snap6-5", &deletePolicy, &defaultSize, &timeNowStamp, false),
			initialSnapshots:  newSnapshotArray("snap6-5", classGold, "", "snapuid6-5", "claim6-5", false, nil, nil, nil),
			expectedSnapshots: newSnapshotArray("snap6-5", classGold, "snapcontent-snapuid6-5", "snapuid6-5", "claim6-5", false, nil, metaTimeNowUnix, getSize(defaultSize)),
			initialClaims:     newClaimArray("claim6-5", "pvc-uid6-5", "1Gi", "volume6-5", v1.ClaimBound, &classEmpty),
			initialVolumes:    newVolumeArray("volume6-5", "pv-uid6-5", "pv-handle6-5", "1Gi", "pvc-uid6-5", "claim6-5", v1.VolumeBound, v1.PersistentVolumeReclaimDelete, classEmpty, mockDriverName, testNamespace),
			expectedCreateCalls: []createCall{
				{
					snapshotName: "snapshot-snapuid6-5",
					volume:       newVolume("volume6-5", "pv-uid6-5", "pv-handle6-5", "1Gi", "pvc-uid6-5", "claim6-5", v1.VolumeBound, v1.PersistentVolumeReclaimDelete, classEmpty, mockDriverName, testNamespace),
					parameters:   map[string]string{"param1": "value1"},
					// information to return
					driverName:   mockDriverName,
					size:         defaultSize,
					snapshotId:   "sid6-5",
					creationTime: timeNow,
					readyToUse:   true,
				},
			},
			errors: noerrors,
			test:   testSyncSnapshot,
		},
		{
			name:              "6-6 - successful create snapshot with status error uploading",
			initialContents:   nocontents,
			expectedContents:  newContentArray("snapcontent-snapuid6-6", classGold, "sid6-6", "pv-uid6-6", "volume6-6", "snapuid6-6", "snap6-6", &deletePolicy, &defaultSize, &timeNowStamp, false),
			initialSnapshots:  newSnapshotArray("snap6-6", classGold, "", "snapuid6-6", "claim6-6", false, nil, nil, nil),
			expectedSnapshots: newSnapshotArray("snap6-6", classGold, "snapcontent-snapuid6-6", "snapuid6-6", "claim6-6", false, nil, metaTimeNowUnix, getSize(defaultSize)),
			initialClaims:     newClaimArray("claim6-6", "pvc-uid6-6", "1Gi", "volume6-6", v1.ClaimBound, &classEmpty),
			initialVolumes:    newVolumeArray("volume6-6", "pv-uid6-6", "pv-handle6-6", "1Gi", "pvc-uid6-6", "claim6-6", v1.VolumeBound, v1.PersistentVolumeReclaimDelete, classEmpty, mockDriverName, testNamespace),
			expectedCreateCalls: []createCall{
				{
					snapshotName: "snapshot-snapuid6-6",
					volume:       newVolume("volume6-6", "pv-uid6-6", "pv-handle6-6", "1Gi", "pvc-uid6-6", "claim6-6", v1.VolumeBound, v1.PersistentVolumeReclaimDelete, classEmpty, mockDriverName, testNamespace),
					parameters:   map[string]string{"param1": "value1"},
					// information to return
					driverName:   mockDriverName,
					size:         defaultSize,
					snapshotId:   "sid6-6",
					creationTime: timeNow,
					readyToUse:   true,
				},
			},
			errors: noerrors,
			test:   testSyncSnapshot,
		},
		{
			name:              "7-1 - fail create snapshot with snapshot class non-existing",
			initialContents:   nocontents,
			expectedContents:  nocontents,
			initialSnapshots:  newSnapshotArray("snap7-1", classNonExisting, "", "snapuid7-1", "claim7-1", false, nil, nil, nil),
			expectedSnapshots: newSnapshotArray("snap7-1", classNonExisting, "", "snapuid7-1", "claim7-1", false, newVolumeError("Failed to create snapshot: failed to get input parameters to create snapshot snap7-1: \"failed to retrieve snapshot class non-existing from the informer: \\\"volumesnapshotclass.snapshot.storage.k8s.io \\\\\\\"non-existing\\\\\\\" not found\\\"\""), nil, nil),
			initialClaims:     newClaimArray("claim7-1", "pvc-uid7-1", "1Gi", "volume7-1", v1.ClaimBound, &classEmpty),
			initialVolumes:    newVolumeArray("volume7-1", "pv-uid7-1", "pv-handle7-1", "1Gi", "pvc-uid7-1", "claim7-1", v1.VolumeBound, v1.PersistentVolumeReclaimDelete, classEmpty, mockDriverName, testNamespace),
			expectedEvents:    []string{"Warning SnapshotCreationFailed"},
			errors:            noerrors,
			test:              testSyncSnapshot,
		},
		{
			name:              "7-2 - fail create snapshot with snapshot class invalid-secret-class",
			initialContents:   nocontents,
			expectedContents:  nocontents,
			initialSnapshots:  newSnapshotArray("snap7-2", invalidSecretClass, "", "snapuid7-2", "claim7-2", false, nil, nil, nil),
			expectedSnapshots: newSnapshotArray("snap7-2", invalidSecretClass, "", "snapuid7-2", "claim7-2", false, newVolumeError("Failed to create snapshot: failed to get input parameters to create snapshot snap7-2: \"failed to get name and namespace template from params: either name and namespace for Snapshotter secrets specified, Both must be specified\""), nil, nil),
			initialClaims:     newClaimArray("claim7-2", "pvc-uid7-2", "1Gi", "volume7-2", v1.ClaimBound, &classEmpty),
			initialVolumes:    newVolumeArray("volume7-2", "pv-uid7-2", "pv-handle7-2", "1Gi", "pvc-uid7-2", "claim7-2", v1.VolumeBound, v1.PersistentVolumeReclaimDelete, classEmpty, mockDriverName, testNamespace),
			expectedEvents:    []string{"Warning SnapshotCreationFailed"},
			errors:            noerrors,
			test:              testSyncSnapshot,
		},
		{
			name:                  "7-3 - fail create snapshot with none snapshot class ",
			initialContents:       nocontents,
			expectedContents:      nocontents,
			initialSnapshots:      newSnapshotArray("snap7-3", "", "", "snapuid7-3", "claim7-3", false, nil, nil, nil),
			expectedSnapshots:     newSnapshotArray("snap7-3", "", "", "snapuid7-3", "claim7-3", false, newVolumeError("Failed to create snapshot: failed to get input parameters to create snapshot snap7-3: \"failed to retrieve snapshot class  from the informer: \\\"volumesnapshotclass.snapshot.storage.k8s.io \\\\\\\"\\\\\\\" not found\\\"\""), nil, nil),
			initialClaims:         newClaimArray("claim7-3", "pvc-uid7-3", "1Gi", "volume7-3", v1.ClaimBound, &classEmpty),
			initialVolumes:        newVolumeArray("volume7-3", "pv-uid7-3", "pv-handle7-3", "1Gi", "pvc-uid7-3", "claim7-3", v1.VolumeBound, v1.PersistentVolumeReclaimDelete, classEmpty, mockDriverName, testNamespace),
			initialStorageClasses: []*storage.StorageClass{diffDriverStorageClass},
			expectedEvents:        []string{"Warning SnapshotCreationFailed"},
			errors:                noerrors,
			test:                  testSyncSnapshot,
		},
		{
			name:              "7-4 - fail create snapshot with no-existing claim",
			initialContents:   nocontents,
			expectedContents:  nocontents,
			initialSnapshots:  newSnapshotArray("snap7-4", classGold, "", "snapuid7-4", "claim7-4", false, nil, nil, nil),
			expectedSnapshots: newSnapshotArray("snap7-4", classGold, "", "snapuid7-4", "claim7-4", false, newVolumeError("Failed to create snapshot: failed to get input parameters to create snapshot snap7-4: \"failed to retrieve PVC claim7-4 from the lister: \\\"persistentvolumeclaim \\\\\\\"claim7-4\\\\\\\" not found\\\"\""), nil, nil),
			initialVolumes:    newVolumeArray("volume7-4", "pv-uid7-4", "pv-handle7-4", "1Gi", "pvc-uid7-4", "claim7-4", v1.VolumeBound, v1.PersistentVolumeReclaimDelete, classEmpty, mockDriverName, testNamespace),
			expectedEvents:    []string{"Warning SnapshotCreationFailed"},
			errors:            noerrors,
			test:              testSyncSnapshot,
		},
		{
			name:              "7-5 - fail create snapshot with no-existing volume",
			initialContents:   nocontents,
			expectedContents:  nocontents,
			initialSnapshots:  newSnapshotArray("snap7-5", classGold, "", "snapuid7-5", "claim7-5", false, nil, nil, nil),
			expectedSnapshots: newSnapshotArray("snap7-5", classGold, "", "snapuid7-5", "claim7-5", false, newVolumeError("Failed to create snapshot: failed to get input parameters to create snapshot snap7-5: \"failed to retrieve PV volume7-5 from the API server: \\\"cannot find volume volume7-5\\\"\""), nil, nil),
			initialClaims:     newClaimArray("claim7-5", "pvc-uid7-5", "1Gi", "volume7-5", v1.ClaimBound, &classEmpty),
			expectedEvents:    []string{"Warning SnapshotCreationFailed"},
			errors:            noerrors,
			test:              testSyncSnapshot,
		},
		{
			name:              "7-6 - fail create snapshot with claim that is not yet bound",
			initialContents:   nocontents,
			expectedContents:  nocontents,
			initialSnapshots:  newSnapshotArray("snap7-6", classGold, "", "snapuid7-6", "claim7-6", false, nil, nil, nil),
			expectedSnapshots: newSnapshotArray("snap7-6", classGold, "", "snapuid7-6", "claim7-6", false, newVolumeError("Failed to create snapshot: failed to get input parameters to create snapshot snap7-6: \"the PVC claim7-6 is not yet bound to a PV, will not attempt to take a snapshot\""), nil, nil),
			initialClaims:     newClaimArray("claim7-6", "pvc-uid7-6", "1Gi", "", v1.ClaimPending, &classEmpty),
			expectedEvents:    []string{"Warning SnapshotCreationFailed"},
			errors:            noerrors,
			test:              testSyncSnapshot,
		},
		{
			name:              "7-7 - fail create snapshot due to csi driver error",
			initialContents:   nocontents,
			expectedContents:  nocontents,
			initialSnapshots:  newSnapshotArray("snap7-7", classGold, "", "snapuid7-7", "claim7-7", false, nil, nil, nil),
			expectedSnapshots: newSnapshotArray("snap7-7", classGold, "", "snapuid7-7", "claim7-7", false, newVolumeError("Failed to create snapshot: failed to take snapshot of the volume, volume7-7: \"mock create snapshot error\""), nil, nil),
			initialClaims:     newClaimArray("claim7-7", "pvc-uid7-7", "1Gi", "volume7-7", v1.ClaimBound, &classEmpty),
			initialVolumes:    newVolumeArray("volume7-7", "pv-uid7-7", "pv-handle7-7", "1Gi", "pvc-uid7-7", "claim7-7", v1.VolumeBound, v1.PersistentVolumeReclaimDelete, classEmpty, mockDriverName, testNamespace),
			expectedCreateCalls: []createCall{
				{
					snapshotName: "snapshot-snapuid7-7",
					volume:       newVolume("volume7-7", "pv-uid7-7", "pv-handle7-7", "1Gi", "pvc-uid7-7", "claim7-7", v1.VolumeBound, v1.PersistentVolumeReclaimDelete, classEmpty, mockDriverName, testNamespace),
					parameters:   map[string]string{"param1": "value1"},
					// information to return
					err: errors.New("mock create snapshot error"),
				},
			},
			errors:         noerrors,
			expectedEvents: []string{"Warning SnapshotCreationFailed"},
			test:           testSyncSnapshot,
		},
		{
			name:              "7-8 - fail create snapshot due to cannot update snapshot status",
			initialContents:   nocontents,
			expectedContents:  nocontents,
			initialSnapshots:  newSnapshotArray("snap7-8", classGold, "", "snapuid7-8", "claim7-8", false, nil, nil, nil),
			expectedSnapshots: newSnapshotArray("snap7-8", classGold, "", "snapuid7-8", "claim7-8", false, newVolumeError("Failed to create snapshot: snapshot controller failed to update default/snap7-8 on API server: mock update error"), nil, nil),
			initialClaims:     newClaimArray("claim7-8", "pvc-uid7-8", "1Gi", "volume7-8", v1.ClaimBound, &classEmpty),
			initialVolumes:    newVolumeArray("volume7-8", "pv-uid7-8", "pv-handle7-8", "1Gi", "pvc-uid7-8", "claim7-8", v1.VolumeBound, v1.PersistentVolumeReclaimDelete, classEmpty, mockDriverName, testNamespace),
			expectedCreateCalls: []createCall{
				{
					snapshotName: "snapshot-snapuid7-8",
					volume:       newVolume("volume7-8", "pv-uid7-8", "pv-handle7-8", "1Gi", "pvc-uid7-8", "claim7-8", v1.VolumeBound, v1.PersistentVolumeReclaimDelete, classEmpty, mockDriverName, testNamespace),
					parameters:   map[string]string{"param1": "value1"},
					// information to return
					driverName:   mockDriverName,
					size:         defaultSize,
					snapshotId:   "sid7-8",
					creationTime: timeNow,
					readyToUse:   true,
				},
			},
			errors: []reactorError{
				// Inject error to the forth client.VolumesnapshotV1alpha1().VolumeSnapshots().Update call.
				// All other calls will succeed.
				{"update", "volumesnapshots", errors.New("mock update error")},
				{"update", "volumesnapshots", errors.New("mock update error")},
				{"update", "volumesnapshots", errors.New("mock update error")},
			},
			expectedEvents: []string{"Warning SnapshotCreationFailed"},
			test:           testSyncSnapshot,
		},
		{
			name:              "7-9 - fail create snapshot due to cannot save snapshot content",
			initialContents:   nocontents,
			expectedContents:  nocontents,
			initialSnapshots:  newSnapshotArray("snap7-9", classGold, "", "snapuid7-9", "claim7-9", false, nil, nil, nil),
			expectedSnapshots: newSnapshotArray("snap7-9", classGold, "", "snapuid7-9", "claim7-9", false, nil, metaTimeNowUnix, getSize(defaultSize)),
			initialClaims:     newClaimArray("claim7-9", "pvc-uid7-9", "1Gi", "volume7-9", v1.ClaimBound, &classEmpty),
			initialVolumes:    newVolumeArray("volume7-9", "pv-uid7-9", "pv-handle7-9", "1Gi", "pvc-uid7-9", "claim7-9", v1.VolumeBound, v1.PersistentVolumeReclaimDelete, classEmpty, mockDriverName, testNamespace),
			expectedCreateCalls: []createCall{
				{
					snapshotName: "snapshot-snapuid7-9",
					volume:       newVolume("volume7-9", "pv-uid7-9", "pv-handle7-9", "1Gi", "pvc-uid7-9", "claim7-9", v1.VolumeBound, v1.PersistentVolumeReclaimDelete, classEmpty, mockDriverName, testNamespace),
					parameters:   map[string]string{"param1": "value1"},
					// information to return
					driverName:   mockDriverName,
					size:         defaultSize,
					snapshotId:   "sid7-9",
					creationTime: timeNow,
					readyToUse:   true,
				},
			},
			errors: []reactorError{
				{"create", "volumesnapshotcontents", errors.New("mock create error")},
				{"create", "volumesnapshotcontents", errors.New("mock create error")},
				{"create", "volumesnapshotcontents", errors.New("mock create error")},
			},
			expectedEvents: []string{"Warning CreateSnapshotContentFailed"},
			test:           testSyncSnapshot,
		},
		{
			name:              "7-10 - fail create snapshot with secret not found",
			initialContents:   nocontents,
			expectedContents:  nocontents,
			initialSnapshots:  newSnapshotArray("snap7-10", validSecretClass, "", "snapuid7-10", "claim7-10", false, nil, nil, nil),
			expectedSnapshots: newSnapshotArray("snap7-10", validSecretClass, "", "snapuid7-10", "claim7-10", false, newVolumeError("Failed to create snapshot: failed to get input parameters to create snapshot snap7-10: \"error getting secret secret in namespace default: cannot find secret secret\""), nil, nil),
			initialClaims:     newClaimArray("claim7-10", "pvc-uid7-10", "1Gi", "volume7-10", v1.ClaimBound, &classEmpty),
			initialVolumes:    newVolumeArray("volume7-10", "pv-uid7-10", "pv-handle7-10", "1Gi", "pvc-uid7-10", "claim7-10", v1.VolumeBound, v1.PersistentVolumeReclaimDelete, classEmpty, mockDriverName, testNamespace),
			initialSecrets:    []*v1.Secret{}, // no initial secret created
			errors:            noerrors,
			test:              testSyncSnapshot,
		},
		{
			name:              "8-1 - fail create snapshot with PVC using unbound PV",
			initialContents:   nocontents,
			expectedContents:  nocontents,
			initialSnapshots:  newSnapshotArray("snap8-1", classGold, "", "snapuid8-1", "claim8-1", false, nil, nil, nil),
			expectedSnapshots: newSnapshotArray("snap8-1", classGold, "", "snapuid8-1", "claim8-1", false, newVolumeError("Failed to create snapshot: failed to get input parameters to create snapshot snap8-1: \"claim in dataSource not bound or invalid\""), nil, nil),
			initialClaims:     newClaimArray("claim8-1", "pvc-uid8-1", "1Gi", "volume8-1", v1.ClaimBound, &classEmpty),
			initialVolumes:    newVolumeArray("volume8-1", "pv-uid8-1", "pv-handle8-1", "1Gi", "pvc-uid8-1", "", v1.VolumeBound, v1.PersistentVolumeReclaimDelete, classEmpty, mockDriverName, testNamespace),
			errors:            noerrors,
			expectedEvents:    []string{"Warning SnapshotCreationFailed"},
			test:              testSyncSnapshot,
		},
		{
			name:              "8-2 - fail create snapshot with PVC using PV bound to another PVC (with wrong UID)",
			initialContents:   nocontents,
			expectedContents:  nocontents,
			initialSnapshots:  newSnapshotArray("snap8-2", classGold, "", "snapuid8-2", "claim8-2", false, nil, nil, nil),
			expectedSnapshots: newSnapshotArray("snap8-2", classGold, "", "snapuid8-2", "claim8-2", false, newVolumeError("Failed to create snapshot: failed to get input parameters to create snapshot snap8-2: \"claim in dataSource not bound or invalid\""), nil, nil),
			initialClaims:     newClaimArray("claim8-2", "pvc-uid8-2", "1Gi", "volume8-2", v1.ClaimBound, &classEmpty),
			initialVolumes:    newVolumeArray("volume8-2", "pv-uid8-2", "pv-handle8-2", "1Gi", "pvc-uid8-2-wrong-UID", "claim8-2", v1.VolumeBound, v1.PersistentVolumeReclaimDelete, classEmpty, mockDriverName, testNamespace),
			errors:            noerrors,
			expectedEvents:    []string{"Warning SnapshotCreationFailed"},
			test:              testSyncSnapshot,
		},
		{
			name:              "8-3 - fail create snapshot with PVC using PV bound to another PVC (with wrong namespace)",
			initialContents:   nocontents,
			expectedContents:  nocontents,
			initialSnapshots:  newSnapshotArray("snap8-3", classGold, "", "snapuid8-3", "claim8-3", false, nil, nil, nil),
			expectedSnapshots: newSnapshotArray("snap8-3", classGold, "", "snapuid8-3", "claim8-3", false, newVolumeError("Failed to create snapshot: failed to get input parameters to create snapshot snap8-3: \"claim in dataSource not bound or invalid\""), nil, nil),
			initialClaims:     newClaimArray("claim8-3", "pvc-uid8-3", "1Gi", "volume8-3", v1.ClaimBound, &classEmpty),
			initialVolumes:    newVolumeArray("volume8-3", "pv-uid8-3", "pv-handle8-3", "1Gi", "pvc-uid8-3", "claim8-3", v1.VolumeBound, v1.PersistentVolumeReclaimDelete, classEmpty, mockDriverName, "wrong-namespace"),
			errors:            noerrors,
			expectedEvents:    []string{"Warning SnapshotCreationFailed"},
			test:              testSyncSnapshot,
		},
		{
			name:              "8-4 - fail create snapshot with PVC using PV bound to another PVC (with wrong name)",
			initialContents:   nocontents,
			expectedContents:  nocontents,
			initialSnapshots:  newSnapshotArray("snap8-4", classGold, "", "snapuid8-4", "claim8-4", false, nil, nil, nil),
			expectedSnapshots: newSnapshotArray("snap8-4", classGold, "", "snapuid8-4", "claim8-4", false, newVolumeError("Failed to create snapshot: failed to get input parameters to create snapshot snap8-4: \"claim in dataSource not bound or invalid\""), nil, nil),
			initialClaims:     newClaimArray("claim8-4", "pvc-uid8-4", "1Gi", "volume8-4", v1.ClaimBound, &classEmpty),
			initialVolumes:    newVolumeArray("volume8-4", "pv-uid8-4", "pv-handle8-4", "1Gi", "pvc-uid8-4", "claim8-4-wrong-name", v1.VolumeBound, v1.PersistentVolumeReclaimDelete, classEmpty, mockDriverName, testNamespace),
			errors:            noerrors,
			expectedEvents:    []string{"Warning SnapshotCreationFailed"},
			test:              testSyncSnapshot,
		},
		{
			name:              "8-5 - fail create snapshot with PVC bound to PV with wrong provisioner",
			initialContents:   nocontents,
			expectedContents:  nocontents,
			initialSnapshots:  newSnapshotArray("snap8-5", classGold, "", "snapuid8-5", "claim8-5", false, nil, nil, nil),
			expectedSnapshots: newSnapshotArray("snap8-5", classGold, "", "snapuid8-5", "claim8-5", false, newVolumeError("Failed to create snapshot: failed to get input parameters to create snapshot snap8-5: \"claim in dataSource not bound or invalid\""), nil, nil),
			initialClaims:     newClaimArray("claim8-5", "pvc-uid8-5", "1Gi", "volume8-5", v1.ClaimBound, &classEmpty),
			initialVolumes:    newVolumeArray("volume8-5", "pv-uid8-5", "pv-handle8-5", "1Gi", "pvc-uid8-5", "claim8-5", v1.VolumeBound, v1.PersistentVolumeReclaimDelete, classEmpty, "wrong-driver", testNamespace),
			errors:            noerrors,
			expectedEvents:    []string{"Warning SnapshotCreationFailed"},
			test:              testSyncSnapshot,
		},
	}
	runSyncTests(t, tests, snapshotClasses)
}
