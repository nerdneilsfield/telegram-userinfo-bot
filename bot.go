// Package main implements a Telegram bot that provides chat and user information.
// 主包实现了一个提供聊天和用户信息的 Telegram 机器人。
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/mymmrac/telego"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
)

// config holds all runtime configuration loaded from environment variables.
// config 保存从环境变量加载的所有运行时配置。
type config struct {
	BotToken   string // TELEGRAM_BOT_TOKEN
	HTTPSProxy string // HTTPS_PROXY（可选）
	MaxWorkers int    // 最大并发处理数，默认 20
	SendTimeout time.Duration // 单次 SendMessage 超时，默认 10s
}

// loadConfig reads configuration from environment variables and applies defaults.
// loadConfig 从环境变量读取配置并应用默认值。
func loadConfig() config {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		slog.Error("Missing required environment variable", "var", "TELEGRAM_BOT_TOKEN")
		os.Exit(1)
	}
	return config{
		BotToken:    token,
		HTTPSProxy:  os.Getenv("HTTPS_PROXY"),
		MaxWorkers:  20,
		SendTimeout: 10 * time.Second,
	}
}

// mdEscape escapes special characters for Telegram MarkdownV2.
// mdEscape 转义 Telegram MarkdownV2 的特殊字符。
func mdEscape(s string) string {
	replacer := strings.NewReplacer(
		`_`, `\_`, `*`, `\*`, `[`, `\[`, `]`, `\]`,
		`(`, `\(`, `)`, `\)`, `~`, `\~`, "`", "\\`",
		`>`, `\>`, `#`, `\#`, `+`, `\+`, `-`, `\-`,
		`=`, `\=`, `|`, `\|`, `{`, `\{`, `}`, `\}`,
		`.`, `\.`, `!`, `\!`,
	)
	return replacer.Replace(s)
}

// boolEmoji converts a bool to a checkmark or cross emoji.
// boolEmoji 将布尔值转换为对勾或叉号 emoji。
func boolEmoji(b bool) string {
	if b {
		return "✅"
	}
	return "❌"
}

// main is the entry point of the application.
// main 是应用程序的入口点。
func main() {
	// 使用 JSON 格式结构化日志，便于日志聚合系统（如 Loki、ELK）解析
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	cfg := loadConfig()

	// --- Initialize Bot / 初始化 Bot ---
	botOptions := []telego.BotOption{}

	if cfg.HTTPSProxy != "" {
		slog.Info("Configuring HTTPS proxy", "proxy", cfg.HTTPSProxy)
		customClient := &fasthttp.Client{
			Dial: fasthttpproxy.FasthttpHTTPDialer(cfg.HTTPSProxy),
		}
		botOptions = append(botOptions, telego.WithFastHTTPClient(customClient))
	}

	bot, err := telego.NewBot(cfg.BotToken, botOptions...)
	if err != nil {
		slog.Error("Cannot create bot", "error", err)
		os.Exit(1)
	}

	botUser, err := bot.GetMe(context.Background())
	if err != nil {
		slog.Error("Cannot get bot information", "error", err)
		os.Exit(1)
	}
	botUsername := strings.ToLower(botUser.Username)
	slog.Info("Bot authorized", "name", botUser.FirstName, "username", botUser.Username)

	// --- Set up bot commands / 设置机器人命令 ---
	err = bot.SetMyCommands(context.Background(), &telego.SetMyCommandsParams{
		Commands: []telego.BotCommand{
			{Command: "userinfo", Description: "Get current chat and user information"},
		},
	})
	if err != nil {
		slog.Warn("Cannot set bot commands", "error", err)
	} else {
		slog.Info("Bot commands set successfully")
	}

	// --- Get updates / 获取更新 ---
	ctx, cancel := context.WithCancel(context.Background())
	// Timeout=30 启用真正的 long polling：Telegram 有消息时立刻推送，无需反复轮询
	updates, err := bot.UpdatesViaLongPolling(ctx, &telego.GetUpdatesParams{Timeout: 30})
	if err != nil {
		slog.Error("Cannot start receiving updates", "error", err)
		os.Exit(1)
	}

	// Worker pool：最多同时处理 maxWorkers 条消息，防止突发流量下 goroutine 无限膨胀
	sem := make(chan struct{}, cfg.MaxWorkers)

	// WaitGroup 追踪所有 in-flight goroutine，确保优雅关闭时不丢消息
	var wg sync.WaitGroup

	// --- Gracefully stop the Bot / 优雅地停止 Bot ---
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigs
		slog.Info("Received stop signal, shutting down", "signal", sig.String())
		cancel() // 停止 long polling，updates channel 随之关闭
	}()

	slog.Info("Bot started, listening for messages")
	for update := range updates {
		if update.Message == nil {
			continue
		}

		message := update.Message

		// 兼容群组里的 /userinfo@botname 格式
		cmd := strings.ToLower(strings.SplitN(message.Text, " ", 2)[0])
		if cmd != "/userinfo" && cmd != "/userinfo@"+botUsername {
			continue
		}

		// 占用一个 worker slot（满载时阻塞，背压）
		sem <- struct{}{}
		wg.Add(1)

		go func(msg *telego.Message) {
			defer wg.Done()
			defer func() { <-sem }()

			chat := msg.Chat
			log := slog.With("chat_id", chat.ID, "msg_id", msg.MessageID)

			// 每条消息最多等 SendTimeout，避免 Telegram API 超时挂死 goroutine
			sendCtx, sendCancel := context.WithTimeout(context.Background(), cfg.SendTimeout)
			defer sendCancel()

			// 构建 Chat 信息块（MarkdownV2 格式）
			usernameStr := ""
			if chat.Username != "" {
				usernameStr = fmt.Sprintf("\nUsername: @%s", mdEscape(chat.Username))
			}
			chatInfo := fmt.Sprintf("💬 *Chat Info*\nType: `%s`\nID: `%d`\nTitle: %s%s\nForum: %s",
				mdEscape(string(chat.Type)),
				chat.ID,
				mdEscape(chat.Title),
				usernameStr,
				boolEmoji(chat.IsForum),
			)

			// 构建 Topic 信息块
			topicInfo := ""
			if msg.IsTopicMessage || (chat.IsForum && msg.MessageThreadID != 0) {
				topicInfo = fmt.Sprintf("\n\n📌 *Topic Info*\nThread ID: `%d`",
					msg.MessageThreadID,
				)
			}

			// 构建 User 信息块
			user := msg.From
			userInfo := ""
			if user != nil {
				fullName := strings.TrimSpace(user.FirstName + " " + user.LastName)
				senderUsername := ""
				if user.Username != "" {
					senderUsername = fmt.Sprintf("\nUsername: @%s", mdEscape(user.Username))
				}
				langStr := ""
				if user.LanguageCode != "" {
					langStr = fmt.Sprintf("\nLang: `%s`", mdEscape(user.LanguageCode))
				}
				userInfo = fmt.Sprintf("\n\n👤 *Sender Info*\nName: %s\nID: `%d`%s\nBot: %s%s",
					mdEscape(fullName),
					user.ID,
					senderUsername,
					boolEmoji(user.IsBot),
					langStr,
				)
			} else {
				userInfo = "\n\n👤 *Sender Info*\n_Cannot get user info \\(anonymous admin or channel\\)_"
			}

			params := &telego.SendMessageParams{
				ChatID:    telego.ChatID{ID: chat.ID},
				Text:      chatInfo + topicInfo + userInfo,
				ParseMode: telego.ModeMarkdownV2,
			}
			if msg.IsTopicMessage || (chat.IsForum && msg.MessageThreadID != 0) {
				params.MessageThreadID = msg.MessageThreadID
			}

			if _, err := bot.SendMessage(sendCtx, params); err != nil {
				log.Error("Failed to send message", "error", err)
			} else {
				log.Info("Replied to /userinfo")
			}
		}(message)
	}

	// 等待所有 in-flight 消息发完再退出
	wg.Wait()
	slog.Info("Bot has been shut down")
}
