package providers

import (
	"fmt"
	"github.com/mailjet/mailjet-apiv3-go"
	"log"
	"yip/src/config"
)

type EmailProvider struct {
	client *mailjet.Client
	config *config.EmailConfig
}

func InitEmailProvider(config *config.EmailConfig) (EmailProvider, error) {
	//mailjetClient := NewMailjetClient(os.Getenv("MJ_APIKEY_PUBLIC"), os.Getenv("MJ_APIKEY_PRIVATE"))
	client := mailjet.NewMailjetClient(config.MailJet.PublicKey, config.MailJet.PrivateKey)
	return EmailProvider{
		client: client,
		config: config,
	}, nil
}

func (ep EmailProvider) SendPinMail(toEmail string, pin string) error {
	fmt.Println("Sending email to ", toEmail, "with Pin", pin)
	messagesInfo := []mailjet.InfoMessagesV31{
		mailjet.InfoMessagesV31{
			From: &mailjet.RecipientV31{
				Email: ep.config.SenderEmail,
				Name:  ep.config.SenderName,
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: toEmail,
					Name:  "Client",
				},
			},
			TemplateID:       ep.config.MailJet.TemplateId,
			TemplateLanguage: true,
			Subject:          "Your PIN",
			Variables: map[string]interface{}{
				"pin": pin,
			},
		},
	}

	messages := mailjet.MessagesV31{Info: messagesInfo}
	res, err := ep.client.SendMailV31(&messages)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	log.Println(res.ResultsV31)
	return nil
}
