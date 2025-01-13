package cal

import (
	"io"
	"os"
  "fmt"
	"time"

	"github.com/apognu/gocal"
	"github.com/lcyvin/gorgeous/pkg/org"
)

type Cal struct {
  Document *org.Document
  parsers []*EventParser
}

func NewCal() *Cal {
  c := &Cal{}
  c.Document = org.New()
  c.Document.Virtual = true

  return c
}

func (c *Cal) ImportFile(p string) (*Cal, error) {
  ep, err := NewEventParser(WithFile(p))
  if err != nil {
    return nil, err
  }

  c.parsers = append(c.parsers, ep)

  return c, nil
}

type DateCounterFormat int

const (
  PREPEND DateCounterFormat = iota
  APPEND
)

type EventParser struct {
  Data io.Reader
  Window *Window
  Config *EventParserConfig
  parser *gocal.Gocal
}

type EventParserConfig struct {
  AddDateCounter bool
  DateCounterFmt DateCounterFormat
}

type Window struct {
  Start time.Time
  End   time.Time
}

type EventParserOpt func(*EventParser) error

func WithFile(path string) EventParserOpt {
  return func(c *EventParser) error {
    f, err := os.Open(path)
    if err != nil {
      return err
    }

    c.Data = f
    return nil
  }
}

func WithData(r io.Reader) EventParserOpt {
  return func(c *EventParser) error {
    c.Data = r
    return nil
  }
}

func WithWindow(start, end time.Time) EventParserOpt {
  return func(ep *EventParser) error {
    ep.Window = &Window{
      Start: start,
      End: end,
    }
    return nil
  }
}

func NewEventParser(opts... EventParserOpt) (*EventParser, error) {
  ce := &EventParser{}

  for _, opt := range opts {
    if err := opt(ce); err != nil {
      return nil, err
    }
  }

  if ce.Data != nil {
    ce.parser = gocal.NewParser(ce.Data)
  }

  return ce, nil
}

func (ep *EventParser) NewWindow(start, end time.Time) *EventParser {
  newEvent, _ := NewEventParser(WithData(ep.Data), WithWindow(start, end))
  return newEvent
}

func (ep *EventParser) Nodes() ([]*org.Node, error) {
  nodes := make([]*org.Node, 0)

  window := ep.Window
  if ep.Window == nil {
    window = NewDefaultWindow()
  }

  evt, _ := NewEventParser(WithData(ep.Data))
  evt.Window = window

  if err := evt.parser.Parse(); err != nil {
    return nodes, err
  }

  for _, item := range evt.parser.Events {
    node, err := ep.newNodes(item)
    if err != nil {
      return []*org.Node{}, err
    }

    nodes = append(nodes, node...)
  }

  return nodes, nil
}

func (ep *EventParser) newNodes(ce gocal.Event) ([]*org.Node, error) {
  n := make([]*org.Node, 0)
  base := org.Node{}
  baseHdg := org.Heading{
    Text: ce.Summary,
    Level: 1,
  }

  baseSct := org.Section{
    Elements: []org.Element{
      &org.Paragraph{Raw: ce.Description},
    },
  }

  base.Section = &baseSct

  timestamps := timestampSet(ce.Start, ce.End)
  days := len(timestamps)
  if days == 1 {
    e := base
    hdg := baseHdg
    hdg.Node = &e
    hdg.Planning = &org.Planning{
      Kind: org.PLANNING_EVENT,
      TimestampRangeOrSexp: timestamps[0],
    }
    e.Heading = &hdg

    return append(n, &e), nil
  }

  for idx, t := range timestamps {
    e := base
    hdg := baseHdg
    hdg.Node = &e
    if ep.Config.AddDateCounter {
      switch ep.Config.DateCounterFmt {
      case PREPEND:
        hdg.Text = fmt.Sprintf("(%d/%d) ", idx+1, days) + hdg.Text
      case APPEND:
        hdg.Text = hdg.Text + fmt.Sprintf(" (%d/%d)", idx+1, days)
      }
    }

    hdg.Planning = &org.Planning{
      Kind: org.PLANNING_EVENT,
      TimestampRangeOrSexp: t,
    }

    e.Heading = &hdg
    n = append(n, &e)
  }

  return n, nil
}

func NewDefaultWindow() *Window {
  today := time.Now().Local()
  start := time.Date(
    today.Year(), today.Month(), today.Day(),
    0, 0, 0, 0,
    today.Location())

  w := &Window{
    Start: start,
    End: start.Add(24*time.Hour),
  }

  return w
}

func timestampSet(start, end *time.Time) []*org.Timestamp {
  out := make([]*org.Timestamp, 0)

  if end == nil {
    return append(out, org.NewTimestamp(*start, org.WithDateOnly()))
  }

  if end.Before(tomorrow(*start)) {
    return append(out, org.NewTimestamp(*start, org.WithEnd(*end)))
  }

  days := 0
  for iter := *start; iter.Before(*end); iter = tomorrow(iter) {
    days += 1
  }

  curStart := *start
  for day := 0; day < days; day++ {
    e := *end
    if day+1 != days {
      e = endOfDay(curStart)
    }

    out = append(out, org.NewTimestamp(curStart, org.WithEnd(e)))
  }

  return out
}

// returns a time.Time with the clock set to 23:59:59 and 999[...] nsec for the
// given day
func endOfDay(d time.Time) time.Time {
  eod := tomorrow(d).Add(-1)

  return eod
}

// returns a time.Time one calendary day after the passed time.Time,
// shifted to exactly midnight
func tomorrow(d time.Time) time.Time {
  t := d.AddDate(0, 0, 1)
  tm := time.Date(
    t.Year(), t.Month(), t.Day(), 
    0, 0, 0, 0, 
    d.Location())

  return tm
}
