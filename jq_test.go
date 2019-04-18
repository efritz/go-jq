package jq

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type JQSuite struct{}

func (s *JQSuite) TestRunBasic(t sweet.T) {
	results, err := Run(".", "foobar")
	Expect(err).To(BeNil())
	Expect(results).To(Equal([]interface{}{"foobar"}))
}

func (s *JQSuite) TestRunArray(t sweet.T) {
	results, err := Run(".[]", []string{"foo", "bar", "baz"})
	Expect(err).To(BeNil())
	Expect(results).To(Equal([]interface{}{"foo", "bar", "baz"}))
}

func (s *JQSuite) TestRunComplex(t sweet.T) {
	results, err := Run(".[] | .[0] + .[1]", [][]int{[]int{1, 2}, []int{11, 22}})
	Expect(err).To(BeNil())
	Expect(results).To(Equal([]interface{}{3, 33}))
}

func (s *JQSuite) TestBadExpression(t sweet.T) {
	_, err := Run(".Item.[]", nil)
	Expect(err).NotTo(BeNil())
	Expect(err.Error()).To(HavePrefix("jq: error: syntax error, unexpected '['"))
}

func (s *JQSuite) TestMarshalNil(t sweet.T) {
	marshaled := marshal(nil)
	Expect(unmarshal(marshaled)).To(BeNil())
	free(marshaled)
}

func (s *JQSuite) TestMarshal(t sweet.T) {
	values := []struct {
		source   interface{}
		expected interface{}
	}{
		{true, true},
		{false, false},
		{"foobar", "foobar"},
		{int(20), 20},
		{uint(20), 20},
		{float32(3.5), 3.5},
		{float64(3.5), 3.5},
		{[]string{"a", "b", "c"}, []interface{}{"a", "b", "c"}},
		{map[string]int{"x": 1, "y": 2, "z": 3}, map[string]interface{}{"x": 1, "y": 2, "z": 3}},
	}

	for _, test := range values {
		marshaled := marshal(test.source)
		Expect(unmarshal(marshaled)).To(Equal(test.expected))
		free(marshaled)
	}
}
