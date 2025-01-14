package org

import (
	"errors"
	"fmt"
	"strings"
)

type Drawer struct {
  Name string
  Elements []Element
}

func (d Drawer) Kind() ElementKind {
  return ELEMENT_DRAWER
}

func (d Drawer) IsGreaterElement() bool {
  return true
}

func (d *Drawer) String() string {
  return strings.Join(d.Strings(), "\n")
}

func (d *Drawer) Strings() []string {
  out := make([]string, 0)
  dOpen := fmt.Sprintf(":%s:", strings.ToUpper(d.Name))
  dClose := fmt.Sprintf(":END:")
  out = append(out, dOpen)

  for _, elem := range d.Elements {
    out = append(out, elem.Strings()...)
  }

  return append(out, dClose)
}

func (d *Drawer) AddElement(e Element) (*Drawer, error) {
  switch e.Kind() {
  case ELEMENT_HEADING:
    fallthrough
  case ELEMENT_DRAWER:
    fallthrough
  case ELEMENT_PROPERTY_DRAWER:
    return nil, errors.New("Invalid element kind for drawer")
  }

  d.Elements = append(d.Elements, e)

  return d, nil
}

type PropertyDrawer struct {
  Node        *Node
  Properties  map[string]*Property
}

func (pd PropertyDrawer) Kind() ElementKind {
  return ELEMENT_PROPERTY_DRAWER
}

func (pd PropertyDrawer) IsGreaterElement() bool {
  return true
}

func (pd *PropertyDrawer) Add(p *Property) *PropertyDrawer {
  // it is unclear if inheritance collisions should be left or right
  // biased, I will need to test. For now I am setting right-bias,
  // meaning new values are set rather than continuing the inheritance
  // chain.
  pd.Properties[p.Key] = p

  // this is unnecessary but I prefer functional-style returns than
  // implicit behaviors
  return pd
}

func (pd *PropertyDrawer) Heritable() *PropertyDrawer {
  ipd := &PropertyDrawer{}
  for _, v := range pd.Properties {
    // again using functional style assignment rather than implicitly
    // modifying the value.
    ipd = ipd.Add(v)
  }

  return ipd
}

func (pd *PropertyDrawer) ValueRestrictions() ([]*Property) {
  out := make([]*Property, 0)
  for _, v := range pd.Properties {
    if v.IsValueRestriction() {
      out = append(out, v)
    }
  }

  return out
}
