package eval_test

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/shashwatrathod/redis-internals/core/eval"
	"github.com/shashwatrathod/redis-internals/core/store"
	mocks "github.com/shashwatrathod/redis-internals/mocks/github.com/shashwatrathod/redis-internals/core/store"
	"github.com/stretchr/testify/mock"
)

func TestEvalSet(t *testing.T) {
	RegisterFailHandler(Fail)
}

var _ = Describe("CommandMap.SET", func() {

	var mockStore *mocks.Store

	BeforeEach(func() {
		mockStore = mocks.NewStore(GinkgoT())
	})

	It("should set a key with a value", func() {
		expectedValue := &store.Value{
			Value:     "value",
			ValueType: store.String,
			Expiry:    nil,
		}

		mockStore.EXPECT().Put("key", expectedValue).Return()

		res := eval.CommandMap[eval.SET].Eval([]string{"key", "value"}, mockStore)
		mockStore.AssertCalled(GinkgoT(), "Put", "key", expectedValue)
		Expect(res.Response).To(Equal([]byte("+OK\r\n")))
		Expect(res.Error).To(BeNil())
	})

	It("should return error for wrong number of arguments", func() {
		res := eval.CommandMap[eval.SET].Eval([]string{"key"}, mockStore)
		Expect(res.Error).To(HaveOccurred())
		Expect(res.Error.Error()).To(Equal("ERR wrong number of arguments for 'set' command"))
		mockStore.AssertNotCalled(GinkgoT(), "Put")
	})

	It("should set a key with an expiry in seconds", func() {
		expectedValue := &store.Value{
			Value:     "value",
			ValueType: store.String,
		}
		mockStore.EXPECT().Put("key", mock.MatchedBy(func(val *store.Value) bool {
			return val.Value == expectedValue.Value &&
				val.ValueType == expectedValue.ValueType &&
				val.Expiry != nil &&
				val.Expiry.GetExpireAtTimestamp().Sub(time.Now().Add(10*time.Second)) < 100*time.Millisecond
		})).Return()

		res := eval.CommandMap[eval.SET].Eval([]string{"key", "value", "EX", "10"}, mockStore)
		mockStore.AssertCalled(GinkgoT(), "Put", "key", mock.MatchedBy(func(val *store.Value) bool {
			return val.Value == expectedValue.Value &&
				val.ValueType == expectedValue.ValueType &&
				val.Expiry != nil &&
				val.Expiry.GetExpireAtTimestamp().Sub(time.Now().Add(10*time.Second)) < 100*time.Millisecond
		}))
		Expect(res.Response).To(Equal([]byte("+OK\r\n")))
		Expect(res.Error).To(BeNil())

	})

	It("should set a key with an expiry in milliseconds", func() {
		expectedValue := &store.Value{
			Value:     "value",
			ValueType: store.String,
		}
		mockStore.EXPECT().Put("key", mock.MatchedBy(func(val *store.Value) bool {
			return val.Value == expectedValue.Value &&
				val.ValueType == expectedValue.ValueType &&
				val.Expiry != nil &&
				val.Expiry.GetExpireAtTimestamp().Sub(time.Now().Add(10000*time.Millisecond)) < 100*time.Millisecond
		})).Return()

		res := eval.CommandMap[eval.SET].Eval([]string{"key", "value", "PX", "10000"}, mockStore)
		mockStore.AssertCalled(GinkgoT(), "Put", "key", mock.MatchedBy(func(val *store.Value) bool {
			return val.Value == expectedValue.Value &&
				val.ValueType == expectedValue.ValueType &&
				val.Expiry != nil &&
				val.Expiry.GetExpireAtTimestamp().Sub(time.Now().Add(10000*time.Millisecond)) < 100*time.Millisecond
		}))
		Expect(res.Response).To(Equal([]byte("+OK\r\n")))
		Expect(res.Error).To(BeNil())
	})

	It("should return error for invalid expiry time", func() {
		res := eval.CommandMap[eval.SET].Eval([]string{"key", "value", "EX", "invalid"}, mockStore)
		Expect(res.Error).To(HaveOccurred())
		Expect(res.Error.Error()).To(Equal("ERR value is not an integer or out of range"))
		mockStore.AssertNotCalled(GinkgoT(), "Put")
	})

	It("should return error for multiple expiry options", func() {
		res := eval.CommandMap[eval.SET].Eval([]string{"key", "value", "EX", "10", "PX", "10000"}, mockStore)
		Expect(res.Error).To(HaveOccurred())
		Expect(res.Error.Error()).To(Equal("ERR syntax error"))
		mockStore.AssertNotCalled(GinkgoT(), "Put")
	})

	It("should return error for unknown option", func() {
		res := eval.CommandMap[eval.SET].Eval([]string{"key", "value", "UNKNOWN"}, mockStore)
		Expect(res.Error).To(HaveOccurred())
		Expect(res.Error.Error()).To(Equal("ERR syntax error"))
		mockStore.AssertNotCalled(GinkgoT(), "Put")
	})
})
