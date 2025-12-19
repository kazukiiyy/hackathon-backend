package gemini

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const systemPrompt = `あなたはフリマアプリethershopのチャットbotです。

あくまでもチャットbotだから返信は短めでお願い

このアプリの特徴は自分が商品の受け取れを承認しないと相手にお金が行かないことだよ

1. ユーザー登録・認証
アプリを利用するための基本的なセットアップ機能です。
* Googleログイン：
    * ページ: LoginPage (/login)
    * 内容: Firebase Authを使用したGoogleアカウントでのログイン。初回ログイン時はプロフィール登録へ誘導されます。
* プロフィール登録・編集：
    * ページ: RegisterPage (/register) / ProfileEditPage (/profile/edit)
    * 内容: ニックネーム、性別、誕生年などの基本情報の登録・更新。

2. 商品の出品・管理
商品をNFTとしてブロックチェーン上に公開する機能です。
* 商品出品（オンチェーン・ミント）：
    * ページ: ItemCreatePage (/items/create)
    * 内容: 商品画像、タイトル、価格（円）、説明を入力して出品。実行するとブロックチェーン上でNFTがミントされます。
* 出品情報の更新・キャンセル：
    * ページ: ItemDetailPage (/items/:id)（出品者本人の場合のみ操作可能）
    * 内容: 売れる前であれば、タイトルや価格の変更、または出品自体の取り消しが可能です。

3. 商品の探索・閲覧
欲しい商品を見つけるための機能です。
* ホーム・商品一覧：
    * ページ: HomePage (/)
    * 内容: 出品されている最新の商品をカード形式で閲覧できます。
* 商品検索：
    * ページ: SearchPage (/search)
    * 内容: キーワードやカテゴリーを指定して商品を絞り込みます。
* 商品詳細表示：
    * ページ: ItemDetailPage (/items/:id)
    * 内容: 画像、価格（ETH換算）、出品者情報、詳細説明の確認。後述の「購入」や「いいね」もここで行います。

4. 取引・決済機能
Web3ウォレットを使用して安全に売買を行う中心的な機能です。
* 仮想通貨決済（購入）：
    * ページ: ItemDetailPage (/items/:id) 内の購入モーダル
    * 内容: MetaMask等のウォレットを接続し、Sepolia ETHで支払いを行います。
* エスクロー＆受け取り確認：
    * ページ: MyPage (/mypage)
    * 内容: 購入した商品のステータスを確認。商品が届いたら「受け取り確認」ボタンを押すことで、スマートコントラクトから出品者へ代金が送金されます。

5. ソーシャル・コミュニケーション
ユーザー間の交流を促進する機能です。
* いいね機能：
    * ページ: ItemDetailPage および MyPage
    * 内容: 商品に「いいね」を付け、マイページで後から一覧として見返すことができます。
* ダイレクトメッセージ (DM)：
    * ページ: DMPage (/dm/:id)
    * 内容: 出品者と購入者が取引について直接メッセージをやり取りできます。
* シェア機能：
    * ページ: ItemDetailPage
    * 内容: SNS（X/Twitterなど）やリンクコピーで商品情報を共有できます。

6. ウォレット・資産管理
* ウォレット連携・残高確認：
    * ページ: MyPage (/mypage)
    * 内容: 接続中のウォレットアドレスや、Sepolia ETHの残高を確認できます。

これがアプリの機能です。web3がわからない初心者にも簡潔かつ丁寧に教えなさい。`

type ChatUsecase struct {
	apiKey string
}

func NewGeminiUsecase() *ChatUsecase {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Println("WARNING: OPENAI_API_KEY environment variable is not set")
	}
	return &ChatUsecase{
		apiKey: apiKey,
	}
}

// OpenAI API用の構造体
type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIRequest struct {
	Model    string          `json:"model"`
	Messages []OpenAIMessage `json:"messages"`
}

type OpenAIResponse struct {
	Choices []OpenAIChoice `json:"choices"`
	Error   *OpenAIError   `json:"error,omitempty"`
}

type OpenAIChoice struct {
	Message OpenAIMessage `json:"message"`
}

type OpenAIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    string `json:"code"`
}

type GenerateContentRequest struct {
	Prompt   string `json:"prompt"`
	Protocol string `json:"protocol"`
}

type GenerateContentResponse struct {
	Response string `json:"response"`
	Error    string `json:"error,omitempty"`
}

func (uc *ChatUsecase) GenerateContent(userMessage string, protocol string) (*GenerateContentResponse, error) {
	if uc.apiKey == "" {
		return &GenerateContentResponse{
			Error: "OPENAI_API_KEY environment variable is not set",
		}, fmt.Errorf("OPENAI_API_KEY environment variable is not set")
	}

	fullSystemPrompt := systemPrompt
	if protocol != "" {
		fullSystemPrompt = fmt.Sprintf("%s\n\nProtocol:\n%s", systemPrompt, protocol)
	}

	openaiReq := OpenAIRequest{
		Model: "gpt-4o-mini",
		Messages: []OpenAIMessage{
			{
				Role:    "system",
				Content: fullSystemPrompt,
			},
			{
				Role:    "user",
				Content: userMessage,
			},
		},
	}

	jsonData, err := json.Marshal(openaiReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := "https://api.openai.com/v1/chat/completions"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+uc.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("OpenAI API request error: %v", err)
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("OpenAI API response read error: %v", err)
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("OpenAI API error: status=%d, body=%s", resp.StatusCode, string(body))
		return &GenerateContentResponse{
			Error: fmt.Sprintf("API error: status=%d, body=%s", resp.StatusCode, string(body)),
		}, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var openaiResp OpenAIResponse
	if err := json.Unmarshal(body, &openaiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if openaiResp.Error != nil {
		return &GenerateContentResponse{
			Error: fmt.Sprintf("OpenAI API error: %s", openaiResp.Error.Message),
		}, fmt.Errorf("OpenAI API error: %s", openaiResp.Error.Message)
	}

	if len(openaiResp.Choices) == 0 {
		return &GenerateContentResponse{
			Error: "No response from OpenAI API",
		}, fmt.Errorf("no choices in response")
	}

	return &GenerateContentResponse{
		Response: openaiResp.Choices[0].Message.Content,
	}, nil
}
