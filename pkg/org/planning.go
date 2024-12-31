package org

import "time"

// Planning elements directly follow headlines with no preceeding newlines
// and are marked by a keyword of "SCHEDULED", "DEADLINE", or "CLOSED", a
// colon, and a timestamp object.
type Planning struct {
  Kind PlanningKind
  Timestamp *TimestampRangeOrSexp
}

type PlanningKind string

const (
  // if we do not see a keyword before a timestamp, it is a simple event
  PLANNING_EVENT PlanningKind = ""
  // We expect to see these keywords depending on which keyword preceedes
  // the timestamp
  PLANNING_SCHEDULED = "SCHEDULED"
  PLANNING_DEADLINE = "DEADLINE"
)

type TimestampRangeOrSexp interface {
  // Should return out of either TIMESTAMP_KIND_TIMESTAMP or TIMESTAMP_KIND_SEXP
  Kind() TimestampKind
  // If defined as a range between two dates without a per-day time definition,
  // the total duration of the time range from the starting time on the starting
  // date, to the ending time on the ending date (defaults 00:00 and 23:59)
  // respectively.
  Duration() (int64, bool, error)
  // Returns a list of all timestamp instances which are valid within the
  // supplied date range. When using the TimestampRange data structure, this
  // will be TimestampRange.StartDate if TimestampRange.EndDate falls outside
  // the window, or both if the window encapsulates the whole range.
  //
  // If the timestamp(s) defined within a TimestampRange use repeat directives,
  // the return will be a single Timestamp. Extrapolation of that timestamp to
  // an agenda-like set of generated Timestamps can then be performed with that
  // return value.
  Instances(start, end time.Time) []*Timestamp
  // Returns the time of day that the timestamp defines as the start, if a start
  // time is present. Else, returns 0, 0, 0
  Time() (int, int, int)
  // Returns the time of day that the timestamp defines as the end of the time
  // range, if a time range is set. Else, returns 0, 0, 0.
  EndTime() (int, int, int)
}

type TimestampKind string

const (
  // sentinel
  TIMESTAMP_KIND_UNKNOWN TimestampKind = ""
  TIMESTAMP_KIND_TIMESTAMP = "timestamp"
  TIMESTAMP_KIND_SEXP = "sexp"
)
