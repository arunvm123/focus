package main

import (
	push "github.com/arunvm/travail-backend/push_notification"
	log "github.com/sirupsen/logrus"

	"github.com/robfig/cron/v3"
)

func (server *server) startCronJobs() error {
	c := cron.New()

	_, err := c.AddFunc("@every 5m", server.notificationForTasksAboutToExpire)
	if err != nil {
		log.WithFields(log.Fields{
			"func":     "startCronJobs",
			"subFunc":  "c.AddFunc",
			"cronFunc": "models.sendPushNotificationForTasksAboutToExpire",
		}).Error(err)
		return err
	}

	c.Start()

	return nil
}

func (server *server) notificationForTasksAboutToExpire() {
	tasks, err := server.db.GetTasksAboutToExpire()
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "sendPushNotificationForTasksAboutToExpire",
			"subFunc": "models.SendPushNotificationForTasksAboutToExpire",
		}).Error(err)
		return
	}

	for i := 0; i < len(*tasks); i++ {
		go func(i int) {
			deviceTokens, err := server.db.GetNotificationTokens((*tasks)[i].UserID)
			if err != nil {
				log.WithFields(log.Fields{
					"func": "SendPushNotificationForTasksAboutToExpire",
					"info": "retrieving users device tokens",
				}).Error(err)
				return
			}

			if len(deviceTokens) == 0 {
				return
			}

			err = server.push.SendPushNotification(deviceTokens, &push.Payload{
				Body:  "'" + (*tasks)[i].Info + "' is due soon",
				Title: (*tasks)[i].Heading,
				Data: map[string]string{
					"link": "/todo/id?=" + (*tasks)[i].ListID,
				},
				ClickAction: "/todo/id?=" + (*tasks)[i].ListID,
			})
			if err != nil {
				log.WithFields(log.Fields{
					"func":    "SendPushNotificationForTasksAboutToExpire",
					"subFunc": "push.SendPushNotification",
				}).Error(err)
				return
			}

		}(i)
	}

	return
}
