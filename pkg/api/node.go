package api

type Node interface {
  Heading() Heading
  Section() Section
  Parent()  Node
}
