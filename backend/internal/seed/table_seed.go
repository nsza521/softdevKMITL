package seed

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"backend/internal/db_model"
)

func seedTableTimeslots(db *gorm.DB) error {
	// Create tables
	var tables []models.Table
	for col := 1; col <= 3; col++ {
		for row := 1; row <= 6; row++ {
			table := models.Table{
				MaxSeats: 6,
				TableRow:      fmt.Sprintf("%c", 'A'+(row-1)),
				TableCol:      fmt.Sprintf("%d", col),
			}
			if err := db.Where("table_row = ? AND table_col = ?", table.TableRow, table.TableCol).
				FirstOrCreate(&table).Error; err != nil {
				return err
			}
			tables = append(tables, table)
		}
	}

	// Create time slots from 10:00 to 15:00 with 15-minute intervals
	loc, _ := time.LoadLocation("Asia/Bangkok")
	now := time.Now().In(loc)
	baseDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)

	start := baseDate.Add(10 * time.Hour)
	end := baseDate.Add(15 * time.Hour)
	duration := 15 * time.Minute

	var timeSlots []models.Timeslot
	for t := start; t.Before(end); t = t.Add(duration) {
		timeSlot := models.Timeslot{
			StartTime: t,
			EndTime:   t.Add(duration),
		}
		if err := db.Where("start_time = ? AND end_time = ?", t, t.Add(duration)).
			FirstOrCreate(&timeSlot).Error; err != nil {
			return fmt.Errorf("error creating timeslot %v-%v: %v", t, t.Add(duration), err)
		}
		timeSlots = append(timeSlots, timeSlot)
	}

	// Create table-timeslot associations
	for _, table := range tables {
		for _, timeSlot := range timeSlots {
			status := "available"
			tableTimeslot := models.TableTimeslot{
				TableID:       table.ID,
				TimeslotID:    timeSlot.ID,
				Status:        status,
				ReservedSeats: 0,
			}
			if err := db.Where("table_id = ? AND timeslot_id = ?", table.ID, timeSlot.ID).
				FirstOrCreate(&tableTimeslot).Error; err != nil {
				return err
			}
		}
	}

	return nil
}