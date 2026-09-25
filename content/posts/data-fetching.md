---
date: 2026-09-23
tags: [tools]
---

# goquery Web Scraping

`goquery`를 활용해 웹 페이지(Hacker News 등)의 헤드라인을 실시간 스크래핑하고 기술 키워드 언급 빈도를 수집한 결과다.

- **파이프라인**: `tools/goquery/main.go` 실행 → `data/sample_goquery.json` 생성 → Hugo `hugo.Data` 로드
- **차트 엔진**: `Chart.js` (Bar Chart)

---

## Chart.js

{{< chartjs data="sample_goquery" key="topic_bar" title="Hacker News 헤드라인 주요 기술 언급 빈도" height="380px" >}}

---

### `tools/goquery/main.go`

```go
doc.Find(".titleline > a").Each(func(i int, s *goquery.Selection) {
    text := strings.ToLower(s.Text())
    for _, t := range topics {
        if strings.Contains(text, strings.ToLower(t)) {
            counts[t]++
        }
    }
})
```

# FRED REST API (Economic Data)

미국 연방준비은행 경제 데이터(FRED API)를 Go 스크립트로 호출하여 최신 실업률(UNRATE) 시계열 데이터를 수집하고 시각화한 결과다.

- **파이프라인**: `tools/rest_api/main.go` 실행 → `data/sample_rest_api.csv` (원천) & `data/sample_rest_api.json` (가공) 생성
- **차트 엔진**: `ECharts` (Time-Series Area Chart)

---

## ECharts

{{< echarts data="sample_rest_api" key="unrate_series" type="time-series" title="미국 실업률 추이 (UNRATE, %)" height="420px" >}}

### `tools/rest_api/main.go`

```go
apiKey := getEnvKey("FRED_API_KEY")
url := fmt.Sprintf("https://api.stlouisfed.org/fred/series/observations?series_id=UNRATE&api_key=%s&file_type=json&sort_order=desc&limit=24", apiKey)
res, err := http.Get(url)
```