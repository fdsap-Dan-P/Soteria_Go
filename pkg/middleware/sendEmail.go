package middleware

import (
	"fmt"
	"net/smtp"
	"soteria_go/pkg/models/response"
	"strconv"
)

func SendMail(receiver_name, receiver_email, subject, body, username, instiCode, appCode, moduleName, methodUsed, endpoint string, requestBodyBytes, responseBodyBytes []byte) response.ReturnModel {
	funcName := "Send Email"

	// Email configuration
	smtpServer := "smtp.gmail.com"
	smtpPort := 587
	senderEmail := "roldan.polintang@cardmri.com"
	senderPassword := "oaqp hlpu doxs dmzy"
	// senderEmail := "dataplatform-support@fortress-asya.com"
	// senderPassword := "4828/iU7&^B7}4Cj"

	// Compose the email message
	to := []string{receiver_email}
	message := "Subject: " + subject + "\r\n" +
		"To: " + receiver_email + "\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"\r\n" + body

	// Connect to the SMTP server
	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpServer)
	sendEmailErr := smtp.SendMail(smtpServer+":"+strconv.Itoa(smtpPort), auth, senderEmail, to, []byte(message))
	if sendEmailErr != nil {
		returnMessage := ResponseData(username, instiCode, appCode, moduleName, funcName, "315", methodUsed, endpoint, requestBodyBytes, []byte(""), "", sendEmailErr)
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}

	ActivityLogger(username, instiCode, appCode, moduleName, funcName, "200", methodUsed, endpoint, []byte(""), []byte(""), "Successful", "", nil)
	return response.ReturnModel{
		RetCode: "200",
		Message: "Successful",
		Data: response.DataModel{
			Message:   fmt.Sprintf("Mail Sent to: %s", receiver_name),
			IsSuccess: true,
			Error:     nil,
		},
	}
}
