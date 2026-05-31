package model

import (
	"errors"
	"fmt"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/logger"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	DistributionRewardModePercent = 1
	DistributionRewardModeFixed   = 2

	CommissionStatusPending = 1
	CommissionStatusSettled = 3
	CommissionStatusInvalid = 4

	CommissionTypePercent      = 1
	CommissionTypeFixed        = 2
	CommissionTypeManualAdjust = 3

	CommissionSourceSettle       = "commission_settle"
	CommissionSourceManualAdjust = "manual_adjust"
)

type DistributionSettings struct {
	Id                     int   `json:"id"`
	Enabled                bool  `json:"enabled" gorm:"default:false"`
	RewardRate             int   `json:"reward_rate" gorm:"type:int;default:10"`
	ReferralLimit          int   `json:"referral_limit" gorm:"type:int;default:0"`
	RewardMode             int   `json:"reward_mode" gorm:"type:int;default:1"`
	FixedRewardAmount      int   `json:"fixed_reward_amount" gorm:"type:int;default:0"`
	FreezeDays             int   `json:"freeze_days" gorm:"type:int;default:7"`
	ReferredUserReward     int   `json:"referred_user_reward" gorm:"type:int;default:0"`
	MinSettlementThreshold int   `json:"min_settlement_threshold" gorm:"type:int;default:0"`
	AdminRechargeTrigger   bool  `json:"admin_recharge_trigger" gorm:"default:true"`
	UpdatedBy              int   `json:"updated_by" gorm:"type:int;default:0"`
	UpdatedAt              int64 `json:"updated_at" gorm:"bigint"`
}

type UserReferral struct {
	Id             int    `json:"id"`
	ReferrerUserId int    `json:"referrer_user_id" gorm:"index"`
	ReferredUserId int    `json:"referred_user_id" gorm:"uniqueIndex"`
	Status         int    `json:"status" gorm:"type:int;default:1"`
	BoundAt        int64  `json:"bound_at" gorm:"bigint;index"`
	Source         string `json:"source" gorm:"type:varchar(50);default:'register'"`
	ReferrerDistId string `json:"referrer_dist_id" gorm:"type:varchar(255);default:''"`
}

type CommissionRecord struct {
	Id               int    `json:"id"`
	ReferrerUserId   int    `json:"referrer_user_id" gorm:"index"`
	ReferredUserId   int    `json:"referred_user_id" gorm:"index"`
	BalanceLogId     int    `json:"balance_log_id" gorm:"type:int;default:0"`
	Source           string `json:"source" gorm:"type:varchar(50);index"`
	SourceId         string `json:"source_id" gorm:"type:varchar(128);index:idx_commission_source,priority:2"`
	IncreasedAmount  int    `json:"increased_amount" gorm:"type:int;default:0"`
	CommissionType   int    `json:"commission_type" gorm:"type:int;default:1"`
	CommissionAmount int    `json:"commission_amount" gorm:"type:int;default:0"`
	Status           int    `json:"status" gorm:"type:int;default:1;index"`
	FreezeUntil      int64  `json:"freeze_until" gorm:"bigint;index"`
	CreatedAt        int64  `json:"created_at" gorm:"bigint;index"`
	SettledAt        int64  `json:"settled_at" gorm:"bigint;default:0"`
	AdminId          int    `json:"admin_id" gorm:"type:int;default:0"`
	Reason           string `json:"reason" gorm:"type:varchar(255);default:''"`
	BeforeBalance    int    `json:"before_balance" gorm:"type:int;default:0"`
	AfterBalance     int    `json:"after_balance" gorm:"type:int;default:0"`
}

func (CommissionRecord) TableName() string {
	return "commission_records"
}

type UserCommissionAccount struct {
	UserId         int   `json:"user_id" gorm:"primaryKey"`
	PendingBalance int   `json:"pending_balance" gorm:"type:int;default:0"`
	SettledBalance int   `json:"settled_balance" gorm:"type:int;default:0"`
	TotalEarned    int   `json:"total_earned" gorm:"type:int;default:0"`
	UpdatedAt      int64 `json:"updated_at" gorm:"bigint"`
}

type CommissionTransfer struct {
	Id             int    `json:"id"`
	UserId         int    `json:"user_id" gorm:"index"`
	TransferAmount int    `json:"transfer_amount" gorm:"type:int;default:0"`
	SourceRecords  string `json:"source_records" gorm:"type:text"`
	Status         int    `json:"status" gorm:"type:int;default:1"`
	CreatedAt      int64  `json:"created_at" gorm:"bigint"`
	CompletedAt    int64  `json:"completed_at" gorm:"bigint"`
}

type ReferralCountTracker struct {
	Id             int   `json:"id"`
	ReferrerUserId int   `json:"referrer_user_id" gorm:"uniqueIndex:idx_referral_count"`
	ReferredUserId int   `json:"referred_user_id" gorm:"uniqueIndex:idx_referral_count"`
	RewardCount    int   `json:"reward_count" gorm:"type:int;default:0"`
	MaxCount       int   `json:"max_count" gorm:"type:int;default:0"`
	UpdatedAt      int64 `json:"updated_at" gorm:"bigint"`
}

type DistributionOverview struct {
	Enabled                bool   `json:"enabled"`
	DistId                 string `json:"dist_id"`
	InviteLink             string `json:"invite_link"`
	TotalEarned            int    `json:"total_earned"`
	PendingCommission      int    `json:"pending_commission"`
	AvailableCommission    int    `json:"available_commission"`
	SettledCommission      int    `json:"settled_commission"`
	InviteCount            int    `json:"invite_count"`
	EffectiveInviteCount   int    `json:"effective_invite_count"`
	MinSettlementThreshold int    `json:"min_settlement_threshold"`
}

type DistributionBalanceIncrease struct {
	UserId       int
	Amount       int
	Source       string
	SourceId     string
	TempDistId   string
	AdminTrigger bool
}

func GetDistributionSettings() (*DistributionSettings, error) {
	settings := &DistributionSettings{}
	err := DB.First(settings, 1).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		settings = &DistributionSettings{
			Id:                     1,
			Enabled:                false,
			RewardRate:             10,
			ReferralLimit:          0,
			RewardMode:             DistributionRewardModePercent,
			FixedRewardAmount:      0,
			FreezeDays:             7,
			ReferredUserReward:     0,
			MinSettlementThreshold: 0,
			AdminRechargeTrigger:   true,
			UpdatedAt:              common.GetTimestamp(),
		}
		err = DB.Create(settings).Error
	}
	return settings, err
}

func UpdateDistributionSettings(settings *DistributionSettings, adminId int) error {
	if settings.RewardRate < 0 || settings.RewardRate > 50 {
		return errors.New("reward_rate must be between 0 and 50")
	}
	if settings.ReferralLimit < 0 || settings.FixedRewardAmount < 0 || settings.FreezeDays < 0 ||
		settings.ReferredUserReward < 0 || settings.MinSettlementThreshold < 0 {
		return errors.New("distribution settings cannot contain negative values")
	}
	if settings.RewardMode != DistributionRewardModePercent && settings.RewardMode != DistributionRewardModeFixed {
		return errors.New("invalid reward_mode")
	}
	settings.Id = 1
	settings.UpdatedBy = adminId
	settings.UpdatedAt = common.GetTimestamp()
	return DB.Save(settings).Error
}

func SyncUserReferral(referredUserId int, source string) {
	if referredUserId == 0 {
		return
	}
	var user User
	if err := DB.Select("id", "inviter_id", "aff_code").First(&user, "id = ?", referredUserId).Error; err != nil {
		return
	}
	if user.InviterId == 0 || user.InviterId == user.Id {
		return
	}
	distId := ""
	_ = DB.Model(&User{}).Where("id = ?", user.InviterId).Select("aff_code").Scan(&distId).Error
	referral := UserReferral{
		ReferrerUserId: user.InviterId,
		ReferredUserId: user.Id,
		Status:         1,
		BoundAt:        common.GetTimestamp(),
		Source:         source,
		ReferrerDistId: distId,
	}
	_ = DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&referral).Error
}

func buildInviteLink(distId string) string {
	if distId == "" {
		return ""
	}
	return "/sign-up?dist_id=" + distId
}

func GetDistributionOverview(userId int) (*DistributionOverview, error) {
	settings, err := GetDistributionSettings()
	if err != nil {
		return nil, err
	}
	user, err := GetUserById(userId, true)
	if err != nil {
		return nil, err
	}
	if user.AffCode == "" {
		user.AffCode = common.GetRandomString(8)
		if err := user.Update(false); err != nil {
			return nil, err
		}
	}
	account, err := getOrCreateCommissionAccount(DB, userId)
	if err != nil {
		return nil, err
	}
	now := common.GetTimestamp()
	var available int64
	if err := DB.Model(&CommissionRecord{}).
		Where("referrer_user_id = ? AND status = ? AND freeze_until <= ? AND commission_amount > 0", userId, CommissionStatusPending, now).
		Select("COALESCE(SUM(commission_amount), 0)").Scan(&available).Error; err != nil {
		return nil, err
	}
	var inviteCount int64
	_ = DB.Model(&UserReferral{}).Where("referrer_user_id = ? AND status = ?", userId, 1).Count(&inviteCount).Error
	var effectiveInviteCount int64
	_ = DB.Model(&CommissionRecord{}).Where("referrer_user_id = ? AND commission_amount > 0", userId).Distinct("referred_user_id").Count(&effectiveInviteCount).Error
	return &DistributionOverview{
		Enabled:                settings.Enabled,
		DistId:                 user.AffCode,
		InviteLink:             buildInviteLink(user.AffCode),
		TotalEarned:            account.TotalEarned,
		PendingCommission:      account.PendingBalance,
		AvailableCommission:    int(available),
		SettledCommission:      account.SettledBalance,
		InviteCount:            int(inviteCount),
		EffectiveInviteCount:   int(effectiveInviteCount),
		MinSettlementThreshold: settings.MinSettlementThreshold,
	}, nil
}

func GrantCommissionForBalanceIncrease(event DistributionBalanceIncrease) {
	if event.UserId == 0 || event.Amount <= 0 || event.Source == CommissionSourceSettle {
		return
	}
	settings, err := GetDistributionSettings()
	if err != nil || !settings.Enabled {
		return
	}
	if event.AdminTrigger && !settings.AdminRechargeTrigger {
		return
	}
	referrerId := 0
	if event.TempDistId != "" {
		referrerId, _ = GetUserIdByAffCode(strings.TrimSpace(event.TempDistId))
	}
	if referrerId == 0 {
		var user User
		if err := DB.Select("id", "inviter_id").First(&user, "id = ?", event.UserId).Error; err != nil {
			return
		}
		referrerId = user.InviterId
	}
	if referrerId == 0 || referrerId == event.UserId {
		return
	}
	commissionAmount := calculateCommissionAmount(event.Amount, settings)
	if commissionAmount <= 0 {
		return
	}
	err = DB.Transaction(func(tx *gorm.DB) error {
		if event.SourceId != "" {
			var count int64
			if err := tx.Model(&CommissionRecord{}).
				Where("source = ? AND source_id = ? AND referrer_user_id = ? AND referred_user_id = ?", event.Source, event.SourceId, referrerId, event.UserId).
				Count(&count).Error; err != nil {
				return err
			}
			if count > 0 {
				return nil
			}
		}
		tracker, err := getOrCreateReferralCountTracker(tx, referrerId, event.UserId, settings.ReferralLimit)
		if err != nil {
			return err
		}
		if settings.ReferralLimit > 0 && tracker.RewardCount >= settings.ReferralLimit {
			return nil
		}
		account, err := getOrCreateCommissionAccount(tx, referrerId)
		if err != nil {
			return err
		}
		before := account.PendingBalance
		commissionType := CommissionTypePercent
		if settings.RewardMode == DistributionRewardModeFixed {
			commissionType = CommissionTypeFixed
		}
		now := common.GetTimestamp()
		record := CommissionRecord{
			ReferrerUserId:   referrerId,
			ReferredUserId:   event.UserId,
			Source:           event.Source,
			SourceId:         event.SourceId,
			IncreasedAmount:  event.Amount,
			CommissionType:   commissionType,
			CommissionAmount: commissionAmount,
			Status:           CommissionStatusPending,
			FreezeUntil:      now + int64(settings.FreezeDays)*86400,
			CreatedAt:        now,
			BeforeBalance:    before,
			AfterBalance:     before + commissionAmount,
		}
		if err := tx.Create(&record).Error; err != nil {
			return err
		}
		if err := tx.Model(&UserCommissionAccount{}).Where("user_id = ?", referrerId).Updates(map[string]interface{}{
			"pending_balance": gorm.Expr("pending_balance + ?", commissionAmount),
			"total_earned":    gorm.Expr("total_earned + ?", commissionAmount),
			"updated_at":      now,
		}).Error; err != nil {
			return err
		}
		if err := tx.Model(&ReferralCountTracker{}).Where("id = ?", tracker.Id).Updates(map[string]interface{}{
			"reward_count": gorm.Expr("reward_count + ?", 1),
			"max_count":    settings.ReferralLimit,
			"updated_at":   now,
		}).Error; err != nil {
			return err
		}
		SyncUserReferral(event.UserId, "balance_increase")
		return nil
	})
	if err != nil {
		common.SysLog("failed to grant distribution commission: " + err.Error())
	}
}

func calculateCommissionAmount(amount int, settings *DistributionSettings) int {
	if settings.RewardMode == DistributionRewardModeFixed {
		return settings.FixedRewardAmount
	}
	return amount * settings.RewardRate / 100
}

func getOrCreateCommissionAccount(tx *gorm.DB, userId int) (*UserCommissionAccount, error) {
	account := &UserCommissionAccount{}
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(account, "user_id = ?", userId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		account = &UserCommissionAccount{UserId: userId, UpdatedAt: common.GetTimestamp()}
		err = tx.Create(account).Error
	}
	return account, err
}

func getOrCreateReferralCountTracker(tx *gorm.DB, referrerId int, referredId int, maxCount int) (*ReferralCountTracker, error) {
	tracker := &ReferralCountTracker{}
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(tracker, "referrer_user_id = ? AND referred_user_id = ?", referrerId, referredId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		tracker = &ReferralCountTracker{ReferrerUserId: referrerId, ReferredUserId: referredId, MaxCount: maxCount, UpdatedAt: common.GetTimestamp()}
		err = tx.Create(tracker).Error
	}
	return tracker, err
}

func ListUserInvites(userId int, pageInfo *common.PageInfo) ([]UserReferral, int64, error) {
	var referrals []UserReferral
	var total int64
	query := DB.Model(&UserReferral{}).Where("referrer_user_id = ?", userId)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Order("id desc").Limit(pageInfo.GetPageSize()).Offset(pageInfo.GetStartIdx()).Find(&referrals).Error
	return referrals, total, err
}

func ListCommissionRecords(referrerId int, pageInfo *common.PageInfo) ([]CommissionRecord, int64, error) {
	var records []CommissionRecord
	var total int64
	query := DB.Model(&CommissionRecord{})
	if referrerId > 0 {
		query = query.Where("referrer_user_id = ?", referrerId)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Order("id desc").Limit(pageInfo.GetPageSize()).Offset(pageInfo.GetStartIdx()).Find(&records).Error
	return records, total, err
}

func SettleDistributionCommission(userId int) (int, error) {
	settings, err := GetDistributionSettings()
	if err != nil {
		return 0, err
	}
	now := common.GetTimestamp()
	var transferAmount int
	var recordIds []int
	err = DB.Transaction(func(tx *gorm.DB) error {
		account, err := getOrCreateCommissionAccount(tx, userId)
		if err != nil {
			return err
		}
		var records []CommissionRecord
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("referrer_user_id = ? AND status = ? AND freeze_until <= ? AND commission_amount > 0", userId, CommissionStatusPending, now).
			Order("id asc").Find(&records).Error; err != nil {
			return err
		}
		for _, record := range records {
			transferAmount += record.CommissionAmount
			recordIds = append(recordIds, record.Id)
		}
		if transferAmount <= 0 {
			return errors.New("no available commission to settle")
		}
		if settings.MinSettlementThreshold > 0 && transferAmount < settings.MinSettlementThreshold {
			return fmt.Errorf("available commission is lower than minimum settlement threshold %s", logger.LogQuota(settings.MinSettlementThreshold))
		}
		idsJson, _ := common.Marshal(recordIds)
		if err := tx.Model(&CommissionRecord{}).Where("id IN ?", recordIds).Updates(map[string]interface{}{
			"status":     CommissionStatusSettled,
			"settled_at": now,
		}).Error; err != nil {
			return err
		}
		if err := tx.Model(&UserCommissionAccount{}).Where("user_id = ?", userId).Updates(map[string]interface{}{
			"pending_balance": gorm.Expr("pending_balance - ?", transferAmount),
			"settled_balance": gorm.Expr("settled_balance + ?", transferAmount),
			"updated_at":      now,
		}).Error; err != nil {
			return err
		}
		if err := tx.Model(&User{}).Where("id = ?", userId).Update("quota", gorm.Expr("quota + ?", transferAmount)).Error; err != nil {
			return err
		}
		transfer := CommissionTransfer{
			UserId:         userId,
			TransferAmount: transferAmount,
			SourceRecords:  string(idsJson),
			Status:         1,
			CreatedAt:      now,
			CompletedAt:    now,
		}
		if err := tx.Create(&transfer).Error; err != nil {
			return err
		}
		account.PendingBalance -= transferAmount
		account.SettledBalance += transferAmount
		return nil
	})
	if err != nil {
		return 0, err
	}
	RecordLog(userId, LogTypeTopup, fmt.Sprintf("分销佣金转入余额 %s", logger.LogQuota(transferAmount)))
	return transferAmount, nil
}

func ManualGrantDistributionCommission(referrerId int, referredId int, amount int, adminId int, reason string) error {
	return AdjustDistributionCommission(referrerId, "increase", amount, adminId, reason, referredId)
}

func AdjustDistributionCommission(userId int, action string, amount int, adminId int, reason string, referredUserId int) error {
	if userId == 0 {
		return errors.New("invalid user id")
	}
	reason = strings.TrimSpace(reason)
	if reason == "" {
		reason = "manual adjustment"
	}
	now := common.GetTimestamp()
	return DB.Transaction(func(tx *gorm.DB) error {
		account, err := getOrCreateCommissionAccount(tx, userId)
		if err != nil {
			return err
		}
		before := account.PendingBalance
		switch action {
		case "increase":
			if amount <= 0 {
				return errors.New("amount must be greater than zero")
			}
			record := CommissionRecord{
				ReferrerUserId:   userId,
				ReferredUserId:   referredUserId,
				Source:           CommissionSourceManualAdjust,
				SourceId:         fmt.Sprintf("admin-%d-%d", adminId, now),
				IncreasedAmount:  0,
				CommissionType:   CommissionTypeManualAdjust,
				CommissionAmount: amount,
				Status:           CommissionStatusPending,
				FreezeUntil:      now,
				CreatedAt:        now,
				AdminId:          adminId,
				Reason:           reason,
				BeforeBalance:    before,
				AfterBalance:     before + amount,
			}
			if err := tx.Create(&record).Error; err != nil {
				return err
			}
			return tx.Model(&UserCommissionAccount{}).Where("user_id = ?", userId).Updates(map[string]interface{}{
				"pending_balance": gorm.Expr("pending_balance + ?", amount),
				"total_earned":    gorm.Expr("total_earned + ?", amount),
				"updated_at":      now,
			}).Error
		case "decrease":
			if amount <= 0 {
				return errors.New("amount must be greater than zero")
			}
			return invalidateCommissionAmount(tx, userId, amount, adminId, reason, now, before)
		case "clear":
			if before <= 0 {
				return nil
			}
			return invalidateCommissionAmount(tx, userId, before, adminId, reason, now, before)
		default:
			return errors.New("invalid adjustment action")
		}
	})
}

func invalidateCommissionAmount(tx *gorm.DB, userId int, amount int, adminId int, reason string, now int64, before int) error {
	remaining := amount
	var records []CommissionRecord
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("referrer_user_id = ? AND status = ? AND commission_amount > 0", userId, CommissionStatusPending).
		Order("freeze_until asc, id asc").Find(&records).Error; err != nil {
		return err
	}
	for _, record := range records {
		if remaining <= 0 {
			break
		}
		deduct := record.CommissionAmount
		if deduct > remaining {
			deduct = remaining
		}
		if deduct == record.CommissionAmount {
			if err := tx.Model(&CommissionRecord{}).Where("id = ?", record.Id).Updates(map[string]interface{}{
				"status":     CommissionStatusInvalid,
				"settled_at": now,
				"admin_id":   adminId,
				"reason":     reason,
			}).Error; err != nil {
				return err
			}
		} else {
			if err := tx.Model(&CommissionRecord{}).Where("id = ?", record.Id).Update("commission_amount", gorm.Expr("commission_amount - ?", deduct)).Error; err != nil {
				return err
			}
		}
		remaining -= deduct
	}
	deducted := amount - remaining
	if deducted <= 0 {
		return errors.New("no unsettled commission to adjust")
	}
	adjustRecord := CommissionRecord{
		ReferrerUserId:   userId,
		Source:           CommissionSourceManualAdjust,
		SourceId:         fmt.Sprintf("admin-%d-%d", adminId, now),
		CommissionType:   CommissionTypeManualAdjust,
		CommissionAmount: -deducted,
		Status:           CommissionStatusSettled,
		FreezeUntil:      now,
		CreatedAt:        now,
		SettledAt:        now,
		AdminId:          adminId,
		Reason:           reason,
		BeforeBalance:    before,
		AfterBalance:     before - deducted,
	}
	if err := tx.Create(&adjustRecord).Error; err != nil {
		return err
	}
	return tx.Model(&UserCommissionAccount{}).Where("user_id = ?", userId).Updates(map[string]interface{}{
		"pending_balance": gorm.Expr("pending_balance - ?", deducted),
		"total_earned":    gorm.Expr("total_earned - ?", deducted),
		"updated_at":      now,
	}).Error
}
