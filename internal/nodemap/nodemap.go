// Package nodemap provides a schema registry index type.
package nodemap

import (
	"github.com/MadBase/go-capnproto2/v2"
	"github.com/MadBase/go-capnproto2/v2/internal/schema"
	"github.com/MadBase/go-capnproto2/v2/schemas"
)

// Map is a lazy index of a registry.
// The zero value is an index of the default registry.
type Map struct {
	reg   *schemas.Registry
	nodes map[uint64]schema.Node
}

func (m *Map) registry() *schemas.Registry {
	if m.reg != nil {
		return m.reg
	}
	return &schemas.DefaultRegistry
}

func (m *Map) UseRegistry(reg *schemas.Registry) {
	m.reg = reg
	m.nodes = make(map[uint64]schema.Node)
}

// Find returns the node for the given ID.
func (m *Map) Find(id uint64) (schema.Node, error) {
	if n := m.nodes[id]; n.IsValid() {
		return n, nil
	}
	data, err := m.registry().Find(id)
	if err != nil {
		return schema.Node{}, err
	}
	msg, err := capnp.Unmarshal(data)
	if err != nil {
		return schema.Node{}, err
	}
	req, err := schema.ReadRootCodeGeneratorRequest(msg)
	if err != nil {
		return schema.Node{}, err
	}
	nodes, err := req.Nodes()
	if err != nil {
		return schema.Node{}, err
	}
	if m.nodes == nil {
		m.nodes = make(map[uint64]schema.Node)
	}
	for i := 0; i < nodes.Len(); i++ {
		n := nodes.At(i)
		m.nodes[n.Id()] = n
	}
	return m.nodes[id], nil
}
