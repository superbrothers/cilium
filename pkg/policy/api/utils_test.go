// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

//go:build !privileged_tests
// +build !privileged_tests

package api

import (
	"testing"

	. "gopkg.in/check.v1"

	"github.com/cilium/cilium/pkg/policy/api/kafka"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) {
	TestingT(t)
}

type PolicyAPITestSuite struct{}

var _ = Suite(&PolicyAPITestSuite{})

func (s *PolicyAPITestSuite) TestHTTPEqual(c *C) {
	rule1 := PortRuleHTTP{Path: "/foo$", Method: "GET", Headers: []string{"X-Test: Foo"}}
	rule2 := PortRuleHTTP{Path: "/bar$", Method: "GET", Headers: []string{"X-Test: Foo"}}
	rule3 := PortRuleHTTP{Path: "/foo$", Method: "GET", Headers: []string{"X-Test: Bar"}}

	c.Assert(rule1.Equal(rule1), Equals, true)
	c.Assert(rule1.Equal(rule2), Equals, false)
	c.Assert(rule1.Equal(rule3), Equals, false)

	rules := L7Rules{
		HTTP: []PortRuleHTTP{rule1, rule2},
	}

	c.Assert(rule1.Exists(rules), Equals, true)
	c.Assert(rule2.Exists(rules), Equals, true)
	c.Assert(rule3.Exists(rules), Equals, false)
}

func (s *PolicyAPITestSuite) TestKafkaEqual(c *C) {
	rule1 := kafka.PortRule{APIVersion: "1", APIKey: "foo", Topic: "topic1"}
	rule2 := kafka.PortRule{APIVersion: "1", APIKey: "bar", Topic: "topic1"}
	rule3 := kafka.PortRule{APIVersion: "1", APIKey: "foo", Topic: "topic2"}

	c.Assert(rule1, Equals, rule1)
	c.Assert(rule1, Not(Equals), rule2)
	c.Assert(rule1, Not(Equals), rule3)

	rules := L7Rules{
		Kafka: []kafka.PortRule{rule1, rule2},
	}

	c.Assert(rule1.Exists(rules.Kafka), Equals, true)
	c.Assert(rule2.Exists(rules.Kafka), Equals, true)
	c.Assert(rule3.Exists(rules.Kafka), Equals, false)
}

func (s *PolicyAPITestSuite) TestL7Equal(c *C) {
	rule1 := PortRuleL7{"Path": "/foo$", "Method": "GET"}
	rule2 := PortRuleL7{"Path": "/bar$", "Method": "GET"}
	rule3 := PortRuleL7{"Path": "/foo$", "Method": "GET", "extra": ""}

	c.Assert(rule1.Equal(rule1), Equals, true)
	c.Assert(rule2.Equal(rule2), Equals, true)
	c.Assert(rule3.Equal(rule3), Equals, true)
	c.Assert(rule1.Equal(rule2), Equals, false)
	c.Assert(rule2.Equal(rule1), Equals, false)
	c.Assert(rule1.Equal(rule3), Equals, false)
	c.Assert(rule3.Equal(rule1), Equals, false)
	c.Assert(rule2.Equal(rule3), Equals, false)
	c.Assert(rule3.Equal(rule2), Equals, false)

	rules := L7Rules{
		L7Proto: "testing",
		L7:      []PortRuleL7{rule1, rule2},
	}

	c.Assert(rule1.Exists(rules), Equals, true)
	c.Assert(rule2.Exists(rules), Equals, true)
	c.Assert(rule3.Exists(rules), Equals, false)
}

func (s *PolicyAPITestSuite) TestValidateL4Proto(c *C) {
	c.Assert(L4Proto("TCP").Validate(), IsNil)
	c.Assert(L4Proto("UDP").Validate(), IsNil)
	c.Assert(L4Proto("ANY").Validate(), IsNil)
	c.Assert(L4Proto("TCP2").Validate(), Not(IsNil))
	c.Assert(L4Proto("t").Validate(), Not(IsNil))
}

func (s *PolicyAPITestSuite) TestParseL4Proto(c *C) {
	p, err := ParseL4Proto("tcp")
	c.Assert(p, Equals, ProtoTCP)
	c.Assert(err, IsNil)

	p, err = ParseL4Proto("Any")
	c.Assert(p, Equals, ProtoAny)
	c.Assert(err, IsNil)

	p, err = ParseL4Proto("")
	c.Assert(p, Equals, ProtoAny)
	c.Assert(err, IsNil)

	_, err = ParseL4Proto("foo2")
	c.Assert(err, Not(IsNil))
}
