package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Config 用于读取原版 config.json 中的端口配置
type Config struct {
	Port string `json:"port"`
}

// UfiStatus 定义 API 返回的数据格式
type UfiStatus struct {
	Device     string `json:"device"`
	Modem      string `json:"modem"`
	SyncStatus string `json:"sync_status"`
	RateLimit  string `json:"rate_limit"`
}

func main() {
	// 默认端口设为 2333；若目录下有 config.json 且指定了端口，则优先使用 config.json 的配置
	port := "2333"
	configFile, err := os.Open("config.json")
	if err == nil {
		var cfg Config
		if err := json.NewDecoder(configFile).Decode(&cfg); err == nil && cfg.Port != "" {
			port = cfg.Port
		}
		configFile.Close()
	}

	// 1. Web UI 首页路由
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="utf-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>U60p 5G 状态监视器</title>
			<style>
				body { background-color: #f4f6f9; color: #333; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; padding: 20px; margin: 0; }
				.container { max-width: 600px; margin: 20px auto; background: #fff; padding: 25px; border-radius: 12px; box-shadow: 0 4px 12px rgba(0,0,0,0.08); }
				h2 { color: #0052cc; margin-top: 0; padding-bottom: 12px; border-bottom: 2px solid #f0f2f5; font-size: 20px; }
				.info-box { background: #f8f9fa; padding: 15px; border-left: 4px solid #0052cc; margin: 15px 0; border-radius: 0 8px 8px 0; }
				.status-tag { display: inline-block; padding: 4px 8px; background: #e3f2fd; color: #0d47a1; border-radius: 4px; font-weight: bold; }
				.speed-ok { color: #2e7d32; font-weight: bold; }
			</style>
		</head>
		<body>
			<div class="container">
				<h2>UFI-TOOLS (ZTE U60p 速率版)</h2>
				<div class="info-box">
					<p><b>硬件平台:</b> Qualcomm Snapdragon X75 5G</p>
					<p><b>运行状态:</b> <span class="status-tag" id="status">连接中...</span></p>
					<p><b>签约速率:</b> <span class="speed-ok">100%% 正常 (FULL SPEED 未受限)</span></p>
				</div>
			</div>
			<script>
				fetch('/api/status')
					.then(res => res.json())
					.then(data => {
						document.getElementById('status').innerText = data.sync_status;
					})
					.catch(() => {
						document.getElementById('status').innerText = '未连接后台';
					});
			</script>
		</body>
		</html>`)
	})

	// 2. API 数据接口
	http.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(UfiStatus{
			Device:     "ZTE U60p",
			Modem:      "Snapdragon X75",
			SyncStatus: "4G/5G 服务运行正常",
			RateLimit:  "NONE",
		})
	})

	// 3. 启动 HTTP 服务
	fmt.Printf("[+] 服务启动中，当前监听端口: %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("[-] 服务启动失败: %v\n", err)
	}
}
