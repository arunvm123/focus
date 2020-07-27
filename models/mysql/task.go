package mysql

import (
	"errors"
	"fmt"
	"time"

	"github.com/arunvm/travail-backend/models"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

func (db *Mysql) CreateTask(args *models.CreateTaskArgs, user *models.User) (*models.Task, error) {
	list, err := db.getList(args.ListID)
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
	var task models.Task

	tx := db.Client.Begin()

	task.ID = uuid.New().String()
	task.Info = args.Info
	task.Order = args.Order
	task.ExpiresAt = args.ExpiresAt
	task.ListID = list.ID
	task.Archived = false
	task.Complete = false
	task.CreatedAt = time.Now().Unix()

	err = tx.Create(&task).Error
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

func (db *Mysql) GetTasks(args *models.GetTasksArgs, user *models.User) (*[]models.TaskInfo, error) {
	var tasks []models.TaskInfo

	var completeFilter string
	if args.Complete != nil {
		completeFilter = fmt.Sprintf(" AND tasks.complete = %v", *args.Complete)
	}

	err := db.Client.Table("lists").Joins("JOIN tasks on tasks.list_id = lists.id").
		Where("lists.archived = false AND tasks.archived = false AND lists.user_id = ? AND lists.id = ?"+completeFilter, user.ID, args.ListID).
		Order("order").
		Select("tasks.*,lists.heading").Find(&tasks).Error
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

func (db *Mysql) UpdateTask(args models.UpdateTaskArgs, user *models.User) error {
	task, err := db.getTaskOfUser(user.ID, args.ID)
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

	err = db.Client.Save(task).Error
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

func (db *Mysql) DeleteTasks(args *models.DeleteTasksArgs, user *models.User) error {
	err := db.Client.Table("tasks JOIN lists on tasks.list_id = lists.id").
		Where("user_id = ? AND tasks.id IN (?)", user.ID, args.TaskIDs).
		UpdateColumn("tasks.archived", true).
		Error
	if err != nil {
		log.WithFields(log.Fields{
			"func":   "DeleteTasks",
			"info":   "deleting tasks specified by id",
			"userID": user.ID,
			"args":   *args,
		}).Error(err)
		return err
	}

	return nil
}

// func createTasks(db *gorm.DB, tasks *[]Task) error {
// 	if len(*tasks) == 0 {
// 		return nil
// 	}
// 	sqlQuery := "INSERT into tasks(list_id,info,created_at,complete,`order`,archived) VALUES "

// 	scope := db.NewScope(Task{})
// 	for i := 0; i < len(*tasks); i++ {
// 		sqlQuery = sqlQuery + "(?,?,?,?,?,?),"
// 		scope.AddToVars((*tasks)[i].ListID)
// 		scope.AddToVars((*tasks)[i].Info)
// 		scope.AddToVars((*tasks)[i].CreatedAt)
// 		scope.AddToVars((*tasks)[i].Complete)
// 		scope.AddToVars((*tasks)[i].Order)
// 		scope.AddToVars((*tasks)[i].Archived)
// 	}

// 	sqlQuery = sqlQuery[:len(sqlQuery)-1]

// 	err := db.Exec(sqlQuery, scope.SQLVars...).Error
// 	if err != nil {
// 		log.Println("Error when creating tasks")
// 		return err
// 	}

// 	return nil
// }

func (db *Mysql) getTaskOfUser(userID int, taskID string) (*models.Task, error) {
	var task models.Task

	err := db.Client.Table("lists").Joins("JOIN tasks on tasks.list_id = lists.id").
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

func (db *Mysql) GetTasksAboutToExpire() (*[]models.TaskInfo, error) {
	startTime := time.Now().Unix()
	endTime := startTime + (5 * 60)

	var tasks []models.TaskInfo
	err := db.Client.Table("tasks").Joins("JOIN lists on tasks.list_id = lists.id").
		Select("tasks.*,lists.user_id,lists.heading").
		Where("expires_at BETWEEN ? AND ? AND tasks.archived = false AND lists.archived = false", startTime, endTime).
		Find(&tasks).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func": "SendPushNotificationForTasksAboutToExpire",
			"info": "retrieving tasks expiring within the next 5 mins",
		}).Error(err)
		return nil, err
	}

	return &tasks, nil
}
