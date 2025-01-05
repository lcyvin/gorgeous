package api

import (
  "time"
)

type RepeatStamp interface {
  // Return the starting time of the timestamp object having a repeater, if one
  // is present.
  Start() time.Time
  // Return the ending time of the timestamp object having a repeater, if one is
  // present.
  End() time.Time
  // For implementors, return a custom Kind for the RepeatStamp for type assertion
  Kind() interface{}
  // string representation of the repeat directive (aka cookie) as would appear
  // within an org document (E.G., +1m)
  Cookie() string
  // Return true if this refers to an active timestamp, in orgmode terms
  Active() bool
  // Return true if the RepeatStamp would occur within the passed window, based
  // on the repeat cookie and implemented repetition logic.
  InWindow(start, end time.Time) bool
}

// May be implemented internal to a RepeatStamp or externally as a handler
type Repeater interface {
  // Perform a single shift on a RepeatStamp based on the cookie it holds.
  Shift() RepeatStamp
  // Perform <n> shifts on a RepeatStamp based on the cookie it holds.
  Shiftn(i int) RepeatStamp
  // Perform as many shifts on a RepeatStamp as needed until it is less than or
  // equal to the time passed to t
  ShiftUntil(t time.Time) RepeatStamp
  // Perform as many shifts on a RepeatStamp as needed to be in the future of
  // the time passed to t
  ShiftUntilAfter(t time.Time) RepeatStamp
}
