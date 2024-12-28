package org

type Drawer struct {
}

func (d Drawer) Kind() ElementKind {
  return ELEMENT_DRAWER
}

func (d Drawer) IsGreaterElement() bool {
  return true
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
    if v.Inheritable {
      // again using functional style assignment rather than implicitly
      // modifying the value.
      ipd = ipd.Add(v)
    }
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
