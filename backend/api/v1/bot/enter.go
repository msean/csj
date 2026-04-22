package bot

import "github.com/msean/csj/backend/service"

var (
	botChatGroupService   = service.ServiceGroupApp.BotServiceGroup.BotChatGroupService
	botBanGroupMemService = service.ServiceGroupApp.BotServiceGroup.BotBanGroupMemService
	taskService           = service.ServiceGroupApp.BotServiceGroup.BotTaskService
	botChannelService     = service.ServiceGroupApp.BotServiceGroup.BotChannelService
	botCmdConfigService   = service.ServiceGroupApp.BotServiceGroup.BotCmdConfigService
	botMsgMassService     = service.ServiceGroupApp.BotServiceGroup.BotMsgMassService
)
