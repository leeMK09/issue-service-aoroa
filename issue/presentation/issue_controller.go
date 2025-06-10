package presentation

import (
	"net/http"
	"strconv"

	"issue-service-aoroa/issue/application"

	"github.com/gin-gonic/gin"
)

type IssueController struct {
	issueService application.IssueService
}

type CreateIssueRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	UserID      *uint  `json:"userId"`
}

type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

func NewIssueController(issueService application.IssueService) *IssueController {
	return &IssueController{
		issueService: issueService,
	}
}

func (c *IssueController) CreateIssue(ctx *gin.Context) {
	var req CreateIssueRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "잘못된 요청 데이터입니다",
			Code:  http.StatusBadRequest,
		})
		return
	}

	issue, err := c.issueService.CreateIssue(req.Title, req.Description, req.UserID)
	if err != nil {
		if err.Error() == "사용자를 찾을 수 없습니다" {
			ctx.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "사용자를 찾을 수 없습니다",
				Code:  http.StatusBadRequest,
			})
			return
		}
		if err.Error() == "제목은 필수입니다" {
			ctx.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "제목은 필수입니다",
				Code:  http.StatusBadRequest,
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "서버 내부 오류입니다",
			Code:  http.StatusInternalServerError,
		})
		return
	}

	ctx.JSON(http.StatusCreated, issue)
}

func (c *IssueController) GetIssues(ctx *gin.Context) {
	status := ctx.Query("status")
	
	if status != "" {
		issues, err := c.issueService.GetIssuesByStatus(status)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "유효하지 않은 상태입니다",
				Code:  http.StatusBadRequest,
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"issues": issues})
		return
	}

	issues := c.issueService.GetAllIssues()
	ctx.JSON(http.StatusOK, gin.H{"issues": issues})
}

func (c *IssueController) GetIssueByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "잘못된 ID 형식입니다",
			Code:  http.StatusBadRequest,
		})
		return
	}

	issue, err := c.issueService.GetIssueByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "서버 내부 오류입니다",
			Code:  http.StatusInternalServerError,
		})
		return
	}

	if issue == nil {
		ctx.JSON(http.StatusNotFound, ErrorResponse{
			Error: "이슈를 찾을 수 없습니다",
			Code:  http.StatusNotFound,
		})
		return
	}

	ctx.JSON(http.StatusOK, issue)
}

func (c *IssueController) UpdateIssue(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "잘못된 ID 형식입니다",
			Code:  http.StatusBadRequest,
		})
		return
	}

	var rawRequest map[string]interface{}
	if err := ctx.ShouldBindJSON(&rawRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "잘못된 요청 데이터입니다",
			Code:  http.StatusBadRequest,
		})
		return
	}

	updates := make(map[string]interface{})
	hasStatusUpdate := false
	
	if title, ok := rawRequest["title"]; ok {
		updates["title"] = title
	}
	if description, ok := rawRequest["description"]; ok {
		updates["description"] = description
	}
	if status, ok := rawRequest["status"]; ok {
		updates["status"] = status
		hasStatusUpdate = true
	}
	if userID, ok := rawRequest["userId"]; ok {
		updates["userId"] = userID
	}
	
	updates["status"] = hasStatusUpdate

	issue, err := c.issueService.UpdateIssue(uint(id), updates)
	if err != nil {
		switch err.Error() {
		case "이슈를 찾을 수 없습니다":
			ctx.JSON(http.StatusNotFound, ErrorResponse{
				Error: "이슈를 찾을 수 없습니다",
				Code:  http.StatusNotFound,
			})
		case "완료되거나 취소된 이슈는 수정할 수 없습니다":
			ctx.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "완료되거나 취소된 이슈는 수정할 수 없습니다",
				Code:  http.StatusBadRequest,
			})
		case "사용자를 찾을 수 없습니다":
			ctx.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "사용자를 찾을 수 없습니다",
				Code:  http.StatusBadRequest,
			})
		case "유효하지 않은 상태입니다":
			ctx.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "유효하지 않은 상태입니다",
				Code:  http.StatusBadRequest,
			})
		case "담당자 없이는 진행중 또는 완료 상태로 변경할 수 없습니다":
			ctx.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "담당자 없이는 진행중 또는 완료 상태로 변경할 수 없습니다",
				Code:  http.StatusBadRequest,
			})
		case "제목은 필수입니다":
			ctx.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "제목은 필수입니다",
				Code:  http.StatusBadRequest,
			})
		default:
			ctx.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "서버 내부 오류입니다",
				Code:  http.StatusInternalServerError,
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, issue)
}