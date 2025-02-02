package resp

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestResp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "RESP Suite")
}

type DecodeTestCase struct {
	Input          string
	ExpectedOutput interface{}
	ExpectError    bool
}

func getDecodeTestCase(input string, expectedOutput interface{}, expectError bool) DecodeTestCase {
	ret := DecodeTestCase{
		Input:          input,
		ExpectedOutput: expectedOutput,
		ExpectError:    expectError,
	}
	return ret
}

func runTestAgainstCase(testCase DecodeTestCase) {
	result, err := Decode([]byte(testCase.Input))

	validateResult(testCase, result, err)
}

func validateResult(testCase DecodeTestCase, result interface{}, err error) {
	if !testCase.ExpectError {
		Expect(err).ToNot(HaveOccurred())
	} else {
		Expect(err).To(HaveOccurred())
	}

	Expect(result).To(Equal(testCase.ExpectedOutput))
}

var _ = Describe("Decode", func() {
	Context("Simple String", func() {
		It("Should decode a simple string", func() {
			test := getDecodeTestCase("+OK\r\n", "OK", false)
			runTestAgainstCase(test)
		})
	})

	Context("Bulk String", func() {
		It("Should decode a bulk string", func() {
			test := getDecodeTestCase("$5\r\nHello\r\n", "Hello", false)
			runTestAgainstCase(test)
		})
	})

	Context("Integer64", func() {
		It("Should decode an integer without any signs", func() {
			test := getDecodeTestCase(":2345\r\n", int64(2345), false)
			runTestAgainstCase(test)
		})

		It("Should decode an integer with + sign", func() {
			test := getDecodeTestCase(":+2345\r\n", int64(2345), false)
			runTestAgainstCase(test)
		})

		It("Should decode an integer with - sign", func() {
			test := getDecodeTestCase(":-2345\r\n", int64(-2345), false)
			runTestAgainstCase(test)
		})
	})

	Context("Simple Error", func() {
		It("Should decode a simple error", func() {
			test := getDecodeTestCase("-Error message\r\n", "Error message", false)
			runTestAgainstCase(test)
		})
	})

	Context("Arrays", func() {
		It("Should decode an empty array", func() {
			test := getDecodeTestCase("*0\r\n", []interface{}{}, false)
			runTestAgainstCase(test)
		})

		It("Should decode an array with simple strings", func() {
			test := getDecodeTestCase("*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n", []interface{}{"hello", "world"}, false)
			runTestAgainstCase(test)
		})

		It("Should decode an array with int64", func() {
			test := getDecodeTestCase("*3\r\n:1\r\n:2\r\n:3\r\n", []interface{}{int64(1), int64(2), int64(3)}, false)
			runTestAgainstCase(test)
		})

		It("Should decode an array with different datatypes", func() {
			test := getDecodeTestCase("*5\r\n:1\r\n:2\r\n:3\r\n:4\r\n$5\r\nhello\r\n",
				[]interface{}{int64(1), int64(2), int64(3), int64(4), "hello"},
				false)
			runTestAgainstCase(test)
		})

		It("Should decode nested arrays", func() {
			test := getDecodeTestCase("*2\r\n*3\r\n:1\r\n:2\r\n:3\r\n*3\r\n+Hello\r\n-World\r\n:-19\r\n",
				[]interface{}{
					[]interface{}{int64(1), int64(2), int64(3)},
					[]interface{}{"Hello", "World", int64(-19)}},
				false)
			runTestAgainstCase(test)
		})
	})
})
