package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type FredResponse struct {
	Observations []struct {
		Date  string `json:"date"`
		Value string `json:"value"`
	} `json:"observations"`
}

// .env 파일에서 키 파싱 (외부 패키지 의존 없이 표준 라이브러리로 처리)
func getEnvKey(key string) string {
	val := os.Getenv(key)
	if val != "" {
		return val
	}

	file, err := os.Open(".env")
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		k := strings.TrimSpace(parts[0])
		v := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
		if k == key {
			return v
		}
	}
	return ""
}

func main() {
	apiKey := getEnvKey("FRED_API_KEY")
	if apiKey == "" {
		log.Fatal("FRED_API_KEY를 .env 파일 또는 환경 변수에서 찾을 수 없습니다.")
	}

	// 미국 실업률(UNRATE) 최신 24개월 데이터 조회
	url := fmt.Sprintf("https://api.stlouisfed.org/fred/series/observations?series_id=UNRATE&api_key=%s&file_type=json&sort_order=desc&limit=24", apiKey)
	res, err := http.Get(url)
	if err != nil {
		log.Fatalf("FRED API 요청 실패: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("FRED API 오류 응답: %d %s", res.StatusCode, res.Status)
	}

	var fredData FredResponse
	if err := json.NewDecoder(res.Body).Decode(&fredData); err != nil {
		log.Fatalf("JSON 파싱 실패: %v", err)
	}

	// 1. sample_rest_api.csv 원천 데이터 저장
	csvFile, err := os.Create(filepath.Join("data", "sample_rest_api.csv"))
	if err != nil {
		log.Fatalf("CSV 생성 실패: %v", err)
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	_ = writer.Write([]string{"date", "unemployment_rate"})

	// 최신순으로 받은 것을 과거->최신(시간순)으로 정렬하여 차트용 JSON 및 CSV 준비
	obs := fredData.Observations
	for i, j := 0, len(obs)-1; i < j; i, j = i+1, j-1 {
		obs[i], obs[j] = obs[j], obs[i]
	}

	var dates []string
	var rates []float64

	for _, o := range obs {
		_ = writer.Write([]string{o.Date, o.Value})
		val, err := strconv.ParseFloat(o.Value, 64)
		if err == nil {
			dates = append(dates, o.Date)
			rates = append(rates, val)
		}
	}

	// 2. Hugo Chart Shortcode용 sample_rest_api.json 저장
	// (ECharts / Chart.js 둘 다 사용 가능한 Time-Series 구조)
	chartData := map[string]interface{}{
		"unrate_series": map[string]interface{}{
			"dates": dates,
			"rates": rates,
			"labels": dates,
			"datasets": []map[string]interface{}{
				{
					"label":           "미국 실업률 (%)",
					"data":            rates,
					"borderColor":     "rgba(239, 68, 68, 1)",
					"backgroundColor": "rgba(239, 68, 68, 0.15)",
					"fill":            true,
					"tension":         0.3,
				},
			},
		},
	}

	jsonBytes, err := json.MarshalIndent(chartData, "", "  ")
	if err != nil {
		log.Fatalf("JSON 직렬화 실패: %v", err)
	}

	jsonPath := filepath.Join("data", "sample_rest_api.json")
	if err := os.WriteFile(jsonPath, jsonBytes, 0644); err != nil {
		log.Fatalf("JSON 저장 실패: %v", err)
	}

	fmt.Printf("성공: data/sample_rest_api.csv 및 %s 생성 완료 (%d개 관측치)\n", jsonPath, len(rates))
}
