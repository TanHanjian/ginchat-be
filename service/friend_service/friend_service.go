package friend_service

import (
	friend_models "ginchat/models/friend_basic"
	user_models "ginchat/models/user_basic"
	"ginchat/utils"

	"github.com/gin-gonic/gin"
)

func convertToFriendInfoList(friendList []user_models.UserBasic) []FriendInfo {
	friendInfoList := make([]FriendInfo, len(friendList))
	for i, friend := range friendList {
		friendInfoList[i] = FriendInfo{
			ID:   friend.ID,
			Name: friend.Name,
		}
	}
	return friendInfoList
}

func GetFriendList(c *gin.Context) {
	user_id, err := utils.GetUserIdFromToken(c)
	if err != nil {
		c.JSON(200, gin.H{
			"message": err.Error(),
		})
		return
	}
	friend_list, err := friend_models.GetFriendList(user_id)
	friend_info_list := convertToFriendInfoList(friend_list)
	if err != nil {
		c.JSON(200, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"data": friend_info_list,
	})
}

func DeleteFriend(c *gin.Context) {
	user_id, err := utils.GetUserIdFromToken(c)
	if err != nil {
		c.JSON(200, gin.H{
			"message": err.Error(),
		})
		return
	}
	delete_dto, err := utils.BodyToModel[DeleteFriendApplyDto](c)
	if err != nil {
		c.JSON(200, gin.H{
			"message": err.Error(),
		})
		return
	}
	err = friend_models.DeleteFriend(user_id, delete_dto.Friend_Id)
	if err != nil {
		c.JSON(200, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "success",
	})
}

func GetFriendApplyToList(c *gin.Context) {
	user_id, err := utils.GetUserIdFromToken(c)
	if err != nil {
		c.JSON(200, gin.H{
			"message": err.Error(),
		})
		return
	}
	apply_list, err := friend_models.GetFriendApplyToList(user_id)
	if err != nil {
		c.JSON(200, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"data": apply_list,
	})
}

func GetFriendApplyFromList(c *gin.Context) {
	user_id, err := utils.GetUserIdFromToken(c)
	if err != nil {
		c.JSON(200, gin.H{
			"message": err.Error(),
		})
		return
	}
	apply_list, err := friend_models.GetFriendApplyFromList(user_id)
	if err != nil {
		c.JSON(200, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"data": apply_list,
	})
}

func CreateFriendApply(c *gin.Context) {
	user_id, err := utils.GetUserIdFromToken(c)
	if err != nil {
		c.JSON(200, gin.H{
			"message": err.Error(),
		})
		return
	}
	to_user_dto, err := utils.BodyToModel[FriendApplyDto](c)
	if err != nil {
		c.JSON(200, gin.H{
			"message": err.Error(),
		})
		return
	}
	if user_id == to_user_dto.Friend_Id {
		c.JSON(200, gin.H{
			"message": "不能添加自己为好友",
		})
		return
	}
	_, err = friend_models.CreateFriendApply(user_id, to_user_dto.Friend_Id, to_user_dto.Reason)
	if err != nil {
		c.JSON(200, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "申请成功",
	})
}

func AgreeFriendApply(c *gin.Context) {
	user_id, err := utils.GetUserIdFromToken(c)
	if err != nil {
		c.JSON(200, gin.H{
			"message": err.Error(),
		})
		return
	}
	agree_dto, err := utils.BodyToModel[AgreeFriendApplyDto](c)
	if err != nil {
		c.JSON(200, gin.H{
			"message": err.Error(),
		})
		return
	}
	err = friend_models.AgreeFriendApply(user_id, agree_dto.Apply_Id)
	if err != nil {
		c.JSON(200, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "同意成功",
	})
}

func RejectFriendApply(c *gin.Context) {
	user_id, err := utils.GetUserIdFromToken(c)
	if err != nil {
		c.JSON(200, gin.H{
			"message": err.Error(),
		})
		return
	}
	reject_dto, err := utils.BodyToModel[RejectFriendApplyDto](c)
	if err != nil {
		c.JSON(200, gin.H{
			"message": err.Error(),
		})
		return
	}
	err = friend_models.RejectFriendApply(user_id, reject_dto.Apply_Id)
	if err != nil {
		c.JSON(200, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "拒绝成功",
	})
}
