package main

import (
	issuePresentation "issue-service-aoroa/issue/presentation"
	issueInfra "issue-service-aoroa/issue/infrastructure"
	issueApp "issue-service-aoroa/issue/application"
	userInfra "issue-service-aoroa/user/infrastructure"

	"github.com/gin-gonic/gin"
)

func main() {
	userRepo := userInfra.NewUserRepository()
	issueRepo := issueInfra.NewIssueRepository()
	issueService := issueApp.NewIssueService(issueRepo, userRepo)
	issueController := issuePresentation.NewIssueController(issueService)

	router := gin.Default()

	router.POST("/issue", issueController.CreateIssue)
	router.GET("/issues", issueController.GetIssues)
	router.GET("/issue/:id", issueController.GetIssueByID)
	router.PATCH("/issue/:id", issueController.UpdateIssue)

	router.Run(":8080")
}

