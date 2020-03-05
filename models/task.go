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
	Info  string `json:"info"`
	Order int    `json:"order"`
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
			Heading:   "Default List",
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
