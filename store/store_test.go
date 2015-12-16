package store_test

import (
	"encoding/json"

	. "github.com/jh-bate/fantail-bot/store"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Store", func() {

	type TestData struct {
		Id   int
		Name string
	}

	var (
		store *RedisStore
		t1    TestData
		t2    TestData
	)

	const test_data = "testing_stuff"

	BeforeEach(func() {

		store = NewRedisStore().Set(STORE_TEST_DB)
		store.Pool.Get().Do("FLUSHDB")
		t1 = TestData{Id: 1, Name: "stuff"}
		t2 = TestData{Id: 2, Name: "moar stuff"}
	})

	Describe("When save", func() {
		It("should not return an error", func() {
			Expect(store.Save(test_data, t1)).To(BeNil())
		})
	})

	Describe("When deleting", func() {
		It("should not return an error", func() {
			store.Save(test_data, t1)
			store.Save(test_data, t2)

			Expect(store.Delete(test_data, t1)).To(BeNil())
		})
	})

	Describe("When getting", func() {
		It("should return all saved data", func() {

			store.Save(test_data, t1)
			store.Save(test_data, t2)

			items, err := store.ReadAll(test_data)

			Expect(err).To(BeNil())

			var data []*TestData

			for i := range items {
				var d TestData
				json.Unmarshal(items[i].([]byte), &d)
				data = append(data, &d)
			}

			Expect(len(data)).To(Equal(2))
		})
	})

})
