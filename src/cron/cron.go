package cron

import (
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
	"pokemonscan-pokeball/src/utils/docker"
)

func init() {
	c := cron.New()
	err := c.AddFunc("0 */30 * * * *", docker.CleanHangContainers)
	if err != nil {
		log.Error(err)
	}
	c.Start()
	log.Info("cron job init success")
}
