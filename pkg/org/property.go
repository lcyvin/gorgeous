package org

import (
	"fmt"
	"strings"
)

type Property struct {
  // The key by which this property's value can be referenced.
  Key         string
  // The value held by this key, or whitespace-separated list
  // of values that are allowed for a corresponding Property
  // having the same name, less the suffix "_All"
  Value       string
}

func (p Property) Kind() ElementKind {
  return ELEMENT_NODE_PROPERTY
}

func (p Property) IsGreaterElement() bool {
  return false
}

// Org syntax defines that properties ending with the suffix `_All`
// should be heritable by nodes down the tree, enforcing the values
// listed (whitespace-separated) as the allowed values for the 
// corresponding property (sans _All).
func (p *Property) IsValueRestriction() bool {
  if p.Key[len(p.Key)-4:] == "_All" {
    return true
  }

  return false
}

// When a property's value is restricted by heritable `_All`-suffixed
// property definitions, it applies to the corresponding property whose
// name matchse the restriction key less the suffix.
func (p *Property) RestrictionKey() string {
  if p.IsValueRestriction() {
    return p.Key[len(p.Key)-4:]
  }

  return p.Key
}

// Property value restrictions set by properties with keys like `Foo_All`
// hold a whitespace-separated list of valid values for the corresponding
// non-suffixed property, with quotation marks being used to group values
// that contain whitespace within themselves.
func (p *Property) RestrictionValues() []string {
  out := make([]string, 0)
  if !p.IsValueRestriction() {
    return out
  }

  brktOpen := false
  val := ""
  for _, c := range p.Value {
    if c == '"' {
      brktOpen = !brktOpen
      continue
    }

    if c == ' ' && !brktOpen {
      out = append(out, val)
      val = ""
      continue
    }

    val += string(c)
  }

  return out
}

// Tests the value of the provided property against the list of valid values
// defined within the value restriction property. Returns nil if validation
// passes, or one of:
//   - InvalidPropertyValueError
//   - NotValueRestrictionPropertyError
// The former is returned when the value in `prop` does not match any of the
// listed values. The latter is returned if this function is called by a
// property that does not set value restrictions. Type assertions may be made
// on the returned error in order to handle each contingency. 
func (p *Property) Validate(prop *Property) error {
  if !p.IsValueRestriction() {
    return NewNotValueRestrictionPropertyError(p)
  }

  for _, v := range p.RestrictionValues() {
    if prop.Value == v {
      return nil
    }
  }

  return NewInvalidPropertyValueError(p, prop)
}

type NotValueRestictionPropertyError struct {
  Property string
}

func NewNotValueRestrictionPropertyError(p *Property) *NotValueRestictionPropertyError {
  nvrpe := &NotValueRestictionPropertyError {
    Property: p.Key,
  }

  return nvrpe
}

func (nvrpe NotValueRestictionPropertyError) Error() string {
  return fmt.Sprintf(
    "Property %s does not implement value restrictions or is misnamed.",
    nvrpe.Property,
    )
}

type InvalidPropertyValueError struct {
  Property          string
  PropertyValue     string
  Restrictor        string
  RestrictorValues  []string
}

func NewInvalidPropertyValueError(p, r *Property) *InvalidPropertyValueError {
  ipve := &InvalidPropertyValueError{
    Property: p.Key,
    PropertyValue: p.Value,
    Restrictor: r.Key,
    RestrictorValues: r.RestrictionValues(),
  }

  return ipve
}

func (ipve InvalidPropertyValueError) Error() string {
  msg := fmt.Sprintf(
    "Invalid property value for property %s: %s.",
    ipve.Property,
    ipve.PropertyValue,
    )

  msg += fmt.Sprintf("%s restricts values to: %s",
    ipve.Restrictor,
    strings.Join(ipve.RestrictorValues, ", "),
    )

  return msg
}
