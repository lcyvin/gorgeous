package org

import (
	"time"
)

var DefaultRepeatConfig RepeatConfig = RepeatConfig{
  ClampToEndOfMonth: false,
  ShiftByDays: false,
  FixedDate: true,
  Location: time.Now().Location(),
}

type RepeatConfig struct {
  // ClampToEndOfMonth modifies the behavior of the Repeater to always set
  // shifted dates to be no later than the last day of the month following the
  // month referenced in the timstamp's current state. Its behavior varies 
  // depending on the values of other RepeatConfig struct items. The possible
  // behaviors are:
  //
  // - ClampToEndOfMonth only is true: date is always shifted to the last day
  //   of the target month.
  //
  // - ShiftByDays is true: date is shifted by 30 days, unless said shift
  //   results in the new timestamp "Jumping" a month. E.G., with this setting
  //   a one month repeat on a date of January 31st would normally shift to 
  //   March 2nd or 3rd, but instead will clamp to the last day of the month,
  //   being Feb. 28th or 29th. 
  //
  // - FixedDate is true: Date will be shifted to the following month, and if
  //   that date does not exist in the month, the resulting shift will be the 
  //   last possible day of the month. Note that Going forward, this date 
  //   becomes the new reference point for all future shifts. 
  //
  // Note: ShiftByDays and FixedDate can not both be true.
  ClampToEndOfMonth bool
  
  // Unless combined with ClampToEndOfMonth, shift by days will always result
  // in the timestamp being shifted by exactly 30 days (the behavior of 
  // time.AddDate)
  ShiftByDays bool

  // Unless combined with ClampToEndOfMonth, shift to the next valid occurence
  // of the given date. This will result in skipping months where the date does
  // not exist, E.G., something which should always repeat on the 29th will 
  // only occur in february if it is a leap year.
  FixedDate bool

  Location *time.Location
}

// RepeatStamp is a meta struct that implements the api.Repeater and
// api.RepeatStamp interfaces, providing handling for repeat directives set in
// a timestamp based on the behavior defined in RepeatConfig. If no
// RepeatConfig is set when a shift, window, etc. operation is called,
// DefaultRepeatConfig is used.
type RepeatStamp struct {
  Timestamp
  RepeatConfig RepeatConfig
}

func NewRepeatStampFromTimestamp(ts *Timestamp, cfg RepeatConfig) *RepeatStamp {
  return &RepeatStamp{
    Timestamp: *ts,
    RepeatConfig: DefaultRepeatConfig,
  }
}

func NewRepeatStamp(start time.Time, cfg RepeatConfig, opts... NewTimestampOpt) *RepeatStamp {
  rs := &RepeatStamp{
    Timestamp: *NewTimestamp(start, opts...),
    RepeatConfig: cfg,
  }

  return rs
}

func (rs *RepeatStamp) InWindow(start, end time.Time) bool {
  
  return false
}

// Implements the Shift() function as required by api.Repeater
// Shifts the timestamp by one interval, based on the configured behavior in
// RepeatStamp.RepeatConfig. Returns a new pointer to a RepeatStamp object.
// This function also considers the cookie held by the underlying timestamp,
// if you wish to perform a shift other than the one specified by the cookie
// (E.G., the cookie is ++7d, but you wish to shift by one week from the
// timestamp's held date), use Shiftn()
func (rs *RepeatStamp) Shift(t time.Time) *RepeatStamp {
  switch rs.Repeat.Kind {
  case REPEAT_KIND_SHIFT:
    return rs.Shiftn(rs.Repeat.IntervalAmount)
  case REPEAT_KIND_SHIFT_FUTURE_FIXED:
    if t.IsZero() {
      t = time.Now()
    }
    return rs.ShiftUntilAfter(t)
  // I'm not certain you can have hourly shifts with a time range, but let's
  // pretend you can and preserve the duration of the range.
  case REPEAT_KIND_SHIFT_FUTURE_RELATIVE:
    nrs := *rs
    now := time.Now()
    if !t.IsZero() {
      now = t
    }

    duration := time.Duration(0)

    if rs.IsRange {
      duration = rs.End.Sub(rs.Start)
    }

    nrs.Start = now
    
    if !rs.End.IsZero() {
      nrs.End = nrs.Start.Add(duration)
    }

    return nrs.Shiftn(1)
  default:
    return nil
  }
}

func (rs *RepeatStamp) Shiftn(i int) *RepeatStamp {
  amt := rs.Repeat.IntervalAmount
  switch rs.Repeat.Interval {
  case REPEAT_INTERVAL_HOUR:
    return rs.shiftByHours(amt*i)
  case REPEAT_INTERVAL_DAY:
    return rs.shiftByDays(amt*i)
  case REPEAT_INTERVAL_WEEK:
    return rs.shiftByWeeks(amt*i)
  case REPEAT_INTERVAL_MONTH:
    o, err := rs.shiftByMonths(amt*i)
    if err != nil {
      panic(err)
    }

    return o
  case REPEAT_INTERVAL_YEAR:
    return rs.shiftByYears(amt)
  default:
    return nil
  }
}

func (rs *RepeatStamp) ShiftUntil(t time.Time) *RepeatStamp {
  nrs := *rs
  one := rs.Shiftn(1)
  next := one.Shiftn(1)
  delta := one.Start.Sub(rs.Start)

  if t.Before(one.Start) && t.Before(one.End) {
    return &nrs
  }

  if t.After(one.Start) && t.Before(next.Start) {
    return one
  }

  // this will not always land, but should get us in the ballpark
  window := t.Sub(rs.Start)
  shiftN := int(window.Hours()/delta.Hours())
  
  try := rs.Shiftn(shiftN)
  if t.Before(try.Start) && t.Before(try.End) {
    halfShift := int(shiftN/2)
    return rs.Shiftn(halfShift).ShiftUntil(t)   
  }

  confirmTry := try.Shiftn(1)
  if t.Before(confirmTry.Start) && t.Before(confirmTry.End) {
    return confirmTry.ShiftUntil(t)
  }

  return try
}

func (rs *RepeatStamp) ShiftUntilAfter(t time.Time) *RepeatStamp {
  before := rs.ShiftUntil(t)
  after := before.Shiftn(1)

  if t.After(after.Start) && t.After(after.End) {
    after.ShiftUntilAfter(t)
  }

  return after
}

func (rs *RepeatStamp) shiftByHours(i int) *RepeatStamp {
  // notice we don't set anything relative to the current time, as these funcs
  // are intended to handle the base shift operation. The relative/fixed/etc.
  // shifting behavior should be implemented by the parent funcs.
  nrs := *rs

  // if the timestamp was originally defined without a specific time, but has
  // an hourly repeat, we need to assume it is incrementing from 00:00 on that
  // calendar day.
  if rs.DateOnly {
    y, m, d := rs.Start.Date()
    rs.Start = time.Date(y, m, d, 0, 0, 0, 0, rs.RepeatConfig.Location)
    nrs.DateOnly = false
  }

  start := rs.Start.Add(time.Duration(i)*time.Hour)
  nrs.Start = start

  if !rs.End.IsZero() {
      nrs.End = rs.End.Add(time.Duration(i)*time.Hour)
  }

  return &nrs
}

func (rs *RepeatStamp) shiftByDays(i int) *RepeatStamp {
  nrs := *rs
  
  nrs.Start = rs.Start.AddDate(0, 0, i)

  if !rs.End.IsZero() {
    nrs.End = rs.End.AddDate(0, 0, i)
  }

  return &nrs
}

func (rs *RepeatStamp) shiftByWeeks(i int) *RepeatStamp {
  nrs := *rs
  
  nrs.Start = rs.Start.AddDate(0, 0, i*7)

  if !rs.End.IsZero() {
    nrs.End = rs.End.AddDate(0, 0, i*7)
  }

  return &nrs
}

// TODO make this way way cleaner than just calling single shifts over and
// over. For most cases it shouldn't matter much, but it's not ideal.
func (rs *RepeatStamp) shiftByMonths(i int) (*RepeatStamp, error) {
  nrs := *rs

  irs := *rs
  for iter := 0; iter < i; iter++ {
    next, err := irs.shiftByMonth()
    if err != nil {
      return nil, err
    }
    irs = *next
  }

  nrs.Start = irs.Start
  if !irs.End.IsZero() {
    nrs.End = irs.End
  }

  return &nrs, nil
}

func (rs *RepeatStamp) shiftByMonth() (*RepeatStamp, error) {
  nrs := *rs

  if rs.RepeatConfig.ShiftByDays && rs.RepeatConfig.FixedDate {
    return nil, NewInvalidRepeatConfigError()
  }

  if rs.RepeatConfig.ClampToEndOfMonth {
    if !rs.RepeatConfig.ShiftByDays && !rs.RepeatConfig.FixedDate {
      loc := rs.RepeatConfig.Location
      sFom := lastDayOfMonth(rs.Start, loc).AddDate(0, 0, 1)
      sEom := lastDayOfMonth(sFom, loc)

      nrs.Start = sEom

      if !rs.End.IsZero() {
        eFom := lastDayOfMonth(rs.End, loc).AddDate(0, 0, 1)
        eEom := lastDayOfMonth(eFom, loc)

        nrs.End = eEom
      }

      return &nrs, nil
    }

    if rs.RepeatConfig.ShiftByDays {
      start := rs.Start.AddDate(0, 0, 30)
      sm := int(rs.Start.Month())

      if sm == 12 {
        sm = 0
      }

      if int(start.Month()) - sm > 1 {
        nm := lastDayOfMonth(rs.Start, rs.RepeatConfig.Location).AddDate(0, 0, 1)
        start = lastDayOfMonth(nm, rs.RepeatConfig.Location)
      }
      
      nrs.Start = start

      if !rs.End.IsZero() {
        year, month, day := start.Date()
        hour, minute, _ := rs.End.Clock()
        nrs.End = time.Date(
          year, month, day,
          hour, minute, 0, 0,
          rs.RepeatConfig.Location,
          )
      }

      return &nrs, nil
    }

    if rs.RepeatConfig.FixedDate {
      nextMonth := lastDayOfMonth(rs.Start, rs.RepeatConfig.Location).AddDate(0, 0, 1)
      monthEnd := lastDayOfMonth(nextMonth, rs.RepeatConfig.Location)
      var start time.Time
      year, month, _ := nextMonth.Date()
      hour, minute, _ := rs.Start.Clock()
      if monthEnd.Day() < rs.Day() {
        start = time.Date(
          year, month, monthEnd.Day(),
          hour, minute, 0, 0,
          rs.RepeatConfig.Location,
          )
      } else {
        start = time.Date(
          year, month, rs.Day(),
          hour, minute, 0, 0,
          rs.RepeatConfig.Location,
          )
      }

      nrs.Start = start
      if !rs.End.IsZero() {
        year, month, day := nrs.Start.Date()
        hour, minute, _ := rs.End.Clock()
        nrs.End = time.Date(
          year, month, day,
          hour, minute, 0, 0,
          rs.RepeatConfig.Location,
          )
      }

      return &nrs, nil
    }
  }

  if rs.RepeatConfig.ShiftByDays {
    nrs.Start = rs.Start.AddDate(0, 0, 30)
    if !nrs.End.IsZero() {
      nrs.End = rs.End.AddDate(0, 0, 30)
    }
  }

  if rs.RepeatConfig.FixedDate {
    nextMonth := lastDayOfMonth(rs.Start, rs.RepeatConfig.Location).AddDate(0, 0, 1)
    lastOfMonth := lastDayOfMonth(nextMonth, rs.RepeatConfig.Location)

    if lastOfMonth.Day() < rs.Day() {
      nextMonth = nextMonth.AddDate(0, 1, 0)
    }

    year, month, _ := nextMonth.Date()
    hour, minute, _ := rs.Start.Clock()
    nrs.Start = time.Date(
      year, month, rs.Day(),
      hour, minute, 0, 0,
      rs.RepeatConfig.Location,
      )

    if !rs.End.IsZero() {
      hour, minute, _ = rs.End.Clock()
      nrs.End = time.Date(
        year, month, rs.Day(),
        hour, minute, 0, 0,
        rs.RepeatConfig.Location,
        )
    }
  }

  return &nrs, nil
}

func (rs *RepeatStamp) shiftByYears(i int) *RepeatStamp {
  nrs := *rs
  nrs.Start = rs.Start.AddDate(i, 0, 0)

  if !rs.End.IsZero() {
    nrs.End = rs.End.AddDate(i, 0, 0)
  }

  return &nrs
}

func firstDayOfMonth(t time.Time, l *time.Location) time.Time {
  t = t.In(l) 
  year, month, _ := t.Date()
  hour, minute, _ := t.Clock()

  return time.Date(
    year, month, 1,
    hour, minute, 0, 0,
    t.Location(),
    )
}

// due to the design of org, we assume all dates to already be in local time
// but we will esnure this with time.In()
func lastDayOfMonth(t time.Time, l *time.Location) time.Time {
  t = t.In(l)
  
  fom := firstDayOfMonth(t, l)

  return fom.AddDate(0, 1, -1)
}

type InvalidRepeatConfigError struct {}

func (irce InvalidRepeatConfigError) Error() string {
  return "Invalid repeat config set: cannot have ShiftByDays AND FixedDate both true."
}

func NewInvalidRepeatConfigError() *InvalidRepeatConfigError {
  return &InvalidRepeatConfigError{}
}
