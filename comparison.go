package linodego

import (
	"fmt"
	"strings"
)

const (
	Eq         = "+eq"
	Neq        = "+neq"
	Gt         = "+gt"
	Gte        = "+gte"
	Lt         = "+lt"
	Lte        = "+lte"
	Contains   = "+contains"
	Ascending  = "asc"
	Descending = "desc"
)

type FilterNode interface {
	JSON() string
	GetChildren() []FilterNode
}

type Filter struct {
	Operator string
	Children []FilterNode
	OrderBy  string
	Order    string
}

func (f *Filter) Add(c *Comp) {
	f.Children = append(f.Children, c)
}

func (f *Filter) GetChildren() []FilterNode {
	return f.Children
}

func (f *Filter) JSON() string {
	children := make([]string, 0, len(f.Children))
	for _, c := range f.Children {
		children = append(children, c.JSON())
	}
	if f.OrderBy != "" {
		orderBy := fmt.Sprintf("\"+order_by\": \"%s\"", f.OrderBy)
		order := fmt.Sprintf("\"+order\": \"%s\"", f.Order)
		if f.Operator == "" {
			return fmt.Sprintf("%s, %s, %s",
				strings.Join(children, ", "), orderBy, order)
		}
		return fmt.Sprintf("\"%s\": [%s], %s, %s", f.Operator,
			strings.Join(children, ", "), orderBy, order)
	}
	if f.Operator == "" {
		return fmt.Sprintf("%s", strings.Join(children, ", "))
	}
	return fmt.Sprintf("\"%s\": [%s]", f.Operator, strings.Join(children, ", "))
}

type Comp struct {
	Column   string
	Operator string
	Value    interface{}
}

func (c *Comp) GetChildren() []FilterNode {
	return []FilterNode{}
}

func (c *Comp) JSON() string {
	if c.Operator == Eq {
		return fmt.Sprintf("{\"%s\": %s}", c.Column,
			getJSONValueString(c.Value))
	}

	return fmt.Sprintf("{\"%s\": {\"%s\": %s}}", c.Column, c.Operator,
		getJSONValueString(c.Value))
}

func Or(order string, orderBy string, nodes ...FilterNode) *Filter {
	return &Filter{"+or", nodes, orderBy, order}
}

func And(order string, orderBy string, nodes ...FilterNode) *Filter {
	return &Filter{"+and", nodes, orderBy, order}
}

func getJSONValueString(value interface{}) string {
	if _, ok := value.(string); ok {
		return fmt.Sprintf("\"%s\"", value)
	}

	return fmt.Sprintf("%v", value)
}
