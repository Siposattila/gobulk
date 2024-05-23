package sync

import (
	"strconv"
	"strings"
	"time"

	"github.com/Siposattila/gobulk/internal/email"
	"github.com/Siposattila/gobulk/internal/interfaces"
	"github.com/Siposattila/gobulk/internal/kill"
	"github.com/Siposattila/gobulk/internal/logger"
	"github.com/robfig/cron/v3"
	"github.com/schollz/progressbar/v3"
	"gorm.io/gorm"
)

type sync struct {
	app      interfaces.AppInterface
	cron     cron.Schedule
	config   interfaces.ConfigInterface
	database interfaces.DatabaseInterface
}

func Init(app interfaces.AppInterface, database interfaces.DatabaseInterface, config interfaces.ConfigInterface) interfaces.SyncInterface {
	cronParser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	schedule, parseError := cronParser.Parse(config.GetSyncCron())
	if parseError != nil {
		logger.Fatal("Your cron expression is invalid or an error occured: " + parseError.Error())
	}

	return &sync{
		app:      app,
		cron:     schedule,
		config:   config,
		database: database,
	}
}

func (s *sync) Start() {
	s.sync()
	s.syncProcess()
}

func (s *sync) syncProcess() {
	timeSignal := time.After(s.cron.Next(time.Now()).Sub(time.Now()))
	logger.Success("Next sync will be at: " + s.cron.Next(time.Now()).String())

	select {
	case <-timeSignal:
		s.sync()
		s.syncProcess()
	case <-kill.KillCtx.Done():
		logger.Warning("Shutdown signal received shutting down sync process.")
	}
}

func (s *sync) sync() {
	logger.Normal("Sync is started. This may take a long time!")
	s.cacheMysqlData()
	diff := s.getDifference()
	var (
		batchInsert []email.Email
		total       int64
	)
	s.database.GetEntityManager().GetGormORM().Find(&email.Cache{}).Count(&total)

	logger.Warning("Synchronizing local database!")
	if len(diff) > 0 {
		bar := progressbar.Default(total)

		for _, d := range diff {
			bar.Add(1)
			var e email.Email
			tx := s.database.GetEntityManager().GetGormORM().First(&e, "name = ? AND email = ?", d.Name, d.Email)
			if tx.Error != nil {
				if tx.Error == gorm.ErrRecordNotFound {
					batchInsert = append(batchInsert, email.Email{Name: d.Name, Email: d.Email})
					logger.LogNormal("Creating " + d.Email + " " + d.Name + ".")
				} else {
					logger.Error("Error executing query: %v", tx.Error)
				}
			} else {
				e.Status = interfaces.EMAIL_STATUS_INACTIVE
				s.database.GetEntityManager().GetGormORM().Save(e)
				logger.LogNormal("Found " + d.Email + " " + d.Name + ". Inactivating!")
			}
		}

		if len(batchInsert) > 0 {
			s.database.GetEntityManager().GetGormORM().CreateInBatches(batchInsert, 100)
		}
	} else {
		logger.Success("Already up to date!")
	}

	logger.Success("Sync finished successfully! Total synced emails: " + strconv.Itoa(int(total)))
}

func (s *sync) cacheMysqlData() {
	logger.Normal("Validating cache...")
	// TODO: will need a logic to know when to delete cache
	// for now this is okay

	logger.Warning("Rebuilding cache!")
	// This is TRUNCATE in sqlite
	s.database.GetEntityManager().GetGormORM().Exec("DELETE FROM caches;")

	var (
		emails      []email.Cache
		batchInsert []email.Cache
		total       int64
	)
	s.database.GetMysqlEntityManager().GetGormORM().Raw(
		"SELECT COUNT(*) FROM (" + s.config.GetMysqlQuery()[:strings.LastIndex(s.config.GetMysqlQuery(), ";")] + ") AS all_users;",
	).Row().Scan(&total)

	bar := progressbar.Default(total)
	s.database.GetMysqlEntityManager().GetGormORM().Raw(
		s.config.GetMysqlQuery(),
	).FindInBatches(&emails, 100, func(tx *gorm.DB, batch int) error {
		for _, mail := range emails {
			bar.Add(1)
			mail.Name = strings.TrimSpace(mail.Name)
			mail.Email = strings.ToLower(strings.TrimSpace(mail.Email))

			if email.IsEmail(&mail.Email) {
				batchInsert = append(batchInsert, mail)
				logger.LogNormal("Create cache " + mail.Email + " " + mail.Name + ".")
			}
		}

		return nil
	})

	if len(batchInsert) > 0 {
		s.database.GetEntityManager().GetGormORM().CreateInBatches(batchInsert, 100)
	}

	logger.Success("Cache building finished successfully! Total email in cache: " + strconv.Itoa(int(total)))
}

func (s *sync) getDifference() []email.Cache {
	var cache []email.Cache
	tx := s.database.GetEntityManager().GetGormORM().Raw(`SELECT email, name FROM emails WHERE (email, name) NOT IN (SELECT email, name FROM caches)
    UNION ALL SELECT email, name FROM caches WHERE (email, name) NOT IN (SELECT email, name FROM emails);`).Scan(&cache)
	if tx.Error != nil {
		logger.Fatal(tx.Error)
	}

	return cache
}
