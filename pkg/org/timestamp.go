package org

import (
	"fmt"
	"time"
)

// TimestampRange is implicitly used in all planning elements, leaving EndDate
// as nil in the case that no range is defined. Helper methods are defined to
// access basic scheduling information without the need to directly use the
// time.Time values held by each Timestamp.
//
// For diary expressions (not yet implemented), the type 
//
// In an org planning element, the timestamp object can represent a range in 
// one of two ways:
//     - with a time range: 
//           <2050-01-01 Sat 00:00-02:00>
// 
//     - with a date time range:
//           <2050-01-01 Sat 00:00>--<2050-01-01 Sat 02:00>
//
// In the above examples, the actual time range is identical, however the first
// is more compact. For representing tasks and events ocurring across multiple 
// days, the latter is required. In this library, either form is held as a 
// TimestampRange. In cases where a date time range form is not used, the
// value of EndDate should be nil.
type TimestampRange struct {
  // The start date will always be present in a TimestampRange and holds the
  // given date and (when relevant) time for an item. If no time is provided,
  // this should default to a time of "00:00", and set the TimeStamp.DateOnly
  // value to "true".
  StartDate *Timestamp
  // The end date is present when a timestamp object contains a date time 
  // range. Implementors of file writers can utilize a nil check on this value
  // to determine if the format of a timestamp should be in
  // a time range or date time range format.
  EndDate *Timestamp
}

// Returns a signed int64 value representng the total number of seconds between
// the StartDate and EndDate, a boolean value representing the validity of the
// range, and an error if the start occurs after the end. The duration is still
// returned in the case of an earlier end date than start date, as a negative
// value. This is done to allow for client-side handling of the data in the 
// event passing error values is infeasible or undesirable. The boolean value 
// will only return `true` if there is no error and the start date occurs
// before the end date. If no EndDate is defined, the returned values are:
// 0, false, nil if the timestamp defined TimestampRange.StartDate is a 
// DateOnly timestamp, and *n*, true, nil if the value of Timestamp.IsRange is
// true, where *n* is the duration returned by Timestamp.Duration()
func (tr *TimestampRange) Duration() (int64, bool, error) {
  if tr.StartDate.IsRange {
    return tr.StartDate.Duration()
  }

  if tr.EndDate == nil {
    
  }

  start := tr.StartDate.Start
  end := tr.EndDate.End

  duration := int64(end.Sub(start).Seconds())

  // start can't be after end
  if start.Compare(end) == 1 {
    err := NewStartTimeAfterEndTimeError(start, end) 
    return duration, false,  err
  }

  return duration, true, nil
}

// HasDuration is a shorthand to access the bool value returned by 
// TimestampRange.Duration. It returns true only when both a start and end date
// are defined, and the duration is non-negative (start occurs before end).
func (tr *TimestampRange) HasDuration() bool {
  _, ok, _ := tr.Duration()

  return ok
}

func NewTimestampRange(start, end *Timestamp) (*TimestampRange, error) {
  if start == nil {
    if end == nil {
      return nil, NewNilTimestampsError()
    }

    return nil, NewNilStartTimeError()
  }

  tr := &TimestampRange{StartDate: start}
  if end != nil {
    tr.EndDate = end
  }

  return tr, nil
}

type Timestamp struct {
  Start time.Time
  End time.Time
  DateOnly bool
  Active bool
  IsRange bool
  Repeat *Repeat
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
func (ts *Timestamp) Within(start, end time.Time) bool {
  startsWithinWindow := ts.Start.After(start) && ts.Start.Before(end)
  endsWithinWindow := ts.End.Before(end) && ts.End.After(start)

  if startsWithinWindow || endsWithinWindow {
    return true
  }

  //TODO implement shift handling
  return false
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

type RepeatKind string

const (
  // sentinel
  REPEAT_KIND_UNKNOWN RepeatKind = ""
  // valid repeat shift markers
  REPEAT_KIND_SHIFT                 = "+"
  REPEAT_KIND_SHIFT_FUTURE_FIXED    = "++"
  REPEAT_KIND_SHIFT_FUTURE_RELATIVE = ".+"
)

func (rk RepeatKind) String() string {
  return string(rk)
}

type RepeatIntervalKind string

const (
  // sentinel
  REPEAT_INTERVAL_UNKNOWN RepeatIntervalKind = ""
  // valid interval markers
  REPEAT_INTERVAL_HOUR  = "h"
  REPEAT_INTERVAL_DAY   = "d"
  REPEAT_INTERVAL_WEEK  = "w"
  REPEAT_INTERVAL_MONTH = "m"
  REPEAT_INTERVAL_YEAR  = "y"
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
  return "Method was called with all Timestamps being nil."
}

func NewNilTimestampsError() NilTimestampsError {
  return NilTimestampsError{}
}
