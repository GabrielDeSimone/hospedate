package routes

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/hospedate/backyard/controllers"
    "github.com/hospedate/backyard/models"
    "github.com/hospedate/backyard/services"
)

func NotificationServiceRoutes(router *gin.Engine, service services.EmailNotificationService, userController controllers.UsersController) {
    router.POST("/notificationService/userHostApplication", func(ctx *gin.Context) {
        var newUserHostApplication models.UserHostApplication

        err := ctx.BindJSON(&newUserHostApplication)
        if err != nil {
            ctx.JSON(http.StatusBadRequest, ResponseErr{ErrBadRequest})
            return
        }

        err = service.SendUserHostApplicationNotification(newUserHostApplication)

        if err != nil {
            if err == services.ErrUserNotFound {
                ctx.JSON(http.StatusBadRequest, ResponseErr{ErrUserNotFound})
            } else {
                logger.Error("Internal server error: ", err.Error())
                ctx.JSON(http.StatusInternalServerError, ResponseErr{ErrInternal})
            }
        } else {
            ctx.JSON(http.StatusOK, ResponseOk{Data: newUserHostApplication})
        }

    })

    router.POST("/notificationService/externalInvitationRequest", func(ctx *gin.Context) {
        var newExternalInvitationRequest models.ExternalInvitationRequest

        err := ctx.BindJSON(&newExternalInvitationRequest)
        if err != nil {
            ctx.JSON(http.StatusBadRequest, ResponseErr{ErrBadRequest})
            return
        }

        err = service.SendIExternalInvitationRequestNotification(newExternalInvitationRequest)
        if err != nil {
            ctx.JSON(http.StatusInternalServerError, ResponseErr{ErrInternal})
            return
        }
        ctx.JSON(http.StatusOK, ResponseOk{Data: newExternalInvitationRequest})
    })
}
