package list

import (
	"ama/api/application/list"
	"ama/api/application/requests"
	"ama/api/application/responses"
	"ama/api/constants"
	"ama/api/interfaces"
	"ama/api/logging"
	"net/http"
	"strings"
)

// PutUserListById godoc
//
//	@Summary		Update a list
//	@Description	Update a list
//	@Tags			list
//	@Accept			json
//	@Produce		json
//	@Param			userId	path		string						true	"User ID"
//	@Param			listId	path		string						true	"List ID"
//	@Param			list	body		requests.PutUserListRequest	true	"List data"
//	@Success		200		{object}	responses.SuccessResponse
//	@Failure		400		{object}	responses.ErrorResponse
//	@Failure		500		{object}	responses.ErrorResponse
//	@Router			/user/{userId}/list/{listId} [put]
func PutUserListById(c interfaces.APIContext, db interfaces.ListUpdater) {
	body := "body"
	logger := logging.GetLogger()
	userId := c.Param(constants.UserIdPathIdentifier)
	listId := c.Param(constants.ListIdPathIdentifier)

	if strings.TrimSpace(userId) == "" || strings.TrimSpace(listId) == "" {
		logger.Error("blank userId or listId",
			constants.UserIdPathIdentifier, userId,
			constants.ListIdPathIdentifier, listId,
			body, c.GetString(body),
		)
		c.IndentedJSON(http.StatusBadRequest, responses.NewError("listId or userId is blank"))
		return
	}

	var listReq requests.PutUserListRequest
	if err := c.BindJSON(&listReq); err != nil {
		logger.Error("unable to bind request body",
			constants.UserIdPathIdentifier, userId,
			constants.ListIdPathIdentifier, listId,
			body, c.GetString(body),
			"error", err,
		)
		c.IndentedJSON(http.StatusBadRequest, responses.NewError("invalid request body"))
		return
	}

	list := list.List{
		ID:   listId,
		Name: listReq.Name,
	}

	if err := db.UpdateList(userId, list); err != nil {
		logger.Error("unable to update list",
			constants.UserIdPathIdentifier, userId,
			constants.ListIdPathIdentifier, listId,
			body, c.GetString(body),
			"error", err,
		)
		c.IndentedJSON(http.StatusInternalServerError, responses.NewError("unable to update list"))
		return
	}

	c.IndentedJSON(http.StatusOK, responses.NewSuccessResponse(true))
}
