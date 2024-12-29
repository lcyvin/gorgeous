package org

// A node represents a discrete collection of elements on the tree consisting
// of, at the very least, a heading element, any elements within the
// section owned by the heading, and ends at the next occurrance of a heading.
type Node struct {
  // Index refers to the node's own position within the list of node objects
  // and is used to simplify top-down walking operations.
  Index       int
  // Heading holds the information defined by the heading's structure,
  // including the level of the heading which is used to determine the node's
  // relatives during a tree-walking operation.
  Heading     *Heading
  // Section refers to the collection of elements following the heading within
  // the node, excepting another heading.
  Section     *Section
  // Properties refers to any properties defined in a property drawer
  // immediately following a heading element, or inherited from parent nodes
  // or document-level properties.
  Properties  []Property
  // Document refers to the document (and thus root tree) which contains this
  // node. For implementations of org features which refer to multiple
  // documents at once, it is necessary to maintain a reference to the specific
  // location of any given node in order to allow for re-filing, sorting, etc.
  Document    *Document
}
