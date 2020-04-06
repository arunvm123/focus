package main

import (
	"github.com/arunvm/travail-backend/models"
	log "github.com/sirupsen/logrus"

	"github.com/robfig/cron/v3"
)

func (server *server) startCronJobs() error {
	c := cron.New()

	_, err := c.AddFunc("@every 30s", server.notificationForTasksAboutToExpire)
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
	err := models.SendPushNotificationForTasksAboutToExpire(server.db, server.pushClient)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "sendPushNotificationForTasksAboutToExpire",
			"subFunc": "models.SendPushNotificationForTasksAboutToExpire",
		}).Error(err)
		return
	}

	return
}
