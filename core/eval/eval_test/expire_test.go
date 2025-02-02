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

func TestEvalExpire(t *testing.T) {
	RegisterFailHandler(Fail)
}

var _ = Describe("CommandMap.EXPIRE", func() {

	var mockStore *mocks.Store

	BeforeEach(func() {
		mockStore = mocks.NewStore(GinkgoT())
	})

	It("should return error for wrong number of arguments", func() {
		res := eval.CommandMap[eval.EXPIRE].Eval([]string{"key"}, mockStore)
		Expect(res.Error).To(HaveOccurred())
		Expect(res.Error.Error()).To(Equal("ERR wrong number of arguments for 'expire' command"))
	})

	It("should return 0 for non-existent key", func() {
		mockStore.On("Get", "nonexistent").Return(nil)

		res := eval.CommandMap[eval.EXPIRE].Eval([]string{"nonexistent", "10"}, mockStore)
		Expect(res.Response).To(Equal([]byte(":0\r\n")))
		Expect(res.Error).To(BeNil())

		mockStore.AssertCalled(GinkgoT(), "Get", "nonexistent")
	})

	It("should return 0 for key without expiry", func() {
		val := &store.Value{
			Value:     "value",
			ValueType: store.String,
			Expiry:    nil,
		}
		mockStore.On("Get", "key").Return(val)

		res := eval.CommandMap[eval.EXPIRE].Eval([]string{"key", "10"}, mockStore)
		Expect(res.Response).To(Equal([]byte(":0\r\n")))
		Expect(res.Error).To(BeNil())

		mockStore.AssertCalled(GinkgoT(), "Get", "key")
	})

	It("should return 0 for expired key", func() {
		expiryTime := utils.FromExpiryInMilliseconds(-1000)
		val := &store.Value{
			Value:     "value",
			ValueType: store.String,
			Expiry:    expiryTime,
		}
		mockStore.On("Get", "key").Return(val)

		res := eval.CommandMap[eval.EXPIRE].Eval([]string{"key", "10"}, mockStore)
		Expect(res.Response).To(Equal([]byte(":0\r\n")))
		Expect(res.Error).To(BeNil())

		mockStore.AssertCalled(GinkgoT(), "Get", "key")
	})

	It("should set expiry for key with existing expiry", func() {
		expiryTime := utils.FromExpiryInSeconds(10)
		val := &store.Value{
			Value:     "value",
			ValueType: store.String,
			Expiry:    expiryTime,
		}
		mockStore.On("Get", "key").Return(val)

		res := eval.CommandMap[eval.EXPIRE].Eval([]string{"key", "20"}, mockStore)
		Expect(res.Response).To(Equal([]byte(":1\r\n")))
		Expect(res.Error).To(BeNil())

		Expect(val.Expiry.GetExpireAtTimestamp().Sub(time.Now().Add(20*time.Second)) < 100*time.Millisecond).To(BeTrue())
	})

	It("should return error for invalid expiry time", func() {
		res := eval.CommandMap[eval.EXPIRE].Eval([]string{"key", "invalid"}, mockStore)
		Expect(res.Error).To(HaveOccurred())
		Expect(res.Error.Error()).To(Equal("ERR unknown command 'EXPIRE', with args beginning with: 'key' 'invalid'"))
		mockStore.AssertNotCalled(GinkgoT(), "Get")
	})
})
