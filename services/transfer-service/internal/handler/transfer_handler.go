package handler

import (
	"net/http"

	"github.com/ikhdamw/e-wallet/transfer-service/internal/model"
	"github.com/ikhdamw/e-wallet/transfer-service/internal/service"
	"github.com/ikhdamw/e-wallet/transfer-service/pkg/response"

	"github.com/gin-gonic/gin"
)

type TransferHandler struct {
	transferService service.TransferService
}

func NewTransferHandler(transferService service.TransferService) *TransferHandler {
	return &TransferHandler{transferService: transferService}
}

func (h *TransferHandler) InternalTransfer(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not found in context")
		return
	}

	var req model.InternalTransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.transferService.InternalTransfer(userID.(string), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, http.StatusOK, "Transfer initiated", result)
}

func (h *TransferHandler) ExternalTransfer(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not found in context")
		return
	}

	var req model.ExternalTransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.transferService.ExternalTransfer(userID.(string), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, http.StatusOK, "External transfer initiated", result)
}

func (h *TransferHandler) GetStatus(c *gin.Context) {
	transferID := c.Param("id")

	transfer, err := h.transferService.GetStatus(transferID)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	if transfer == nil {
		response.NotFound(c, "Transfer not found")
		return
	}

	response.Success(c, http.StatusOK, transfer)
}
