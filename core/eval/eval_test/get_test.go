package eval_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/shashwatrathod/redis-internals/core/eval"
	"github.com/shashwatrathod/redis-internals/core/resp"
	"github.com/shashwatrathod/redis-internals/core/store"
	mocks "github.com/shashwatrathod/redis-internals/mocks/github.com/shashwatrathod/redis-internals/core/store"
	"github.com/shashwatrathod/redis-internals/utils"
)

func TestEvalGet(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "eval.GET Suite")
}

var _ = Describe("CommandMap.GET", func() {

	var mockStore *mocks.Store

	BeforeEach(func() {
		mockStore = mocks.NewStore(GinkgoT())
	})

	It("should return nil for non-existent key", func() {
		mockStore.EXPECT().Get("nonexistent").Return(nil)

		res := eval.CommandMap[eval.GET].Eval([]string{"nonexistent"}, mockStore)
		Expect(res.Response).To(Equal([]byte("$-1\r\n")))
	})

	It("should return nil for expired key", func() {
		expiredTime := utils.FromExpiryInMilliseconds(-100)
		expiredValue := &store.Value{
			Value:     "expired",
			ValueType: store.String,
			Expiry:    expiredTime, // Expired time
		}
		mockStore.EXPECT().Get("expired").Return(expiredValue)

		res := eval.CommandMap[eval.GET].Eval([]string{"expired"}, mockStore)
		Expect(res.Response).To(Equal([]byte("$-1\r\n")))
	})

	It("should return error for wrong number of arguments", func() {
		res := eval.CommandMap[eval.GET].Eval([]string{}, mockStore)
		Expect(res.Error).To(HaveOccurred())
		Expect(res.Error.Error()).To(Equal("ERR wrong number of arguments for 'get' command"))
	})

	It("should return value for key with no expiry", func() {
		value := &store.Value{
			Value:     "noexpiry",
			ValueType: store.String,
			Expiry:    nil,
		}
		mockStore.EXPECT().Get("noexpiry").Return(value)

		res := eval.CommandMap[eval.GET].Eval([]string{"noexpiry"}, mockStore)
		Expect(res.Response).To(Equal(resp.Encode(value.Value, false)))
	})
})
