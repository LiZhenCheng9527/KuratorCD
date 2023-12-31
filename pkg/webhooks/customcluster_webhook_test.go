/*
Copyright Kurator Authors.

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

package webhooks

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/yaml"

	"kurator.dev/kurator/pkg/apis/infra/v1alpha1"
)

func TestCustomClusterValidation(t *testing.T) {
	// read configuration from examples directory
	caseNames := []string{"../../examples/infra/customcluster/cc-customcluster.yaml"}

	wh := &CustomClusterWebhook{}
	for _, tt := range caseNames {
		t.Run(tt, func(t *testing.T) {
			g := NewWithT(t)
			c, err := readCustomCluster(tt)
			g.Expect(err).NotTo(HaveOccurred())

			err = wh.validate(c)
			g.Expect(err).NotTo(HaveOccurred())
		})
	}
}

func TestInvalidCustomClusterValidation(t *testing.T) {
	r := path.Join("testdata", "customcluster")
	caseNames := make([]string, 0)
	err := filepath.WalkDir(r, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		caseNames = append(caseNames, path)

		return nil
	})
	assert.NoError(t, err)

	wh := &CustomClusterWebhook{}
	for _, tt := range caseNames {
		t.Run(tt, func(t *testing.T) {
			g := NewWithT(t)
			c, err := readCustomCluster(tt)
			g.Expect(err).NotTo(HaveOccurred())

			err = wh.validate(c)
			g.Expect(err).To(HaveOccurred())
			t.Logf("%v", err)
		})
	}
}

func readCustomCluster(filename string) (*v1alpha1.CustomCluster, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c := &v1alpha1.CustomCluster{}
	if err := yaml.Unmarshal(b, c); err != nil {
		return nil, err
	}

	return c, nil
}
