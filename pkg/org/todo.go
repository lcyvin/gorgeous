package org

import (
	"fmt"

	"github.com/lcyvin/gorgeous/internal/util"
)

// TodoKeywordKind enums are used for handling the output of
// TodoSequence.GetKeywordKind(keyword). Necessary for handling state changes
// on scheduled and repeated todos for proper agenda filtering and repeat
// triggering.
type TodoKeywordKind string

const (
  TODO_KEYWORD_KIND_UNKNOWN TodoKeywordKind = ""
  TODO_KEYWORD_KIND_PROCESS = "process"
  TODO_KEYWORD_KIND_DONE = "done"
)

// Returns the string value of the TodoKeywordKind
func (tkk TodoKeywordKind) String() string {
  return string(tkk)
}

// Returns a TodoKeywordKind enum based on the passed string (unknown values
// are coerced to TODO_KEYWORD_KIND_UNKNOWN)
func (tkk TodoKeywordKind) FromString(s string) TodoKeywordKind {
  switch tkk {
  case "process":
    return TODO_KEYWORD_KIND_PROCESS
  case "done":
    return TODO_KEYWORD_KIND_DONE
  default:
    return TODO_KEYWORD_KIND_UNKNOWN
  }
}

// TodoSequence represents a single defined keyword sequence, E.G.,:
//    TODO(t) INPROGRESS(p) REVIEWING BLOCKED | DONE(d) CANCELLED(c)
//
// fast-access keys are stored in the FastAccessMap, for convenience,
// the keyword and its kind can be retrieved with the GetAccessKeyword function
type TodoSequence struct {
  // ProcessKeywords refers to any keyword occurring before a PIPE ("|")
  // character in a todo keyword sequence definition. These states do not
  // trigger a repetition.
  ProcessKeywords []string
  
  // DoneKeywords refers to any keywords occuring after a PIPE ("|") in a todo
  // keyword sequence definition. These states trigger a repetition.
  DoneKeywords    []string

  // FastAccessMap refers to any fast access keys defined for a keyword within
  // a todo keyword sequence definition (E.G., TODO(t))
  FastAccessMap   map[string]string
  
  // Kind refers to the sequence kind being defined. The valid kinds are:
  // - TODO_SEQUENCE_STATE
  // - TODO_SEQUENCE_TYPE
  //
  // TODO_SEQUENCE_TYPE exists for backwards compatibility purposes,
  // it is recommended to utilize tags to refer to these values than to set
  // them as todo keywords.
  Kind            TodoSequenceKind
}

// Returns the string value of the keyword referenced by the fast access key,
// and the TodoKeywordKind of the keyword referenced.
func (ts *TodoSequence) GetAccessKeyword(k string) (string, TodoKeywordKind) {
  var out string

  out, ok := ts.FastAccessMap[k]
  if !ok {
    out = ""
  }

  return out, ts.GetKeywordKind(out)
}

// Returns the TodoKeywordKind of the keyword referenced by k. GetKeywordKind
// does not refer to the FastAccessMap, for that use GetAccessKeyword.
func (ts *TodoSequence) GetKeywordKind(k string) TodoKeywordKind {
  if util.In(k, ts.ProcessKeywords) {
    return TODO_KEYWORD_KIND_PROCESS
  }

  if util.In(k, ts.DoneKeywords) {
    return TODO_KEYWORD_KIND_DONE
  }

  return TODO_KEYWORD_KIND_UNKNOWN
}

// Returns the fast access key defined for the passed Keyword k. E.G.,
// if a keyword is defined as "TODO(t)", GetFastAccessKey("TODO") returns "t"
func (ts *TodoSequence) GetFastAccessKey(k string) string {
  for key, v := range ts.FastAccessMap {
    if v == k {
      return key
    }
  }

  return ""
}

func (ts *TodoSequence) keywords() []string {
  return append(ts.ProcessKeywords, ts.DoneKeywords...)
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
  fastAccessKeySequence map[string]*TodoSequence
}

func (ts *TodoSettings) validate(seq []string) error {
  for _, k := range seq {
    if _, ok := ts.keywords[k]; ok {
      return NewTodoSequenceKeyCollisionError(k)
    }
  }

  return nil
}

func (ts *TodoSettings) fAdd(seq *TodoSequence) (*TodoSettings, error) {
  if seq.Kind.String() == TODO_SEQUENCE_UNKNOWN.String() {
    return nil, NewTodoSequenceKindInvalidError()
  }

  for _, s := range ts.Sequences {
    if nok, c := ts.fMapIntersects(s.FastAccessMap, seq.FastAccessMap); nok {
      return nil, NewTodoFastAccessKeyCollisionError(
        seq.GetFastAccessKey(c),
        c,
        seq.FastAccessMap[s.GetFastAccessKey(c)],
        )
    }
  }

  if nok, collision := ts.fIntersectsAny(seq, ts.Sequences); nok {
    return nil, NewTodoSequenceKeyCollisionError(collision)
  }

  if seq.Kind == TODO_SEQUENCE_STATE {
    return &TodoSettings{
      StateSequences: append(ts.StateSequences, seq),
      Sequences: append(ts.Sequences, seq),
    }, nil
  }

  return &TodoSettings{
    TypeSequences: append(ts.TypeSequences, seq),
    Sequences: append(ts.Sequences, seq),
  }, nil
}

func (ts *TodoSettings) fMapIntersects(left, right map[string]string) (bool, string) {
  for k := range right {
    if v, ok := left[k]; ok {
      return true, v
    }
  }

  return false, ""
}

func (ts *TodoSettings) fIntersectsAny(left *TodoSequence, right []*TodoSequence) (bool, string) {
  for _, v := range right {
    if ok, collision := ts.fIntersects(left, v); !ok {
      return true, collision
    }
  }

  return false, ""
}

func (ts *TodoSettings) fIntersects(left, right *TodoSequence) (bool, string) {
  for _, v := range left.keywords() {
    if ok, collision := ts.fExists(v, right); !ok {
      return true, collision
    }
  }

  return false, ""
}

func (ts *TodoSettings) fExists(q string, s *TodoSequence) (bool, string) {
  if collision := s.GetKeywordKind(q) != TODO_KEYWORD_KIND_UNKNOWN; collision {
    return false, q
  }

  return true, ""
}

func (ts *TodoSettings) GetFastAccessKey(k string) string {
  for _, seq := range ts.Sequences {
    if v := seq.GetFastAccessKey(k); v != "" {
      return v
    }
  }

  return "" 
}

// Adds a todo sequence to the 
func (ts *TodoSettings) Add(seq *TodoSequence) (*TodoSettings, error) {
  nts, err := ts.fAdd(seq)
  if err != nil {
    return nts, err
  }

  ts = nts
  return ts, err
}

type TodoSequenceKind int

const (
  TODO_SEQUENCE_UNKNOWN TodoSequenceKind = iota
  TODO_SEQUENCE_TYPE
  TODO_SEQUENCE_STATE
)

func (tsk TodoSequenceKind) String() string {
  switch int(tsk) {
  case 1:
    return "type"
  case 2:
    return "state"
  default:
    return "unknown"
  }
}

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

type TodoFastAccessKeyCollisionError struct {
  Key string
  Exist string
  New string
}

func (tfakce TodoFastAccessKeyCollisionError) Error() string {
  msg := "Fast access key %s already defined for todo keyword %s, "
  msg += "redefined by keyword %s"

  return fmt.Sprintf(msg, tfakce.Key, tfakce.Exist, tfakce.New)
}

func NewTodoFastAccessKeyCollisionError(k, e, n string) *TodoFastAccessKeyCollisionError {
  return &TodoFastAccessKeyCollisionError{
    Key: k,
    Exist: e,
    New: n,
  }
}
