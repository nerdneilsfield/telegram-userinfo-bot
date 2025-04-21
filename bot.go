// Package main implements a Telegram bot that provides chat and user information.
// 主包实现了一个提供聊天和用户信息的 Telegram 机器人。
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/mymmrac/telego"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
	// "github.com/mymmrac/telego/telegoext" // Optional, for more convenient handlers / 可选，用于更方便的 handler
)

// main is the entry point of the application.
// main 是应用程序的入口点。
func main() {
	// 1. Get Bot Token from environment variable / 从环境变量获取 Bot Token
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("Error: Please set the TELEGRAM_BOT_TOKEN environment variable / 错误：请设置环境变量 TELEGRAM_BOT_TOKEN")
	}

	// 2. (Optional) Get HTTPS proxy address from environment variable / （可选）从环境变量获取 HTTPS 代理地址
	httpsProxy := os.Getenv("HTTPS_PROXY")

	// --- Initialize Bot / 初始化 Bot ---
	// Create Bot options list / 创建 Bot 选项列表
	botOptions := []telego.BotOption{}

	// If HTTPS_PROXY is set, configure proxy / 如果设置了 HTTPS_PROXY，则配置代理
	if httpsProxy != "" {
		log.Printf("HTTPS_PROXY detected: %s, configuring proxy... / 检测到 HTTPS_PROXY: %s，正在配置代理...", httpsProxy, httpsProxy)
		// proxyURL variable is no longer needed, parsing step removed / proxyURL 变量不再需要，移除解析步骤
		// proxyURL, err := url.Parse(httpsProxy)
		// if err != nil {
		// 	log.Fatalf("Unable to parse HTTPS_PROXY URL '%s': %v / 无法解析 HTTPS_PROXY URL '%s': %v", httpsProxy, err)
		// }

		// Create fasthttp client with proxy / 创建使用代理的 fasthttp client
		customClient := &fasthttp.Client{
			// Configure Dial method with fasthttpproxy / 使用 fasthttpproxy 配置 Dial 方法
			Dial: fasthttpproxy.FasthttpHTTPDialer(httpsProxy),
			// Can configure other fasthttp.Client parameters as needed, such as timeout / 可以根据需要配置其他 fasthttp.Client 参数，例如超时
			// ReadTimeout:  5 * time.Second,
			// WriteTimeout: 5 * time.Second,
		}

		// Add custom client to Bot options / 将自定义 client 添加到 Bot 选项中
		botOptions = append(botOptions, telego.WithFastHTTPClient(customClient))
		log.Println("Fasthttp client configured to use proxy. / 已配置 fasthttp 客户端以使用代理。")
	}

	// Optional: Add other Bot options, such as logger / 可选：添加其他 Bot 选项，例如 logger
	// botOptions = append(botOptions, telego.WithLogger(&telego.LeveledLogger{MinLevel: telego.LogLevelDebug}))

	// Create Bot instance with Token and options / 使用 Token 和选项创建 Bot 实例
	bot, err := telego.NewBot(botToken, botOptions...)
	if err != nil {
		log.Fatalf("Cannot create bot: %s / 无法创建 bot: %s", err, err)
	}

	// Get Bot information / 获取 Bot 信息
	botUser, err := bot.GetMe(context.Background())
	if err != nil {
		log.Fatalf("Cannot get bot information: %s / 无法获取 bot 信息: %s", err, err)
	}
	log.Printf("Bot authorized as: %s (@%s) / 机器人已授权为: %s (@%s)", botUser.FirstName, botUser.Username, botUser.FirstName, botUser.Username)

	// --- Set up bot commands / 设置机器人命令 ---
	commands := []telego.BotCommand{
		{
			Command:     "userinfo",
			Description: "Get current chat and user information / 获取当前聊天和用户信息", // Command description / 命令描述
		},
		// Add other commands here if needed / 如果有其他命令，可以在这里添加
	}
	setCmdParams := &telego.SetMyCommandsParams{
		Commands: commands,
		// Scope:    telego.BotCommandScopeDefault{}, // Can specify scope, default is all private chats / 可以指定范围，默认为所有私聊
		// LanguageCode: "", // Can specify language / 可以指定语言
	}
	err = bot.SetMyCommands(context.Background(), setCmdParams)
	if err != nil {
		log.Printf("Warning: Cannot set bot commands: %s / 警告：无法设置机器人命令: %s", err, err) // Non-fatal error, just log it / 非致命错误，记录日志即可
	} else {
		log.Println("Bot commands successfully set. / 机器人命令已成功设置。")
	}

	// --- Get updates / 获取更新 ---
	ctx, cancel := context.WithCancel(context.Background())
	updates, err := bot.UpdatesViaLongPolling(ctx, nil)
	if err != nil {
		log.Fatalf("Cannot start receiving updates: %s / 无法开始接收更新: %s", err, err)
	}

	// --- Gracefully stop the Bot / 优雅地停止 Bot ---
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs

		log.Println("Received stop signal, shutting down... / 收到停止信号，正在关闭...")
		cancel()
		log.Println("Long polling cancel function called. / Long polling cancel 函数已调用。")
	}()

	log.Println("Bot started, listening for messages... / 机器人已启动，正在监听消息...")
	for update := range updates {
		if update.Message == nil {
			continue
		}

		message := update.Message

		if !strings.HasPrefix(message.Text, "/userinfo") {
			continue
		}

		chat := message.Chat
		chatInfo := fmt.Sprintf("--- Chat Information (telego) / 聊天信息 (telego) ---\n"+
			"Type / 类型: %s\n"+
			"ID: %d\n"+
			"Title / 标题: %s\n"+
			"Username / 用户名: @%s\n"+
			"Is Forum/Topics Enabled / 是否启用话题: %s\n",
			chat.Type,
			chat.ID,
			chat.Title,
			chat.Username,
			strconv.FormatBool(chat.IsForum),
		)

		topicInfo := ""
		if message.IsTopicMessage || (chat.IsForum && message.MessageThreadID != 0) {
			topicInfo = fmt.Sprintf("\n--- Topic Info / 话题信息 ---\n"+
				"Is Topic Message / 是否在话题内: %s\n"+
				"Thread ID / 话题 ID: %d\n"+
				"Note: Topic name needs to be obtained through other means. / 注意：话题名称需通过其他方式获取。\n",
				strconv.FormatBool(message.IsTopicMessage),
				message.MessageThreadID,
			)
		}

		user := message.From
		userInfo := ""
		if user != nil {
			userInfo = fmt.Sprintf("\n--- Sender Information / 发起用户信息 ---\n"+
				"User ID / 用户 ID: %d\n"+
				"First Name / 名字: %s\n"+
				"Last Name / 姓氏: %s\n"+
				"Username / 用户名: @%s\n"+
				"Is Bot / 是否机器人: %s\n"+
				"Language Code / 语言代码: %s\n",
				user.ID,
				user.FirstName,
				user.LastName,
				user.Username,
				strconv.FormatBool(user.IsBot),
				user.LanguageCode,
			)
		} else {
			userInfo = "\n--- Sender Information / 发起用户信息 ---\n" +
				"Cannot get user information (possibly anonymous admin or channel message) / " +
				"无法获取用户信息 (可能是匿名管理员或频道消息)\n"
		}

		responseText := chatInfo + topicInfo + userInfo

		params := &telego.SendMessageParams{
			ChatID: telego.ChatID{ID: chat.ID},
			Text:   responseText,
		}

		if message.IsTopicMessage || (chat.IsForum && message.MessageThreadID != 0) {
			params.MessageThreadID = message.MessageThreadID
		}

		_, err := bot.SendMessage(context.Background(), params)
		if err != nil {
			log.Printf("Failed to send message (ChatID: %d): %s / 发送消息失败 (ChatID: %d): %s", chat.ID, err, chat.ID, err)
		}
	}

	log.Println("Bot has been shut down. / 机器人已关闭。")
}
