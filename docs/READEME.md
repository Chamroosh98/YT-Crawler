<div align="center">

  <h1>🎬 YouTube Crawler</h1>

</div>

<p align="center">
  <strong>🕊️ Remembering the IRAN Massacre on Jan 8–9, 2026</strong>
</p>

---

<p align="center">
  <a href="https://github.com/Chamroosh98/YT-Crawler/releases">
    <img src="https://img.shields.io/badge/v1.0.0-181717?style=for-the-badge&logo=github&logoColor=white" alt="Version 1.0.0">
  </a>
  <a href="https://www.youtube.com/">
    <img src="https://img.shields.io/badge/YouTube-FF0000?style=for-the-badge&logo=youtube&logoColor=white" alt="YouTube">
  </a>
  <a href="https://go.dev/">
    <img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go">
  </a>
  <a href="https://github.com/features/actions">
    <img src="https://img.shields.io/badge/GitHub%20Actions-2088FF?style=for-the-badge&logo=githubactions&logoColor=white" alt="GitHub Actions">
  </a>
  <a href="https://opensource.org/licenses/MIT">
    <img src="https://img.shields.io/badge/MIT-000000?style=for-the-badge&logo=opensourceinitiative&logoColor=white" alt="MIT License">
  </a>
  <a href="https://t.me/Chamroosh98">
    <img src="https://img.shields.io/badge/Telegram-26A5E4?style=for-the-badge&logo=telegram&logoColor=white" alt="Telegram">
  </a>
</p>

<p align="center">
  <a href="i18n/README_Fa.md"><strong>Persian</strong></a>
</p>

---

### 📖 Overview

**YouTube Crawler** is a clean and modular YouTube crawler written in **Go**.

It searches for videos using customizable queries, stores newly discovered videos in a **SQLite** database, and sends notifications through **Telegram** whenever new content is found.

---

### ✨ Features

* 💾 **SQLite storage** with duplicate prevention
* 🔔 **Telegram notifications** for newly discovered videos
* ⏰ **Automated execution** with GitHub Actions
* 🔎 **Custom search queries** for targeted video discovery

---

### 🔑 Required Credentials

Before running the crawler, you'll need the following three values:

* `YOUTUBE_API_KEY`
* `TELEGRAM_BOT_TOKEN`
* `TELEGRAM_CHAT_ID`

#### 🎬 Get a YouTube API Key

The crawler uses the **YouTube Data API** to search for videos.

First, create or select a project in Google Cloud Console:

[Google Cloud — Credentials](https://console.cloud.google.com/apis/credentials)

If **YouTube Data API v3** isn't enabled for your project yet, enable it here:

[YouTube Data API v3 — Google Cloud](https://console.cloud.google.com/apis/library/youtube.googleapis.com)

Then go to:

**Google Cloud → APIs & Services → Credentials → Create Credentials → API Key**

---

#### 🤖 Get a Telegram Bot Token

To create a Telegram bot, open **@BotFather**:

[Telegram BotFather](https://t.me/BotFather)

Add the token to your `.env` file :

```bash 
TELEGRAM_BOT_TOKEN=your_bot_token
```

---

#### 💬 Get the Telegram Chat ID

The `TELEGRAM_CHAT_ID` identifies the chat where the crawler will send notifications.

Add the Chat ID to your `.env` file :

```bash 
TELEGRAM_CHAT_ID=your_chat_id
```
---

### 🚀 Option 1: Quick Start

Once you have the required credentials, clone the repository and create your `.env` file :

```bash 
git clone https://github.com/Chamroosh98/YT-Crawler.git
cd YT-Crawler
cp .env.example .env
```

Open `.env` and add your credentials :

```bash 
YOUTUBE_API_KEY=your_youtube_api_key
TELEGRAM_BOT_TOKEN=your_telegram_bot_token
TELEGRAM_CHAT_ID=your_telegram_chat_id
```

Then start the crawler :

```bash
go run ./cmd/main.go
```
---

### 🐙 Option 2: GitHub Actions

If you want the crawler to run automatically without keeping your own machine online, you can use **GitHub Actions**.

1. **Fork** the YT-Crawler repository .

2. Then from your GitHub repository, go to :

    **Settings → Secrets and variables → Actions → Secrets**

3. Finally Create these three repository secrets :

* `YOUTUBE_API_KEY`
* `TELEGRAM_BOT_TOKEN`
* `TELEGRAM_CHAT_ID`

and Set each secret to its corresponding value!

---

### 📅 Automation

The project is configured for automated execution through **GitHub Actions**.

The workflow runs every **6 hours** and sends a Telegram notification only when **new videos** are discovered.

