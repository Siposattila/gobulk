package bulk

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Siposattila/gobulk/internal/config"
	"github.com/Siposattila/gobulk/internal/console"
	"github.com/Siposattila/gobulk/internal/email"
	"github.com/Siposattila/gobulk/internal/gorm"
	g "gorm.io/gorm"
)

type Bulk struct {
	Config      *config.Config
	EM          *gorm.EntityManager
	EmailClient email.ClientInterface
}

func Init() *Bulk {
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

	bulk := &Bulk{
		EM: gorm.Gorm(),
	}
	bulk.Config = gorm.GetConfig(bulk.EM.GormORM)

	body := email.EmailBody{
		Subject:     &subject,
		Greeting:    &greeting,
		Message:     &message,
		Farewell:    &farewell,
		Company:     &bulk.Config.CompanyName,
		Unsubscribe: &bulk.Config.UnsubscribeEndpoint,
	}
	bulk.EmailClient = email.NewClient(&bulk.Config.EmailDSN, &body)

	return bulk
}

func InitForServer() *Bulk {
	bulk := &Bulk{
		EM: gorm.Gorm(),
	}
	bulk.Config = gorm.GetConfig(bulk.EM.GormORM)

	return bulk
}

func (b *Bulk) Start() {
	console.Normal("Bulk email sending is starting now. This may take a long time!!!")

	last := b.getLast()
	offset := 0
	if last != nil {
		offset = int(last.EmailID) - 1
	}

	var results []email.Email
	b.EM.GormORM.Where("valid = ? AND status = ?", email.EMAIL_VALID, email.EMAIL_STATUS_ACTIVE).Offset(offset).FindInBatches(&results, (60*1000/int(b.Config.SendDelay))*2, func(tx *g.DB, batch int) error {
		for _, result := range results {
			select {
			case <-email.ShutdownChan:
				last := email.NewLast(&result.ID, email.LAST_PROCESS_SEND)
				b.EM.GormORM.Create(last)
				console.Warning("Unexpected shutdown while sending emails. Saving last progress...")

				os.Exit(1)
			default:
				time.Sleep(time.Duration(b.Config.SendDelay) * time.Millisecond)
				b.EmailClient.Send(&result)
			}
		}

		// Returning an error will stop further batch processing
		return nil
	})
	console.Success("Bulk email sending is done!")
}

func (b *Bulk) getLast() *email.Last {
	var last email.Last
	tx := b.EM.GormORM.First(&last, "process_id = ?", email.LAST_PROCESS_SEND)

	if tx.Error == nil {
		b.EM.GormORM.Delete(last)
	}

	return &last
}
