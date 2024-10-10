package friend_service

type FriendApplyDto struct {
	Friend_Id uint   `json:"friendId" binding:"required"`
	Reason    string `json:"reason"`
}

type AgreeFriendApplyDto struct {
	Apply_Id uint `json:"applyId" binding:"required"`
}

type RejectFriendApplyDto struct {
	Apply_Id uint `json:"applyId" binding:"required"`
}

type DeleteFriendApplyDto struct {
	Friend_Id uint `json:"friendId" binding:"required"`
}
