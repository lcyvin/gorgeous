package org

// Headings represent logical separations between sections
// within a document and are components of an outline (or node)
// tree within a document. A heading occuring at a "lower" (higher
// int value) level than the preceeding heading are considered to
// be children of the preceeding heading. A child heading will always
// inherit tags from its parents as well as any set with the `FILETAGS`
// property. Other heritable items (e.g., properties) are handled at the
// node level.
type Heading struct {
  Text        string
  Priority    HeadingPriority
  IsComment   bool
  Tags        []string
  Level       int
  Node        *Node
}


