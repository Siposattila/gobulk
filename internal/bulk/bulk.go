package bulk

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Siposattila/gobulk/internal/console"
	"github.com/Siposattila/gobulk/internal/email"
	"github.com/Siposattila/gobulk/internal/interfaces"
	"github.com/Siposattila/gobulk/internal/kill"
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

	last := b.getLast()
	offset := 0
	if last != nil {
		offset = int(last.Offset)
	}

	var results []email.Email
	b.database.GetEntityManager().GetGormORM().Where(
		"valid = ? AND status = ?",
		email.EMAIL_VALID,
		email.EMAIL_STATUS_ACTIVE,
	).Offset(offset).FindInBatches(&results, (60*1000/int(b.config.GetSendDelay()))*2, func(tx *gorm.DB, batch int) error {
		for _, result := range results {
			select {
			case <-kill.KillCtx.Done():
				last := email.NewLast(int64(offset), email.LAST_PROCESS_SEND)
				b.database.GetEntityManager().GetGormORM().Create(last)
				console.Warning("Unexpected shutdown while sending emails. Saving last progress...")
			default:
				offset += 1
				time.Sleep(time.Duration(b.config.GetSendDelay()) * time.Millisecond)
				b.emailClient.Send(&result)
			}
		}

		// Returning an error will stop further batch processing
		return nil
	})
	console.Success("Bulk email sending is done!")
}

// FIXME: code dup (validate.go)
func (b *Bulk) getLast() *email.Last {
	var last email.Last
	tx := b.database.GetEntityManager().GetGormORM().First(&last, "process_id = ?", email.LAST_PROCESS_SEND)

	if tx.Error == nil {
		b.database.GetEntityManager().GetGormORM().Delete(last)
	}

	return &last
}
