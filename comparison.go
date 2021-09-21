package linodego

import (
	"encoding/json"
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
	GetChildren() []FilterNode
	Key() string
	JSONValueSegment() interface{}
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

func (f *Filter) MarshalJSON() ([]byte, error) {
	result := make(map[string]interface{})

	if f.OrderBy != "" {
		result["+order_by"] = f.OrderBy
	}

	if f.Order != "" {
		result["+order"] = f.Order
	}

	if f.Operator == "" {
		for _, c := range f.Children {
			result[c.Key()] = c.JSONValueSegment()
		}

		return json.Marshal(result)
	}

	fields := make([]map[string]interface{}, len(f.Children))
	for i, c := range f.Children {
		fields[i] = map[string]interface{}{
			c.Key(): c.JSONValueSegment(),
		}
	}

	result[f.Operator] = fields

	return json.Marshal(result)
}

type Comp struct {
	Column   string
	Operator string
	Value    interface{}
}

func (c *Comp) GetChildren() []FilterNode {
	return []FilterNode{}
}

func (c *Comp) Key() string {
	return c.Column
}

func (c *Comp) JSONValueSegment() interface{} {
	if c.Operator == Eq {
		return c.Value
	}

	return map[string]interface{}{
		c.Operator: c.Value,
	}
}

func Or(order string, orderBy string, nodes ...FilterNode) *Filter {
	return &Filter{"+or", nodes, orderBy, order}
}

func And(order string, orderBy string, nodes ...FilterNode) *Filter {
	return &Filter{"+and", nodes, orderBy, order}
}
