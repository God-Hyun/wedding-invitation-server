package sqldb

import (
	"fmt"
	"time"
)

func initializeAttendanceTable() error {
	_, err := sqlDb.Exec(`
		CREATE TABLE IF NOT EXISTS attendance (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			side VARCHAR(10),
			name VARCHER(20),
			meal VARCHAR(20),
			count INTEGER,
			timestamp INTEGER
		)
	`)
	return err
}

type AttendanceRow struct {
	Id        int
	Side      string
	Name      string
	Meal      string
	Count     int
	Timestamp int64
}

func GetAllAttendance() ([]AttendanceRow, error) {
	rows, err := sqlDb.Query(`
		SELECT id, side, name, meal, count, timestamp
		FROM attendance
		ORDER BY timestamp DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []AttendanceRow
	for rows.Next() {
		var row AttendanceRow
		err := rows.Scan(&row.Id, &row.Side, &row.Name, &row.Meal, &row.Count, &row.Timestamp)
		if err != nil {
			return nil, err
		}
		result = append(result, row)
	}
	return result, nil
}

func CreateAttendance(side, name, meal string, count int) error {
	_, err := sqlDb.Exec(`
		INSERT INTO attendance (side, name, meal, count, timestamp)
		VALUES (?, ?, ?, ?, ?)
	`, side, name, meal, count, time.Now().Unix())
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
