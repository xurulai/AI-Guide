package model

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// DifyWorkflowRequest Dify工作流请求
type DifyWorkflowRequest struct {
	Inputs       map[string]interface{} `json:"inputs"`
	ResponseMode string                 `json:"response_mode"` // blocking/streaming
	User         string                 `json:"user"`
}

// DifyWorkflowResponse Dify工作流响应
type DifyWorkflowResponse struct {
	WorkflowRunID string              `json:"workflow_run_id"`
	TaskID        string              `json:"task_id"`
	Data          DifyWorkflowRunData `json:"data"`
}

// DifyWorkflowRunData Dify工作流运行数据
type DifyWorkflowRunData struct {
	ID          string                 `json:"id"`
	WorkflowID  string                 `json:"workflow_id"`
	Status      string                 `json:"status"`
	Outputs     map[string]interface{} `json:"outputs"`
	Error       string                 `json:"error"`
	ElapsedTime float64                `json:"elapsed_time"`
	TotalTokens int                    `json:"total_tokens"`
	TotalSteps  int                    `json:"total_steps"`
	CreatedAt   int64                  `json:"created_at"`
	FinishedAt  int64                  `json:"finished_at"`
}

// PlannerResult Planner解析结果
type PlannerResult struct {
	//RealShoppingIntentionClear bool                   `json:"real_shopping_intention_clear"`
	//RealShoppingIntentionItem  string                 `json:"real_shopping_intention_item"`
	Tool string `json:"tool"`
	//ToolInput                  string                 `json:"tool_input"`
	//Metadata                   map[string]interface{} `json:"metadata"`
}

// ExecutorResult Executor解析结果
type ExecutorResult struct {
	Response            string                 `json:"response"`
	RecommendedProducts []RecommendedProduct   `json:"recommended_products,omitempty"`
	Metadata            map[string]interface{} `json:"metadata"`
}

// Tool类型常量
const (
	ToolProductRecommendation = "PRODUCT_RECOMMENDATION_MODULE"
	ToolShoppingGuide         = "SHOPPING_GUIDE_AND_INTENT_MINING_MODULE"
	ToolQAAssistant           = "ECOMMERCE_QA_ASSISTANT_MODULE"
)

// StreamChunk 流式响应块
type StreamChunk struct {
	Event string      `json:"event"` // planner/chunk/products/done/error
	Data  interface{} `json:"data"`
}

var plannerBaseUrl = "http://jp02-inf-k8s00.jp02.baidu.com:8266/v1"
var plannerToken = "app-pYmNiYe8zZdoaaYHW7aaUHGG"

var guideBaseUrl = "http://jp02-inf-k8s00.jp02.baidu.com:8266/v1"
var guideToken = "app-Lgn1rQ2RvybwlFto7158FqB1"

func PlannerChat(ctx context.Context, request DifyWorkflowRequest) (string, error) {

	client := &http.Client{
		Timeout: 600 * time.Second,
	}

	// 构造请求数据
	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// 打印请求信息（调试用）
	fmt.Printf("\n=== Dify API Request ===\n")
	fmt.Printf("URL: %s\n", plannerBaseUrl+"/workflows/run")
	fmt.Printf("Request Body: %s\n", string(jsonData))
	fmt.Printf("========================\n\n")

	// 创建 HTTP 请求
	req, err := http.NewRequestWithContext(ctx, "POST", plannerBaseUrl+"/workflows/run", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+plannerToken)

	// 发送请求
	response, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer response.Body.Close()

	// 读取响应
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// 打印响应信息（调试用）
	fmt.Printf("\n=== Dify API Response ===\n")
	fmt.Printf("Status Code: %d\n", response.StatusCode)
	fmt.Printf("Response Body: %s\n", string(body))
	fmt.Printf("=========================\n\n")

	// 检查HTTP状态码
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("dify api error: status=%d, body=%s", response.StatusCode, string(body))
	}

	// 解析响应
	var result DifyWorkflowResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w, body=%s", err, string(body))
	}

	// 检查工作流执行状态
	if result.Data.Status != "succeeded" {
		return "", fmt.Errorf("workflow failed: status=%s, error=%s", result.Data.Status, result.Data.Error)
	}

	// 尝试提取 outputs 中的结果字段（可能是 text 或 result）
	var text string
	var ok bool

	// 先尝试 text 字段（Planner 通常返回 text）
	if text, ok = result.Data.Outputs["text"].(string); !ok {
		// 如果没有 text，尝试 result 字段
		if text, ok = result.Data.Outputs["result"].(string); !ok {
			return "", fmt.Errorf("neither 'text' nor 'result' field found in outputs, outputs=%+v", result.Data.Outputs)
		}
	}

	fmt.Printf("✅ Dify workflow succeeded: tokens=%d, elapsed=%.2fs\n", result.Data.TotalTokens, result.Data.ElapsedTime)

	return text, nil
}

func GuideChat(ctx context.Context, request DifyWorkflowRequest) (string, error) {

	client := &http.Client{
		Timeout: 600 * time.Second,
	}

	//request.Inputs["user_portrait"] = "{\"age\": 25, \"gender\": \"male\", \"interests\": [\"fashion\", \"technology\", \"travel\"]}"

	// 构造请求数据
	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// 打印请求信息（调试用）
	fmt.Printf("\n=== Guide Chat API Request ===\n")
	fmt.Printf("URL: %s\n", guideBaseUrl+"/workflows/run")
	fmt.Printf("Request Body: %s\n", string(jsonData))
	fmt.Printf("==============================\n\n")

	req, err := http.NewRequestWithContext(ctx, "POST", guideBaseUrl+"/workflows/run", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+guideToken)

	response, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// 打印响应信息（调试用）
	fmt.Printf("\n=== Guide Chat API Response ===\n")
	fmt.Printf("Status Code: %d\n", response.StatusCode)
	fmt.Printf("Response Body: %s\n", string(body))
	fmt.Printf("================================\n\n")

	// 检查HTTP状态码
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("guide chat api error: status=%d, body=%s", response.StatusCode, string(body))
	}

	result := DifyWorkflowResponse{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w, body=%s", err, string(body))
	}

	// 检查工作流执行状态
	if result.Data.Status != "succeeded" {
		return "", fmt.Errorf("workflow failed: status=%s, error=%s", result.Data.Status, result.Data.Error)
	}

	// 尝试提取 outputs 中的结果字段（可能是 text 或 result）
	var text string
	var ok bool

	// 先尝试 result 字段
	if text, ok = result.Data.Outputs["result"].(string); !ok {
		// 如果没有 result，尝试 text 字段
		if text, ok = result.Data.Outputs["text"].(string); !ok {
			return "", fmt.Errorf("neither 'result' nor 'text' field found in outputs, outputs=%+v", result.Data.Outputs)
		}
	}

	fmt.Printf("✅ Guide Chat workflow succeeded: tokens=%d, elapsed=%.2fs\n", result.Data.TotalTokens, result.Data.ElapsedTime)
	fmt.Printf("📄 Response content: %s\n\n", text)

	return text, nil

}
