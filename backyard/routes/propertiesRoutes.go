package routes

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/hospedate/backyard/controllers"
    "github.com/hospedate/backyard/log"
    "github.com/hospedate/backyard/models"
)

func PropertiesRoutes(router *gin.Engine, controller controllers.PropertiesController) {

    logger = log.GetOrCreateLogger("RoutesLogger", "INFO")

    router.POST("/properties", func(ctx *gin.Context) {
        var newPropertyRequest models.NewPropertyRequest

        err := ctx.BindJSON(&newPropertyRequest)
        if err != nil {
            ctx.JSON(http.StatusBadRequest, ResponseErr{ErrBadRequest})
            return
        }

        property, err := controller.Create(&newPropertyRequest)

        if err != nil {
            if err == controllers.ErrDuplicateKey {
                ctx.JSON(http.StatusBadRequest, ResponseErr{ErrPropertyAlreadyTaken})
            } else if err == controllers.ErrUserDoesNotExist {
                ctx.JSON(http.StatusBadRequest, ResponseErr{NewErrBadRequest("The user id provided does not exist")})
            } else {
                logger.Error("Internal server error: ", err.Error())
                ctx.JSON(http.StatusInternalServerError, ResponseErr{ErrInternal})
            }
        } else {
            ctx.JSON(http.StatusCreated, ResponseOk{Data: property})
        }
    })

    router.GET("/properties/:id", func(ctx *gin.Context) {
        id := ctx.Param("id")
        property := controller.GetById(id)

        if property == nil {
            ctx.JSON(http.StatusNotFound, ResponseErr{ErrNotFound})
            return
        }

        ctx.JSON(http.StatusOK, ResponseOk{Data: property})
    })

    router.GET("/properties/search", func(ctx *gin.Context) {
        queryParams := ctx.Request.URL.Query()
        searchParams, err := models.NewPropSearchParams(queryParams)
        if err != nil {
            ctx.JSON(
                http.StatusBadRequest,
                ResponseErr{NewErrBadRequest(err.Error())},
            )
            return
        } else if searchParams.HasOnlyOneDate() {
            ctx.JSON(
                http.StatusBadRequest,
                ResponseErr{NewErrBadRequest(
                    "If a date is provided, both start and end are mandatory",
                )},
            )
            return
        } else if searchParams.HasDates() && (searchParams.GetDateStart().After(searchParams.GetDateEnd())) {
            ctx.JSON(
                http.StatusBadRequest,
                ResponseErr{NewErrBadRequest("Dates provided are not valid!")},
            )
            return
        }

        properties := controller.SearchProperties(searchParams)
        if properties == nil {
            logger.Error("Internal server error searching properties ")
            ctx.JSON(http.StatusInternalServerError, ResponseErr{ErrInternal})
            return
        }

        ctx.JSON(http.StatusOK, ResponseOk{Data: properties})
    })

    router.PUT("/properties/:id", func(ctx *gin.Context) {
        id := ctx.Param("id")

        var propertyEditRequest models.PropertyEditRequest

        err := ctx.BindJSON(&propertyEditRequest)
        if err != nil {
            ctx.JSON(http.StatusBadRequest, ResponseErr{NewErrBadRequest(err.Error())})
            return
        }
        // We assign the id from the URL to the struct
        // because we don't actually expect users to include this attribute in the body request
        propertyEditRequest.Id = id

        propertyEdited, err := controller.Edit(propertyEditRequest)

        if err != nil {
            if err == controllers.ErrPropertyArchivedStatus {
                ctx.JSON(http.StatusBadRequest, ResponseErr{ErrPropertyArchivedStatus})
            } else if err == controllers.ErrPropertyHasActiveOrders {
                ctx.JSON(http.StatusBadRequest, ResponseErr{ErrPropertyHasActiveOrders})
            } else if err == controllers.ErrPropertyDoesNotExist {
                ctx.JSON(http.StatusNotFound, ResponseErr{ErrNotFound})
            } else {
                logger.Error("Internal server error: ", err.Error())
                ctx.JSON(http.StatusInternalServerError, ResponseErr{ErrInternal})
            }
        } else {
            ctx.JSON(http.StatusOK, ResponseOk{Data: propertyEdited})
        }
    })

    router.POST("/properties/:id/blocks", func(ctx *gin.Context) {
        property_id := ctx.Param("id")

        var newBlockRequest models.NewBlockRequest
        err := ctx.BindJSON(&newBlockRequest)
        if err != nil {
            // validation failed
            ctx.JSON(http.StatusBadRequest, ResponseErr{ErrBadRequest})
            return
        }

        if !newBlockRequest.DateEnd.After(&newBlockRequest.DateStart) {
            ctx.JSON(http.StatusBadRequest, ResponseErr{NewErrBadRequest("Dates provided are not valid!")})
            return
        }

        // We assign the property_id from the URL to the struct
        //because we don't actually expect users to include this attribute in the body request
        newBlockRequest.PropertyId = property_id

        block, err := controller.CreateBlock(&newBlockRequest)

        if err != nil {
            if err == controllers.ErrCollision {
                ctx.JSON(http.StatusBadRequest, ResponseErr{ErrCollision})
            } else {
                logger.Error("Internal server error: ", err.Error())
                ctx.JSON(http.StatusInternalServerError, ResponseErr{ErrInternal})
            }
        } else {
            ctx.JSON(http.StatusCreated, ResponseOk{Data: block})
        }
    })
    router.GET("/properties/:id/blocks", func(ctx *gin.Context) {
        id := ctx.Param("id")
        blocks := controller.GetBlocksByPropertyId(id)

        if blocks == nil {
            logger.Error("Internal server error fetching property ")
            ctx.JSON(http.StatusInternalServerError, ResponseErr{ErrInternal})
            return
        }

        ctx.JSON(http.StatusOK, ResponseOk{Data: blocks})
    })

    router.DELETE("/properties/:id/blocks/:id_block", func(ctx *gin.Context) {
        id_property := ctx.Param("id")
        id_block := ctx.Param("id_block")
        rows, err := controller.DeleteBlockById(id_property, id_block)

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
