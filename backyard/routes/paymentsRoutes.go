package routes

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/hospedate/backyard/controllers"
    "github.com/hospedate/backyard/log"
    "github.com/hospedate/backyard/models"
)

func PaymentsRoutes(router *gin.Engine, controller controllers.PaymentsController) {

    logger = log.GetOrCreateLogger("RoutesLogger", "INFO")

    router.GET("/payments/:id", func(ctx *gin.Context) {
        id := ctx.Param("id")
        payment := controller.GetById(id)

        if payment == nil {
            ctx.JSON(http.StatusNotFound, ResponseErr{ErrNotFound})
            return
        }

        ctx.JSON(http.StatusOK, ResponseOk{Data: payment})
    })

    router.GET("/payments/search", func(ctx *gin.Context) {
        queryParams := ctx.Request.URL.Query()
        searchParams, err := models.NewPaymentsSearchParams(queryParams)
        if err != nil {
            ctx.JSON(http.StatusBadRequest, ResponseErr{ErrBadRequest})
            return
        }

        payments := controller.Search(searchParams)
        if payments == nil {
            logger.Error("Internal server error when searching payments")
            ctx.JSON(http.StatusInternalServerError, ResponseErr{ErrInternal})
            return
        }

        ctx.JSON(http.StatusOK, ResponseOk{Data: payments})
    })

    router.POST("/payments", func(ctx *gin.Context) {
        var newPaymentRequest models.NewPaymentRequest

        err := ctx.BindJSON(&newPaymentRequest)
        if err != nil {
            ctx.JSON(http.StatusBadRequest, ResponseErr{ErrBadRequest})
            return
        }
        payment, err := controller.Create(&newPaymentRequest)

        if err != nil {
            if err == controllers.ErrOrderDoesNotExist {
                ctx.JSON(http.StatusBadRequest, ResponseErr{NewErrBadRequest("The order id provided does not exist")})
            } else {
                logger.Error("Internal server error: ", err.Error())
                ctx.JSON(http.StatusInternalServerError, ResponseErr{ErrInternal})
            }
        } else {
            ctx.JSON(http.StatusCreated, ResponseOk{Data: payment})
        }
    })

    router.PUT("/payments/:id", func(ctx *gin.Context) {
        id := ctx.Param("id")
        payment := controller.GetById(id)

        if payment == nil {
            ctx.JSON(http.StatusNotFound, ResponseErr{ErrNotFound})
            return
        }
        paymentEditRequest, err := models.NewPaymentEditRequest(ctx)
        if err != nil {
            ctx.JSON(http.StatusBadRequest, ResponseErr{NewErrBadRequest(err.Error())})
            return
        }
        // We assign the id from the URL to the struct
        //because we don't actually expect users to include this attribute in the body request
        paymentEditRequest.Id = id

        paymentEdited, err := controller.Edit(*paymentEditRequest)

        if err != nil {
            logger.Error("Internal server error: ", err.Error())
            ctx.JSON(http.StatusInternalServerError, ResponseErr{ErrInternal})
            return
        } else {
            ctx.JSON(http.StatusOK, ResponseOk{Data: paymentEdited})
        }
    })

    router.DELETE("/payments/:id", func(ctx *gin.Context) {
        id := ctx.Param("id")
        rows, err := controller.DeleteById(id)

        if err != nil {
            logger.Error("Internal server error: ", err.Error())
            ctx.JSON(http.StatusInternalServerError, ResponseErr{ErrInternal})
            return
        } else {
            rows_deleted := DeletedRows{Nrows: NumberRowsDeleted(rows)}
            ctx.JSON(http.StatusOK, ResponseOk{Data: rows_deleted})
        }
    })

}
