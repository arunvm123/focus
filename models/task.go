package models

import (
	"errors"
	"time"

	"firebase.google.com/go/messaging"
	push "github.com/arunvm/travail-backend/push_notification"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
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
	ListID    int    `json:"listID" binding:"required"`
	Info      string `json:"info"`
	Order     int    `json:"order"`
	ExpiresAt *int64 `json:"expiresAt"`
}

type GetTasksArgs struct {
	ListID int `json:"listID" binding:"required"`
}

type UpdateTaskArgs struct {
	ID        int     `json:"id" binding:"required"`
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

func (user *User) CreateTask(db *gorm.DB, args *CreateTaskArgs) (*Task, error) {
	list, err := getList(db, args.ListID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("No list with specified id")
		}
		log.WithFields(log.Fields{
			"func":    "CreateTask",
			"subFunc": "getList",
			"userID":  user.ID,
			"listID":  args.ListID,
		}).Error(err)
		return nil, err
	}

	if list.UserID != user.ID {
		return nil, errors.New("Not user's list")
	}
	var task Task

	tx := db.Begin()

	task.Info = args.Info
	task.Order = args.Order
	task.ExpiresAt = args.ExpiresAt
	task.ListID = list.ID
	task.Archived = false
	task.Complete = false
	task.CreatedAt = time.Now().Unix()

	err = task.Create(tx)
	if err != nil {
		tx.Rollback()
		log.WithFields(log.Fields{
			"func":    "CreateTask",
			"subFunc": "task.Create",
			"userID":  user.ID,
			"listID":  task.ListID,
		}).Error(err)
		return nil, err
	}

	tx.Commit()
	return &task, nil
}

func (user *User) GetTasks(db *gorm.DB, args *GetTasksArgs) (*[]Task, error) {
	var tasks []Task

	err := db.Table("lists").Joins("JOIN tasks on tasks.list_id = lists.id").
		Where("lists.archived = false AND tasks.archived = false AND lists.user_id = ? AND lists.id = ?", user.ID, args.ListID).
		Order("order").
		Select("tasks.*").Find(&tasks).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func":   "GetTasks",
			"info":   "retrieving task details",
			"userID": user.ID,
			"listID": args.ListID,
		}).Error(err)
		return nil, err
	}

	return &tasks, nil
}

func (user *User) UpdateTask(db *gorm.DB, args UpdateTaskArgs) error {
	task, err := getTaskOfUser(db, user.ID, args.ID)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "UpdateTask",
			"subFunc": "getTasksOfUser",
			"userID":  user.ID,
			"taskID":  args.ID,
		}).Error(err)
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
		log.WithFields(log.Fields{
			"func":    "UpdateTask",
			"subFunc": "task.Save",
			"userID":  user.ID,
			"taskID":  args.ID,
		}).Error(err)
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
		log.WithFields(log.Fields{
			"func":   "getTaskOfUser",
			"info":   "retrieving task info",
			"taskID": taskID,
			"userID": userID,
		}).Error(err)
		return nil, err
	}

	return &task, nil
}

func SendPushNotificationForTasksAboutToExpire(db *gorm.DB, pushClient *messaging.Client) error {
	startTime := time.Now().Unix()
	endTime := startTime + (5 * 60)

	var userIDs []int
	err := db.Table("tasks").Joins("JOIN lists on tasks.list_id = lists.id").
		Where("expires_at BETWEEN ? AND ? AND tasks.archived = false AND lists.archived = false", startTime, endTime).
		Pluck("DISTINCT(user_id)", &userIDs).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func": "SendPushNotificationForTasksAboutToExpire",
			"info": "retrieving users who have tasks expiring within the next 5 mins",
		}).Error(err)
		return err
	}

	if len(userIDs) == 0 {
		return nil
	}

	var deviceTokens []string
	err = db.Table("fcm_notification_tokens").Where("user_id IN (?)", userIDs).
		Pluck("token", &deviceTokens).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func": "SendPushNotificationForTasksAboutToExpire",
			"info": "retrieving users device tokens",
		}).Error(err)
		return err
	}

	if len(deviceTokens) == 0 {
		return nil
	}

	err = push.SendPushNotification(pushClient, deviceTokens, "You have tasks that are about to expire")
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "SendPushNotificationForTasksAboutToExpire",
			"subFunc": "push.SendPushNotification",
		}).Error(err)
		return err
	}

	return nil
}
