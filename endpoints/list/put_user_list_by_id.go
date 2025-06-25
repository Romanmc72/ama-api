package list

import (
	"ama/api/application/list"
	"ama/api/application/requests"
	"ama/api/application/responses"
	"ama/api/interfaces"
	"ama/api/logging"
	"net/http"
	"time"
)

func PutUserListById(c interfaces.APIContext, db interfaces.ListUpdater) {
	logger := logging.GetLogger()
	userId := c.Param("userId")
	listId := c.Param("listId")

	if userId == "" || listId == "" {
		logger.Error("blank userId or listId",
			"userId", userId,
			"listId", listId,
			"body", c.GetString("body"),
		)
		c.IndentedJSON(http.StatusBadRequest, "userId or listId is empty")
		return
	}

	var listReq requests.PutListRequest
	if err := c.BindJSON(&listReq); err != nil {
		logger.Error("unable to bind request body",
			"userId", userId,
			"listId", listId,
			"body", c.GetString("body"),
			"error", err,
		)
		c.IndentedJSON(http.StatusBadRequest, "invalid request body")
		return
	}

	list := list.List{
		ID:   listId,
		Name: listReq.Name,
	}

	if err := db.UpdateList(userId, list); err != nil {
		logger.Error("unable to update list",
			"userId", userId,
			"listId", listId,
			"body", c.GetString("body"),
			"error", err,
		)
		c.IndentedJSON(http.StatusInternalServerError, "unable to update list")
		return
	}

	c.IndentedJSON(http.StatusOK, responses.NewDeleteUserResponse(true, time.Now().Unix()))
}
