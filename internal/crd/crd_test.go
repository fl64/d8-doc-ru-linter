package crd

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("d8-doc-ru-linter:crd", func() {

	Context("CRD should be normalized", func() {
		It("CRD should be normalized", func() {
			var err error
			var srcCRD CRD
			err = srcCRD.Load("../../testdata/src.yaml")
			Expect(err).NotTo(HaveOccurred())

			resultYAML, err := os.ReadFile("../../testdata/result-normalized.yaml")
			Expect(err).NotTo(HaveOccurred())

			srcCRDYAML, err := srcCRD.Marshal()
			Expect(err).NotTo(HaveOccurred())
			Expect(srcCRDYAML).To(MatchYAML(resultYAML))
		})
	})

	Context("CRDs should merged", func() {
		It("CRDs should merged", func() {
			var err error
			var srcCRD CRD
			err = srcCRD.Load("../../testdata/src.yaml")
			Expect(err).NotTo(HaveOccurred())

			var dstCRD CRD
			err = dstCRD.Load("../../testdata/dst.yaml")
			Expect(err).NotTo(HaveOccurred())

			mergedCRDData, ops := srcCRD.CompareWith(dstCRD)
			mergedCRDYAML, err := mergedCRDData.Marshal()
			Expect(err).NotTo(HaveOccurred())

			resultYAML, err := os.ReadFile("../../testdata/result-merged.yaml")
			Expect(err).NotTo(HaveOccurred())

			Expect(mergedCRDYAML).To(MatchYAML(resultYAML))

			opsJSONReport, err := ops.MarshalJSONReport()
			Expect(err).NotTo(HaveOccurred())

			Expect(opsJSONReport).To(MatchJSON(`{
        "count": 4,
        "operations": [
          {
            "path": "/spec/versions/v2",
            "op": "add"
          },
          {
            "path": "/spec/versions/v1/schema/openAPIV3Schema/properties/p1-object/properties/p1-1-object/properties/only-in-src",
            "op": "add"
          },
          {
            "path": "/spec/versions/v1/schema/openAPIV3Schema/properties/p1-object/properties/p1-1-object/properties/only-in-dst",
            "op": "delete"
          },
          {
            "path": "/spec/versions/v0",
            "op": "delete"
          }
        ]
      }`))

		})
	})

})

func TestFencingAgent(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "d8-doc-ru-linter Suite")
}
