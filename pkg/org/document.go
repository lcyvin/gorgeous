package org

type Document struct {
  // A list containing each node in order of parsing where the int value
  // refers to the list index of the node's parent. For a first-level Node,
  // this will always be `0`. For Any node at a lower level than 1, it will
  // be the most recent node occurring at a higher level. Id est, in a Document
  // wherein the first heading occurs at a level of `2`, the most recent
  // higher-level node would still be the zero-th node.
  Nodes          []map[int]*Node
  BufferSettings *BufferSettings
}

// Instantiate a new blank document with no contents
func New() *Document {
  d := &Document{}
  d.Nodes = []map[int]*Node{
    {0: &Node{}},
  }

  return d
}
