# Telegram UserInfo Bot / Telegram 用户信息机器人

[![Docker Hub](https://img.shields.io/docker/pulls/nerdneils/telegram-userinfo-bot.svg)](https://hub.docker.com/r/nerdneils/telegram-userinfo-bot)
[![GitHub Container Registry](https://ghcr-badge.egpl.dev/nerdneilsfield/telegram-userinfo-bot/size?color=green)](https://github.com/nerdneilsfield/telegram-userinfo-bot/pkgs/container/telegram-userinfo-bot)

[English](#english) | [中文](#中文)

---

## English

A simple Telegram bot that replies with chat and user information when you send the `/userinfo` command.

### Features

* Responds to `/userinfo` command.
* Displays information about the current chat (type, ID, title, username, forum status).
* Displays information about the message sender (ID, first name, last name, username, bot status, language code).
* Supports topic messages in forum-enabled groups.
* Supports HTTPS proxy via `HTTPS_PROXY` environment variable.

### Usage

1. **Add the bot to your chat or start a private chat with it.**
2. **Send the `/userinfo` command.**
3. The bot will reply with the current chat and your user information.

### Setup

#### Prerequisites

* [Telegram Bot Token](https://core.telegram.org/bots#how-do-i-create-a-bot): Get one from [@BotFather](https://t.me/BotFather).

#### Using Docker (Recommended)

Replace `<YOUR_BOT_TOKEN>` with your actual Telegram Bot Token.

* **From Docker Hub:**

    ```bash
    docker run -d --name userinfo-bot \
      -e TELEGRAM_BOT_TOKEN=<YOUR_BOT_TOKEN> \
      nerdneils/telegram-userinfo-bot
    ```

* **From GitHub Container Registry (GHCR):**

    ```bash
    docker run -d --name userinfo-bot \
      -e TELEGRAM_BOT_TOKEN=<YOUR_BOT_TOKEN> \
      ghcr.io/nerdneilsfield/telegram-userinfo-bot
    ```

* **(Optional) Using an HTTPS Proxy:**
    Set the `HTTPS_PROXY` environment variable. Replace `<YOUR_PROXY_URL>` with your proxy address (e.g., `http://user:pass@host:port` or `socks5://user:pass@host:port`).

  * **From Docker Hub:**

        ```bash
        docker run -d --name userinfo-bot \
          -e TELEGRAM_BOT_TOKEN=<YOUR_BOT_TOKEN> \
          -e HTTPS_PROXY=<YOUR_PROXY_URL> \
          nerdneils/telegram-userinfo-bot
        ```

  * **From GHCR:**

        ```bash
        docker run -d --name userinfo-bot \
          -e TELEGRAM_BOT_TOKEN=<YOUR_BOT_TOKEN> \
          -e HTTPS_PROXY=<YOUR_PROXY_URL> \
          ghcr.io/nerdneilsfield/telegram-userinfo-bot
        ```

#### Using Docker Compose

Create a `docker-compose.yml` file:

* **Basic (Docker Hub):**
    ```yaml
    services:
      userinfo-bot:
        image: nerdneils/telegram-userinfo-bot
        restart: unless-stopped
        environment:
          - TELEGRAM_BOT_TOKEN=<YOUR_BOT_TOKEN>
        container_name: userinfo-bot
    ```

* **Basic (GHCR):**
    ```yaml
    services:
      userinfo-bot:
        image: ghcr.io/nerdneilsfield/telegram-userinfo-bot
        restart: unless-stopped
        environment:
          - TELEGRAM_BOT_TOKEN=<YOUR_BOT_TOKEN>
        container_name: userinfo-bot
    ```

* **(Optional) With HTTPS Proxy (Docker Hub):**
    ```yaml
    version: '3.8'
    services:
      userinfo-bot:
        image: nerdneils/telegram-userinfo-bot
        restart: unless-stopped
        environment:
          - TELEGRAM_BOT_TOKEN=<YOUR_BOT_TOKEN>
          - HTTPS_PROXY=<YOUR_PROXY_URL>
        container_name: userinfo-bot
    ```

* **(Optional) With HTTPS Proxy (GHCR):**
    ```yaml
    version: '3.8'
    services:
      userinfo-bot:
        image: ghcr.io/nerdneilsfield/telegram-userinfo-bot
        restart: unless-stopped
        environment:
          - TELEGRAM_BOT_TOKEN=<YOUR_BOT_TOKEN>
          - HTTPS_PROXY=<YOUR_PROXY_URL>
        container_name: userinfo-bot
    ```

Replace `<YOUR_BOT_TOKEN>` and (optionally) `<YOUR_PROXY_URL>` with your actual values.

Then run:
```bash
docker-compose up -d
```

#### Building from Source

1. Clone the repository:

    ```bash
    git clone https://github.com/nerdneilsfield/telegram-userinfo-bot.git
    cd telegram-userinfo-bot
    ```

2. Set environment variables:

    ```bash
    export TELEGRAM_BOT_TOKEN=<YOUR_BOT_TOKEN>
    # Optional:
    # export HTTPS_PROXY=<YOUR_PROXY_URL>
    ```

3. Build and run:

    ```bash
    go build .
    ./telegram-userinfo-bot
    ```

### Source Code

[https://github.com/nerdneilsfield/telegram-userinfo-bot](https://github.com/nerdneilsfield/telegram-userinfo-bot)

---

## 中文

一个简单的 Telegram 机器人，当你发送 `/userinfo` 命令时，它会回复聊天和用户信息。

### 功能

* 响应 `/userinfo` 命令。
* 显示当前聊天信息（类型、ID、标题、用户名、是否启用话题）。
* 显示消息发送者信息（ID、名字、姓氏、用户名、是否机器人、语言代码）。
* 支持启用话题的群组中的话题消息。
* 通过 `HTTPS_PROXY` 环境变量支持 HTTPS 代理。

### 使用方法

1. **将机器人添加到您的聊天中或与它开始私聊。**
2. **发送 `/userinfo` 命令。**
3. 机器人将回复当前的聊天和您的用户信息。

### 设置

#### 先决条件

* [Telegram Bot Token](https://core.telegram.org/bots#how-do-i-create-a-bot)：从 [@BotFather](https://t.me/BotFather) 获取。

#### 使用 Docker（推荐）

将 `<YOUR_BOT_TOKEN>` 替换为你的实际 Telegram Bot Token。

* **从 Docker Hub:**

    ```bash
    docker run -d --name userinfo-bot \
      -e TELEGRAM_BOT_TOKEN=<YOUR_BOT_TOKEN> \
      nerdneils/telegram-userinfo-bot
    ```

* **从 GitHub Container Registry (GHCR):**

    ```bash
    docker run -d --name userinfo-bot \
      -e TELEGRAM_BOT_TOKEN=<YOUR_BOT_TOKEN> \
      ghcr.io/nerdneilsfield/telegram-userinfo-bot
    ```

* **(可选) 使用 HTTPS 代理:**
    设置 `HTTPS_PROXY` 环境变量。将 `<YOUR_PROXY_URL>` 替换为你的代理地址（例如 `http://user:pass@host:port` 或 `socks5://user:pass@host:port`）。

  * **从 Docker Hub:**

        ```bash
        docker run -d --name userinfo-bot \
          -e TELEGRAM_BOT_TOKEN=<YOUR_BOT_TOKEN> \
          -e HTTPS_PROXY=<YOUR_PROXY_URL> \
          nerdneils/telegram-userinfo-bot
        ```

  * **从 GHCR:**

        ```bash
        docker run -d --name userinfo-bot \
          -e TELEGRAM_BOT_TOKEN=<YOUR_BOT_TOKEN> \
          -e HTTPS_PROXY=<YOUR_PROXY_URL> \
          ghcr.io/nerdneilsfield/telegram-userinfo-bot
        ```

#### 使用 Docker Compose

创建一个 `docker-compose.yml` 文件：

* **基础 (Docker Hub):**
    ```yaml
    version: '3.8'
    services:
      userinfo-bot:
        image: nerdneils/telegram-userinfo-bot
        restart: unless-stopped
        environment:
          - TELEGRAM_BOT_TOKEN=<YOUR_BOT_TOKEN>
        container_name: userinfo-bot
    ```

* **基础 (GHCR):**
    ```yaml
    version: '3.8'
    services:
      userinfo-bot:
        image: ghcr.io/nerdneilsfield/telegram-userinfo-bot
        restart: unless-stopped
        environment:
          - TELEGRAM_BOT_TOKEN=<YOUR_BOT_TOKEN>
        container_name: userinfo-bot
    ```

* **(可选) 使用 HTTPS 代理 (Docker Hub):**
    ```yaml
    version: '3.8'
    services:
      userinfo-bot:
        image: nerdneils/telegram-userinfo-bot
        restart: unless-stopped
        environment:
          - TELEGRAM_BOT_TOKEN=<YOUR_BOT_TOKEN>
          - HTTPS_PROXY=<YOUR_PROXY_URL>
        container_name: userinfo-bot
    ```

* **(可选) 使用 HTTPS 代理 (GHCR):**
    ```yaml
    version: '3.8'
    services:
      userinfo-bot:
        image: ghcr.io/nerdneilsfield/telegram-userinfo-bot
        restart: unless-stopped
        environment:
          - TELEGRAM_BOT_TOKEN=<YOUR_BOT_TOKEN>
          - HTTPS_PROXY=<YOUR_PROXY_URL>
        container_name: userinfo-bot
    ```

将 `<YOUR_BOT_TOKEN>` 和（可选的）`<YOUR_PROXY_URL>` 替换为你的实际值。

然后运行：
```bash
docker-compose up -d
```

#### 从源代码构建

1. 克隆仓库：

    ```bash
    git clone https://github.com/nerdneilsfield/telegram-userinfo-bot.git
    cd telegram-userinfo-bot
    ```

2. 设置环境变量：

    ```bash
    export TELEGRAM_BOT_TOKEN=<YOUR_BOT_TOKEN>
    # 可选:
    # export HTTPS_PROXY=<YOUR_PROXY_URL>
    ```

3. 构建并运行：

    ```bash
    go build .
    ./telegram-userinfo-bot
    ```

### 源代码

[https://github.com/nerdneilsfield/telegram-userinfo-bot](https://github.com/nerdneilsfield/telegram-userinfo-bot)
