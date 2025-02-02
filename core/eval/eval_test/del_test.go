package eval_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/shashwatrathod/redis-internals/core/eval"
	mocks "github.com/shashwatrathod/redis-internals/mocks/github.com/shashwatrathod/redis-internals/core/store"
)

func TestEvalDel(t *testing.T) {
	RegisterFailHandler(Fail)
}

var _ = Describe("CommandMap.DEL", func() {

	var mockStore *mocks.Store

	BeforeEach(func() {
		mockStore = mocks.NewStore(GinkgoT())
	})

	It("should delete a single existing key", func() {
		mockStore.EXPECT().Delete("key1").Return(true)

		res := eval.CommandMap[eval.DEL].Eval([]string{"key1"}, mockStore)
		Expect(res.Response).To(Equal([]byte(":1\r\n")))
		Expect(res.Error).To(BeNil())

		mockStore.AssertCalled(GinkgoT(), "Delete", "key1")
	})

	It("should delete multiple existing keys", func() {
		mockStore.EXPECT().Delete("key1").Return(true)
		mockStore.EXPECT().Delete("key2").Return(true)

		res := eval.CommandMap[eval.DEL].Eval([]string{"key1", "key2"}, mockStore)
		Expect(res.Response).To(Equal([]byte(":2\r\n")))
		Expect(res.Error).To(BeNil())

		mockStore.AssertCalled(GinkgoT(), "Delete", "key1")
		mockStore.AssertCalled(GinkgoT(), "Delete", "key2")
	})

	It("should return 0 for non-existent keys", func() {
		mockStore.EXPECT().Delete("nonexistent").Return(false)

		res := eval.CommandMap[eval.DEL].Eval([]string{"nonexistent"}, mockStore)
		Expect(res.Response).To(Equal([]byte(":0\r\n")))
		Expect(res.Error).To(BeNil())

		mockStore.AssertCalled(GinkgoT(), "Delete", "nonexistent")
	})

	It("should return the correct count for a mix of existing and non-existent keys", func() {
		mockStore.EXPECT().Delete("key1").Return(true)
		mockStore.EXPECT().Delete("nonexistent").Return(false)

		res := eval.CommandMap[eval.DEL].Eval([]string{"key1", "nonexistent"}, mockStore)
		Expect(res.Response).To(Equal([]byte(":1\r\n")))
		Expect(res.Error).To(BeNil())

		mockStore.AssertCalled(GinkgoT(), "Delete", "key1")
		mockStore.AssertCalled(GinkgoT(), "Delete", "nonexistent")
	})

	It("should return an error for no arguments", func() {
		res := eval.CommandMap[eval.DEL].Eval([]string{}, mockStore)
		Expect(res.Error).To(HaveOccurred())
		Expect(res.Error.Error()).To(Equal("ERR wrong number of arguments for 'del' command"))
		mockStore.AssertNotCalled(GinkgoT(), "Delete")
	})
})
