package eval_test

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/shashwatrathod/redis-internals/core/eval"
	"github.com/shashwatrathod/redis-internals/core/store"
	mocks "github.com/shashwatrathod/redis-internals/mocks/github.com/shashwatrathod/redis-internals/core/store"
	"github.com/shashwatrathod/redis-internals/utils"
)

func TestEvalTtl(t *testing.T) {
	RegisterFailHandler(Fail)
}

var _ = Describe("CommandMap.TTL", func() {

	var mockStore *mocks.Store

	BeforeEach(func() {
		mockStore = mocks.NewStore(GinkgoT())
	})

	It("should return -2 for non-existent key", func() {
		mockStore.On("Get", "nonexistent").Return(nil)

		res := eval.CommandMap[eval.TTL].Eval([]string{"nonexistent"}, mockStore)
		Expect(res.Response).To(Equal([]byte(":-2\r\n")))
		Expect(res.Error).To(BeNil())

		mockStore.AssertCalled(GinkgoT(), "Get", "nonexistent")
	})

	It("should return -1 for key without expiry", func() {
		val := &store.Value{
			Value:     "value",
			ValueType: store.String,
			Expiry:    nil,
		}
		mockStore.On("Get", "key").Return(val)

		res := eval.CommandMap[eval.TTL].Eval([]string{"key"}, mockStore)
		Expect(res.Response).To(Equal([]byte(":-1\r\n")))
		Expect(res.Error).To(BeNil())

		mockStore.AssertCalled(GinkgoT(), "Get", "key")
	})

	It("should return -2 for expired key", func() {
		expiryTime := utils.FromExpiryInMilliseconds(-1000)
		val := &store.Value{
			Value:     "value",
			ValueType: store.String,
			Expiry:    expiryTime,
		}
		mockStore.On("Get", "key").Return(val)

		res := eval.CommandMap[eval.TTL].Eval([]string{"key"}, mockStore)
		Expect(res.Response).To(Equal([]byte(":-2\r\n")))
		Expect(res.Error).To(BeNil())

		mockStore.AssertCalled(GinkgoT(), "Get", "key")
	})

	It("should return remaining TTL for key with expiry", func() {
		expiryTime := utils.FromExpiryInSeconds(10)
		val := &store.Value{
			Value:     "value",
			ValueType: store.String,
			Expiry:    expiryTime,
		}
		mockStore.On("Get", "key").Return(val)

		res := eval.CommandMap[eval.TTL].Eval([]string{"key"}, mockStore)
		Expect(res.Response).To(Equal([]byte(":10\r\n")))
		Expect(res.Error).To(BeNil())
		mockStore.AssertCalled(GinkgoT(), "Get", "key")

		time.Sleep(2 * time.Second)

		res = eval.CommandMap[eval.TTL].Eval([]string{"key"}, mockStore)
		Expect(res.Response).To(Equal([]byte(":8\r\n")))
		Expect(res.Error).To(BeNil())
		mockStore.AssertCalled(GinkgoT(), "Get", "key")
	})

	It("should return remaining TTL for key with expiry rounded to the nearest ceil integer", func() {
		expiryTime := utils.FromExpiryInSeconds(10)
		val := &store.Value{
			Value:     "value",
			ValueType: store.String,
			Expiry:    expiryTime,
		}
		mockStore.On("Get", "key").Return(val)

		res := eval.CommandMap[eval.TTL].Eval([]string{"key"}, mockStore)
		Expect(res.Response).To(Equal([]byte(":10\r\n")))
		Expect(res.Error).To(BeNil())
		mockStore.AssertCalled(GinkgoT(), "Get", "key")

		// Advance time by 1.2s.
		time.Sleep(1200 * time.Millisecond)

		res = eval.CommandMap[eval.TTL].Eval([]string{"key"}, mockStore)
		Expect(res.Response).To(Equal([]byte(":9\r\n"))) // Nearest time is still 9s.
		Expect(res.Error).To(BeNil())
		mockStore.AssertCalled(GinkgoT(), "Get", "key")
	})

	It("should return remaining TTL for key with expiry rounded to the nearest floor integer", func() {
		expiryTime := utils.FromExpiryInSeconds(10)
		val := &store.Value{
			Value:     "value",
			ValueType: store.String,
			Expiry:    expiryTime,
		}
		mockStore.On("Get", "key").Return(val)

		res := eval.CommandMap[eval.TTL].Eval([]string{"key"}, mockStore)
		Expect(res.Response).To(Equal([]byte(":10\r\n")))
		Expect(res.Error).To(BeNil())
		mockStore.AssertCalled(GinkgoT(), "Get", "key")

		// Advance time by 1.8s.
		time.Sleep(1800 * time.Millisecond)

		res = eval.CommandMap[eval.TTL].Eval([]string{"key"}, mockStore)
		Expect(res.Response).To(Equal([]byte(":8\r\n"))) // Nearest time is still 9s.
		Expect(res.Error).To(BeNil())
		mockStore.AssertCalled(GinkgoT(), "Get", "key")
	})

	It("should return error for wrong number of arguments", func() {
		res := eval.CommandMap[eval.TTL].Eval([]string{}, mockStore)
		Expect(res.Error).To(HaveOccurred())
		Expect(res.Error.Error()).To(Equal("ERR wrong number of arguments for 'ttl' command"))
	})
})
