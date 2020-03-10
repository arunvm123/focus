package models

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
)

// Task model
type Task struct {
	ID        int    `json:"id" gorm:"primary_key"`
	ListID    int    `json:"listID"`
	Info      string `json:"info"`
	CreatedAt int64  `json:"createdAt"`
	UpdatedAt *int64 `json:"updatedAt"`
	ExpiresAt *int64 `json:"expiresAt"`
	Complete  bool   `json:"complete"`
	Order     int    `json:"order"`
	Archived  bool   `json:"archived"`
}

// Create is a helper function to create a new task
func (t *Task) Create(db *gorm.DB) error {
	return db.Create(&t).Error
}

// Save is a helper function to update existing task
func (t *Task) Save(db *gorm.DB) error {
	return db.Save(&t).Error
}

type CreateTaskArgs struct {
	Info      string `json:"info"`
	Order     int    `json:"order"`
	ExpiresAt *int64 `json:"expiresAt"`
}

type UpdateTaskArgs struct {
	ID        int     `json:"id"`
	Info      *string `json:"info"`
	Order     *int    `json:"order"`
	Complete  *bool   `json:"complete"`
	Archived  *bool   `json:"archived"`
	ExpiresAt *int64  `json:"expiresAt"`
}

func (user *User) CreateTasks(db *gorm.DB, args *[]CreateTaskArgs) error {
	list, err := getListOfUser(db, user.ID)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Printf("Error when fetching users list\n%v", err)
			return err
		}
	}

	tx := db.Begin()
	if list == nil {
		list = &List{
			UserID:    user.ID,
			Archived:  false,
			CreatedAt: time.Now().Unix(),
			Heading:   "Default List",
		}

		err = list.Create(tx)
		if err != nil {
			tx.Rollback()
			log.Printf("Error when creating list\n%v", err)
			return err
		}
	}

	var tasks []Task
	for i := range *args {
		task := Task{
			ListID:    list.ID,
			Archived:  false,
			Complete:  false,
			Info:      (*args)[i].Info,
			Order:     (*args)[i].Order,
			CreatedAt: time.Now().Unix(),
		}

		tasks = append(tasks, task)
	}

	err = createTasks(tx, &tasks)
	if err != nil {
		tx.Rollback()
		log.Printf("Error when creating tasks\n%v", err)
		return err
	}

	tx.Commit()
	return nil
}

func (user *User) CreateTask(db *gorm.DB, task *Task) error {
	list, err := getListOfUser(db, user.ID)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Printf("Error when fetching users list\n%v", err)
			return err
		}
	}

	tx := db.Begin()
	if list == nil {
		list = &List{
			UserID:    user.ID,
			Archived:  false,
			CreatedAt: time.Now().Unix(),

			Heading: "Default List",
		}

		err = list.Create(tx)
		if err != nil {
			tx.Rollback()
			log.Printf("Error when creating list\n%v", err)
			return err
		}
	}

	task.ListID = list.ID
	task.Archived = false
	task.Complete = false
	task.CreatedAt = time.Now().Unix()

	err = task.Create(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (user *User) GetTasks(db *gorm.DB) (*[]Task, error) {
	var tasks []Task

	err := db.Table("lists").Joins("JOIN tasks on tasks.list_id = lists.id").
		Where("lists.archived = false AND tasks.archived = false AND lists.user_id = ?", user.ID).
		Select("tasks.*").Find(&tasks).Error
	if err != nil {
		return nil, err
	}

	return &tasks, nil
}

func (user *User) UpdateTask(db *gorm.DB, args UpdateTaskArgs) error {
	task, err := getTaskOfUser(db, user.ID, args.ID)
	if err != nil {
		log.Printf("Error fetching task\n%v", err)
		return err
	}

	if args.Archived != nil {
		task.Archived = *args.Archived
	}
	if args.Complete != nil {
		task.Complete = *args.Complete
	}
	if args.Info != nil {
		task.Info = *args.Info
	}
	if args.Order != nil {
		task.Order = *args.Order
	}
	if args.ExpiresAt != nil {
		task.ExpiresAt = args.ExpiresAt
	}

	updatedAt := time.Now().Unix()
	task.UpdatedAt = &updatedAt

	err = task.Save(db)
	if err != nil {
		log.Printf("Error updating task\n%v", err)
		return err
	}

	return nil
}

func createTasks(db *gorm.DB, tasks *[]Task) error {
	if len(*tasks) == 0 {
		return nil
	}
	sqlQuery := "INSERT into tasks(list_id,info,created_at,complete,`order`,archived) VALUES "

	scope := db.NewScope(Task{})
	for i := 0; i < len(*tasks); i++ {
		sqlQuery = sqlQuery + "(?,?,?,?,?,?),"
		scope.AddToVars((*tasks)[i].ListID)
		scope.AddToVars((*tasks)[i].Info)
		scope.AddToVars((*tasks)[i].CreatedAt)
		scope.AddToVars((*tasks)[i].Complete)
		scope.AddToVars((*tasks)[i].Order)
		scope.AddToVars((*tasks)[i].Archived)
	}

	sqlQuery = sqlQuery[:len(sqlQuery)-1]

	err := db.Exec(sqlQuery, scope.SQLVars...).Error
	if err != nil {
		log.Println("Error when creating tasks")
		return err
	}

	return nil
}

func getTaskOfUser(db *gorm.DB, userID int, taskID int) (*Task, error) {
	var task Task

	err := db.Table("lists").Joins("JOIN tasks on tasks.list_id = lists.id").
		Where("lists.archived = false AND tasks.archived = false AND lists.user_id = ? AND tasks.id = ?", userID, taskID).
		Select("tasks.*").Find(&task).Error
	if err != nil {
		return nil, err
	}

	return &task, nil
}
