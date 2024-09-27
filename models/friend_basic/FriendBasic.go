package friend_models

import (
	"errors"
	user_models "ginchat/models/user_basic"
	"ginchat/mydb"

	"gorm.io/gorm"
)

type FriendRelation struct {
	User_Id   uint
	Friend_Id uint
}

type FriendApply struct {
	gorm.Model
	From_User_Id uint
	To_Friend_Id uint
	Status       uint
	Reason       string
}

const (
	ApplyPending  = 1
	ApplyApproved = 2
	ApplyRejected = 3
)

func IsFriend(user_id, friend_id uint) (bool, error) {
	var friend_basic FriendRelation
	err := mydb.DB.Table("friend_basic").Where("user_id = ? AND friend_id = ?", user_id, friend_id).First(&friend_basic).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func GetFriendList(user_id uint) ([]user_models.UserBasic, error) {
	var friends []user_models.UserBasic
	err := mydb.DB.Table("user_basic").
		Select("user_basic.id, user_basic.name, user_basic.created_at, user_basic.updated_at").
		Joins("JOIN friend_basic ON user_basic.id = friend_basic.friend_id").
		Where("friend_basic.user_id = ?", user_id).
		Find(&friends).Error
	if err != nil {
		return nil, err
	}
	return friends, nil
}

func GetFriendApplyToList(user_id uint) ([]FriendApply, error) {
	var list []FriendApply
	err := mydb.DB.Raw(`
		SELECT * FROM friend_approve WHERE to_friend_id = ?
	`, user_id).Scan(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func GetFriendApplyFromList(user_id uint) ([]FriendApply, error) {
	var list []FriendApply
	err := mydb.DB.Raw(`
		SELECT * FROM friend_approve WHERE from_friend_id = ?
	`, user_id).Scan(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func CreateFriendApply(from_user_id, to_friend_id uint, reason string) (FriendApply, error) {
	is_friend, err := IsFriend(from_user_id, to_friend_id)
	if err != nil {
		return FriendApply{}, err
	}
	if is_friend {
		return FriendApply{}, errors.New("already friends")
	}
	approve := FriendApply{
		From_User_Id: from_user_id,
		To_Friend_Id: to_friend_id,
		Status:       ApplyPending,
		Reason:       reason,
	}
	return approve, mydb.DB.Table("friend_apply").Create(&approve).Error
}

func GetFriendApply(from_user_id, to_friend_id uint) (FriendApply, error) {
	var apply FriendApply
	err := mydb.DB.Table("friend_apply").Where("user_id = ? AND friend_id = ?", from_user_id, to_friend_id).First(&apply).Error
	if err != nil {
		return FriendApply{}, err
	}
	return apply, nil
}

func AgreeFriendApply(from_user_id, to_friend_id uint) error {
	is_friend, err := IsFriend(from_user_id, to_friend_id)
	if err != nil {
		return err
	}
	if is_friend {
		return errors.New("already friends")
	}
	apply, err := GetFriendApply(from_user_id, to_friend_id)
	if err != nil {
		return err
	}
	if apply.Status != ApplyPending {
		return errors.New("apply status not pending")
	}
	tx := mydb.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	res := mydb.DB.Model(&apply).Where("id = ?", apply.ID).Update("status", ApplyApproved)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}
	from_relation := FriendRelation{
		User_Id:   apply.From_User_Id,
		Friend_Id: apply.To_Friend_Id,
	}
	from_res := mydb.DB.Table("friend_basic").Create(&from_relation)
	if from_res.Error != nil {
		tx.Rollback()
		return from_res.Error
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func RejectFriendApply(user_id, apply_id uint) error {
	var apply FriendApply
	res := mydb.DB.Model(&apply).Where("id = ?", apply_id).Update("status", ApplyRejected)
	if res.Error != nil {
		return res.Error
	}
	if apply.To_Friend_Id != user_id {
		return errors.New("to_user_id not match")
	}
	return nil
}
