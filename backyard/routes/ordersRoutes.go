package routes

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/hospedate/backyard/controllers"
    "github.com/hospedate/backyard/log"
    "github.com/hospedate/backyard/models"
)

func OrdersRoutes(router *gin.Engine, controller controllers.OrdersController) {

    logger = log.GetOrCreateLogger("RoutesLogger", "INFO")

    router.GET("/orders/:id", func(ctx *gin.Context) {
        id := ctx.Param("id")
        order := controller.GetById(id)

        if order == nil {
            ctx.JSON(http.StatusNotFound, ResponseErr{ErrNotFound})
            return
        }

        ctx.JSON(http.StatusOK, ResponseOk{Data: order})
    })

    router.GET("/orders/search", func(ctx *gin.Context) {
        queryParams := ctx.Request.URL.Query()
        searchParams, err := models.NewOrdersSearchParams(queryParams)
        if err != nil {
            ctx.JSON(http.StatusBadRequest, ResponseErr{ErrBadRequest})
            return
        }

        orders := controller.Search(searchParams)
        if orders == nil {
            logger.Error("Internal server error searching orders ")
            ctx.JSON(http.StatusInternalServerError, ResponseErr{ErrInternal})
            return
        }

        ctx.JSON(http.StatusOK, ResponseOk{Data: orders})
    })

    router.POST("/orders", func(ctx *gin.Context) {
        var newOrderRequest models.NewOrderRequest

        err := ctx.BindJSON(&newOrderRequest)
        if err != nil {
            ctx.JSON(http.StatusBadRequest, ResponseErr{ErrBadRequest})
            return
        } else if !newOrderRequest.DateEnd.After(&newOrderRequest.DateStart) {
            ctx.JSON(http.StatusBadRequest, ResponseErr{NewErrBadRequest("Dates provided are not valid!")})
            return
        }
        order, err := controller.Create(&newOrderRequest)

        if err != nil {
            if err == controllers.ErrUserDoesNotExist {
                ctx.JSON(http.StatusBadRequest, ResponseErr{NewErrBadRequest("The user id provided does not exist")})
            } else if err == controllers.ErrPropertyDoesNotExist {
                ctx.JSON(http.StatusBadRequest, ResponseErr{NewErrBadRequest("The property id provided does not exist")})
            } else if err == controllers.ErrCollision {
                ctx.JSON(http.StatusBadRequest, ResponseErr{NewErrBadRequest("The property is blocked!")})
            } else if err == controllers.ErrOwnerBookingProperty {
                ctx.JSON(http.StatusBadRequest, ResponseErr{ErrOwnerBookingProperty})
            } else if err == controllers.ErrForTooManyGuests {
                ctx.JSON(http.StatusBadRequest, ResponseErr{ErrGuestsNumberExceeded})
            } else if err == controllers.ErrPropertyArchivedStatus {
                ctx.JSON(http.StatusBadRequest, ResponseErr{ErrPropertyArchivedStatus})
            } else {
                logger.Error("Internal server error: ", err.Error())
                ctx.JSON(http.StatusInternalServerError, ResponseErr{ErrInternal})
            }
        } else {
            ctx.JSON(http.StatusCreated, ResponseOk{Data: order})
        }
    })

    router.PUT("/orders/:id", func(ctx *gin.Context) {
        id := ctx.Param("id")
        order := controller.GetById(id)

        if order == nil {
            ctx.JSON(http.StatusNotFound, ResponseErr{ErrNotFound})
            return
        }

        var orderEditRequest models.OrderEditRequest

        err := ctx.BindJSON(&orderEditRequest)
        if err != nil {
            ctx.JSON(http.StatusBadRequest, ResponseErr{ErrBadRequest})
            return
        }
        if !isValidOrderEditRequest(orderEditRequest) {
            ctx.JSON(http.StatusBadRequest, ResponseErr{ErrNoValuesToUpdate})
            return
        }

        // We assign the id from the URL to the struct
        // because we don't actually expect users to include this attribute in the body request
        orderEditRequest.Id = id

        order_edited, err := controller.Edit(orderEditRequest)

        if err != nil {
            logger.Error("Internal server error: ", err.Error())
            ctx.JSON(http.StatusInternalServerError, ResponseErr{ErrInternal})
            return
        } else {
            ctx.JSON(http.StatusOK, ResponseOk{Data: order_edited})
        }
    })

    router.DELETE("/orders/:id", func(ctx *gin.Context) {
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

func isValidOrderEditRequest(editRequest models.OrderEditRequest) bool {
    if (editRequest.CanceledBy != nil) || (editRequest.Status != nil) {
        return true
    } else {
        return false
    }
}
