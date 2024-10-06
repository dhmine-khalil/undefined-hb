package routes

import (
	"habitat-server/utils"

	"github.com/kataras/iris/v12"
)

func TestMessageNotification(ctx iris.Context) {
	data := map[string]string{
		"url": "exp://192.168.100.3:19000/--/messages/2/TestNotification",
	}

	err := utils.SendNotification(
		"ExponentPushToken[Xxxxxxxxxxxxxxxxxxxxxx]",
		"Push Title", "Push body is this message", data)
	if err != nil {
		utils.CreateInternalServerError(ctx)
		return
	}

	ctx.JSON(iris.Map{
		"sent": true,
	})
}
