package config_test

import (
	"github.com/tscolari/memcached-broker/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {

	Describe("#Load", func() {
		It("parses the file correctly", func() {
			config, err := config.Load("./assets/valid.config.yml")
			Expect(err).ToNot(HaveOccurred())

			Expect(len(config.Catalog.Services)).To(Equal(1))
		})

		Context("when the file doesn't exist", func() {
			It("fails", func() {
				_, err := config.Load("./assets/not-here.config.yml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("no such file or directory"))
			})
		})

		Context("when the is not a valid yaml", func() {
			It("fails", func() {
				_, err := config.Load("./assets/invalid.config.yml")
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("#Parse", func() {

		It("parses the data correctly", func() {
			data := `---
catalog:
  services:
  - id: service-id`

			config, err := config.Parse([]byte(data))
			Expect(err).ToNot(HaveOccurred())

			Expect(len(config.Catalog.Services)).To(Equal(1))
			Expect(config.Catalog.Services[0].Id).To(Equal("service-id"))
		})

		Context("when the data is not a valid yaml", func() {
			It("fails", func() {
				data := "not-yaml"
				_, err := config.Parse([]byte(data))
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("cannot unmarshal"))
			})
		})
	})
})
