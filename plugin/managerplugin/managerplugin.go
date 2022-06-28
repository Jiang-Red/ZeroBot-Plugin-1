// Package managerplugin 自定义群管插件
package managerplugin

import (
	"strconv"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/math"
)

func init() {
	engine := control.Register("managerplugin", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  true,
		Help:              "自定义的群管插件\n - 开启全员禁言 群号\n - 解除全员禁言 群号\n - 踢出并拉黑 QQ号\n - 反\"XX召唤术\"\n",
		PrivateDataFolder: "managerplugin",
	})
	// 指定开启某群全群禁言 Usage: 开启全员禁言123456
	engine.OnRegex(`^开启全员禁言.*?(\d+)`, zero.SuperUserPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SetGroupWholeBan(
				math.Str2Int64(ctx.State["regex_matched"].([]string)[1]),
				true,
			)
			ctx.SendChain(message.Text("全员自闭开始"))
		})
	// 指定解除某群全群禁言 Usage: 解除全员禁言123456
	engine.OnRegex(`^解除全员禁言.*?(\d+)`, zero.SuperUserPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SetGroupWholeBan(
				math.Str2Int64(ctx.State["regex_matched"].([]string)[1]),
				false,
			)
			ctx.SendChain(message.Text("全员自闭结束"))
		})
	engine.OnRegex(`^踢出并拉黑.*?(\d+)`, zero.OnlyGroup, zero.AdminPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			uid := math.Str2Int64(ctx.State["regex_matched"].([]string)[1])
			ctx.SetGroupKick(ctx.Event.GroupID, uid, true)
			nickname := ctx.CardOrNickName(uid)
			ctx.SendChain(message.Text("已将", nickname, "流放到边界外~"))
		})
	engine.OnRegex(`^\[CQ:xml`, zero.OnlyGroup, zero.KeywordRule("serviceID=\"60\"")).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			nickname := ctx.CardOrNickName(ctx.Event.UserID)
			ctx.SetGroupKick(ctx.Event.GroupID, ctx.Event.UserID, false)
			ctx.SetGroupBan(ctx.Event.GroupID, ctx.Event.UserID, 7*24*60*60)
			ctx.SendChain(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("检测到 ["+nickname+"]("+strconv.FormatInt(ctx.Event.UserID, 10)+") 发送了干扰性消息,已处理"))...)
			ctx.DeleteMessage(message.NewMessageIDFromInteger(ctx.Event.MessageID.(int64)))
		})
}
