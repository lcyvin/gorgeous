package org

import "fmt"

type TodoSequence struct {
  ProcessKeywords []string
  DoneKeywords    []string
  Kind            TodoSequenceKind
}

// Todo keywords can be defined as a sequence of either states, represented
// by all-caps strings containing only alphabet characters, or for backwards
// compatibility as types, represented by strings of only alphabet characters
// and starting with an uppercase Character.
// In a definition, the sequence "|" occurring in the list of states or types
// indicates that all following items represent "done" states. If none is
// present, the final element in the list is considered the "done" state.
//
// Multiple todo sequences may be defined in a file, for instance, for
// different workflows. The same key should not be defined in multiple places
// in order to allow for proper sequence cycling. E.G.:
//
//     #+TODO: BACKLOG TODO BLOCKED STARTED REVIEWING STAGED | DONE CANCELLED
//     #+TODO: IDEA SPIKE RFC | PLANNED DENIED
//     #+TODO: WAITING MEETING CALL EMAIL | HANDLED NOOP
//
// It is recommended to use tags in favor of types where relevant.
type TodoSettings struct {
  TypeSequences   []*TodoSequence
  StateSequences  []*TodoSequence
  Sequences       []*TodoSequence
  keywords        map[string]struct{}
}

func (ts *TodoSettings) validate(seq []string) error {
  for _, k := range seq {
    if _, ok := ts.keywords[k]; ok {
      return NewTodoSequenceKeyCollisionError(k)
    }
  }

  return nil
}

func (ts *TodoSettings) Add(seq *TodoSequence) (*TodoSettings, error) {
  newKeys := make(map[string]struct{}, 0)

  allSeqKeys := make([]string, 0)
  allSeqKeys = append(allSeqKeys, seq.ProcessKeywords...)
  allSeqKeys = append(allSeqKeys, seq.DoneKeywords...)

  for _, v := range allSeqKeys {
    _, ok := newKeys[v]; if ok {
      return nil, NewTodoSequenceKeyCollisionError(v)
    }
    newKeys[v] = struct{}{}
  }

  if err := ts.validate(allSeqKeys); err != nil {
    return nil, err
  }

  switch seq.Kind {
  case TODO_SEQUENCE_TYPE:
    ts.TypeSequences = append(ts.TypeSequences, seq)
  case TODO_SEQUENCE_STATE:
    ts.StateSequences = append(ts.StateSequences, seq)
  default:
    return nil, NewTodoSequenceKindInvalidError()
  }

  ts.Sequences = append(ts.Sequences, seq)

  for k, v := range newKeys {
    ts.keywords[k] = v
  }

  return ts, nil
}

type TodoSequenceKind int

const (
  TODO_SEQUENCE_TYPE TodoSequenceKind = iota
  TODO_SEQUENCE_STATE
)

type TodoSequenceKeyCollisionError struct {
  Key string
}

func (tskce TodoSequenceKeyCollisionError) Error() string {
  return fmt.Sprintf("Todo keyword collision: keyword %s is used in another sequence or occurs twice.", tskce.Key)
}

func NewTodoSequenceKeyCollisionError(k string) *TodoSequenceKeyCollisionError {
  e := &TodoSequenceKeyCollisionError{Key: k}
  return e
}

type TodoSequenceKindInvalidError struct {}

func (TodoSequenceKindInvalidError) Error() string {
  return "Unknown todo sequence kind passed to handler"
}

func NewTodoSequenceKindInvalidError() *TodoSequenceKindInvalidError {
  return &TodoSequenceKindInvalidError{}
}
