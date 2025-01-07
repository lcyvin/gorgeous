package api

type Timing interface {
  // Start returns a Timestamp interface representing the first or only
  // timestamp in a planning element. A Timestamp may represent a single point
  // in time which an event or task begines, or a timestamp range defining the
  // start and end times of the task or event's duration.
  Start()           Timestamp

  // If the timing element is a date range or date time range, End retruns the
  // Last day on which the event or task occurs as currently scheduled. E.G.:
  //    - <2050-01-01 Wed 00:00-02:00>--<2050-01-01 Wed 00:00-02:00>
  //      an event which spans only one day for a specific range of time. This
  //      construction is overly long, but is valid.
  //      This is equivalent in function to a simple Time Range.
  //    - <2050-01-01 Wed>--<2050-01-03 Fri>
  //      an event which spans from the first of January until the third of
  //      January. This is a Date Range
  //    - <2050-01-01 Wed 00:00-02:00>--<2050-01-03 Fri 00:00-02:00>
  //      an event which repeats from the 1st to the 3rd, from midnight until 2
  //      each day. This is a Date Time Range.
  // 
  // The End Timestamp is only useful if IsDateRange() or IsDateTimeRange() are
  // true.
  End()             Timestamp

  // If the timestamp is defined as being "active" (enclosed in angle brackets,
  // in an actual org document), Active() returns true
  Active()          bool

  // If the timing element has a repeat cookie set, this returns true
  IsRepeat()        bool

  // If the timing element is a DateRange or DateTimeRange, this returns true
  IsDateRange()     bool

  // If the timing element is a DateTimeRange, this returns true
  IsDateTimeRange() bool

  // If the timing element has a Repeat cookie set, this returns a new Timing
  // element based on the configured shifting behavior set by the cookie. 
  Shift()           Timing
}

type Timestamp interface {
  Date()      [3]int
  StartTime() [3]int
  EndTime()   [3]int
  DateOnly()  bool
  IsRange()   bool
  Cookie()    string
}
