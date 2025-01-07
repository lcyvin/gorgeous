package api

type NodeTree interface {
  // Returns the parent nodeTree of this subtree, or nil if there is no parent.
  Parent() NodeTree

  // Returns a Node interface for the Node at this level
  Node() Node

  // Returns a slice of NodeTree subtrees owned by this node tree
  Children() []NodeTree
}
