package api

type Heading interface {
  Level()     int
  Text()      string
  Priority()  Priority
  Tags()      []string
  Planning()  Planning
  IsComment() bool
  Node()      Node
}
