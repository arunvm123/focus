package main

import (
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
	err := server.db.SendPushNotificationForTasksAboutToExpire(server.push)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "sendPushNotificationForTasksAboutToExpire",
			"subFunc": "models.SendPushNotificationForTasksAboutToExpire",
		}).Error(err)
		return
	}

	return
}
