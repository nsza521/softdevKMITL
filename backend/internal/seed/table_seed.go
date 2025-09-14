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
				PeopleNum: 6,
				Row:      fmt.Sprintf("%c", 'A'+(row-1)),
				Col:      fmt.Sprintf("%d", col),
			}
			if err := db.Create(&table).Error; err != nil {
				return err
			}
			tables = append(tables, table)
		}
	}

	// Create time slots from 10:01 to 13:00 with 15-minute intervals
	loc, _ := time.LoadLocation("Asia/Bangkok")
	now := time.Now().In(loc)
	baseDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)

	gap := 1 * time.Minute
	start := baseDate.Add(10 * time.Hour + gap)   // 10:01 in Thailand
	end := baseDate.Add(13 * time.Hour)     // 13:00 in Thailand
	duration := 14 * time.Minute
	var timeSlots []models.Timeslot

	for t := start; t.Before(end); t = t.Add(duration + gap) {
		timeSlot := models.Timeslot{
			StartTime: t,
			EndTime:   t.Add(duration),
		}
		if err := db.Create(&timeSlot).Error; err != nil {
			return err
		}
		timeSlots = append(timeSlots, timeSlot)
	}


	// Create table-time slot associations
	for _, table := range tables {
		for _, timeSlot := range timeSlots {

			status := "available"
			if timeSlot.EndTime.In(loc).After(time.Now().In(loc)) {
				status = "expired"
			}

			tableTimeslot := models.TableTimeslot{
				TableID:    table.ID,
				TimeslotID: timeSlot.ID,
				Status:     status,
			}
			if err := db.Create(&tableTimeslot).Error; err != nil {
				return err
			}
		}
	}
	return nil
}