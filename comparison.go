package linodego

/**
 * Pagination and Filtering types and helpers
 */

import (
	"fmt"
	"strings"
)

type ComparisonOperator int

const (
	Eq = iota
	Neq

	Gt
	Gte

	Lt
	Lte

	Contains
)

func (c ComparisonOperator) String() string {
	switch c {
	case Eq:
		return "+eq"
	case Neq:
		return "+neq"
	case Gt:
		return "+gt"
	case Gte:
		return "+gte"
	case Lt:
		return "+lt"
	case Lte:
		return "+lte"
	case Contains:
		return "+contains"
	default:
		return "Unknown ComparisonOperator"
	}
}

type LogicalOperator int

const (
	LogicalAnd = iota
	LogicalOr
)

func (l LogicalOperator) String() string {
	switch l {
	case LogicalAnd:
		return "+and"
	case LogicalOr:
		return "+or"
	default:
		return "Unknown LogicalOperator"
	}
}

type FilterNode interface {
	GetChildren() []FilterNode
	JSON() string
}

type Filter struct {
	Operator LogicalOperator
	Children []FilterNode
}

func (f *Filter) GetChildren() []FilterNode {
	return f.Children
}

func (f *Filter) JSON() string {
	children := make([]string, 0, len(f.Children))
	for _, c := range f.Children {
		children = append(children, c.JSON())
	}
	return fmt.Sprintf("\"%s\": [%s]", f.Operator, strings.Join(children, ", "))
}

type Comparison struct {
	Column   string
	Operator ComparisonOperator
	Value    interface{}
}

func (c *Comparison) GetChildren() []FilterNode {
	return []FilterNode{}
}

func (c *Comparison) JSON() string {
	if c.Operator == Eq {
		return fmt.Sprintf("{\"%s\": %s}", c.Column, getJSONValueString(c.Value))
	}

	return fmt.Sprintf("{\"%s\": {\"%s\": %s}",
		c.Column, c.Operator, getJSONValueString(c.Value))
}

func And(nodes ...FilterNode) *Filter {
	return &Filter{LogicalAnd, nodes}
}

func Or(nodes ...FilterNode) *Filter {
	return &Filter{LogicalOr, nodes}
}

func getJSONValueString(value interface{}) string {
	if _, ok := value.(string); ok {
		return fmt.Sprintf("\"%s\"", value)
	}

	return fmt.Sprintf("%v", value)
}
