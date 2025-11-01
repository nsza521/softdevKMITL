// internal/db_model/types.go
package models

import (
	"database/sql/driver"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type OnlyTime struct{ time.Time }

// generic data type
func (OnlyTime) GormDataType() string { return "time" }

// DB-specific (GORM v2 expects this signature)
func (OnlyTime) GormDBDataType(db *gorm.DB, f *schema.Field) string {
	return "TIME"
}

// write to DB as HH:MM:SS
func (t OnlyTime) Value() (driver.Value, error) {
	return t.Time.Format("15:04:05"), nil
}

// read from DB (TIME -> string -> time.Time)
func (t *OnlyTime) Scan(value any) error {
	switch v := value.(type) {
	case []byte:
		tt, err := time.Parse("15:04:05", string(v))
		if err != nil { return err }
		t.Time = tt
		return nil
	case string:
		tt, err := time.Parse("15:04:05", v)
		if err != nil { return err }
		t.Time = tt
		return nil
	default:
		return fmt.Errorf("unsupported type %T for OnlyTime", value)
	}
}
