package controllers

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"superhoneypotguard/utils"

	"github.com/gin-gonic/gin"
)

type HFishController struct{}

func NewHFishController() *HFishController {
	return &HFishController{}
}

const (
	HFishBaseURL = "https://115.190.62.202:4433/api/v1"
	APIKey       = "sOPdmLBemXeWqPmizIvkMqKrRIgdnqkqbzOMciukucEiBFVAhcotwVLoLnsGgyNa"
)

type AttackIPResponse struct {
	Success bool       `json:"success"`
	Message string     `json:"message"`
	Data    []AttackIP `json:"data"`
}

type AttackIP struct {
	ID        string `json:"id"`
	IP        string `json:"ip"`
	Count     int    `json:"count"`
	FirstSeen string `json:"first_seen"`
	LastSeen  string `json:"last_seen"`
}

type AttackDetailResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Data    []AttackDetail `json:"data"`
}

type AttackDetail struct {
	ID          string `json:"id"`
	IP          string `json:"ip"`
	AttackType  string `json:"attack_type"`
	Protocol    string `json:"protocol"`
	Port        int    `json:"port"`
	Payload     string `json:"payload"`
	RequestTime string `json:"request_time"`
	Account     string `json:"account"`
}

type AccountInfoResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Data    []AccountInfo `json:"data"`
}

type AccountInfo struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Protocol    string `json:"protocol"`
	IP          string `json:"ip"`
	AttackCount int    `json:"attack_count"`
}

type SysInfoResponse struct {
	Success bool    `json:"success"`
	Message string  `json:"message"`
	Data    SysInfo `json:"data"`
}

type SysInfo struct {
	TotalHoneypots  int    `json:"total_honeypots"`
	ActiveHoneypots int    `json:"active_honeypots"`
	TotalAttacks    int    `json:"total_attacks"`
	LastAttackTime  string `json:"last_attack_time"`
	SystemStatus    string `json:"system_status"`
}

var httpClient = &http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	},
}

func (ctrl *HFishController) GetAttackIPs(c *gin.Context) {
	url := fmt.Sprintf("%s/attack/ip?api_key=%s", HFishBaseURL, APIKey)
	log.Printf("调用 HFish API: %s", url)

	resp, err := httpClient.Post(url, "application/json", nil)
	if err != nil {
		log.Printf("调用 HFish API 失败: %v", err)
		utils.ErrorResponse(c, 500, "调用 HFish API 失败: "+err.Error())
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Printf("HFish API 响应状态码: %d, 响应内容: %s", resp.StatusCode, string(body))

	var result AttackIPResponse
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("解析 HFish API 响应失败: %v", err)
		utils.ErrorResponse(c, 500, "解析 HFish API 响应失败")
		return
	}

	if !result.Success {
		log.Printf("HFish API 返回错误: %s", result.Message)
		utils.ErrorResponse(c, 500, result.Message)
		return
	}

	utils.SuccessResponse(c, result.Data)
}

func (ctrl *HFishController) GetAttackDetails(c *gin.Context) {
	url := fmt.Sprintf("%s/attack/detail?api_key=%s", HFishBaseURL, APIKey)
	log.Printf("调用 HFish API: %s", url)

	resp, err := httpClient.Post(url, "application/json", nil)
	if err != nil {
		log.Printf("调用 HFish API 失败: %v", err)
		utils.ErrorResponse(c, 500, "调用 HFish API 失败: "+err.Error())
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Printf("HFish API 响应状态码: %d, 响应内容: %s", resp.StatusCode, string(body))

	var result AttackDetailResponse
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("解析 HFish API 响应失败: %v", err)
		utils.ErrorResponse(c, 500, "解析 HFish API 响应失败")
		return
	}

	if !result.Success {
		log.Printf("HFish API 返回错误: %s", result.Message)
		utils.ErrorResponse(c, 500, result.Message)
		return
	}

	utils.SuccessResponse(c, result.Data)
}

func (ctrl *HFishController) GetAccountInfo(c *gin.Context) {
	url := fmt.Sprintf("%s/attack/account?api_key=%s", HFishBaseURL, APIKey)
	log.Printf("调用 HFish API: %s", url)

	resp, err := httpClient.Post(url, "application/json", nil)
	if err != nil {
		log.Printf("调用 HFish API 失败: %v", err)
		utils.ErrorResponse(c, 500, "调用 HFish API 失败: "+err.Error())
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Printf("HFish API 响应状态码: %d, 响应内容: %s", resp.StatusCode, string(body))

	var result AccountInfoResponse
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("解析 HFish API 响应失败: %v", err)
		utils.ErrorResponse(c, 500, "解析 HFish API 响应失败")
		return
	}

	if !result.Success {
		log.Printf("HFish API 返回错误: %s", result.Message)
		utils.ErrorResponse(c, 500, result.Message)
		return
	}

	utils.SuccessResponse(c, result.Data)
}

func (ctrl *HFishController) GetSysInfo(c *gin.Context) {
	url := fmt.Sprintf("%s/hfish/sys_info?api_key=%s", HFishBaseURL, APIKey)
	log.Printf("调用 HFish API: %s", url)

	resp, err := httpClient.Get(url)
	if err != nil {
		log.Printf("调用 HFish API 失败: %v", err)
		utils.ErrorResponse(c, 500, "调用 HFish API 失败: "+err.Error())
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Printf("HFish API 响应状态码: %d, 响应内容: %s", resp.StatusCode, string(body))

	var result SysInfoResponse
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("解析 HFish API 响应失败: %v", err)
		utils.ErrorResponse(c, 500, "解析 HFish API 响应失败")
		return
	}

	if !result.Success {
		log.Printf("HFish API 返回错误: %s", result.Message)
		utils.ErrorResponse(c, 500, result.Message)
		return
	}

	utils.SuccessResponse(c, result.Data)
}

func (ctrl *HFishController) BlockIP(c *gin.Context) {
	var req struct {
		IP     string `json:"ip" binding:"required"`
		Reason string `json:"reason"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("参数验证失败: %v", err)
		utils.ErrorResponse(c, 400, "参数验证失败")
		return
	}

	url := fmt.Sprintf("%s/attack/ip/block?api_key=%s", HFishBaseURL, APIKey)
	log.Printf("调用 HFish API 封禁 IP: %s, IP: %s, Reason: %s", url, req.IP, req.Reason)

	jsonData, _ := json.Marshal(req)
	resp, err := httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("调用 HFish API 失败: %v", err)
		utils.ErrorResponse(c, 500, "调用 HFish API 失败: "+err.Error())
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Printf("HFish API 响应状态码: %d, 响应内容: %s", resp.StatusCode, string(body))

	var result map[string]interface{}
	json.Unmarshal(body, &result)

	if result["success"] == true {
		log.Printf("封禁 IP 成功: %s", req.IP)
		utils.SuccessResponse(c, nil)
	} else {
		log.Printf("封禁 IP 失败: %s", result["message"])
		utils.ErrorResponse(c, 500, fmt.Sprintf("封禁 IP 失败: %v", result["message"]))
	}
}
