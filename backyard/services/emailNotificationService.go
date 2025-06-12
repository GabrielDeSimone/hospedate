package services

import (
    "context"
    "fmt"
    "os/exec"
    "time"

    "github.com/hospedate/backyard/log"
    "github.com/hospedate/backyard/models"
    "github.com/hospedate/backyard/repositories"
)

type emailNotificationEvent struct {
    Subject          string
    SenderAddress    string
    RecipientAddress string
    HtmlBody         string
}

type EmailNotificationService interface {
    Start()
    SendUserHostApplicationNotification(newUserHostApplication models.UserHostApplication) error
    SendIExternalInvitationRequestNotification(newExternalInvitationRequest models.ExternalInvitationRequest) error
    SendOwnerFundsNotification(orderId string, amountCents uint, recipient string) error
    SendUserCreditNotificationSourceTraveler(amount float32, recipientID string, userInvitedName string) error
    SendUserCreditNotificationSourceOwner(amount float32, recipientID string, userInvitedName string) error
    SendOrderCanceledNotification(orderID string, recipient string) error
    SendPaymentReceivedNotification(orderID string, recipient string) error
    SendOrderPendingOwnerNotification(orderID string, recipientID string) error
    SendOrderPendingTravelerNotification(orderID string, recipientID string) error
    SendOrderConfirmedTravelerNotification(orderID string, recipientID string) error
}

type EmailNotificationServiceImp struct {
    emailChannel     chan emailNotificationEvent
    logger           log.Logger
    usersRepo        repositories.UsersRepository
    propertiesRepo   repositories.PropertiesRepository
    SenderAddress    string
    TimeOutEmailSend time.Duration
    disabled         bool
    pathPrefix       string
}

func NewEmailNotificationService(
    usersRepo repositories.UsersRepository,
    propertiesRepo repositories.PropertiesRepository,
    SenderAddress string,
    EMAILS_CHANNEL_CAPACITY int32,
    TimeOutEmailSend time.Duration,
    disabled bool,
    pathPrefix string,
) EmailNotificationService {
    logger := log.GetOrCreateLogger("EmailNotificationService", "INFO")
    return &EmailNotificationServiceImp{
        logger:           logger,
        emailChannel:     make(chan emailNotificationEvent, EMAILS_CHANNEL_CAPACITY),
        usersRepo:        usersRepo,
        propertiesRepo:   propertiesRepo,
        SenderAddress:    SenderAddress,
        TimeOutEmailSend: TimeOutEmailSend,
        disabled:         disabled,
        pathPrefix:       pathPrefix,
    }
}
func (es *EmailNotificationServiceImp) sendEmail(subject, body, recipientEmail string) {
    e := emailNotificationEvent{
        Subject:          subject,
        SenderAddress:    es.SenderAddress,
        RecipientAddress: recipientEmail,
        HtmlBody:         body,
    }
    es.emailChannel <- e
}

func (es *EmailNotificationServiceImp) SendUserHostApplicationNotification(newUserHostApplication models.UserHostApplication) error {
    user := es.usersRepo.GetById(newUserHostApplication.UserId)
    if user == nil {
        return ErrUserNotFound
    }
    body := es.emailTempTitleDesc(
        "Promoción a Anfitrión",
        fmt.Sprintf("El usuario con ID %s ha solicitado ser promovido a Anfitrión.", newUserHostApplication.UserId),
    )
    subject := fmt.Sprintf("Promoción a Anfitrión ID %s ", newUserHostApplication.UserId)
    es.sendEmail(subject, body, es.SenderAddress)
    return nil
}

func (es *EmailNotificationServiceImp) SendIExternalInvitationRequestNotification(newExternalInvitationRequest models.ExternalInvitationRequest) error {
    body := es.emailTempTitleContent(
        "Nueva solicitud de invitación",
        fmt.Sprintf(
            "<div><p>Nombre: %v</p><p>Email: %v</p><p>Mensaje:</p><p>%v</p></div>",
            newExternalInvitationRequest.Name,
            newExternalInvitationRequest.Email,
            newExternalInvitationRequest.Body,
        ),
    )
    subject := fmt.Sprintf("Nueva solicitud de invitación de %v", newExternalInvitationRequest.Name)
    es.sendEmail(subject, body, es.SenderAddress)
    return nil
}

func (es *EmailNotificationServiceImp) SendOwnerFundsNotification(orderId string, amountCents uint, recipientID string) error {
    recipientUser := es.usersRepo.GetById(recipientID)
    if recipientUser == nil {
        return ErrUserNotFound
    }
    body := es.emailTempTitleDesc(
        "Hemos actualizado tu saldo disponible",
        fmt.Sprintf("Se han acreditado USDT $ %.2f a tu saldo disponible en concepto de la reserva protegida con ID %v.", float32(amountCents)/100, orderId),
    )
    es.sendEmail("Acreditación de fondos", body, recipientUser.Email)
    return nil
}

func (es *EmailNotificationServiceImp) SendUserCreditNotificationSourceTraveler(amount float32, recipientID string, userInvitedName string) error {
    recipientUser := es.usersRepo.GetById(recipientID)
    if recipientUser == nil {
        return ErrUserNotFound
    }
    body := es.emailTempTitleDesc(
        "Has recibido crédito 🎉",
        fmt.Sprintf("Gracias a que %v completó su primera orden protegida en Hospedate, se han acreditado USDT $ %.2f como crédito disponible a tu cuenta!", userInvitedName, amount),
    )
    es.sendEmail("Acreditación de fondos", body, recipientUser.Email)
    return nil
}

func (es *EmailNotificationServiceImp) SendUserCreditNotificationSourceOwner(amount float32, recipientID string, userInvitedName string) error {
    recipientUser := es.usersRepo.GetById(recipientID)
    if recipientUser == nil {
        return ErrUserNotFound
    }
    body := es.emailTempTitleDesc(
        "Has recibido crédito 🎉",
        fmt.Sprintf("Gracias a que %v verificó una propiedad por primera vez en Hospedate, se han acreditado USDT $ %.2f como crédito disponible a tu cuenta!", userInvitedName, amount),
    )
    es.sendEmail("Acreditación de fondos", body, recipientUser.Email)
    return nil
}

func (es *EmailNotificationServiceImp) SendOrderCanceledNotification(orderID, recipientID string) error {
    recipientUser := es.usersRepo.GetById(recipientID)
    if recipientUser == nil {
        return ErrUserNotFound
    }
    body := es.emailTempTitleDesc(
        "Reserva cancelada",
        fmt.Sprintf("La reserva con ID %s ha sido cancelada.", orderID),
    )
    es.sendEmail(fmt.Sprintf("Reserva cancelada %s", orderID), body, recipientUser.Email)
    return nil
}

func (es *EmailNotificationServiceImp) SendOrderPendingOwnerNotification(orderID, propertyID string) error {
    propertyRecipient := es.propertiesRepo.GetById(propertyID)
    if propertyRecipient == nil {
        return ErrUserNotFound
    }
    recipientUser := es.usersRepo.GetById(propertyRecipient.UserId)
    if recipientUser == nil {
        return ErrUserNotFound
    }
    body := es.emailTempTitleDesc(
        "Nueva reserva pendiente",
        fmt.Sprintf("Tenés una reserva pendiente con ID %s. Entrá a la plataforma para confirmarla! Será cancelada automáticamente si no se confirma en 5 horas.", orderID),
    )
    es.sendEmail(fmt.Sprintf("Reserva pendiente %s", orderID), body, recipientUser.Email)
    return nil
}

func (es *EmailNotificationServiceImp) SendOrderPendingTravelerNotification(orderID, recipientID string) error {
    recipientUser := es.usersRepo.GetById(recipientID)
    if recipientUser == nil {
        return ErrUserNotFound
    }
    body := es.emailTempTitleDesc(
        "Tu reserva está pendiente de confirmación",
        fmt.Sprintf("Tu reserva con ID %s está pendiente de confirmación por parte del anfitrión", orderID),
    )
    es.sendEmail(fmt.Sprintf("Reserva pendiente %s", orderID), body, recipientUser.Email)
    return nil
}

func (es *EmailNotificationServiceImp) SendOrderConfirmedTravelerNotification(orderID, recipientID string) error {
    recipientUser := es.usersRepo.GetById(recipientID)
    if recipientUser == nil {
        return ErrUserNotFound
    }
    body := es.emailTempTitleDesc(
        "Reserva confirmada!",
        fmt.Sprintf("Tu reserva con ID %s fue confirmada! Prepará las valijas!", orderID),
    )
    es.sendEmail(fmt.Sprintf("Reserva confirmada %s", orderID), body, recipientUser.Email)
    return nil
}

func (es *EmailNotificationServiceImp) SendPaymentReceivedNotification(orderID, recipientID string) error {
    recipientUser := es.usersRepo.GetById(recipientID)
    if recipientUser == nil {
        return ErrUserNotFound
    }
    body := es.emailTempTitleDesc(
        "Recibimos tu pago!",
        fmt.Sprintf("Recibimos tu pago asociado a la reserva con ID %s! Solo queda esperar la confirmación del anfitrión", orderID),
    )
    es.sendEmail(fmt.Sprintf("Pago recibido %s", orderID), body, recipientUser.Email)
    return nil
}

func (es *EmailNotificationServiceImp) Start() {
    if es.disabled {
        es.logger.Info("Skipping EmailNotificationService initialization (disabled = true)")
    } else {
        go es.mainSender() // mainSender runs in a separate go routine
        es.logger.Info("EmailNotificationService initialized")
    }
}

func (es *EmailNotificationServiceImp) mainSender() {
    for { // infinite loop
        select {
        case mail := <-es.emailChannel:
            es.logger.Infof("Received an email to Send Subject: \"%v\"", mail.Subject)
            es.sendNotification(mail)
        }
    }
}

func (es *EmailNotificationServiceImp) sendNotification(e emailNotificationEvent) {
    timeout := es.TimeOutEmailSend
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()
    cmdPath := fmt.Sprintf("%v/email-sender/main.py", es.pathPrefix)
    pythonArgs := []string{
        cmdPath,
        e.Subject,
        e.SenderAddress,
        e.RecipientAddress,
        e.HtmlBody,
    }
    cmd := exec.CommandContext(ctx, "python", pythonArgs...)

    err := cmd.Run()
    if err != nil {
        if ctx.Err() == context.DeadlineExceeded {
            es.logger.Error("Sending email notification timed out.")
            return
        }
        es.logger.Error("Error sending email notification:", err.Error())
        return
    }
}

func (es *EmailNotificationServiceImp) emailTempTitleContent(title string, content string) string {
    htmlContent := fmt.Sprintf(`<h1>%v</h1>%v`, title, content)
    return es.emailTemplateRaw(htmlContent)
}

func (es *EmailNotificationServiceImp) emailTempTitleDesc(title string, description string) string {
    htmlContent := fmt.Sprintf(`<h1>%v</h1><p>%v</p>`, title, description)
    return es.emailTemplateRaw(htmlContent)
}

func (es *EmailNotificationServiceImp) emailTemplateRaw(htmlContent string) string {
    template := `
<html lang="es">
  <head>
    <meta charset="utf-8">
    <link rel="preconnect" href="https://fonts.googleapis.com" />
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;700;800&display=swap" rel="stylesheet"/>
    <style>
      body {font-size:18px;background-color:white;font-family:Inter,sans-serif;}
      a {color:#2A8BD2}
      #content {max-width:640px;margin: 0 auto; border:1px solid #e4e4e4;border-radius:20px;padding:5px 30px}
      h1{font-size:32px;font-weight:500;}
      p{font-size:18px}
      #footer, #footer p{font-size:12px;}
      .bar{height:0;border-top: 1px solid #e4e4e4;margin: 35px 0;}
    </style>
  </head>
  <body>
    <div id="content">
      <div>
	<img src="https://hospedate.app/hospedatelogosmall.png" style="height:65px" />
      </div>
      <div>
	    %v
      </div>
      <div class="bar"></div>
      <div id="footer">
	<p>¡Estamos para ayudarte! Comunicate con nuestro equipo de soporte por email a <a href="mailto:info@hospedate.app">info@hospedate.app</a> o por nuestro WhatsApp al <a target="_blank" href="https://wa.me/17869406354">+1 (786) 940 6354</a></p>
      </div>
    </div>
  </body>
</html>
`
    htmlFinal := fmt.Sprintf(template, htmlContent)
    return htmlFinal
}
