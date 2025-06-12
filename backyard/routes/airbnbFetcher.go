package routes

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/hospedate/backyard/models"
    "github.com/hospedate/backyard/services"
)

func AirbnbFetcherRoutes(router *gin.Engine, service services.AirbnbFetcher) {
    router.GET("/airbnbFetcher/status", func(ctx *gin.Context) {
        ctx.JSON(
            http.StatusOK,
            ResponseOk{Data: models.NewAirbnbFetcherStatus{
                Status: service.GetStatus(),
            }},
        )
    })

    router.POST("/airbnbFetcher/status", func(ctx *gin.Context) {
        var newStatusRequest models.NewAirbnbFetcherStatus
        err := ctx.BindJSON(&newStatusRequest)

        if err != nil {
            ctx.JSON(http.StatusBadRequest, ResponseErr{NewErrBadRequest(err.Error())})
            return
        }

        if newStatusRequest.Status == services.STATUS_STOPPED {
            service.Stop()
            ctx.JSON(
                http.StatusOK,
                ResponseOk{Data: models.NewAirbnbFetcherStatus{
                    Status: services.STATUS_STOPPED,
                }},
            )
        } else if newStatusRequest.Status == services.STATUS_STARTED {
            service.Start()
            ctx.JSON(
                http.StatusOK,
                ResponseOk{Data: models.NewAirbnbFetcherStatus{
                    Status: services.STATUS_STARTED,
                }},
            )
        } else {
            ctx.JSON(http.StatusBadRequest, ResponseErr{ErrBadRequest})
        }
    })
}