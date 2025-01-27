package org

import (
	"fmt"
	"strings"
	"time"
)

// TimestampRange is implicitly used in all planning elements, leaving EndDate
// as nil in the case that no range is defined. Helper methods are defined to
// access basic scheduling information without the need to directly use the
// time.Time values held by each Timestamp.
//
// For diary expressions (not yet implemented), the type
//
// In an org planning element, the timestamp object can represent a time range
// in one of two ways:
//     - with a time range:
//           <2050-01-01 Sat 00:00-02:00>
//
//     - with a date time range:
//           <2050-01-01 Sat 00:00>--<2050-01-01 Sat 02:00>
//
// In the above examples, the actual time range is identical, however the first
// is more compact. For representing tasks and events ocurring across multiple
// days, the latter is required.
//
// Additionally, for an event which occurs on multiple days at a specific time
// period or time of day, the format used is:
//
//    <2050-01-01 Sat 00:00-02:00>--<2050-01-03 Mon 00:00-02:00>
//
// In this library, either form is held as a
// TimestampRange. In cases where a date time range form is not used, the
// value of EndDate should be nil.
type TimestampRange struct {
  // The start date will always be present in a TimestampRange and holds the
  // given date and (when relevant) time for an item. If no time is provided,
  // this should default to a time of "00:00", and set the TimeStamp.DateOnly
  // value to "true".
  StartDate *Timestamp

  // The end date is present when a timestamp object contains a date time 
  // range. Implementers of file writers can utilize a nil check on this value
  // to determine if the format of a timestamp should be in
  // a time range or date time range format.
  EndDate *Timestamp

  // Allow out-of-band behavior for client compatibility with various clients
  // that implement some org mode parsing and writing but may not exactly match
  // vanilla orgmode syntax.
  Compatibility bool
}

type TimestampRangeOpt func(*TimestampRange)

func WithCompatibility() TimestampRangeOpt {
  return func(tr *TimestampRange) {
    tr.Compatibility = true
  }
}

func NewTimestampRange(start, end *Timestamp, opts... TimestampRangeOpt) (*TimestampRange, error) {
  if start == nil {
    if end == nil {
      return nil, NewNilTimestampsError()
    }

    return nil, NewNilStartTimeError()
  }

  tr := &TimestampRange{
    StartDate: start,
    Compatibility: false,
  }

  if end == nil {
    return nil, NewNilTimestampsError()
  }

  tr.EndDate = end

  for _, opt := range opts {
    opt(tr)
  }

  return tr, nil
}

func (tr *TimestampRange) String() string {
  strs := tr.Strings()
  out := strings.Join(strs, "--")

  return out
}

func (tr *TimestampRange) Strings() []string {
  return append(tr.StartDate.Strings(), tr.EndDate.String())
}

// Returns true if the timestamps held by TimestampRange represent a date/time
// range, E.G.: <2050-01-01 Sat 00:00-02:00>--<2050-01-03 Mon 00:00-02:00>
func (tr *TimestampRange) IsRecurringRange() bool {
  return tr.StartDate.IsRange && tr.EndDate.IsRange
}

// Returns true if the StartDate.Repeat is non-nil
func (tr *TimestampRange) IsRepeating() bool {
  return tr.StartDate.Repeat != nil
}

// Returns true if TimestampRange.StartDate.Active is true
func (tr *TimestampRange) IsActive() bool {
  return tr.StartDate.Active
}

func (tr *TimestampRange) Time() (int, int, int) {
  return tr.StartDate.Time()
}

func (tr *TimestampRange) EndTime() (int, int, int) {
  endExists := tr.EndDate != nil
  endIsRange := tr.EndDate.IsRange
  // in theory this should never be true if endIsRange is false, however,
  // out-of-band behaviors from plugins and varying client implementations 
  // mean it is worth adding this sanity check.
  endIsDateOnly := tr.EndDate.DateOnly

  // when the range was defined as a date/time range
  // <YYYY-MM-DD Day HH:MM-HH:MM>--<YYYY-MM-DD Day HH:MM-HH:MM>
  if endExists && endIsRange {
    return tr.EndDate.EndTime()
  }

  if tr.Compatibility {
    // when the range was defined as a date range with fixed times on the start 
    // and end. I don't believe this is a default supported behavior in org, but
    // other clients (E.G., orgzly-style android org clients) may implement
    // nonstandard behavior, so we will attempt to handle the obvious edge cases.
    if endExists && !endIsDateOnly {
      return tr.EndDate.Time()
    }

    // As described above, to handle out-of-band implementations and possible
    // alternate syntaxes or plugin values, we will return an end-of-the-day time
    // for a date range that hypothetically looked something like:
    // <YYYY-MM-DD Day HH:MM>--<YYYY-MM-DD Day>
    if endExists && endIsDateOnly {
      return 23, 59, 59
    }
  }

  return tr.StartDate.EndTime()
}

// Flips the value of the StartDate.Active and EndDate.Active (if EndDate != nil)
func (tr *TimestampRange) ToggleActive() *TimestampRange {
  newState := !tr.StartDate.Active

  tr.StartDate.Active = newState
  if tr.EndDate != nil {
    tr.EndDate.Active = newState
  }

  return tr
}

func (tr *TimestampRange) Kind() TimestampKind {
  return TIMESTAMP_KIND_TIMESTAMP_RANGE
}

// Returns true if the event defined within the TimestampRange occurs within
// the provided window. This value does not consider active/inactive timestamp
// status, nor agenda view delay settings. These factors should be handled
// down-stream, as there are cases where one may want to query events including
// those whose visibility is normally hidden in a standard agenda view.
func (tr *TimestampRange) InWindow(start, end time.Time) bool {
  sWin := tr.StartDate.InWindow(start, end)
  eWin := tr.EndDate.InWindow(start, end)

  return sWin || eWin
}

type NewTimestampOpt func(*Timestamp)

func WithEnd(e time.Time) NewTimestampOpt {
  return func(t *Timestamp) {
    t.End = e
  }
}

func WithRepeat(r *Repeat) NewTimestampOpt {
  return func(t *Timestamp) {
    t.Repeat = r
    rawCookieStr := "%s%d%s"
    shiftStr := r.Kind.String()
    amt := r.IntervalAmount
    shiftInterval := r.Interval.String()
    t.RawCookie = fmt.Sprintf(rawCookieStr, shiftStr, amt, shiftInterval)
  }
}

func WithInactive() NewTimestampOpt {
  return func(t *Timestamp) {
    t.Active = false
  }
}

func WithDateOnly() NewTimestampOpt {
  return func(t *Timestamp) {
    t.DateOnly = true
  }
}

type Timestamp struct {
  Start time.Time
  End time.Time
  DateOnly bool
  Active bool
  IsRange bool
  Repeat *Repeat
  RawCookie string
}

func (t *Timestamp) String() string {
  out := fmt.Sprintf("%d-%d-%d %s", t.Year(), t.Month(), t.Day(), t.Weekday())
  if !t.DateOnly {
    out += fmt.Sprintf(" %02d:%02d", t.Start.Hour(), t.Start.Minute())
  }

  if !t.DateOnly && !t.End.IsZero() {
    out += fmt.Sprintf("-%02d:%02d", t.End.Hour(), t.End.Minute()) 
  }

  enclose := "<%s>"
  if !t.Active {
    enclose = "[%s]"
  }

  return fmt.Sprintf(enclose, out)
}

func (t *Timestamp) Strings() []string {
  return []string{t.String()}
}

func NewTimestamp(start time.Time, opts... NewTimestampOpt) *Timestamp {
  ts := &Timestamp{Start: start, IsRange: true, Active: true}
  
  for _, opt := range(opts) {
    opt(ts)
  }

  if ts.End.IsZero() {
    ts.IsRange = false
  }

  return ts
}

// Returns a string representation of the day of the week as expected by org
// when parsing timestamp objects. For a the full name of the day, use the
// standard go Weekday methods available from Timestamp.DateTime
func (ts *Timestamp) Weekday() string {
  return ts.Start.Weekday().String()[:3]
}

// Returns the integer value of the day of the month held by Timestamp.DateTime
func (ts *Timestamp) Day() int {
  return ts.Start.Day()
}

// Returns the integer value of the month of the year held by Timestamp.DateTime
func (ts *Timestamp) Month() int {
  return int(ts.Start.Month())
}

// Returns the integer value of the year held by Timestamp.DateTime
func (ts *Timestamp) Year() int {
  return ts.Start.Year()
}

// Returns all 0s if the timestamp is defined as DateOnly, else returns the
// hour, minute, and second values for the time held by `Timestamp.Start` as
// representable on a clock.
func (ts *Timestamp) Time() (int, int, int) {
  if ts.DateOnly {
    return 0, 0, 0
  }

  return ts.Start.Clock()
}

// returns 0 if the timestamp is a date only definition, or is not a time range 
// timestamp. Else returns the hour, minute, and second values for the time 
// held by Timestamp.End as representable on a clock.
func (ts *Timestamp) EndTime() (int, int, int) {
  if ts.DateOnly || !ts.IsRange {
    return 0, 0, 0
  }

  return ts.End.Clock()
}

// Returns the duration in seconds, validity of range definition, and any errors
// for the timestamp's Start and End values. If Timestamp.IsRange is false, 
// the return values are: 0, false, nil.
//
// If a range is present, but the start occurs after the end, the returned values
// will be: a negative duration, false, StartTimeAfterEndTimeError
// This is done to allow flexibility with implementors' handling of invalid values.
func (ts *Timestamp) Duration() (int64, bool, error) {
  if !ts.IsRange {
    return 0, false, nil
  }

  duration := int64(ts.End.Sub(ts.Start))

  if ts.Start.After(ts.End) {
    return duration, false, NewStartTimeAfterEndTimeError(ts.Start, ts.End) 
  }

  return duration, true, nil
}

// Returns true if the timestamp or one of its repetitions occurs within the
// provided window. Note that depending on the kind of repetition, the
// validity of this response is only valid until the timestamp is updated on
// a state change.
func (ts *Timestamp) InWindow(start, end time.Time) bool {
  startsWithinWindow := ts.Start.After(start) && ts.Start.Before(end)
  endsWithinWindow := ts.End.Before(end) && ts.End.After(start)

  return startsWithinWindow || endsWithinWindow
}

func (ts Timestamp) Kind() TimestampKind {
  return TIMESTAMP_KIND_TIMESTAMP
}

func (ts *Timestamp) Cookie() string {
  if ts.Repeat == nil {
    return ""
  }

  if ts.RawCookie == "" {
    return fmt.Sprintf("%s%d%s", 
      ts.Repeat.Kind.String(),
      ts.Repeat.IntervalAmount,
      ts.Repeat.Interval.String(),
      )
  }

  return ts.RawCookie
}

// Repeat holds the specific information related to a repeating task or
// deadline as needed for filtering and displaying agenda views, as well as
// handling timestamp shifts when a repeating agenda item is marked with a 
// todo done-type state.
type Repeat struct {
  // Controls how repeating items are shifted when their state is changed.
  //    - REPEAT_KIND_SHIFT: add the repeat interval to the timestamp
  //    - REPEAT_KIND_SHIFT_FUTURE_FIXED: add repeat interval to the timestamp
  //      until the date is in the future (minimum shift of 1 interval)
  //    - REPEAT_KIND_SHIFT_FUTURE_RELATIVE: add repeat interval to current
  //      time until it is in the future, does not preserve day-of-week.
  //      (minimum shift of 1 interval)
  Kind RepeatKind

  // Controls the number of times the interval is applied on repetition, E.G.:
  //    +1m == 1*REPEAT_INTERVAL_MONTH
  //    +3d == 3*REPEAT_INTERVAL_DAY
  IntervalAmount int

  // Controls the total amount of time an interval represents. One of:
  //    - REPEAT_INTERVAL_HOUR
  //    - REPEAT_INTERVAL_DAY
  //    - REPEAT_INTERVAL_WEEK
  //    - REPEAT_INTERVAL_MONTH
  //    - REPEAT_INTERVAL_YEAR
  //
  // Note: REPEAT_INTERVAL_MONTH (e.g., +1m) behavior is not consistent across
  // different org clients. The specific behavior of a 1 month shift in this
  // library has two modes, set by RelativeMonth.
  Interval RepeatIntervalKind

  // Sets the point at which an agenda item either appears in the agenda view.
  // Behaves different dependent on the kind of planning element it belongs to.
  //    - SCHEDULED: delays the appearance of the item in the agenda view by
  //      the duration set
  //    - DEADLINE: displays the warning ahead of the deadline by the duration
  //      amount. 
  AgendaWindow time.Duration

  // Because of inconsistencies with month intervals, this is an internally
  // available flag to control how an interval of +1m is set. When false,
  // the shift will always be 30 days (golang's normalized month value). With
  // RelativeMonth set to true, the value is incremented by no more than one
  // calendar month. So, an event occurring on January 30th, with a 1 month 
  // interval, would shift to February 28th or 29th the following month, rather
  // than a date in March. Similarly, an item scheduled for the 15th will
  // always be shifted to the 15th of the following month. For best accuracy,
  // it is recommended to ensure the time.Time held by Timestamp's Start and
  // End values have TimeZone values to refer to. By default, all times are
  // assumed to be UTC.
  RelativeMonth bool
}

func (r *Repeat) String() string {
  return fmt.Sprintf("%s%d%s", r.Kind.String(), r.IntervalAmount, r.Interval.String())
}

type RepeatKind string

const (
  // sentinel
  REPEAT_KIND_UNKNOWN RepeatKind = ""
  // valid repeat shift markers
  REPEAT_KIND_SHIFT                 RepeatKind = "+"
  REPEAT_KIND_SHIFT_FUTURE_FIXED    RepeatKind = "++"
  REPEAT_KIND_SHIFT_FUTURE_RELATIVE RepeatKind = ".+"
)

func (rk RepeatKind) String() string {
  return string(rk)
}

type RepeatIntervalKind string

const (
  // sentinel
  REPEAT_INTERVAL_UNKNOWN RepeatIntervalKind = ""
  // valid interval markers
  REPEAT_INTERVAL_HOUR  RepeatIntervalKind = "h"
  REPEAT_INTERVAL_DAY   RepeatIntervalKind = "d"
  REPEAT_INTERVAL_WEEK  RepeatIntervalKind = "w"
  REPEAT_INTERVAL_MONTH RepeatIntervalKind = "m"
  REPEAT_INTERVAL_YEAR  RepeatIntervalKind = "y"
)

func (rik RepeatIntervalKind) String() string {
  return string(rik)
}

type StartTimeAfterEndTimeError struct {
  startTime time.Time
  endTime time.Time
}

func NewStartTimeAfterEndTimeError(start, end time.Time) *StartTimeAfterEndTimeError {
  return &StartTimeAfterEndTimeError{
    startTime: start,
    endTime: end,
  }
}

func (saee StartTimeAfterEndTimeError) Error() string {
  msg := "Start time [%s] occurs after End time [%s]"
  startDt := saee.startTime.Format(time.DateTime)
  endDt := saee.endTime.Format(time.DateTime)
  return fmt.Sprintf(msg, startDt, endDt)
}

type NilStartTimeError struct {}

func (NilStartTimeError) Error() string {
  return "Attempt was made to instantiate a timestamp with a nil start time."
}

func NewNilStartTimeError() NilStartTimeError {
  return NilStartTimeError{}
}

type NilTimestampsError struct {}

func (NilTimestampsError) Error() string {
  return "Method was called with one or all Timestamps being nil."
}

func NewNilTimestampsError() NilTimestampsError {
  return NilTimestampsError{}
}
