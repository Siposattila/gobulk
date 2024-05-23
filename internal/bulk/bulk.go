package bulk

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Siposattila/gobulk/internal/console"
	"github.com/Siposattila/gobulk/internal/email"
	"github.com/Siposattila/gobulk/internal/interfaces"
	"github.com/Siposattila/gobulk/internal/kill"
	"github.com/schollz/progressbar/v3"
	"gorm.io/gorm"
)

type Bulk struct {
	emailClient email.ClientInterface
	database    interfaces.DatabaseInterface
	config      interfaces.ConfigInterface
}

func Init(database interfaces.DatabaseInterface, config interfaces.ConfigInterface) *Bulk {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Please give a subject:")
	subject, err := reader.ReadString('\n')
	if err != nil {
		console.Fatal(err)
	}
	subject = strings.TrimSpace(subject)

	fmt.Println("Please give a greeting:")
	greeting, err := reader.ReadString('\n')
	if err != nil {
		console.Fatal(err)
	}
	greeting = strings.TrimSpace(greeting)

	fmt.Println("Please give a content:")
	message, err := reader.ReadString('\n')
	if err != nil {
		console.Fatal(err)
	}
	message = strings.TrimSpace(message)

	fmt.Println("Please give a farewell:")
	farewell, err := reader.ReadString('\n')
	if err != nil {
		console.Fatal(err)
	}
	farewell = strings.TrimSpace(farewell)

	return &Bulk{
		emailClient: email.NewClient(
			config.GetEmailDSN(),
			&email.EmailBody{
				Subject:     subject,
				Greeting:    greeting,
				Message:     message,
				Farewell:    farewell,
				Company:     config.GetCompanyName(),
				Unsubscribe: config.GetUnsubscribeEndpoint(),
			},
		),
		database: database,
		config:   config,
	}
}

func InitForServer(database interfaces.DatabaseInterface, config interfaces.ConfigInterface) *Bulk {
	return &Bulk{
		database: database,
		config:   config,
	}
}

func (b *Bulk) Start() {
	console.Normal("Bulk email sending is starting now. This may take a long time!!!")

	var (
		emails []email.Email
		total  int64
	)
	b.database.GetEntityManager().GetGormORM().Find(
		&email.Email{},
		"valid = ? AND status = ? AND send_status = ?",
		email.EMAIL_VALID,
		email.EMAIL_STATUS_ACTIVE,
		email.EMAIL_SEND_STATUS_NOT_SENT,
	).Count(&total)
	bar := progressbar.Default(total)

	b.database.GetEntityManager().GetGormORM().Where(
		"valid = ? AND status = ?",
		email.EMAIL_VALID,
		email.EMAIL_STATUS_ACTIVE,
	).FindInBatches(&emails, (60*1000/int(b.config.GetSendDelay()))*2, func(tx *gorm.DB, batch int) error {
		for _, mail := range emails {
			bar.Add(1)
			select {
			case <-kill.KillCtx.Done():
				console.Warning("Unexpected shutdown while sending emails.")

				return errors.New("Shutdown")
			default:
				time.Sleep(time.Duration(b.config.GetSendDelay()) * time.Millisecond)
				b.emailClient.Send(&mail)
				b.database.GetEntityManager().GetGormORM().Save(mail)
			}
		}

		return nil
	})

	console.Success("Bulk email sending is done!")
}
