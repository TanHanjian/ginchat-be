package chatroom_service

type CreateSingleRoomDto struct {
	Name     string
	FriendID uint
}

type CreateMultiRoomDto struct {
	Name      string
	FriendIDs []uint `json:"friendIDs" binding:"required"`
}

type JoinChatroomDto struct {
	ChatroomID uint `json:"chatroomID" binding:"required"`
}

type QuitChatroomDto struct {
	ChatroomID uint `json:"chatroomID" binding:"required"`
}
