package db

import "fmt"

// ID хранится строкой, чтобы не терять точность в JSON на стороне JS
type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func GetTask(id string) (*Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	row := db.QueryRow(query, id)

	t := &Task{}
	var dbID int64
	if err := row.Scan(&dbID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
		return nil, fmt.Errorf("задача не найдена")
	}
	t.ID = fmt.Sprint(dbID)
	return t, nil
}

func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date=?, title=?, comment=?, repeat=? WHERE id=?`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}

func UpdateDate(id, date string) error {
	query := `UPDATE scheduler SET date=? WHERE id=?`
	res, err := db.Exec(query, date, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}

func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id=?`
	res, err := db.Exec(query, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}

func Tasks(limit int) ([]*Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`
	rows, err := db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		t := &Task{}
		var dbID int64
		if err := rows.Scan(&dbID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, err
		}
		t.ID = fmt.Sprint(dbID)
		tasks = append(tasks, t)
	}
	if tasks == nil {
		tasks = []*Task{}
	}
	return tasks, rows.Err()
}

func TasksBySearch(search string, limit int) ([]*Task, error) {
	pattern := "%" + search + "%"
	query := `SELECT id, date, title, comment, repeat FROM scheduler
	          WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?`
	rows, err := db.Query(query, pattern, pattern, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		t := &Task{}
		var dbID int64
		if err := rows.Scan(&dbID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, err
		}
		t.ID = fmt.Sprint(dbID)
		tasks = append(tasks, t)
	}
	if tasks == nil {
		tasks = []*Task{}
	}
	return tasks, rows.Err()
}

func TasksByDate(date string, limit int) ([]*Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date LIMIT ?`
	rows, err := db.Query(query, date, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		t := &Task{}
		var dbID int64
		if err := rows.Scan(&dbID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, err
		}
		t.ID = fmt.Sprint(dbID)
		tasks = append(tasks, t)
	}
	if tasks == nil {
		tasks = []*Task{}
	}
	return tasks, rows.Err()
}
