package admin

import (
	"FanCode/controller/admin"
	"github.com/gin-gonic/gin"
)

func SetupProblemRoutes(r *gin.Engine, problemController admin.ProblemManagementController) {
	//题目相关路由
	problem := r.Group("/manage/problem")
	{
		problem.GET("/code/check/:number", problemController.CheckProblemNumber)
		problem.POST("", problemController.InsertProblem)
		problem.PUT("", problemController.UpdateProblem)
		problem.DELETE("/:id", problemController.DeleteProblem)
		problem.GET("/list", problemController.GetProblemList)
		problem.GET("/:id", problemController.GetProblemByID)
		problem.POST("/enable", problemController.UpdateProblemEnable)
	}
}
