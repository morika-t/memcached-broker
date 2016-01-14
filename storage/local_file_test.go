package storage_test

import (
	"fmt"
	"math/rand"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tscolari/memcached-broker/config"
	"github.com/tscolari/memcached-broker/storage"
)

var _ = Describe("LocalFile", func() {
	var localFile *storage.LocalFile
	var filename string

	BeforeEach(func() {
		filename = fmt.Sprintf("/tmp/state-%d", rand.Int63())
		var err error
		localFile, err = storage.NewLocalFile(filename)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		os.Remove(filename)
	})

	Describe("NewLocalFile", func() {
		It("creates the file if it doesn't exist", func() {
			_, err := os.Stat(filename)
			Expect(os.IsNotExist(err)).To(BeFalse())
		})
	})

	Describe("Get and Put state", func() {
		It("stores and retrieves the state from LocalFile", func() {
			originalState := config.State{
				Capacity: 100,
			}

			localFile.PutState(originalState)
			state := localFile.GetState()

			Expect(state).To(Equal(originalState))
		})
	})

	Describe("Save", func() {
		It("saves the current state to a file", func() {
			originalState := config.State{
				Capacity:  100,
				Instances: map[string]config.Instance{},
			}
			localFile.PutState(originalState)
			err := localFile.Save()
			Expect(err).ToNot(HaveOccurred())

			anotherLocalFile, err := storage.NewLocalFile(filename)
			Expect(err).ToNot(HaveOccurred())
			currentState := anotherLocalFile.GetState()

			Expect(currentState).To(Equal(originalState))
		})
	})

	Describe("Reload", func() {
		It("reloads the localfile based on the local file", func() {
			notInSyncLocalFile, err := storage.NewLocalFile(filename)

			originalState := config.State{
				Capacity:  100,
				Instances: map[string]config.Instance{},
			}
			localFile.PutState(originalState)
			err = localFile.Save()
			Expect(err).ToNot(HaveOccurred())

			Expect(notInSyncLocalFile.GetState()).ToNot(Equal(originalState))
			err = notInSyncLocalFile.Reload()
			Expect(err).ToNot(HaveOccurred())

			Expect(notInSyncLocalFile.GetState()).To(Equal(originalState))
		})
	})
})
