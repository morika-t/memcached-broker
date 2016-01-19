package config_test

import (
	"github.com/tscolari/memcached-broker/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("State", func() {
	var state config.State

	BeforeEach(func() {
		state = config.NewState(1)
	})

	Describe("InstanceExists", func() {
		Context("when the instance exists", func() {
			BeforeEach(func() {
				state.Instances["instance-id"] = config.Instance{}
			})

			It("returns true", func() {
				Expect(state.InstanceExists("instance-id")).To(BeTrue())
			})
		})

		Context("when the instance doesn't exist", func() {
			It("returns false", func() {
				Expect(state.InstanceExists("not-here-id")).To(BeFalse())
			})
		})
	})

	Describe("Instance", func() {
		It("returns a pointer to the instance", func() {
			instance := config.Instance{
				Host: "127.0.0.1",
				Port: "11111",
			}
			state.Instances["instance-id"] = instance

			fetchedInstance, err := state.Instance("instance-id")
			Expect(err).ToNot(HaveOccurred())
			Expect(fetchedInstance).To(Equal(&instance))
		})

		Context("when the instance doesn't exist", func() {
			It("returns an error", func() {
				_, err := state.Instance("instance-id")
				Expect(err).To(MatchError("Instance not found"))
			})
		})
	})

	Describe("AddInstance", func() {
		It("adds the instance to the state", func() {
			instance := config.Instance{
				Host: "127.0.0.1",
				Port: "11111",
			}

			err := state.AddInstance("instance-id", instance)
			Expect(err).ToNot(HaveOccurred())

			fetchedInstance, err := state.Instance("instance-id")
			Expect(err).ToNot(HaveOccurred())

			Expect(fetchedInstance).To(Equal(&instance))
		})

		It("updates the capacity", func() {
			capacityBefore := state.Capacity
			err := state.AddInstance("instance-id", config.Instance{})
			Expect(err).ToNot(HaveOccurred())

			capacityNow := state.Capacity
			Expect(capacityNow).To(Equal(capacityBefore - 1))
		})

		Context("when there's no capacity left", func() {
			BeforeEach(func() {
				state.Capacity = 0
			})

			It("returns an error", func() {
				err := state.AddInstance("instance-id", config.Instance{})
				Expect(err).To(MatchError("Can't allocate instance, no capacity"))
			})
		})

		Context("when the instance id is taken", func() {
			BeforeEach(func() {
				state.Instances["instance-id"] = config.Instance{}
			})

			It("returns an error", func() {
				err := state.AddInstance("instance-id", config.Instance{})
				Expect(err).To(MatchError("Instance ID is taken"))
			})
		})
	})

	Describe("UpdateInstance", func() {
		BeforeEach(func() {
			instance := config.Instance{
				Host: "127.0.0.1",
				Port: "11111",
			}

			err := state.AddInstance("instance-id", instance)
			Expect(err).ToNot(HaveOccurred())
		})

		It("updates the instance", func() {
			newInstance := config.Instance{
				Host: "0.0.0.0",
				Port: "2222",
			}

			err := state.UpdateInstance("instance-id", newInstance)
			Expect(err).ToNot(HaveOccurred())

			fetchedInstance, err := state.Instance("instance-id")
			Expect(err).ToNot(HaveOccurred())
			Expect(fetchedInstance).To(Equal(&newInstance))
		})

		Context("when the instance is not found", func() {
			It("returns an error", func() {
				err := state.UpdateInstance("instance-id-2", config.Instance{})
				Expect(err).To(MatchError("Instance not found"))
			})
		})
	})

	Describe("DeleteInstance", func() {
		BeforeEach(func() {
			state.Instances["instance-id"] = config.Instance{}
		})

		It("removes the instance", func() {
			err := state.DeleteInstance("instance-id")
			Expect(err).ToNot(HaveOccurred())

			_, err = state.Instance("instance-id")
			Expect(err).To(MatchError("Instance not found"))
		})

		It("updates the capacity", func() {
			capacityBefore := state.Capacity
			err := state.DeleteInstance("instance-id")
			Expect(err).ToNot(HaveOccurred())

			capacityNow := state.Capacity
			Expect(capacityNow).To(Equal(capacityBefore + 1))
		})

		Context("when there's no instance with the given id", func() {
			It("returns an error", func() {
				err := state.DeleteInstance("instance-id-2")
				Expect(err).To(MatchError("Instance not found"))
			})
		})
	})
})
