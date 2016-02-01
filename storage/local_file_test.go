package storage_test

import (
	"fmt"
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tscolari/cf-broker-api/common/repository"
	"github.com/tscolari/memcached-broker/storage"
)

var _ = Describe("LocalFile", func() {
	var localFile *storage.LocalFile
	var tempFileName string

	BeforeEach(func() {
		dir, err := ioutil.TempDir("/tmp/", "local-file")
		Expect(err).ToNot(HaveOccurred())

		tempFileName = fmt.Sprintf("%s/state.yml", dir)
		localFile, err = storage.NewLocalFile(tempFileName, 1)
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("InstanceExists", func() {
		Context("when the instance exists", func() {
			BeforeEach(func() {
				localFile.AddInstance(repository.Instance{ID: "instance-id"})
			})

			It("returns true", func() {
				Expect(localFile.InstanceExists("instance-id")).To(BeTrue())
			})
		})

		Context("when the instance doesn't exist", func() {
			It("returns false", func() {
				Expect(localFile.InstanceExists("not-here-id")).To(BeFalse())
			})
		})
	})

	Describe("Instance", func() {
		It("returns a pointer to the instance", func() {
			instance := repository.Instance{
				ID:   "instance-id",
				Host: "127.0.0.1",
				Port: "11111",
			}
			localFile.AddInstance(instance)

			fetchedInstance, err := localFile.Instance("instance-id")
			Expect(err).ToNot(HaveOccurred())
			Expect(fetchedInstance).To(Equal(&instance))
		})

		Context("when the instance doesn't exist", func() {
			It("returns an error", func() {
				_, err := localFile.Instance("instance-id")
				Expect(err).To(MatchError("Instance not found"))
			})
		})
	})

	Describe("AddInstance", func() {
		It("adds the instance to the object", func() {
			instance := repository.Instance{
				ID:   "instance-id",
				Host: "127.0.0.1",
				Port: "11111",
			}

			err := localFile.AddInstance(instance)
			Expect(err).ToNot(HaveOccurred())

			fetchedInstance, err := localFile.Instance("instance-id")
			Expect(err).ToNot(HaveOccurred())

			Expect(fetchedInstance).To(Equal(&instance))
		})

		It("updates the capacity", func() {
			capacityBefore := localFile.AvailableInstances()
			err := localFile.AddInstance(repository.Instance{ID: "instance-id"})
			Expect(err).ToNot(HaveOccurred())

			capacityNow := localFile.AvailableInstances()
			Expect(capacityNow).To(Equal(capacityBefore - 1))
		})

		It("persists the change on disk", func() {
			instance := repository.Instance{
				ID:       "instance-id",
				Host:     "127.0.0.1",
				Port:     "11111",
				Bindings: []string{},
			}

			err := localFile.AddInstance(instance)
			Expect(err).ToNot(HaveOccurred())

			newLocalFile, err := storage.NewLocalFile(tempFileName, -10)
			Expect(err).ToNot(HaveOccurred())

			Expect(newLocalFile.AvailableInstances()).To(Equal(localFile.AvailableInstances()))

			fetchedInstance, err := newLocalFile.Instance("instance-id")
			Expect(err).ToNot(HaveOccurred())
			Expect(fetchedInstance).To(Equal(&instance))
		})

		Context("when there's no capacity left", func() {
			BeforeEach(func() {
				var err error
				localFile, err = storage.NewLocalFile(tempFileName, 0)
				Expect(err).ToNot(HaveOccurred())
			})

			It("returns an error", func() {
				err := localFile.AddInstance(repository.Instance{ID: "instance-id"})
				Expect(err).To(MatchError("Can't allocate instance, no capacity"))
			})
		})

		Context("when the instance id is taken", func() {
			BeforeEach(func() {
				var err error
				localFile, err = storage.NewLocalFile(tempFileName, 5)
				Expect(err).ToNot(HaveOccurred())
				err = localFile.AddInstance(repository.Instance{ID: "instance-id"})
				Expect(err).ToNot(HaveOccurred())
			})

			It("returns an error", func() {
				err := localFile.AddInstance(repository.Instance{ID: "instance-id"})
				Expect(err).To(MatchError("Instance ID is taken"))
			})
		})
	})

	Describe("UpdateInstance", func() {
		BeforeEach(func() {
			instance := repository.Instance{
				ID:   "instance-id",
				Host: "127.0.0.1",
				Port: "11111",
			}

			err := localFile.AddInstance(instance)
			Expect(err).ToNot(HaveOccurred())
		})

		It("updates the instance", func() {
			newInstance := repository.Instance{
				ID:   "instance-id",
				Host: "0.0.0.0",
				Port: "2222",
			}

			err := localFile.UpdateInstance(newInstance)
			Expect(err).ToNot(HaveOccurred())

			fetchedInstance, err := localFile.Instance("instance-id")
			Expect(err).ToNot(HaveOccurred())
			Expect(fetchedInstance).To(Equal(&newInstance))
		})

		It("persists the change on disk", func() {
			newInstance := repository.Instance{
				ID:       "instance-id",
				Host:     "0.0.0.0",
				Port:     "2222",
				Bindings: []string{},
			}

			err := localFile.UpdateInstance(newInstance)
			Expect(err).ToNot(HaveOccurred())

			newLocalFile, err := storage.NewLocalFile(tempFileName, -10)
			Expect(err).ToNot(HaveOccurred())

			fetchedInstance, err := newLocalFile.Instance("instance-id")
			Expect(err).ToNot(HaveOccurred())
			Expect(fetchedInstance).To(Equal(&newInstance))
		})

		Context("when the instance is not found", func() {
			It("returns an error", func() {
				err := localFile.UpdateInstance(repository.Instance{ID: "instance-id-2"})
				Expect(err).To(MatchError("Instance not found"))
			})
		})
	})

	Describe("DeleteInstance", func() {
		BeforeEach(func() {
			localFile.AddInstance(repository.Instance{ID: "instance-id"})
		})

		It("removes the instance", func() {
			err := localFile.DeleteInstance("instance-id")
			Expect(err).ToNot(HaveOccurred())

			_, err = localFile.Instance("instance-id")
			Expect(err).To(MatchError("Instance not found"))
		})

		It("updates the capacity", func() {
			capacityBefore := localFile.AvailableInstances()
			err := localFile.DeleteInstance("instance-id")
			Expect(err).ToNot(HaveOccurred())

			capacityNow := localFile.AvailableInstances()
			Expect(capacityNow).To(Equal(capacityBefore + 1))
		})

		It("persists the change on disk", func() {
			capacityBefore := localFile.AvailableInstances()
			err := localFile.DeleteInstance("instance-id")
			Expect(err).ToNot(HaveOccurred())

			newLocalFile, err := storage.NewLocalFile(tempFileName, -10)
			Expect(err).ToNot(HaveOccurred())

			fetchedCapacity := newLocalFile.AvailableInstances()
			Expect(fetchedCapacity).To(Equal(capacityBefore + 1))

			Expect(newLocalFile.InstanceExists("instance-id")).To(BeFalse())
		})

		Context("when there's no instance with the given id", func() {
			It("returns an error", func() {
				err := localFile.DeleteInstance("instance-id-2")
				Expect(err).To(MatchError("Instance not found"))
			})
		})
	})

	Describe("InstanceBindingExists", func() {
		Context("when the binding exists for the given instance", func() {
			BeforeEach(func() {
				instance := repository.Instance{
					ID:       "instance-id",
					Bindings: []string{"binding-1"},
				}
				localFile.AddInstance(instance)
			})

			It("returns true", func() {
				Expect(localFile.InstanceBindingExists("instance-id", "binding-1")).To(BeTrue())
			})
		})

		Context("when the instance doesn't exist", func() {
			It("returns false", func() {
				Expect(localFile.InstanceBindingExists("instance-id", "binding-1")).To(BeFalse())
			})
		})

		Context("when the instance binding doesn't exist", func() {
			BeforeEach(func() {
				localFile.AddInstance(repository.Instance{ID: "instance-id"})
			})

			It("returns false", func() {
				Expect(localFile.InstanceBindingExists("instance-id", "binding-1")).To(BeFalse())
			})
		})
	})

	Describe("AddInstanceBinding", func() {
		BeforeEach(func() {
			instance := repository.Instance{
				ID:       "instance-id",
				Host:     "127.0.0.1",
				Port:     "11111",
				Bindings: []string{"existing-binding"},
			}

			err := localFile.AddInstance(instance)
			Expect(err).ToNot(HaveOccurred())
		})

		It("adds the binding to the instance", func() {
			err := localFile.AddInstanceBinding("instance-id", "binding-1")
			Expect(err).ToNot(HaveOccurred())

			Expect(localFile.InstanceBindingExists("instance-id", "binding-1")).To(BeTrue())
		})

		Context("when the instance doesn't exist", func() {
			It("returns an error", func() {
				err := localFile.AddInstanceBinding("instance-id-2", "binding-1")
				Expect(err).To(MatchError("Instance not found"))
			})
		})

		Context("when the binding id is already taken", func() {
			It("returns an error", func() {
				err := localFile.AddInstanceBinding("instance-id", "existing-binding")
				Expect(err).To(MatchError("Binding ID is taken"))
			})
		})
	})

	Describe("DeleteInstanceBinding", func() {
		BeforeEach(func() {
			instance := repository.Instance{
				ID:       "instance-id",
				Host:     "127.0.0.1",
				Port:     "11111",
				Bindings: []string{"existing-binding"},
			}

			err := localFile.AddInstance(instance)
			Expect(err).ToNot(HaveOccurred())
		})

		It("deletes the binding from the instance", func() {
			err := localFile.DeleteInstanceBinding("instance-id", "existing-binding")
			Expect(err).ToNot(HaveOccurred())

			Expect(localFile.InstanceBindingExists("instance-id", "existing-binding")).To(BeFalse())
		})

		Context("when the instance doesn't exist", func() {
			It("returns an error", func() {
				err := localFile.DeleteInstanceBinding("instance-id-2", "binding-1")
				Expect(err).To(MatchError("Instance not found"))
			})
		})

		Context("when the binding id doesn't exist", func() {
			It("returns an error", func() {
				err := localFile.DeleteInstanceBinding("instance-id", "binding-1")
				Expect(err).To(MatchError("Binding not found"))
			})
		})
	})
})
