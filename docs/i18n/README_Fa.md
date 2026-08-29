<div align="center">

  <h1>🎬 YouTube Crawler</h1>

</div>

<p align="center">   <strong>🕊️ یادبود کشتار فجیعانه ایران در ۱۸ و ۱۹ دی ماه ۱۴۰۴</strong> </p>

---

<p align="center">
  <a href="https://github.com/Chamroosh98/YT-Crawler/releases"><img src="https://img.shields.io/badge/v1.0.0-181717?style=for-the-badge&logo=github&logoColor=white" alt="Version 1.0.0"></a>
  <a href="https://www.youtube.com/"><img src="https://img.shields.io/badge/YouTube-FF0000?style=for-the-badge&logo=youtube&logoColor=white" alt="YouTube"></a>
  <a href="https://go.dev/"><img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go"></a>
  <a href="https://github.com/features/actions"><img src="https://img.shields.io/badge/GitHub%20Actions-2088FF?style=for-the-badge&logo=githubactions&logoColor=white" alt="GitHub Actions"></a>
  <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/MIT-000000?style=for-the-badge&logo=opensourceinitiative&logoColor=white" alt="MIT License"></a>
  <a href="https://t.me/Chamroosh98"><img src="https://img.shields.io/badge/Telegram-26A5E4?style=for-the-badge&logo=telegram&logoColor=white" alt="Telegram"></a>
</p>

<p align="center">   <a href="../README.md"><strong>English</strong></a> </p>

---

### 📖 درباره پروژه

ابزار **YouTube Crawler** یه خزنده‌ی سبک و ماژولار برای یوتیوبه که با **Go** نوشته شده.

این ابزار با بهره از کلیدواژه های دلخواه، ویدیوهای موردنظرت رو پیدا می‌کنه، ویدیوهای جدید رو داخل پایگاه‌داده‌ی **SQLite** ذخیره می‌کنه و هر زمان دیتای جدیدی پیدا بشه، به وسیله **Telegram** بهت خبر می‌ده!

---

### ✨ امکانات

* 💾 ذخیره‌ی دیتا در **SQLite** همراه با جلوگیری از ذخیره ویدیوهای تکراری
* 🔔 فرستادن نوتیف در **Telegram** برای ویدیوهای جدید
* ⏰ اجرای خودکار با بهره از **GitHub Actions**
* 🔎 جست‌وجوی ویدیوها بر اساس کلیدواژه های دلخواه

---

### 🔑 دریافت داده موردنیاز

برای اجرای پروژه به سه تا توکن نیاز داری :

* `YOUTUBE_API_KEY`
* `TELEGRAM_BOT_TOKEN`
* `TELEGRAM_CHAT_ID`

#### 🎬 دریافت YouTube API Key

برای بهره از **YouTube Data API** باس یک پروژه در Google Cloud داشته باشی و برای اون API Key بسازی.

از صفحه‌ی رسمی Credentials گوگل وارد شو :

[Google Cloud — Credentials](https://console.cloud.google.com/apis/credentials)

اگه هنوز **YouTube Data API v3** رو برای پروژه اکتیو نکردی، از این صفحه اکتیوش کن:

[YouTube Data API v3 — Google Cloud](https://console.cloud.google.com/apis/library/youtube.googleapis.com)

پس از اکتیو سازی :

**Google Cloud → APIs & Services → Credentials → Create Credentials → API Key**

---

#### 🤖 دریافت Telegram Bot Token

برای ساخت ربات تلگرام، وارد **@BotFather** شو :

[Telegram BotFather](https://t.me/BotFather)

توکن دریافت‌شده رو قرار بده :

```bash
TELEGRAM_BOT_TOKEN=your_bot_token
```
---

#### 💬 دریافت Telegram Chat ID

برای `TELEGRAM_CHAT_ID` باید شناسه‌ی چتی که قراره نوتیف ها داخلش فرستاده بشن رو داشته باشی.

مقدار نهایی رو در `.env` قرار بده :

```bash
TELEGRAM_CHAT_ID=your_chat_id
```
---

### 🚀 روش اول : اجرای سریع

پس از اینکه دیتای موردنیاز رو آماده کردی، ریپو رو دریافت کن و فایل پیکره بندی رو بساز :

```bash
git clone https://github.com/Chamroosh98/YT-Crawler.git
cd YT-Crawler
cp .env.example .env
```

حالا `.env` رو باز کن و توکن های زیر رو وارد کن :

```bash
YOUTUBE_API_KEY=your_youtube_api_key
TELEGRAM_BOT_TOKEN=your_telegram_bot_token
TELEGRAM_CHAT_ID=your_telegram_chat_id
```

سپس برنامه رو اجرا کن :

```bash
go run ./cmd/main.go
```
---

### 🐙 روش دوم : GitHub Actions

اگه می‌خوای پروژه بدون نیاز به اجرای دستی روی سیستم خودت کار کنه، می‌تونی اجرای اون رو به **GitHub Actions** بسپری.

۱. ابتدا پروژه رو Fork کن .

۲. سپس از صفحه‌ی اصلی ریپو در GitHub برو به :

**Settings → Secrets and variables → Actions → Secrets**

۳. در نهایت سه Secret زیر رو بساز :

* `YOUTUBE_API_KEY`
* `TELEGRAM_BOT_TOKEN`
* `TELEGRAM_CHAT_ID`

و مقدار هرکدوم رو هم در Secret مربوط به خودش قرار بده!

---

### 📅 اجرای خودکار

این پروژه برای اجرای خودکار با **GitHub Actions** آماده شده.
Workflow برنامه هر **۶ ساعت یک‌بار** اجرا می‌شه و تنها زمانی ویدیوی **جدیدی** پیدا بشه، به وسیله Telegram بهت خبر می‌ده!
