package email

import (
	"fmt"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"go.uber.org/zap"
)

func sendEmail(mailHelper *mail.SGMailV3, logger *zap.Logger) error {

	request := sendgrid.GetRequest(os.Getenv("SENDGRID_API_KEY"), os.Getenv("SENDGRID_URL_PATH"), os.Getenv("SENDGRID_URL"))
	request.Method = "POST"
	var Body = mail.GetRequestBody(mailHelper)
	request.Body = Body
	response, err := sendgrid.API(request)
	if err != nil {
		logger.Error(err.Error())
	}
	logger.Info(response.Body)
	return nil
}

func VerificationEmail(email string, firstName string, secret string, logger *zap.Logger) error {
	emailURL := fmt.Sprintf("%s?ctx=%s&email=%s", os.Getenv("WEBSITE_URL"), secret, email)

	mailHelper := mail.NewV3Mail()

	fromEmailAddress := mail.NewEmail("Confirm Email", os.Getenv("SENDGRID_EMAIL"))
	receiversEmailAddress := mail.NewEmail("user's email", email)

	mailHelper.SetFrom(fromEmailAddress)
	mailHelper.SetTemplateID(os.Getenv("SENDGRID_TEMPLATE_ID"))

	personalize := mail.NewPersonalization()
	tos := []*mail.Email{receiversEmailAddress}
	personalize.AddTos(tos...)
	personalize.SetDynamicTemplateData("first_name", firstName)
	personalize.SetDynamicTemplateData("verify_url", emailURL)

	mailHelper.AddPersonalizations(personalize)

	err := sendEmail(mailHelper, logger)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}

func PasswordResetEmail(email string, secret string, logger *zap.Logger) error {

	emailURL := fmt.Sprintf("%s?ctx=%s&email=%s", os.Getenv("WEBSITE_URL"), secret, email)

	mailHelper := mail.NewV3Mail()

	fromEmailAddress := mail.NewEmail("Reset Password", os.Getenv("SENDGRID_EMAIL"))
	receiversEmailAddress := mail.NewEmail(email, email)

	mailHelper.SetFrom(fromEmailAddress)
	mailHelper.SetTemplateID(os.Getenv("SENDGRID_PASSWORD_TEMPLATE_ID"))

	personalize := mail.NewPersonalization()
	tos := []*mail.Email{receiversEmailAddress}
	personalize.AddTos(tos...)
	personalize.SetDynamicTemplateData("password_reset_url", emailURL)

	mailHelper.AddPersonalizations(personalize)

	err := sendEmail(mailHelper, logger)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}
