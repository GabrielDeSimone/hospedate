package routes

import (
    "github.com/gin-gonic/gin"
    "github.com/hospedate/backyard/controllers"
    "github.com/hospedate/backyard/log"
    "github.com/hospedate/backyard/models"

    "net/http"
)

var logger log.Logger

func UsersRoutes(router *gin.Engine, controller controllers.UsersController, invitationsController controllers.InvitationsController) {

    logger = log.GetOrCreateLogger("RoutesLogger", "INFO")

    router.GET("/users/:id", func(ctx *gin.Context) {
        id := ctx.Param("id")
        user := controller.GetById(id)

        if user == nil {
            ctx.JSON(http.StatusNotFound, ResponseErr{ErrNotFound})
            return
        }

        ctx.JSON(http.StatusOK, ResponseOk{Data: user})
    })

    router.GET("/users/search", func(ctx *gin.Context) {
        email := ctx.Query("email")
        password := ctx.Query("password")

        if email == "" || password == "" {
            ctx.JSON(http.StatusBadRequest, ResponseErr{NewErrBadRequest("You need to provide email and password")})
            return
        }

        user := controller.Get(email, password)
        if user == nil {
            ctx.JSON(http.StatusOK, ResponseOk{Data: []*models.User{}})
            return
        }

        ctx.JSON(http.StatusOK, ResponseOk{Data: []*models.User{user}})
    })

    router.POST("/users", func(ctx *gin.Context) {
        var newUserRequest models.NewUserRequest

        err := ctx.BindJSON(&newUserRequest)
        if err == models.ErrInvitationNotValid {
            ctx.JSON(http.StatusBadRequest, ResponseErr{ErrInvitationNotValid})
            return
        } else if err != nil {
            ctx.JSON(http.StatusBadRequest, ResponseErr{ErrBadRequest})
            return
        }

        user, err := controller.Create(&newUserRequest, invitationsController)

        if err != nil {
            if err == controllers.ErrDuplicateKey {
                ctx.JSON(http.StatusBadRequest, ResponseErr{ErrEmailOrPhoneAlreadyTaken})
            } else if err == controllers.ErrInvitationDoesNotExist {
                ctx.JSON(http.StatusBadRequest, ResponseErr{ErrInvitationDoesNotExist})
            } else if err == controllers.ErrInvitationAlreadyUsed {
                ctx.JSON(http.StatusBadRequest, ResponseErr{ErrInvitationAlreadyUsed})
            } else {
                logger.Info("Internal server error: ", err.Error())
                ctx.JSON(http.StatusInternalServerError, ResponseErr{ErrInternal})
            }
        } else {
            ctx.JSON(http.StatusCreated, ResponseOk{Data: user})
        }
    })

    router.PUT("/users/:id", func(ctx *gin.Context) {
        id := ctx.Param("id")
        user := controller.GetById(id)

        if user == nil {
            ctx.JSON(http.StatusNotFound, ResponseErr{ErrNotFound})
            return
        }

        var userEditRequest models.UserEditRequest

        err := ctx.BindJSON(&userEditRequest)
        if err != nil {
            ctx.JSON(http.StatusBadRequest, ResponseErr{NewErrBadRequest(err.Error())})
            return
        }

        // We assign the id from the URL to the struct
        // because we don't actually expect users to include this attribute in the body request
        userEditRequest.Id = id

        userEdited, err := controller.Edit(userEditRequest)

        if err != nil {
            logger.Error("Internal server error: ", err.Error())
            ctx.JSON(http.StatusInternalServerError, ResponseErr{ErrInternal})
            return
        } else {
            ctx.JSON(http.StatusOK, ResponseOk{Data: userEdited})
        }
    })

    router.GET("/users/:id/balance", func(ctx *gin.Context) {
        id := ctx.Param("id")
        balance, err := controller.GetBalance(id)

        if err != nil {
            if err == controllers.ErrUserDoesNotExist {
                ctx.JSON(http.StatusBadRequest, ResponseErr{ErrNotFound})
            } else {
                logger.Info("Internal server error: ", err.Error())
                ctx.JSON(http.StatusInternalServerError, ResponseErr{ErrInternal})
            }
        } else {
            ctx.JSON(http.StatusOK, ResponseOk{Data: balance})
        }
    })

    router.GET("/users/:id/withdrawals", func(ctx *gin.Context) {
        id := ctx.Param("id")
        withdrawals, err := controller.GetWithdrawals(id)
        if err != nil {
            if err == controllers.ErrUserDoesNotExist {
                ctx.JSON(http.StatusBadRequest, ResponseErr{ErrNotFound})
            } else {
                logger.Info("Internal server error: ", err.Error())
                ctx.JSON(http.StatusInternalServerError, ResponseErr{ErrInternal})
            }
        } else {
            ctx.JSON(http.StatusOK, ResponseOk{Data: withdrawals})
        }
    })

    router.POST("/users/:id/withdrawals", func(ctx *gin.Context) {
        var newWithdrawalRequest models.NewWithdrawalRequest
        user_id := ctx.Param("id")

        err := ctx.BindJSON(&newWithdrawalRequest)
        if err != nil {
            ctx.JSON(http.StatusBadRequest, ResponseErr{ErrBadRequest})
            return
        }
        newWithdrawalRequest.UserId = user_id

        withdrawal, err := controller.CreateWithdrawal(&newWithdrawalRequest)

        if err != nil {
            if err == controllers.ErrUserDoesNotExist {
                ctx.JSON(http.StatusBadRequest, ResponseErr{NewErrBadRequest("The user id provided does not exist")})
            } else if err == controllers.ErrBalanceNotEnough {
                ctx.JSON(http.StatusBadRequest, ResponseErr{ErrBalanceNotEnough})
            } else {
                logger.Error("Internal server error: ", err.Error())
                ctx.JSON(http.StatusInternalServerError, ResponseErr{ErrInternal})
            }
        } else {
            ctx.JSON(http.StatusCreated, ResponseOk{Data: withdrawal})
        }
    })

    router.PUT("/users/withdrawals/:id", func(ctx *gin.Context) {
        id := ctx.Param("id")
        withdraw := controller.GetWithdrawalById(id)

        if withdraw == nil {
            ctx.JSON(http.StatusNotFound, ResponseErr{ErrNotFound})
            return
        }

        var withdrawalEditRequest models.WithdrawalEditRequest

        err := ctx.BindJSON(&withdrawalEditRequest)
        if err != nil {
            ctx.JSON(http.StatusBadRequest, ResponseErr{ErrBadRequest})
            return
        }

        // We assign the id from the URL to the struct
        // because we don't actually expect users to include this attribute in the body request
        withdrawalEditRequest.Id = id

        withdrawal_edited, err := controller.EditWithdrawal(withdrawalEditRequest)

        if err != nil {
            logger.Error("Internal server error: ", err.Error())
            ctx.JSON(http.StatusInternalServerError, ResponseErr{ErrInternal})
            return
        } else {
            ctx.JSON(http.StatusOK, ResponseOk{Data: withdrawal_edited})
        }
    })

    router.GET("/users/:id/earnings", func(ctx *gin.Context) {
        id := ctx.Param("id")
        earnings, err := controller.GetOwnerEarnedInstances(id)
        if err != nil {
            if err == controllers.ErrUserDoesNotExist {
                ctx.JSON(http.StatusBadRequest, ResponseErr{ErrNotFound})
            } else {
                logger.Info("Internal server error: ", err.Error())
                ctx.JSON(http.StatusInternalServerError, ResponseErr{ErrInternal})
            }
        } else {
            ctx.JSON(http.StatusOK, ResponseOk{Data: earnings})
        }
    })

    router.GET("/users/:id/credit", func(ctx *gin.Context) {
        id := ctx.Param("id")
        credits, err := controller.GetUserCreditInstances(id)
        if err != nil {
            if err == controllers.ErrUserDoesNotExist {
                ctx.JSON(http.StatusBadRequest, ResponseErr{ErrNotFound})
            } else {
                logger.Info("Internal server error: ", err.Error())
                ctx.JSON(http.StatusInternalServerError, ResponseErr{ErrInternal})
            }
        } else {
            ctx.JSON(http.StatusOK, ResponseOk{Data: credits})
        }
    })

    router.POST("/users/:id/credit", func(ctx *gin.Context) {
        var newUserCreditInstanceRequest models.NewUserCreditInstanceRequest
        user_id := ctx.Param("id")

        err := ctx.BindJSON(&newUserCreditInstanceRequest)
        if err != nil {
            ctx.JSON(http.StatusBadRequest, ResponseErr{ErrBadRequest})
            return
        }
        newUserCreditInstanceRequest.UserId = user_id

        userCreditInstance, err := controller.SendCreditToUser(&newUserCreditInstanceRequest)

        if err != nil {
            if err == controllers.ErrUserDoesNotExist {
                ctx.JSON(http.StatusBadRequest, ResponseErr{NewErrBadRequest("The user id provided does not exist")})
            } else {
                logger.Error("Internal server error: ", err.Error())
                ctx.JSON(http.StatusInternalServerError, ResponseErr{ErrInternal})
            }
        } else {
            ctx.JSON(http.StatusCreated, ResponseOk{Data: userCreditInstance})
        }
    })
}
