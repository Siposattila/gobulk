package sync

import (
	"time"

	"github.com/Siposattila/gobulk/internal/config"
	"github.com/Siposattila/gobulk/internal/console"
	"github.com/Siposattila/gobulk/internal/email"
	"github.com/Siposattila/gobulk/internal/gorm"
	"github.com/robfig/cron/v3"
	g "gorm.io/gorm"
)

type Sync struct {
	Config      *config.Config
	Cron        cron.Schedule
	EM          *gorm.EntityManager
	MEM         *gorm.EntityManager
	stopChannel chan bool
}

func Init() *Sync {
	sync := &Sync{
		EM: gorm.Gorm(),
	}
	sync.Config = gorm.GetConfig(sync.EM.GormORM)
	sync.MEM = gorm.GormExternal(&sync.Config.MysqlDSN)

	cronParser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	schedule, parseError := cronParser.Parse(sync.Config.SyncCron)
	if parseError != nil {
		console.Fatal("Your cron expression is invalid or an error occured: " + parseError.Error())
	}
	sync.Cron = schedule

	return sync
}

func (s *Sync) Start() {
	s.syncProcess()
}

func (s *Sync) Stop() {
	s.stopChannel <- true
}

func (s *Sync) syncProcess() {
	s.sync()
	timeSignal := time.After(s.Cron.Next(time.Now()).Sub(time.Now()))
	console.Success("Next sync will be at: " + s.Cron.Next(time.Now()).String())

	select {
	case <-timeSignal:
		s.sync()
		s.syncProcess()
	case <-s.stopChannel:
		console.Normal("Stopping sync process...")
	}
}

func (s *Sync) sync() {
	console.Normal("Sync is started. This may take a long time!!!")
	s.cacheMysqlData()
	diff := s.getDifference()
	for _, d := range diff {
		var e email.Email
		tx := s.EM.GormORM.First(&e, "name = ? AND email = ?", d.Name, d.Email)
		if tx.Error != nil {
			if tx.Error == g.ErrRecordNotFound {
				s.EM.GormORM.Create(&email.Email{Name: d.Name, Email: d.Email})
				console.Normal("Create record " + d.Name + " " + d.Email)
			} else {
				console.Error("Error executing query: %v", tx.Error)
			}
		} else {
			e.Status = email.EMAIL_STATUS_INACTIVE
			s.EM.GormORM.Save(&e)
			console.Normal("Found record " + d.Name + " " + d.Email)
		}
	}
	console.Success("Sync finished successfully!")
}

func (s *Sync) cacheMysqlData() {
	// This is TRUNCATE in sqlite
	s.EM.GormORM.Exec("DELETE FROM caches;")

	var results []email.Cache
	s.MEM.GormORM.Raw(s.Config.MysqlQuery).FindInBatches(&results, 100, func(tx *g.DB, batch int) error {
		for _, result := range results {
			tx := s.EM.GormORM.Create(result)
			if tx.Error != nil {
				console.Error(tx.Error)

				return tx.Error
			}
		}

		// Returning an error will stop further batch processing
		return nil
	})
}

func (s *Sync) getDifference() []email.Cache {
	var cache []email.Cache
	tx := s.EM.GormORM.Raw(`SELECT email, name FROM emails WHERE (email, name) NOT IN (SELECT email, name FROM caches)
    UNION ALL SELECT email, name FROM caches WHERE (email, name) NOT IN (SELECT email, name FROM emails);`).Scan(&cache)
	if tx.Error != nil {
		console.Fatal(tx.Error)
	}

	return cache
}
