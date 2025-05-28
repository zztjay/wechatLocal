package controller

import (
	"feiyu.com/wx/api/model"
	"feiyu.com/wx/api/service"
	"feiyu.com/wx/api/vo"
	"feiyu.com/wx/protobuf/wechat"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// GetContactListApi 获取全部联系人
func GetContactListApi(ctx *gin.Context) {
	reqModel := new(model.GetContactListModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.GetContactListService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// GetContactAllListApi 获取所有联系人
func GetContactAllListApi(ctx *gin.Context) {
	reqModel := new(model.GetContactAllListModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.GetContactAllListService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// GetContactAllListApi 同步获取所有联系人
func SyncContactAllListApi(ctx *gin.Context) {
	reqModel := new(model.GetContactAllListModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	// ContactList集合[]string
	ContactList := make([]string, 0)

	i := 0
	query := &model.GetContactListModel{
		CurrentWxcontactSeq:       0,
		CurrentChatRoomContactSeq: 0,
		ContactUsernameList:       nil,
	}
	for {
		result := service.GetContactListService(queryKey, *query)
		query.CurrentWxcontactSeq = *result.Data.(gin.H)["ContactList"].(*wechat.InitContactResp).CurrentWxcontactSeq
		query.CurrentChatRoomContactSeq = *result.Data.(gin.H)["ContactList"].(*wechat.InitContactResp).CurrentChatRoomContactSeq
		query.ContactUsernameList = result.Data.(gin.H)["ContactList"].(*wechat.InitContactResp).GetContactUsernameList()
		//i = query.ContactUsernameList长度
		i = len(query.ContactUsernameList)
		// 将query.ContactUsernameList添加到ContactList
		ContactList = append(ContactList, query.ContactUsernameList...)
		// 如果i没有或者i小于0，break
		if i < 100 {
			break
		}
	}
	// 将ContactList每二十个调用一次service.GetContactContactService(queryKey, *reqModel)并将结果添加到userList

	reqListModel := &model.BatchGetContactModel{
		UserNames:    make([]string, 0),
		RoomWxIDList: make([]string, 0),
	}
	userList := make([]*wechat.ModContact, 0)

	for i := 0; i < len(ContactList); i += 20 {
		end := i + 20

		// 防止索引越界
		if end > len(ContactList) {
			end = len(ContactList)
		}

		reqListModel.UserNames = ContactList[i:end]
		result := service.GetContactContactService(queryKey, *reqListModel)
		if contactList, ok := result.Data.(*wechat.GetContactResponse); ok {
			userList = append(userList, contactList.ContactList...)
		} else {
			// handle error
		}
	}

	ctx.JSON(http.StatusOK, userList)
}

//// GetFriendListApi 获取好友列表
//func GetFriendListApi(ctx *gin.Context) {
//	queryKey, isExist := ctx.GetQuery("key")
//	if !isExist || strings.Trim(queryKey, "") == "" {
//		//确保每次都有Key
//		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
//		return
//	}
//
//	result := service.GetFriendListService(queryKey)
//	ctx.JSON(http.StatusOK, result)
//}
//
//
//// GetGHListApi 获取好友列表
//func GetGHListApi(ctx *gin.Context) {
//	queryKey, isExist := ctx.GetQuery("key")
//	if !isExist || strings.Trim(queryKey, "") == "" {
//		//确保每次都有Key
//		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
//		return
//	}
//
//	result := service.GetGHListService(queryKey)
//	ctx.JSON(http.StatusOK, result)
//}

// FollowGHApi 关注公众号
func FollowGHApi(ctx *gin.Context) {

	reqModel := new(model.FollowGHModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}

	result := service.FollowGHService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// UploadMContactApi
func UploadMContactApi(ctx *gin.Context) {
	reqModel := new(model.UploadMContactModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}

	result := service.UploadMContactService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// GetMFriendApi
func GetMFriendApi(ctx *gin.Context) {
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	result := service.GetMFriendService(queryKey)
	ctx.JSON(http.StatusOK, result)
}

// 获取联系人详情
func GetContactContactApi(ctx *gin.Context) {
	reqModel := new(model.BatchGetContactModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.GetContactContactService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 获取好友关系
func GetFriendRelationApi(ctx *gin.Context) {
	reqModel := new(model.GetFriendRelationModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.GetFriendRelationService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 获取好友关系
func GetFriendRelationsApi(ctx *gin.Context) {
	reqModel := new(model.GetFriendRelationModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}
	if !validateData(ctx, &reqModel) {
		return
	}
	result := service.GetFriendRelationsService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// SearchContactRequestApi 搜索联系人
func SearchContactRequestApi(ctx *gin.Context) {
	reqModel := new(model.SearchContactRequestModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}

	result := service.SearchContactRequestService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// VerifyUserRequestApi 验证用户
// / <summary>
// / v1 v2操作
// / </summary>
// / <param name="opCode">1关注公众号2打招呼 主动添加好友3通过好友请求</param>
// / <param name="Content">内容</param>
// / <param name="antispamTicket">v2</param>
// / <param name="value">v1</param>
// / <param name="sceneList">1来源QQ2来源邮箱3来源微信号14群聊15手机号18附近的人25漂流瓶29摇一摇30二维码13来源通讯录</param>
// / <returns></returns>
func VerifyUserRequestApi(ctx *gin.Context) {
	reqModel := new(model.VerifyUserRequestModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}

	result := service.VerifyUserRequestService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}

// 同意好友请求
func AgreeAddApi(ctx *gin.Context) {
	reqModel := new(model.VerifyUserRequestModel)
	queryKey, isExist := ctx.GetQuery("key")
	if !isExist || strings.Trim(queryKey, "") == "" {
		//确保每次都有Key
		ctx.JSON(http.StatusOK, vo.NewFailUUId(""))
		return
	}

	if !validateData(ctx, &reqModel) {
		return
	}
	if reqModel.Scene == 0 {
		reqModel.Scene = 0x06
	}
	result := service.VerifyUserRequestService(queryKey, *reqModel)
	ctx.JSON(http.StatusOK, result)
}
