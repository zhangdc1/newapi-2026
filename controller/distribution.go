package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"

	"github.com/gin-gonic/gin"
)

func GetDistributionOverview(c *gin.Context) {
	overview, err := model.GetDistributionOverview(c.GetInt("id"))
	if err != nil {
		common.ApiError(c, err)
		return
	}
	common.ApiSuccess(c, overview)
}

func GenerateDistributionInviteLink(c *gin.Context) {
	overview, err := model.GetDistributionOverview(c.GetInt("id"))
	if err != nil {
		common.ApiError(c, err)
		return
	}
	common.ApiSuccess(c, gin.H{
		"dist_id":     overview.DistId,
		"invite_link": overview.InviteLink,
	})
}

func GetDistributionInvites(c *gin.Context) {
	pageInfo := common.GetPageQuery(c)
	referrals, total, err := model.ListUserInvites(c.GetInt("id"), pageInfo)
	if err != nil {
		common.ApiError(c, err)
		return
	}
	pageInfo.SetTotal(int(total))
	pageInfo.SetItems(referrals)
	common.ApiSuccess(c, pageInfo)
}

func GetUserDistributionCommissionRecords(c *gin.Context) {
	pageInfo := common.GetPageQuery(c)
	records, total, err := model.ListCommissionRecords(c.GetInt("id"), pageInfo)
	if err != nil {
		common.ApiError(c, err)
		return
	}
	pageInfo.SetTotal(int(total))
	pageInfo.SetItems(records)
	common.ApiSuccess(c, pageInfo)
}

func SettleDistributionCommission(c *gin.Context) {
	amount, err := model.SettleDistributionCommission(c.GetInt("id"))
	if err != nil {
		common.ApiError(c, err)
		return
	}
	common.ApiSuccess(c, gin.H{"amount": amount})
}

func AdminGetDistributionSettings(c *gin.Context) {
	settings, err := model.GetDistributionSettings()
	if err != nil {
		common.ApiError(c, err)
		return
	}
	common.ApiSuccess(c, settings)
}

func AdminUpdateDistributionSettings(c *gin.Context) {
	var settings model.DistributionSettings
	if err := common.DecodeJson(c.Request.Body, &settings); err != nil {
		common.ApiError(c, err)
		return
	}
	if err := model.UpdateDistributionSettings(&settings, c.GetInt("id")); err != nil {
		common.ApiError(c, err)
		return
	}
	common.ApiSuccess(c, settings)
}

func AdminGetDistributionCommissionRecords(c *gin.Context) {
	pageInfo := common.GetPageQuery(c)
	referrerId, _ := strconv.Atoi(c.Query("referrer_user_id"))
	records, total, err := model.ListCommissionRecords(referrerId, pageInfo)
	if err != nil {
		common.ApiError(c, err)
		return
	}
	pageInfo.SetTotal(int(total))
	pageInfo.SetItems(records)
	common.ApiSuccess(c, pageInfo)
}

func AdminGetDistributionReferrals(c *gin.Context) {
	pageInfo := common.GetPageQuery(c)
	referrerId, _ := strconv.Atoi(c.Query("referrer_user_id"))
	if referrerId <= 0 {
		common.ApiError(c, errors.New("referrer_user_id is required"))
		return
	}
	referrals, total, err := model.ListUserInvites(referrerId, pageInfo)
	if err != nil {
		common.ApiError(c, err)
		return
	}
	pageInfo.SetTotal(int(total))
	pageInfo.SetItems(referrals)
	common.ApiSuccess(c, pageInfo)
}

type manualGrantDistributionRequest struct {
	ReferrerUserId int    `json:"referrer_user_id"`
	ReferredUserId int    `json:"referred_user_id"`
	Amount         int    `json:"amount"`
	Reason         string `json:"reason"`
}

func AdminManualGrantDistributionCommission(c *gin.Context) {
	var req manualGrantDistributionRequest
	if err := common.DecodeJson(c.Request.Body, &req); err != nil {
		common.ApiError(c, err)
		return
	}
	if err := model.ManualGrantDistributionCommission(req.ReferrerUserId, req.ReferredUserId, req.Amount, c.GetInt("id"), req.Reason); err != nil {
		common.ApiError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "", "data": nil})
}

type adjustDistributionCommissionRequest struct {
	UserId int    `json:"user_id"`
	Action string `json:"action"`
	Amount int    `json:"amount"`
	Reason string `json:"reason"`
}

func AdminAdjustDistributionCommission(c *gin.Context) {
	var req adjustDistributionCommissionRequest
	if err := common.DecodeJson(c.Request.Body, &req); err != nil {
		common.ApiError(c, err)
		return
	}
	if err := model.AdjustDistributionCommission(req.UserId, req.Action, req.Amount, c.GetInt("id"), req.Reason, 0); err != nil {
		common.ApiError(c, err)
		return
	}
	common.ApiSuccess(c, nil)
}
