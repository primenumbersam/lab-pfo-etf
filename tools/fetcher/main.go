package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type HoldingInfo struct {
	Symbol string  `json:"symbol"`
	Name   string  `json:"name"`
	Weight float64 `json:"weight"`
}

type ReturnMetrics struct {
	Ret1M     float64 `json:"1m"`
	Ret3M     float64 `json:"3m"`
	RetReal3M float64 `json:"real_3m"`
	RetYTD    float64 `json:"ytd"`
	Ret1Y     float64 `json:"1y"`
}

type ETFInfo struct {
	Symbol        string         `json:"symbol"`
	ISIN          string         `json:"isin,omitempty"`
	Name          string         `json:"name"`
	Provider      string         `json:"provider"`
	Index         string         `json:"index"`
	Currency      string         `json:"currency"`
	InceptionDate string         `json:"inception_date"`
	ExpenseRatio  float64        `json:"expense_ratio"`
	AUM           float64        `json:"aum"`
	AssetClass    string         `json:"asset_class,omitempty"`
	AssetType     string         `json:"asset_type,omitempty"`
	TopHoldings   []HoldingInfo  `json:"top_holdings"`
	Returns       *ReturnMetrics `json:"returns,omitempty"`
}

var etfCatalog = map[string][2]string{
	"SPY":  {"US Equities", "Stocks-US-Large"},
	"VOO":  {"US Equities", "Stocks-US-Large"},
	"QQQ":  {"US Equities", "Stocks-US-Large"},
	"DIA":  {"US Equities", "Stocks-US-Large"},
	"VBR":  {"US Equities", "Stocks-US-Small-Value"},
	"VB":   {"US Equities", "Stocks-US-Small"},
	"IJR":  {"US Equities", "Stocks-US-Small"},
	"SCHD": {"US Equities", "Stocks-US-Large-Value"},
	"VIG":  {"US Equities", "Stocks-US-Large-Value"},
	"SPYV": {"US Equities", "Stocks-US-Large-Value"},
	"SPYG": {"US Equities", "Stocks-US-Large-Growth"},
	"VUG":  {"US Equities", "Stocks-US-Large-Growth"},
	"IWF":  {"US Equities", "Stocks-US-Large-Growth"},
	"VBK":  {"US Equities", "Stocks-US-Small-Growth"},
	"VEA":  {"Global Equities", "Stocks-Developed-Large"},
	"SCHF": {"Global Equities", "Stocks-Developed-Large"},
	"VT":   {"Global Equities", "Stocks-Developed-Large"},
	"URTH": {"Global Equities", "Stocks-Developed-Large"},
	"VWO":  {"Global Equities", "Stocks-Emerging-Large"},
	"IEMG": {"Global Equities", "Stocks-Emerging-Large"},
	"VSS":  {"Global Equities", "Stocks-Developed-Small"},
	"EFV":  {"Global Equities", "Stocks-Developed-Large-Value"},
	"VGIT": {"Fixed Income & Cash", "Bond-US-Intermediate"},
	"GOVT": {"Fixed Income & Cash", "Bond-US-Intermediate"},
	"SGOV": {"Fixed Income & Cash", "Bond-US-Bills"},
	"VBIL": {"Fixed Income & Cash", "Bond-US-Bills"},
	"BIL":  {"Fixed Income & Cash", "Bond-US-Bills"},
	"VGLT": {"Fixed Income & Cash", "Bond-US-Long"},
	"TLT":  {"Fixed Income & Cash", "Bond-US-Long"},
	"VGSH": {"Fixed Income & Cash", "Bond-US-Short"},
	"SHY":  {"Fixed Income & Cash", "Bond-US-Short"},
	"IGOV": {"Fixed Income & Cash", "Bond-Developed-Intermediate"},
	"VNQ":  {"Commodities & Real Estate", "REITs-US"},
	"SCHH": {"Commodities & Real Estate", "REITs-US"},
	"IAU":  {"Commodities & Real Estate", "Gold-Global"},
	"GLD":  {"Commodities & Real Estate", "Gold-Global"},
	"GSG":  {"Commodities & Real Estate", "Commodities-Global"},
}

var etfIndexCatalog = map[string]string{
	"SPY":  "S&P 500 Index",
	"VOO":  "S&P 500 Index",
	"QQQ":  "NASDAQ-100 Index",
	"DIA":  "Dow Jones Industrial Average",
	"VBR":  "CRSP US Small Cap Value Index",
	"VB":   "CRSP US Small Cap Index",
	"IJR":  "S&P SmallCap 600 Index",
	"SCHD": "Dow Jones U.S. Dividend 100 Index",
	"VIG":  "S&P U.S. Dividend Growers Index",
	"SPYV": "S&P 500 Value Index",
	"SPYG": "S&P 500 Growth Index",
	"VUG":  "CRSP US Large Cap Growth Index",
	"IWF":  "Russell 1000 Growth Index",
	"VBK":  "CRSP US Small Cap Growth Index",
	"VEA":  "FTSE Developed All Cap ex US Index",
	"SCHF": "FTSE Developed ex US Index",
	"VT":   "FTSE Global All Cap Index",
	"URTH": "MSCI World Index",
	"VWO":  "FTSE Emerging Markets All Cap China A Inclusion Index",
	"IEMG": "MSCI Emerging Markets Investable Market Index",
	"VSS":  "FTSE Global All Cap ex US Small Cap Index",
	"EFV":  "MSCI EAFE Value Index",
	"VGIT": "Bloomberg U.S. 3-10 Year Government Float Adjusted Index",
	"GOVT": "ICE U.S. Treasury Core Bond Index",
	"SGOV": "ICE 0-3 Month US Treasury Bill Index",
	"VBIL": "Bloomberg 1-3 Month U.S. Treasury Bill Index",
	"BIL":  "Bloomberg 1-3 Month U.S. Treasury Bill Index",
	"VGLT": "Bloomberg U.S. Long Government Float Adjusted Index",
	"TLT":  "ICE U.S. Treasury 20+ Year Bond Index",
	"VGSH": "Bloomberg U.S. 1-5 Year Government Float Adjusted Index",
	"SHY":  "ICE U.S. Treasury 1-3 Year Bond Index",
	"IGOV": "S&P International Sovereign Ex-U.S. 1-3 Year Bond Index",
	"VNQ":  "MSCI US Investable Market Real Estate 25/50 Index",
	"SCHH": "S&P United States REIT Index",
	"IAU":  "LBMA Gold Price PM",
	"GLD":  "LBMA Gold Price PM",
	"GSG":  "S&P GSCI Commodity-Indexed Trust",
}

func scrapeIndexWithGoquery(client *http.Client, sym string) string {
	url := fmt.Sprintf("https://finance.yahoo.com/quote/%s/profile/", sym)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ""
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		if resp != nil {
			resp.Body.Close()
		}
		return ""
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return ""
	}

	var matched string
	doc.Find("section p, div p").Each(func(i int, s *goquery.Selection) {
		if matched != "" {
			return
		}
		t := s.Text()
		if strings.Contains(strings.ToLower(t), "index") {
			for _, official := range etfIndexCatalog {
				if strings.Contains(strings.ToLower(t), strings.ToLower(official)) {
					matched = official
					break
				}
			}
		}
	})
	return matched
}

func resolveBenchmarkIndex(client *http.Client, sym string) string {
	if scraped := scrapeIndexWithGoquery(client, sym); scraped != "" {
		return scraped
	}
	if official, ok := etfIndexCatalog[sym]; ok {
		return official
	}
	return ""
}

func main() {
	targetFile := "data/etfs.json"
	fmt.Printf("[Fetcher] 37개 ETF 메타데이터 동기화: %s\n", targetFile)

	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar:     jar,
		Timeout: 15 * time.Second,
	}

	// 1. Get cookie
	req, _ := http.NewRequest("GET", "https://fc.yahoo.com", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	resp, _ := client.Do(req)
	if resp != nil {
		resp.Body.Close()
	}

	// 2. Get Crumb
	req, _ = http.NewRequest("GET", "https://query2.finance.yahoo.com/v1/test/getcrumb", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	cResp, err := client.Do(req)
	crumb := ""
	if err == nil && cResp.StatusCode == 200 {
		cBytes, _ := io.ReadAll(cResp.Body)
		crumb = strings.TrimSpace(string(cBytes))
		cResp.Body.Close()
		fmt.Printf("[Fetcher] Crumb 세션 획득 성공: %s\n", crumb)
	}

	// Load existing or init
	etfMap := make(map[string]ETFInfo)
	if data, err := os.ReadFile(targetFile); err == nil {
		_ = json.Unmarshal(data, &etfMap)
	}

	for sym, meta := range etfCatalog {
		item, exists := etfMap[sym]
		if !exists {
			item = ETFInfo{Symbol: sym}
		}
		item.AssetClass = meta[0]
		item.AssetType = meta[1]
		item.Index = resolveBenchmarkIndex(client, sym)

		if crumb != "" {
			url := fmt.Sprintf("https://query2.finance.yahoo.com/v10/finance/quoteSummary/%s?crumb=%s&modules=fundProfile,topHoldings,summaryDetail,price", sym, crumb)
			qReq, _ := http.NewRequest("GET", url, nil)
			qReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
			qResp, qErr := client.Do(qReq)
			if qErr == nil && qResp.StatusCode == 200 {
				bodyBytes, _ := io.ReadAll(qResp.Body)
				qResp.Body.Close()

				var payload struct {
					QuoteSummary struct {
						Result []struct {
							Price struct {
								LongName  string `json:"longName"`
								ShortName string `json:"shortName"`
								Currency  string `json:"currency"`
							} `json:"price"`
							FundProfile struct {
								Family                 string `json:"family"`
								FeesExpensesInvestment struct {
									AnnualReportExpenseRatio struct {
										Raw float64 `json:"raw"`
									} `json:"annualReportExpenseRatio"`
								} `json:"feesExpensesInvestment"`
							} `json:"fundProfile"`
							SummaryDetail struct {
								TotalAssets struct {
									Raw float64 `json:"raw"`
								} `json:"totalAssets"`
							} `json:"summaryDetail"`
							TopHoldings struct {
								Holdings []struct {
									Symbol         string `json:"symbol"`
									HoldingName    string `json:"holdingName"`
									HoldingPercent struct {
										Raw float64 `json:"raw"`
									} `json:"holdingPercent"`
								} `json:"holdings"`
							} `json:"topHoldings"`
						} `json:"result"`
					} `json:"quoteSummary"`
				}

				if err := json.Unmarshal(bodyBytes, &payload); err == nil && len(payload.QuoteSummary.Result) > 0 {
					r := payload.QuoteSummary.Result[0]
					if r.Price.LongName != "" {
						item.Name = r.Price.LongName
					} else if r.Price.ShortName != "" {
						item.Name = r.Price.ShortName
					}
					item.Currency = r.Price.Currency
					item.Provider = r.FundProfile.Family
					item.ExpenseRatio = r.FundProfile.FeesExpensesInvestment.AnnualReportExpenseRatio.Raw
					item.AUM = r.SummaryDetail.TotalAssets.Raw

					var hList []HoldingInfo
					for _, h := range r.TopHoldings.Holdings {
						hList = append(hList, HoldingInfo{
							Symbol: h.Symbol,
							Name:   h.HoldingName,
							Weight: h.HoldingPercent.Raw,
						})
					}
					item.TopHoldings = hList
				}
			}
		}

		// 3. Fetch Returns & Inception Date
		if ret, incep := fetchReturns(client, sym); ret != nil {
			item.Returns = ret
			if incep != "" {
				item.InceptionDate = incep
			}
			fmt.Printf(" [%s] Index: %s | Inception: %s | 3M Real: %.2f%%\n", sym, item.Index, item.InceptionDate, ret.RetReal3M)
		}

		etfMap[sym] = item
		time.Sleep(50 * time.Millisecond)
	}

	outBytes, _ := json.MarshalIndent(etfMap, "", "  ")
	_ = os.WriteFile(targetFile, outBytes, 0644)
	fmt.Printf("[Fetcher] %s 갱신 완료 (총 %d개 ETF 등록)\n", targetFile, len(etfMap))
}

func fetchReturns(client *http.Client, sym string) (*ReturnMetrics, string) {
	url := fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s?range=2y&interval=1d", sym)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, ""
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		if resp != nil {
			resp.Body.Close()
		}
		return nil, ""
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ""
	}

	var payload struct {
		Chart struct {
			Result []struct {
				Meta struct {
					FirstTradeDate int64 `json:"firstTradeDate"`
				} `json:"meta"`
				Timestamp  []int64 `json:"timestamp"`
				Indicators struct {
					Quote []struct {
						Close []float64 `json:"close"`
					} `json:"quote"`
					AdjClose []struct {
						AdjClose []float64 `json:"adjclose"`
					} `json:"adjclose"`
				} `json:"indicators"`
			} `json:"result"`
		} `json:"chart"`
	}

	if err := json.Unmarshal(bodyBytes, &payload); err != nil || len(payload.Chart.Result) == 0 {
		return nil, ""
	}

	res := payload.Chart.Result[0]
	inceptionDate := ""
	if res.Meta.FirstTradeDate > 0 {
		inceptionDate = time.Unix(res.Meta.FirstTradeDate, 0).UTC().Format("2006-01-02")
	}

	timestamps := res.Timestamp
	var prices []float64
	if len(res.Indicators.AdjClose) > 0 && len(res.Indicators.AdjClose[0].AdjClose) > 0 {
		prices = res.Indicators.AdjClose[0].AdjClose
	} else if len(res.Indicators.Quote) > 0 && len(res.Indicators.Quote[0].Close) > 0 {
		prices = res.Indicators.Quote[0].Close
	}

	type point struct {
		t int64
		p float64
	}
	var series []point
	for i := 0; i < len(timestamps) && i < len(prices); i++ {
		p := prices[i]
		if p > 0 {
			series = append(series, point{t: timestamps[i], p: p})
		}
	}
	if len(series) < 2 {
		return nil, inceptionDate
	}

	latest := series[len(series)-1]
	latestTime := time.Unix(latest.t, 0).UTC()

	t1m := latest.t - 30*86400
	t3m := latest.t - 91*86400
	tYtd := time.Date(latestTime.Year(), 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	t1y := latest.t - 365*86400

	findClosest := func(targetT int64) float64 {
		closest := series[0].p
		minDiff := targetT - series[0].t
		if minDiff < 0 {
			minDiff = -minDiff
		}
		for i := 1; i < len(series); i++ {
			diff := targetT - series[i].t
			if diff < 0 {
				diff = -diff
			}
			if diff < minDiff {
				minDiff = diff
				closest = series[i].p
			}
		}
		return closest
	}

	p1m := findClosest(t1m)
	p3m := findClosest(t3m)
	pYtd := findClosest(tYtd)
	p1y := findClosest(t1y)

	calcRet := func(base float64) float64 {
		if base <= 0 {
			return 0
		}
		ret := ((latest.p - base) / base) * 100
		return math.Round(ret*100) / 100
	}

	ret1m := calcRet(p1m)
	ret3m := calcRet(p3m)
	retYtd := calcRet(pYtd)
	ret1y := calcRet(p1y)

	// 최근 3개월 CPI 분기 추정치 약 0.6%
	cpi3m := 0.6
	retReal3m := math.Round((((1+ret3m/100)/(1+cpi3m/100))-1)*10000) / 100

	return &ReturnMetrics{
		Ret1M:     ret1m,
		Ret3M:     ret3m,
		RetReal3M: retReal3m,
		RetYTD:    retYtd,
		Ret1Y:     ret1y,
	}, inceptionDate
}
