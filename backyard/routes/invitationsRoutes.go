package routes

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/hospedate/backyard/controllers"
    "github.com/hospedate/backyard/log"
    "github.com/hospedate/backyard/models"
)

func InvitationsRoutes(router *gin.Engine, controller controllers.InvitationsController) {

    logger = log.GetOrCreateLogger("RoutesLogger", "INFO")

    router.GET("/invitations/:id", func(ctx *gin.Context) {
        id := ctx.Param("id")
        invitationId, err := models.NewInvitationIdFromStr(id)
        if err == models.ErrInvitationNotValid {
            ctx.JSON(http.StatusBadRequest, ResponseErr{ErrInvitationNotValid})
            return
        } else if err != nil {
            ctx.JSON(http.StatusBadRequest, ResponseErr{ErrBadRequest})
            return
        }
        invitation := controller.GetById(invitationId)
        if invitation == nil {
            ctx.JSON(http.StatusNotFound, ResponseErr{ErrNotFound})
            return
        }

        ctx.JSON(http.StatusOK, ResponseOk{Data: invitation})
    })

    router.GET("/invitations/search", func(ctx *gin.Context) {
        queryParams := ctx.Request.URL.Query()
        searchParams, err := models.NewInvitationsSearchParams(queryParams)
        if err != nil {
            ctx.JSON(http.StatusBadRequest, ResponseErr{NewErrBadRequest(err.Error())})
            return
        }

        invitations := controller.Search(searchParams)
        if invitations == nil {
            logger.Error("Internal server error when searching invitations")
            ctx.JSON(http.StatusInternalServerError, ResponseErr{ErrInternal})
            return
        }

        ctx.JSON(http.StatusOK, ResponseOk{Data: invitations})
    })

    router.POST("/invitations", func(ctx *gin.Context) {
        var newInvitationRequest models.NewInvitationRequest

        err := ctx.BindJSON(&newInvitationRequest)
        if err != nil {
            ctx.JSON(http.StatusBadRequest, ResponseErr{ErrBadRequest})
            return
        }

        invitation, err := controller.Create(&newInvitationRequest)

        if err != nil {
            logger.Info("Internal server error: ", err.Error())
            ctx.JSON(http.StatusInternalServerError, ResponseErr{ErrInternal})
        } else {
            ctx.JSON(http.StatusCreated, ResponseOk{Data: invitation})
        }
    })
}
