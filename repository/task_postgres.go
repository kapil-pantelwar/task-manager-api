package repository

import (
   "database/sql"
    "task-manager/domain"
)

type TaskPostgresRepo struct {
	db *sql.DB
}

func NewTaskPostgresRepo(db *sql.DB) *TaskPostgresRepo {
	repo := &TaskPostgresRepo{db: db}
	repo.initDB()
	return repo
}

func (r *TaskPostgresRepo) initDB() {
	query := `
	CREATE TABLE IF NOT EXISTS tasks (
	id SERIAL PRIMARY KEY,
            title TEXT NOT NULL,
            description TEXT,
            status TEXT NOT NULL
            );`
  _,err := r.db.Exec(query)
  if err != nil {
	panic("Failed to initialize database: "+err.Error())
  }
}

func (r *TaskPostgresRepo) Create(task domain.Task)(domain.Task, error){
	err := r.db.QueryRow("INSERT INTO tasks (title, description, status) VALUES ($1,$2,$3) RETURNING id",task.Title,task.Description,task.Status,).Scan(&task.ID)
    return task,err
}

func (r *TaskPostgresRepo) GetAll() ([]domain.Task, error){
	rows,err := r.db.Query("SELECT id, title, description, status FROM tasks")
	if err != nil {
		return nil,err
	}
	defer rows.Close()

	var tasks []domain.Task
	for rows.Next() {
		var task domain.Task
		err := rows.Scan(&task.ID,&task.Title,&task.Description,&task.Status)

		if err != nil {
			return nil,err	
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *TaskPostgresRepo) GetByID(id int) (domain.Task, error){
	var task domain.Task
    err := r.db.QueryRow("SELECT id,title,description, status FROM tasks WHERE id = $1",id).Scan(&task.ID,&task.Title,&task.Description,&task.Status)
    if err == sql.ErrNoRows {
        return domain.Task{},sql.ErrNoRows
    }
    return task, err
}

func (r *TaskPostgresRepo) Update(task domain.Task) (domain.Task, error) {
   var t domain.Task
   err := r.db.QueryRow(`
   UPDATE tasks SET
    title = COALESCE(NULLIF($1,''),title),
    description = COALESCE(NULLIF($2,''),description),
    status = COALESCE(NULLIF($3,''),status)
    WHERE id = $4 RETURNING id, title, description, status`,task.Title,task.Description,task.Status,task.ID,).Scan(&t.ID,&t.Title,&t.Description,&t.Status)

    return t,err

}
func (r *TaskPostgresRepo) Delete(id int) error {
    result, err := r.db.Exec("DELETE FROM tasks WHERE id = $1", id)
    if err != nil {
        return err
    }
    rowsAffected, err := result.RowsAffected()
   
    if rowsAffected == 0 {
        return sql.ErrNoRows
    }
    return err
}