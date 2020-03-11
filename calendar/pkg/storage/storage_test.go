package storage

import (
	"fmt"
	_ "golang_home_work/calendar/pkg/config"
	"golang_home_work/work_11/config"
	"os"
	"testing"
	"time"

	"go.uber.org/zap"
)

var (
	tableTest = [...]struct {
		start, end time.Time
		topic      string
	}{
		{time.Now().Add(-3 * time.Hour), time.Now().Add(-2 * time.Hour), `test1`},
		{time.Now().Add(-8 * time.Hour), time.Now().Add(-6 * time.Hour), `test2`},
		{time.Now().Add(2 * time.Hour), time.Now().Add(6 * time.Hour), `test3`},
		{time.Now().Add(7 * time.Hour), time.Now().Add(8 * time.Hour), `test4`},
		{time.Now().Add(3 * time.Minute), time.Now().Add(6 * time.Minute), `test5`},
		{time.Now().Add(-10 * time.Minute), time.Now().Add(1 * time.Minute), `test6`},
		{time.Now().Add(30 * time.Minute), time.Now().Add(45 * time.Minute), `test7`},
		{time.Now().Add(65 * time.Minute), time.Now().Add(85 * time.Minute), `test8`},
	}
)

func createCalendarForTest() (c *Calendar) {
	c = CreateCalendar()
	for _, t := range tableTest {
		if _, err := c.AddEvent(t.start.Unix(), t.end.Unix(), t.topic); err != nil {
			zap.S().Error(err)
			return c
		}
	}
	return c
}

func TestInsertEvents(t *testing.T) {

	var calendar Storage = CreateCalendar()
	var err error
	var lastAddEventID int
	for _, t := range tableTest {
		if lastAddEventID, err = calendar.AddEvent(t.start.Unix(), t.end.Unix(), t.topic); err != nil {
			zap.S().Error(err)
		} else {
			zap.S().Infof("successfully added event with id =%d", lastAddEventID)
		}
	}
	fmt.Fprint(os.Stdout, "Event list\n")
	calendar.ShowEvents(os.Stdout)
	fmt.Fprint(os.Stdout, "Upcoming event sheet\n")
	calendar.ShowUpcomingEvents(os.Stdout)
}

func TestDeleteEvents(t *testing.T) {

	var calendar Storage = createCalendarForTest()
	var err error

	fmt.Fprint(os.Stdout, "Event list before delete\n")
	calendar.ShowEvents(os.Stdout)
	fmt.Fprintf(os.Stdout, "Delete event with start date %s\n", tableTest[3].start.Format("2006-01-02 15:04:05"))
	if err = calendar.DeleteEvent(tableTest[3].start.Unix()); err != nil {
		zap.S().Error(err)
	}
	fmt.Fprintf(os.Stdout, "Delete event with IDs %d %d %d\n", 1, 5, 6)
	if _, err = calendar.DeleteEvents(1, 5, 6); err != nil {
		zap.S().Error(err)
	}
	fmt.Fprint(os.Stdout, "Event list after delete\n")
	calendar.ShowEvents(os.Stdout)
}

func TestEditEvents(t *testing.T) {

	var calendar Storage = createCalendarForTest()
	var err error

	fmt.Fprint(os.Stdout, "Event list before edit\n")
	calendar.ShowEvents(os.Stdout)
	fmt.Fprintf(os.Stdout, "Edit event with start date %s\n", tableTest[5].start.Format("2006-01-02 15:04:05"))
	if err = calendar.EditEvent(tableTest[5].start.Unix(), tableTest[5].start.Add(-1*time.Hour).Unix(), 0, `test 6 was mod`); err != nil {
		zap.S().Error(err)
	}
	fmt.Fprintf(os.Stdout, "Edit event with start ID %d\n", 4)
	if err = calendar.EditEventByID(4, tableTest[4].start.Add(-1*time.Hour).Unix(), tableTest[4].end.Add(-1*time.Hour).Unix(), `test 4 was mod`); err != nil {
		zap.S().Error(err)
	}
	fmt.Fprint(os.Stdout, "Event list after edit\n")
	calendar.ShowEvents(os.Stdout)
}

func TestGetErrors(t *testing.T) {

	var calendar Storage = createCalendarForTest()
	var err error

	fmt.Fprint(os.Stdout, "Event list before edit\n")
	calendar.ShowEvents(os.Stdout)
	fmt.Fprintf(os.Stdout, "Add event with start date %s for test error '%s'\n", time.Now().Add(30*time.Minute).Format("2006-01-02 15:04:05"), ErrDateBusy)
	if _, err = calendar.AddEvent(time.Now().Add(30*time.Minute).Unix(), 0, `Test event bysy Time`); err != ErrDateBusy {
		zap.S().Errorf("failed to get error '%s'\n", ErrDateBusy)
	}
	fmt.Fprintf(os.Stdout, "Add event with start date %s for test error '%s'\n", time.Now().Add(7*time.Hour+30*time.Minute).Format("2006-01-02 15:04:05"), ErrDateBusy)
	if _, err = calendar.AddEvent(time.Now().Add(7*time.Hour+30*time.Minute).Unix(), 0, `Test event bysy Time`); err != ErrDateBusy {
		zap.S().Errorf("failed to get error '%s'\n", ErrDateBusy)
	}
	fmt.Fprintf(os.Stdout, "Add event with empty topic for test error '%s'\n", ErrEmptyTopic)
	if _, err = calendar.AddEvent(time.Now().Add(10*time.Hour).Unix(), 0, ``); err != ErrEmptyTopic {
		zap.S().Errorf("failed to get error '%s'\n", ErrEmptyTopic)
	}
	fmt.Fprintf(os.Stdout, "Add event with start %s  end %s for test error '%s'\n", time.Now().Add(10*time.Hour).Format("2006-01-02 15:04:05"),
		time.Now().Add(9*time.Hour).Format("2006-01-02 15:04:05"), ErrWrongTime)
	if _, err = calendar.AddEvent(time.Now().Add(10*time.Hour).Unix(), time.Now().Add(9*time.Hour).Unix(), `Test event wrong time`); err != ErrWrongTime {
		zap.S().Errorf("failed to get error '%s'\n", ErrWrongTime)
	}
	fmt.Fprintf(os.Stdout, "Delete event with startDate %s  for test error '%s'\n", time.Now().Add(10*time.Hour).Format("2006-01-02 15:04:05"), ErrNotFound)
	if err = calendar.DeleteEvent(time.Now().Add(10 * time.Hour).Unix()); err != ErrNotFound {
		zap.S().Errorf("failed to get error '%s'\n", ErrNotFound)
	}
	fmt.Fprintf(os.Stdout, "Delete event with IDs %d %d %d  for test error '%s'\n", 2, 5, 12, ErrNotFound)
	if _, err = calendar.DeleteEvents(2, 5, 12); err != ErrNotFound {
		zap.S().Errorf("failed to get error '%s'\n", ErrNotFound)
	}
	calendar.ShowEvents(os.Stdout)
	fmt.Fprintf(os.Stdout, "Edit event with old data %s start %s  end %s for test error '%s'\n", tableTest[7].start.Format("2006-01-02 15:04:05"),
		tableTest[7].end.Format("2006-01-02 15:04:05"),
		tableTest[7].start.Format("2006-01-02 15:04:05"), ErrWrongTime)
	if err = calendar.EditEvent(tableTest[7].start.Unix(), tableTest[7].end.Unix(), tableTest[7].start.Unix(), tableTest[7].topic); err != ErrWrongTime {
		zap.S().Errorf("failed to get error '%s'\n", ErrWrongTime)
	}
	fmt.Fprintf(os.Stdout, "Edit event with id %d - change topic to empty for test error '%s'\n", 4, ErrEmptyTopic)
	if err = calendar.EditEventByID(4, tableTest[7].start.Unix(), tableTest[7].end.Unix(), ``); err != ErrEmptyTopic {
		zap.S().Errorf("failed to get error '%s'\n", ErrEmptyTopic)
	}
	fmt.Fprintf(os.Stdout, "Edit event with id %d - change topic to empty for test error '%s'\n", 15, ErrNotFound)
	if err = calendar.EditEventByID(15, tableTest[7].start.Unix(), tableTest[7].end.Unix(), `Test get error edit - not found`); err != ErrNotFound {
		zap.S().Errorf("failed to get error '%s'\n", ErrNotFound)
	}
	fmt.Fprintf(os.Stdout, "Edit event with old data %s start %s  end %s for test error '%s'\n", tableTest[7].start.Add(5*time.Minute).Format("2006-01-02 15:04:05"),
		tableTest[7].start.Format("2006-01-02 15:04:05"),
		tableTest[7].end.Format("2006-01-02 15:04:05"), ErrNotFound)
	if err = calendar.EditEvent(tableTest[7].start.Add(5*time.Minute).Unix(), tableTest[7].start.Unix(), tableTest[7].end.Unix(), tableTest[7].topic); err != ErrNotFound {
		zap.S().Errorf("failed to get error '%s'\n", ErrNotFound)
	}
	fmt.Fprintf(os.Stdout, "Edit event with id %d, to start %s end %s for test error '%s'\n", 4, tableTest[3].start.Add(-2*time.Hour).Format("2006-01-02 15:04:05"),
		tableTest[3].end.Add(-2*time.Hour).Format("2006-01-02 15:04:05"), ErrDateBusy)
	if err = calendar.EditEventByID(4, tableTest[3].start.Add(-2*time.Hour).Unix(), tableTest[3].end.Add(-2*time.Hour).Unix(),
		`Test get error edit - time is bysy`); err != ErrDateBusy {
		zap.S().Errorf("failed to get error '%s'\n", ErrDateBusy)
	}
}

func init() {
	var logger *zap.Logger
	var err error
	if logger, err = config.GetStandartLogger(`development`); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
		return
	}
	zap.ReplaceGlobals(logger)
}
