package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
)

type CompositionItem struct {
	Element string  `json:"element"`
	Weight  float64 `json:"weight"`
}

type PortfolioData struct {
	Slug         string            `json:"slug"`
	Name         string            `json:"name"`
	Category     string            `json:"category"`
	CategorySlug string            `json:"category_slug"`
	Composition  []CompositionItem `json:"composition"`
	AssetNotes   string            `json:"asset_notes"`
	Author       string            `json:"author"`
	Overview     string            `json:"overview"`
}

type AssetReturnsPayload struct {
	Years     []int                `json:"years"`
	Inflation []float64            `json:"inflation"`
	Assets    map[string][]float64 `json:"assets"`
}

type PortfolioScoreItem struct {
	Slug                 string  `json:"slug"`
	Name                 string  `json:"name"`
	Category             string  `json:"category"`
	CategorySlug         string  `json:"category_slug"`
	GeometricReturn      float64 `json:"geometric_return"`
	OmegaRatio           float64 `json:"omega_ratio"`
	RmseDD               float64 `json:"rmse_dd"`
	StartDateSensitivity float64 `json:"start_date_sensitivity"`
	CompositeScore       float64 `json:"composite_score"`
	Rank                 int     `json:"rank"`
}

type ScoresPayload struct {
	Items       []PortfolioScoreItem `json:"items"`
	LastUpdated string               `json:"last_updated"`
}

type SeriesDataPayload struct {
	Slug          string               `json:"slug"`
	Name          string               `json:"name"`
	Years         []int                `json:"years"`
	AnnualReturns []float64            `json:"annual_returns"`
	Rolling10Y    map[string]float64   `json:"rolling_10y"`
	FunnelByN     []map[string]float64 `json:"funnel_by_n"`
	Heatmap       [][]float64          `json:"heatmap"` // [startYearIdx][holdingPeriodN]
}

func percentile(sortedVals []float64, p float64) float64 {
	if len(sortedVals) == 0 {
		return 0
	}
	idx := p * float64(len(sortedVals)-1)
	lower := int(math.Floor(idx))
	upper := int(math.Ceil(idx))
	if lower == upper {
		return sortedVals[lower]
	}
	weight := idx - float64(lower)
	return sortedVals[lower]*(1-weight) + sortedVals[upper]*weight
}

func main() {
	pfoFile := "data/portfolios.json"
	returnsFile := "data/series/asset_returns.json"
	scoresFile := "data/scores.json"

	pfoBytes, err := os.ReadFile(pfoFile)
	if err != nil {
		log.Fatalf("[Error] %s 읽기 실패: %v\n", pfoFile, err)
	}
	var pfoMap map[string]PortfolioData
	if err := json.Unmarshal(pfoBytes, &pfoMap); err != nil {
		log.Fatalf("[Error] %s 파싱 실패: %v\n", pfoFile, err)
	}

	retBytes, err := os.ReadFile(returnsFile)
	if err != nil {
		log.Fatalf("[Error] %s 읽기 실패: %v\n", returnsFile, err)
	}
	var retPayload AssetReturnsPayload
	if err := json.Unmarshal(retBytes, &retPayload); err != nil {
		log.Fatalf("[Error] %s 파싱 실패: %v\n", returnsFile, err)
	}

	numYears := len(retPayload.Years)
	fmt.Printf("[Calculator] 18개 자산군 시계열 로드 완료 (%d년: %d ~ %d)\n", numYears, retPayload.Years[0], retPayload.Years[numYears-1])

	// T-Bill threshold for Omega Ratio
	tbills := retPayload.Assets["Bond-US-Bills"]
	sumTBill := 0.0
	for _, b := range tbills {
		sumTBill += b
	}
	avgTBill := sumTBill / float64(len(tbills))

	var scoreItems []PortfolioScoreItem

	for slug, pfo := range pfoMap {
		// 1. Compute portfolio annual real returns (Annual Rebalancing)
		pfoReturns := make([]float64, numYears)
		for t := 0; t < numYears; t++ {
			retT := 0.0
			for _, comp := range pfo.Composition {
				assetSeries, exists := retPayload.Assets[comp.Element]
				if exists && t < len(assetSeries) {
					retT += comp.Weight * assetSeries[t]
				}
			}
			pfoReturns[t] = retT
		}

		// 2. Metric 1: 10-year rolling Geometric Return (15th percentile)
		window := 10
		var cagr10List []float64
		for s := 0; s <= numYears-window; s++ {
			cum := 1.0
			for t := s; t < s+window; t++ {
				cum *= (1.0 + pfoReturns[t])
			}
			cagr := math.Pow(cum, 1.0/float64(window)) - 1.0
			cagr10List = append(cagr10List, cagr)
		}
		sortedCagr10 := make([]float64, len(cagr10List))
		copy(sortedCagr10, cagr10List)
		sort.Float64s(sortedCagr10)
		geom15th := percentile(sortedCagr10, 0.15) * 100.0

		// 3. Metric 2: Omega Ratio (threshold = avgTBill)
		posSum := 0.0
		negSum := 0.0
		for _, r := range pfoReturns {
			diff := r - avgTBill
			if diff > 0 {
				posSum += diff
			} else {
				negSum += math.Abs(diff)
			}
		}
		omega := 1.0
		if negSum > 0 {
			omega = posSum / negSum
		}

		// 4. Metric 3: RMSE of Drawdown (Ulcer Index)
		wealth := 1.0
		peak := 1.0
		sumDD2 := 0.0
		for _, r := range pfoReturns {
			wealth *= (1.0 + r)
			if wealth > peak {
				peak = wealth
			}
			dd := (peak - wealth) / peak * 100.0
			sumDD2 += dd * dd
		}
		rmseDD := math.Sqrt(sumDD2 / float64(numYears))

		// 5. Metric 4: Start Date Sensitivity (Std dev of 10Y CAGRs across start dates)
		meanCagr10 := 0.0
		for _, c := range cagr10List {
			meanCagr10 += c
		}
		meanCagr10 /= float64(len(cagr10List))
		varSum := 0.0
		for _, c := range cagr10List {
			diff := c - meanCagr10
			varSum += diff * diff
		}
		sensitivity := math.Sqrt(varSum/float64(len(cagr10List))) * 100.0

		// 6. Generate detailed series for ECharts
		rolling10Y := make(map[string]float64)
		for s := 0; s <= numYears-10; s++ {
			cum := 1.0
			for t := s; t < s+10; t++ {
				cum *= (1.0 + pfoReturns[t])
			}
			endYear := retPayload.Years[s+9]
			rolling10Y[fmt.Sprintf("%d", endYear)] = math.Round((math.Pow(cum, 0.1)-1.0)*1000) / 10.0
		}

		// Funnel distribution along N (N = 1 to 30)
		var funnelByN []map[string]float64
		for n := 1; n <= 30 && n <= numYears; n++ {
			var nCagrs []float64
			for s := 0; s <= numYears-n; s++ {
				cum := 1.0
				for t := s; t < s+n; t++ {
					cum *= (1.0 + pfoReturns[t])
				}
				nCagrs = append(nCagrs, (math.Pow(cum, 1.0/float64(n))-1.0)*100.0)
			}
			sort.Float64s(nCagrs)
			funnelByN = append(funnelByN, map[string]float64{
				"n":      float64(n),
				"min":    math.Round(nCagrs[0]*10) / 10,
				"p15":    math.Round(percentile(nCagrs, 0.15)*10) / 10,
				"median": math.Round(percentile(nCagrs, 0.50)*10) / 10,
				"p85":    math.Round(percentile(nCagrs, 0.85)*10) / 10,
				"max":    math.Round(nCagrs[len(nCagrs)-1]*10) / 10,
			})
		}

		// Heatmap matrix: start year s x holding period N (N = 1 to numYears-s)
		var heatmap [][]float64
		for s := 0; s < numYears; s++ {
			var row []float64
			for n := 1; n <= numYears; n++ {
				if s+n <= numYears {
					cum := 1.0
					for t := s; t < s+n; t++ {
						cum *= (1.0 + pfoReturns[t])
					}
					cagrVal := (math.Pow(cum, 1.0/float64(n)) - 1.0) * 100.0
					row = append(row, math.Round(cagrVal*10)/10)
				} else {
					row = append(row, 0.0)
				}
			}
			heatmap = append(heatmap, row)
		}

		annualRetPct := make([]float64, numYears)
		for i, r := range pfoReturns {
			annualRetPct[i] = math.Round(r*1000) / 10.0
		}

		// Write portfolio series JSON
		seriesPayload := SeriesDataPayload{
			Slug:          slug,
			Name:          pfo.Name,
			Years:         retPayload.Years,
			AnnualReturns: annualRetPct,
			Rolling10Y:    rolling10Y,
			FunnelByN:     funnelByN,
			Heatmap:       heatmap,
		}
		sBytes, _ := json.MarshalIndent(seriesPayload, "", "  ")
		_ = os.WriteFile(fmt.Sprintf("data/series/%s.json", slug), sBytes, 0644)

		scoreItems = append(scoreItems, PortfolioScoreItem{
			Slug:                 slug,
			Name:                 pfo.Name,
			Category:             pfo.Category,
			CategorySlug:         pfo.CategorySlug,
			GeometricReturn:      math.Round(geom15th*10) / 10,
			OmegaRatio:           math.Round(omega*100) / 100,
			RmseDD:               math.Round(rmseDD*10) / 10,
			StartDateSensitivity: math.Round(sensitivity*10) / 10,
		})
	}

	// Normalization & Composite Score calculation
	minG, maxG := math.MaxFloat64, -math.MaxFloat64
	minO, maxO := math.MaxFloat64, -math.MaxFloat64
	minD, maxD := math.MaxFloat64, -math.MaxFloat64
	minS, maxS := math.MaxFloat64, -math.MaxFloat64

	for _, item := range scoreItems {
		if item.GeometricReturn < minG {
			minG = item.GeometricReturn
		}
		if item.GeometricReturn > maxG {
			maxG = item.GeometricReturn
		}
		if item.OmegaRatio < minO {
			minO = item.OmegaRatio
		}
		if item.OmegaRatio > maxO {
			maxO = item.OmegaRatio
		}
		if item.RmseDD < minD {
			minD = item.RmseDD
		}
		if item.RmseDD > maxD {
			maxD = item.RmseDD
		}
		if item.StartDateSensitivity < minS {
			minS = item.StartDateSensitivity
		}
		if item.StartDateSensitivity > maxS {
			maxS = item.StartDateSensitivity
		}
	}

	norm := func(val, min, max float64, invert bool) float64 {
		if max == min {
			return 50.0
		}
		res := (val - min) / (max - min) * 100.0
		if invert {
			return 100.0 - res
		}
		return res
	}

	for i := range scoreItems {
		item := &scoreItems[i]
		s1 := norm(item.GeometricReturn, minG, maxG, false)
		s2 := norm(item.OmegaRatio, minO, maxO, false)
		s3 := norm(item.RmseDD, minD, maxD, true)
		s4 := norm(item.StartDateSensitivity, minS, maxS, true)

		// Equal 25% weights baseline
		comp := (s1 + s2 + s3 + s4) / 4.0
		item.CompositeScore = math.Round(comp*10) / 10
	}

	sort.Slice(scoreItems, func(i, j int) bool {
		return scoreItems[i].CompositeScore > scoreItems[j].CompositeScore
	})

	for i := range scoreItems {
		scoreItems[i].Rank = i + 1
	}

	outScores := ScoresPayload{
		Items:       scoreItems,
		LastUpdated: "2026-09 (Annual Rebalanced 1970-2025)",
	}

	outBytes, _ := json.MarshalIndent(outScores, "", "  ")
	_ = os.WriteFile(scoresFile, outBytes, 0644)
	fmt.Printf("[Calculator] %s 갱신 완료 (21개 모델 연산 및 data/series/*.json 개별 차트 시계열 도출)\n", scoresFile)
}
