package storage

import (
	"errors"
	"fmt"
	"io"
	"sync"
	"time"
)

const (
	defatuldEventDurationInSeconds = 300
)

var (
	ErrNotFound   = errors.New("not found")
	ErrDateBusy   = errors.New("this time is already taken")
	ErrWrongTime  = errors.New("the time is incorrect")
	ErrEmptyTopic = errors.New("event topic not specified")
)

// Storage - interface for interacting with calendars
type Storage interface {
	AddEvent(timeStart, timeEnd int64, topic string) (eventID int, err error)
	DeleteEvent(timeStart int64) (err error)
	DeleteEvents(eventIDs ...int) (deletedEventIDs []int, err error)
	EditEvent(oldTimeStart, timeStart, timeEnd int64, topic string) (err error)
	EditEventByID(eventID int, timeStart, timeEnd int64, topic string) (err error)
	ShowEvents(io.Writer)
	ShowUpcomingEvents(io.Writer)
}

// Calendar - an object for working with a set of events
type Calendar struct {
	Events              map[int]*Event
	EventsByStartTime   map[int64]*Event
	eventsOrderedByTime []*Event
	mu                  *sync.Mutex
	lastEventID         int
}

// Event - record with description, start and end date
type Event struct {
	ID                 int
	TimeStart, TimeEnd time.Time
	Topic              string
}

func CreateCalendar() *Calendar {
	return &Calendar{
		Events:              make(map[int]*Event, 256),
		EventsByStartTime:   make(map[int64]*Event, 256),
		eventsOrderedByTime: make([]*Event, 0),
		mu:                  &sync.Mutex{},
		lastEventID:         0,
	}
}

func (e *Event) CreateCopy() (eCopy *Event) {
	eCopy = &Event{
		ID:        e.ID,
		TimeStart: e.TimeStart,
		TimeEnd:   e.TimeEnd,
		Topic:     e.Topic,
	}
	return eCopy
}

func InsertNewEventInTimeLine(event *Event, timeLine []*Event) (newTimeLine []*Event, err error) {
	if len(timeLine) == 0 {
		timeLine = append(timeLine, event)
		return timeLine, nil
	}
	for i, e := range timeLine {
		switch {
		case e.TimeEnd.Before(event.TimeStart):
			continue
		case e.TimeStart.After(event.TimeEnd) || e.TimeStart.Equal(event.TimeEnd):
			if (i == 0) || (timeLine[i-1].TimeEnd.Before(event.TimeStart) || timeLine[i-1].TimeEnd.Equal(event.TimeStart)) {
				timeLine = append(timeLine[:i], append([]*Event{event}, timeLine[i:]...)...)
				return timeLine, nil
			} else {
				return nil, ErrDateBusy
			}
		}
	}
	if timeLine[len(timeLine)-1].TimeEnd.Before(event.TimeStart) || timeLine[len(timeLine)-1].TimeEnd.Equal(event.TimeStart) {
		timeLine = append(timeLine, event)
		return timeLine, nil
	}
	return nil, ErrDateBusy
}

func DeleteEventFromTimeLine(id int, timeLine []*Event) (newTimeLine []*Event, err error) {
	for i, e := range timeLine {
		if e.ID == id {
			if i < len(timeLine)-1 {
				copy(timeLine[i:], timeLine[i+1:])
			}
			timeLine = timeLine[:len(timeLine)-1]
			return timeLine, nil
		}
	}
	return timeLine, ErrNotFound
}

func ChangeEventInTimeLine(oldE, newE *Event, timeLine []*Event) (newTimeLine []*Event, err error) {
	newTimeLine = make([]*Event, len(timeLine))
	copy(newTimeLine, timeLine)
	if newTimeLine, err = DeleteEventFromTimeLine(oldE.ID, newTimeLine); err != nil {
		return nil, err
	}
	if newTimeLine, err = InsertNewEventInTimeLine(newE, newTimeLine); err != nil {
		return nil, err
	}
	return newTimeLine, nil
}

func (c *Calendar) AddEvent(timeStart, timeEnd int64, topic string) (eventID int, err error) {

	if timeStart == 0 || (timeEnd != 0 && timeStart >= timeEnd) {
		return 0, ErrWrongTime
	}
	if topic == "" {
		return 0, ErrEmptyTopic
	}
	if timeEnd == 0 {
		timeEnd = timeStart + defatuldEventDurationInSeconds
	}
	// in the future, ID will be taken from the database
	c.mu.Lock()
	var timeLine []*Event
	newEvent := &Event{
		ID:        c.lastEventID + 1,
		TimeStart: time.Unix(timeStart, 0),
		TimeEnd:   time.Unix(timeEnd, 0),
		Topic:     topic,
	}
	if timeLine, err = InsertNewEventInTimeLine(newEvent, c.eventsOrderedByTime); err != nil {
		c.mu.Unlock()
		return 0, err
	}
	c.lastEventID = newEvent.ID
	c.eventsOrderedByTime = timeLine
	c.Events[newEvent.ID] = newEvent
	c.EventsByStartTime[newEvent.TimeStart.Unix()] = newEvent
	c.mu.Unlock()
	return newEvent.ID, nil
}

func (c *Calendar) DeleteEvent(timeStart int64) (err error) {
	c.mu.Lock()
	var e *Event
	var ok bool
	if e, ok = c.EventsByStartTime[timeStart]; !ok {
		c.mu.Unlock()
		return ErrNotFound
	}
	var timeLine []*Event
	if timeLine, err = DeleteEventFromTimeLine(e.ID, c.eventsOrderedByTime); err != nil {
		c.mu.Unlock()
		return err
	}
	c.eventsOrderedByTime = timeLine
	delete(c.EventsByStartTime, timeStart)
	delete(c.Events, e.ID)
	c.mu.Unlock()
	return nil

}

func (c *Calendar) DeleteEvents(idsForDelete ...int) (deletedIDs []int, err error) {
	c.mu.Lock()
	var timeLine []*Event
	var e *Event
	var ok bool
	deletedIDs = make([]int, 0, len(idsForDelete))
	for _, id := range idsForDelete {
		if e, ok = c.Events[id]; !ok {
			c.mu.Unlock()
			return deletedIDs, ErrNotFound
		}
		if timeLine, err = DeleteEventFromTimeLine(e.ID, c.eventsOrderedByTime); err != nil {
			c.mu.Unlock()
			return deletedIDs, err
		}
		c.eventsOrderedByTime = timeLine
		delete(c.EventsByStartTime, e.TimeStart.Unix())
		delete(c.Events, e.ID)
		deletedIDs = append(deletedIDs, e.ID)

	}
	c.mu.Unlock()
	return deletedIDs, nil
}

func (c *Calendar) EditEvent(oldTimeStart, timeStart, timeEnd int64, topic string) (err error) {
	if timeStart == 0 || (timeEnd != 0 && timeStart >= timeEnd) {
		return ErrWrongTime
	}
	if topic == "" {
		return ErrEmptyTopic
	}
	if timeEnd == 0 {
		timeEnd = timeStart + defatuldEventDurationInSeconds
	}
	c.mu.Lock()
	var e *Event
	var ok bool
	if e, ok = c.EventsByStartTime[oldTimeStart]; !ok {
		c.mu.Unlock()
		return ErrNotFound
	}
	newE := e.CreateCopy()
	newE.Topic = topic
	if e.TimeStart.Unix() != timeStart || e.TimeEnd.Unix() != timeEnd {
		newE.TimeStart = time.Unix(timeStart, 0)
		newE.TimeEnd = time.Unix(timeEnd, 0)
		var timeLine []*Event
		if timeLine, err = ChangeEventInTimeLine(e, newE, c.eventsOrderedByTime); err != nil {
			c.mu.Unlock()
			return err
		}
		c.eventsOrderedByTime = timeLine
	}
	if !e.TimeStart.Equal(newE.TimeStart) {
		delete(c.EventsByStartTime, oldTimeStart)
		c.EventsByStartTime[timeStart] = newE
	}
	c.mu.Unlock()
	return nil
}

func (c *Calendar) EditEventByID(eventID int, timeStart, timeEnd int64, topic string) (err error) {
	var e *Event
	var ok bool
	if e, ok = c.Events[eventID]; !ok {
		return ErrNotFound
	}
	if err := c.EditEvent(e.TimeStart.Unix(), timeStart, timeEnd, topic); err != nil {
		return err
	}
	return nil
}

func (c *Calendar) ShowEvents(w io.Writer) {
	for _, e := range c.eventsOrderedByTime {
		fmt.Fprintf(w, "%s - %s id %d :'%s'\n", e.TimeStart.Format("2006-01-02 15:04:05"), e.TimeEnd.Format("2006-01-02 15:04:05"), e.ID, e.Topic)
	}
}

func (c *Calendar) ShowUpcomingEvents(w io.Writer) {
	var j int
	for j = range c.eventsOrderedByTime {
		if c.eventsOrderedByTime[j].TimeEnd.After(time.Now()) {
			break
		}
	}
	for i := j; i < len(c.eventsOrderedByTime); i++ {
		fmt.Fprintf(w, "%s - %s id %d :'%s'\n", c.eventsOrderedByTime[i].TimeStart.Format("2006-01-02 15:04:05"),
			c.eventsOrderedByTime[i].TimeEnd.Format("2006-01-02 15:04:05"), c.eventsOrderedByTime[i].ID, c.eventsOrderedByTime[i].Topic)
	}

}
