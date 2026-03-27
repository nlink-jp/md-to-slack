# md-to-slack

Markdown を [Slack Block Kit](https://api.slack.com/block-kit) JSON に変換するフィルターツール。

標準入力から Markdown を読み込み、`{"blocks": [...]}` 形式の JSON を標準出力に書き出します。
Slack API や [slackcat](https://github.com/bcicen/slackcat) などのツールと組み合わせて使用できます。

## 機能

- GFM（GitHub Flavored Markdown）対応
- H1/H2 → Slack ヘッダーブロック（大きいフォント）
- H3–H6 → 太字セクションブロック
- 段落、ブロック引用、順序付き・順序なしリスト（ネスト対応）
- フェンスおよびインデント付きコードブロック
- GFM テーブル → 整形済みプレーンテキストのコードブロック
- 単独画像 → Slack image ブロック
- インライン画像 → リンクにフォールバック
- 取り消し線、太字、斜体、インラインコード、リンク、自動リンク
- 水平罫線 → Slack divider ブロック
- 生 HTML は破棄（Slack では表示不可）

## インストール

```bash
go install github.com/nlink-jp/md-to-slack/cmd/md-to-slack@latest
```

または [Releases](https://github.com/nlink-jp/md-to-slack/releases) ページからビルド済みバイナリをダウンロードしてください。

## 使い方

```bash
md-to-slack < README.md
echo "# こんにちは **世界**" | md-to-slack
```

### Slack への送信例

```bash
md-to-slack < message.md | curl -s \
  -X POST https://slack.com/api/chat.postMessage \
  -H "Authorization: Bearer $SLACK_TOKEN" \
  -H "Content-Type: application/json" \
  -d @- -d '{"channel":"#general"}'
```

### フラグ

| フラグ | 説明 |
|--------|------|
| `--version`, `-V` | バージョンを表示して終了 |
| `--help`, `-h` | ヘルプを表示して終了 |

## ビルド

```bash
make build       # 現在のプラットフォーム向けにビルド
make build-all   # 全プラットフォーム向けにクロスコンパイル
make test        # テスト実行
make check       # vet + lint + test + build
```

## 出力形式

Slack Block Kit ペイロードを出力します:

```json
{
  "blocks": [
    { "type": "header", "text": { "type": "plain_text", "text": "タイトル", "emoji": true } },
    { "type": "section", "text": { "type": "mrkdwn", "text": "本文テキスト。" } }
  ]
}
```
