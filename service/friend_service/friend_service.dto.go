package friend_service

type FriendApplyDto struct {
	Friend_Id uint   `json:"friend_id" binding:"required"`
	Reason    string `json:"reason"`
}

type AgreeFriendApplyDto struct {
	Apply_Id uint `json:"apply_id" binding:"required"`
}

type RejectFriendApplyDto struct {
	Apply_Id uint `json:"apply_id" binding:"required"`
}
