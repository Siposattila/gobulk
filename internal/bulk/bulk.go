package bulk

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Siposattila/gobulk/internal/email"
	"github.com/Siposattila/gobulk/internal/interfaces"
	"github.com/Siposattila/gobulk/internal/kill"
	"github.com/Siposattila/gobulk/internal/logger"
	"github.com/schollz/progressbar/v3"
	"gorm.io/gorm"
)

type bulk struct {
	emailClient interfaces.EmailClientInterface
	database    interfaces.DatabaseInterface
	config      interfaces.ConfigInterface
}

func Init(database interfaces.DatabaseInterface, config interfaces.ConfigInterface) interfaces.BulkInterface {
	return &bulk{
		emailClient: email.NewClient(
			config.GetEmailDSN(),
			nil,
		),
		database: database,
		config:   config,
	}
}

func (b *bulk) StartConsole() {
	b.emailClient.SetEmailBody(email.NewBodyConsole(b.config.GetCompanyName(), b.config.GetUnsubscribeEndpoint()))
	b.bulkSend()
}

func (b *bulk) Start(subject string, greeting string, message string, farewell string) {
	b.emailClient.SetEmailBody(email.NewBody(subject, greeting, message, farewell, b.config.GetCompanyName(), b.config.GetUnsubscribeEndpoint()))
	b.bulkSend()
}

func (b *bulk) bulkSend() {
	b.checkStatusOfBulk()
	logger.Normal("Bulk email sending is starting now. This may take a long time!")

	var (
		emails []email.Email
		total  int64
	)
	b.database.GetEntityManager().GetGormORM().Find(
		&email.Email{},
		"valid = ? AND status = ? AND send_status = ?",
		interfaces.EMAIL_VALID,
		interfaces.EMAIL_STATUS_ACTIVE,
		interfaces.EMAIL_SEND_STATUS_NOT_SENT,
	).Count(&total)
	logger.Normal("Sending " + strconv.Itoa(int(total)) + " emails.")
	bar := progressbar.Default(total)

	b.database.GetEntityManager().GetGormORM().Where(
		"valid = ? AND status = ? AND send_status = ?",
		interfaces.EMAIL_VALID,
		interfaces.EMAIL_STATUS_ACTIVE,
		interfaces.EMAIL_SEND_STATUS_NOT_SENT,
	).FindInBatches(&emails, (60*1000/int(b.config.GetSendDelay()))*2, func(tx *gorm.DB, batch int) error {
		for _, mail := range emails {
			bar.Add(1)
			select {
			case <-kill.KillCtx.Done():
				logger.Warning("Shutdown signal received shutting down bulk sending process.")

				return errors.New("Shutdown")
			default:
				time.Sleep(time.Duration(b.config.GetSendDelay()) * time.Millisecond)
				b.emailClient.Send(&mail)
				b.database.GetEntityManager().GetGormORM().Save(mail)
			}
		}

		return nil
	})

	logger.Success("Bulk email sending is done!")
}

func (b *bulk) checkStatusOfBulk() {
	var totalNotSent int64
	b.database.GetEntityManager().GetGormORM().Find(
		&email.Email{},
		"valid = ? AND status = ? AND send_status = ?",
		interfaces.EMAIL_VALID,
		interfaces.EMAIL_STATUS_ACTIVE,
		interfaces.EMAIL_SEND_STATUS_NOT_SENT,
	).Count(&totalNotSent)

	if totalNotSent != 0 {
		reader := bufio.NewReader(os.Stdin)

		fmt.Println("The last bulk sending session was interrupted do you want to continue? [y/n]")
		for {
			answer, err := reader.ReadString('\n')
			if err != nil {
				logger.Fatal(err)
			}
			answer = strings.TrimSpace(answer)

			if answer == "y" || answer == "n" {
				if answer == "n" {
					b.database.GetEntityManager().GetGormORM().Model(email.Email{}).Where("1=1").Updates(email.Email{SendStatus: interfaces.EMAIL_SEND_STATUS_NOT_SENT})
				}
				break
			}
		}
	}

	if totalNotSent == 0 {
		logger.Debug(totalNotSent)
		b.database.GetEntityManager().GetGormORM().Model(email.Email{}).Where("1=1").Updates(email.Email{SendStatus: interfaces.EMAIL_SEND_STATUS_NOT_SENT})
	}
}
