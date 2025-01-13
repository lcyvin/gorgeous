package org

// A section belongs to a heading (excepting the zero-th
// section occuring before the first headline) and contains
// all elements except another Heading
type Section struct {
  // The heading that owns this section. The zero-th Section
  // will have a `nil` heading. No other section should be
  // able to have a `nil` heading.
  Heading   *Heading
  // Collection of all elements within the section. In order
  // to handle the various element types, type inference will 
  // be necessary based on the `Element.Kind()` output.
  Elements  []Element
  // Raw byte array containing the data used to parse this section,
  // when built using a parser. This is useful for debugging
  // purposes when implementing parser rules or when parsing mishandles
  // an element within the section. When the document is constructed
  // without parsing, this can be blank.
  Raw       []byte
}
