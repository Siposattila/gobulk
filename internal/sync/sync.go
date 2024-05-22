package sync

import (
	"strings"
	"time"

	"github.com/Siposattila/gobulk/internal/console"
	"github.com/Siposattila/gobulk/internal/email"
	"github.com/Siposattila/gobulk/internal/interfaces"
	"github.com/Siposattila/gobulk/internal/kill"
	"github.com/robfig/cron/v3"
	"github.com/schollz/progressbar/v3"
	"gorm.io/gorm"
)

type Sync struct {
	cron     cron.Schedule
	config   interfaces.ConfigInterface
	database interfaces.DatabaseInterface
}

func Init(database interfaces.DatabaseInterface, config interfaces.ConfigInterface) *Sync {
	cronParser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	schedule, parseError := cronParser.Parse(config.GetSyncCron())
	if parseError != nil {
		console.Fatal("Your cron expression is invalid or an error occured: " + parseError.Error())
	}

	return &Sync{
		cron:     schedule,
		config:   config,
		database: database,
	}
}

func (s *Sync) Start() {
	s.syncProcess()
}

func (s *Sync) syncProcess() {
	s.sync()
	timeSignal := time.After(s.cron.Next(time.Now()).Sub(time.Now()))
	console.Success("Next sync will be at: " + s.cron.Next(time.Now()).String())

	select {
	case <-timeSignal:
		s.sync()
		s.syncProcess()
	case <-kill.KillCtx.Done():
		console.Warning("Shutdown signal received shutting down sync process.")
	}
}

func (s *Sync) sync() {
	console.Normal("Sync is started. This may take a long time!!!")
	s.cacheMysqlData()
	diff := s.getDifference()
	var (
		batchInsert []email.Email
		total       int64
	)
	s.database.GetEntityManager().GetGormORM().Find(&email.Cache{}).Count(&total)

	console.Warning("Synchronizing local database!")
	if len(diff) > 0 {
		bar := progressbar.Default(total)

		for _, d := range diff {
			bar.Add(1)
			var e email.Email
			tx := s.database.GetEntityManager().GetGormORM().First(&e, "name = ? AND email = ?", d.Name, d.Email)
			if tx.Error != nil {
				if tx.Error == gorm.ErrRecordNotFound {
					batchInsert = append(batchInsert, email.Email{Name: d.Name, Email: d.Email})
				} else {
					console.Error("Error executing query: %v", tx.Error)
				}
			} else {
				e.Status = email.EMAIL_STATUS_INACTIVE
				s.database.GetEntityManager().GetGormORM().Save(e)
			}
		}

		if len(batchInsert) > 0 {
			s.database.GetEntityManager().GetGormORM().CreateInBatches(batchInsert, 100)
		}
	} else {
		console.Success("Already up to date!")
	}

	console.Success("Sync finished successfully!")
}

func (s *Sync) cacheMysqlData() {
	console.Normal("Validating cache...")
	// TODO: will need a logic to know when to delete cache
	// for now this is okay

	console.Warning("Rebuilding cache!")
	// This is TRUNCATE in sqlite
	s.database.GetEntityManager().GetGormORM().Exec("DELETE FROM caches;")

	var (
		results     []email.Cache
		batchInsert []email.Cache
		total       int64
	)
	s.database.GetMysqlEntityManager().GetGormORM().Raw(
		"SELECT COUNT(*) FROM (" + s.config.GetMysqlQuery()[:strings.LastIndex(s.config.GetMysqlQuery(), ";")] + ") AS all_users;",
	).Row().Scan(&total)

	bar := progressbar.Default(total)
	s.database.GetMysqlEntityManager().GetGormORM().Raw(
		s.config.GetMysqlQuery(),
	).FindInBatches(&results, 100, func(tx *gorm.DB, batch int) error {
		for _, result := range results {
			bar.Add(1)
			result.Name = strings.TrimSpace(result.Name)
			result.Email = strings.ToLower(strings.TrimSpace(result.Email))

			if email.IsEmail(&result.Email) {
				batchInsert = append(batchInsert, result)
			}
		}

		return nil
	})

	if len(batchInsert) > 0 {
		s.database.GetEntityManager().GetGormORM().CreateInBatches(batchInsert, 100)
	}
}

func (s *Sync) getDifference() []email.Cache {
	var cache []email.Cache
	tx := s.database.GetEntityManager().GetGormORM().Raw(`SELECT email, name FROM emails WHERE (email, name) NOT IN (SELECT email, name FROM caches)
    UNION ALL SELECT email, name FROM caches WHERE (email, name) NOT IN (SELECT email, name FROM emails);`).Scan(&cache)
	if tx.Error != nil {
		console.Fatal(tx.Error)
	}

	return cache
}
