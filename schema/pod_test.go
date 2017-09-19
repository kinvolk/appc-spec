// Copyright 2015 The appc Authors
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

package schema

import (
	"testing"

	"github.com/kinvolk/spec/schema/types"
)

func TestPodManifestMerge(t *testing.T) {
	pmj := `{}`
	pm := &PodManifest{}

	if pm.UnmarshalJSON([]byte(pmj)) == nil {
		t.Fatal("Manifest JSON without acKind and acVersion unmarshalled successfully")
	}

	pm = BlankPodManifest()

	err := pm.UnmarshalJSON([]byte(pmj))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestAppList(t *testing.T) {
	ri := RuntimeImage{
		ID: *types.NewHashSHA512([]byte{}),
	}
	al := AppList{
		RuntimeApp{
			Name:  "foo",
			Image: ri,
		},
		RuntimeApp{
			Name:  "bar",
			Image: ri,
		},
	}
	if _, err := al.MarshalJSON(); err != nil {
		t.Errorf("want err=nil, got %v", err)
	}
	dal := AppList{
		RuntimeApp{
			Name:  "foo",
			Image: ri,
		},
		RuntimeApp{
			Name:  "bar",
			Image: ri,
		},
		RuntimeApp{
			Name:  "foo",
			Image: ri,
		},
	}
	if _, err := dal.MarshalJSON(); err == nil {
		t.Errorf("want err, got nil")
	}
}

func TestPodManifestUniqueVolumeNames(t *testing.T) {
	pm := BlankPodManifest()
	volA, err := types.VolumeFromString("simple,kind=empty")
	if err != nil {
		t.Errorf("expected nil error, got %q", err)
	}
	volB, err := types.VolumeFromString("simple,kind=host,source=/tmp")
	if err != nil {
		t.Errorf("expected nil error, got %q", err)
	}
	pm.Volumes = append(pm.Volumes, *volA, *volB)
	if _, err := pm.MarshalJSON(); err == nil {
		t.Errorf("expected duplicate volume name error, got nil")
	}

	pmj := `
{
        "acVersion": "0.8.11",
        "acKind": "PodManifest",
        "apps": [
        ],
        "volumes": [
                {
                        "name": "simplename",
                        "kind": "empty"
                },
                {
                        "name": "simplename",
                        "kind": "host",
                        "source": "/tmp"
                }
        ]
}
`
	err = pm.UnmarshalJSON([]byte(pmj))
	if err == nil {
		t.Errorf("expected duplicate volume name error, got nil")
	}
}
