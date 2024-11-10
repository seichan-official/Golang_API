# 🎶 Spotube 🎥

**Spotube**は、Spotifyの再生履歴をもとに、関連するYouTubeのミュージックビデオを検索できるWebアプリです！  
バックエンドはGo言語、フロントエンドはReactで実装されています。

## 🎬 デモ動画
[app-demo-video.mp4](リンクを挿入)

---

## ✨ 特徴
- **Spotify履歴の自動取得**  
   Spotify APIを使ってユーザーの再生履歴を取得します。
- **YouTubeでミュージックビデオを検索**  
   履歴に基づきYouTube APIで関連動画を検索し、リンクを提供します。
- **高速なバックエンド**  
   Go言語で構築された軽量バックエンドにより、快適な動作を実現！

---

## 🛠️ 技術スタック

### 🌐 フロントエンド
- React
- HTML5
- CSS

### 🖥️ バックエンド
- Go 1.23.2

---

## 🚀 セットアップ手順

### 1️⃣ 前提条件
- **SpotifyとYouTube APIのAPIキー**が必要です。

### 2️⃣ APIキーの取得と設定
- **Spotify APIキー**を[Spotify for Developers](https://developer.spotify.com)から取得
- **YouTube APIキー**を[Google Cloud Console](https://console.cloud.google.com/)から取得
- プロジェクトルート直下に`.env`ファイルを作成し、以下を記入してください：
    ```plaintext
    SPOTIFY_CLIENT_ID=<spotify_client_id>
    SPOTIFY_CLIENT_SECRET=<spotify_client_secret>
    YOUTUBE_API_KEY=<youtube_api_key>
    ```

### 3️⃣ アプリケーションの起動

#### 🖥️ バックエンドの起動
```bash
cd back
go run main.go spotify.go youtube.go
```
バックエンドはデフォルトで`http://localhost:8080`で起動します。

#### 🌐 フロントエンドの起動
1. フロントエンドディレクトリに移動
2. 以下のコマンドで依存関係をインストールし、サーバーを起動
    ```bash
    npm install
    npm start
    ```
フロントエンドはデフォルトで`http://localhost:3000`で起動します。

---

## 🎵 機能説明

### ミュージックビデオ検索
- 🎶 Spotifyの再生履歴から楽曲情報を取得し、その楽曲に関連するYouTubeのミュージックビデオを表示します！
 

---

## 👤 著者
- **ctake099**