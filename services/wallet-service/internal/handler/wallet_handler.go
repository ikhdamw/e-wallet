package handler

import (
	"net/http"
	"strconv"

	"github.com/ikhdamw/e-wallet/wallet-service/internal/model"
	"github.com/ikhdamw/e-wallet/wallet-service/internal/service"
	"github.com/ikhdamw/e-wallet/wallet-service/pkg/response"

	"github.com/gin-gonic/gin"
)

type WalletHandler struct {
	walletService service.WalletService
}

func NewWalletHandler(walletService service.WalletService) *WalletHandler {
	return &WalletHandler{walletService: walletService}
}

func (h *WalletHandler) GetBalance(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not found in context")
		return
	}

	balance, err := h.walletService.GetBalance(userID.(string))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, balance)
}

func (h *WalletHandler) TopUp(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not found in context")
		return
	}

	var req model.TopUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	transaction, err := h.walletService.TopUp(userID.(string), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, http.StatusOK, "Top up successful", transaction)
}

func (h *WalletHandler) GetHistory(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not found in context")
		return
	}

	// Parse query params
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	history, err := h.walletService.GetHistory(userID.(string), page, limit)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, history)
}
