---
date: 2026-09-24
tags: [logic]
BookPostThumbnail: "images/character-samsoon.jpg"
---

`portfolios.json` 스키마에 이 5가지 분류 체계를 `category` 필드로 매핑

## 5 Categories of Portfolio Allocation

### Value-Weighted (시장 자본 가중형)

공분산 제어나 자산군 강제 균등 배분 없이, 시장 시가총액과 전통적 2자산(주식/채권) 자본 비중($\mathbf{w}$)에 집중하는 정적 모델입니다.

* **포트폴리오**: [Classic 60-40 Portfolio](https://portfoliocharts.com/portfolios/classic-60-40-portfolio/?utm_source=gemini), [Three-Fund Portfolio](https://portfoliocharts.com/portfolios/three-fund-portfolio/?utm_source=gemini), [Total Stock Market Portfolio](https://portfoliocharts.com/portfolios/total-stock-market-portfolio/?utm_source=gemini), [Global Market Portfolio](https://portfoliocharts.com/portfolios/global-market-portfolio/?utm_source=gemini), [Core Four Portfolio](https://portfoliocharts.com/portfolios/core-four-portfolio/?utm_source=gemini)
* **수학적 특성**: 자본 투입 비율 고정, 위험은 주식 베타($\beta_{\text{equity}}$)가 지배($\text{TRC}_{\text{stock}} > 80\%$).


### Equal-Weighted (자산군 균등 분할 및 기금 모델)

예일대 기금 모델(Swensen)이나 균등 가중(1/N) 철학에 기반하여 주식, 채권, 부동산(REITs), 원자재 등 서로 다른 성격의 자산군 바스켓을 균등하거나 준균등하게 배분하는 모델입니다.

* **포트폴리오**: [7Twelve Portfolio](https://portfoliocharts.com/portfolios/7twelve-portfolio/?utm_source=gemini), [Ivy Portfolio](https://portfoliocharts.com/portfolios/ivy-portfolio/?utm_source=gemini), [Swensen Portfolio](https://portfoliocharts.com/portfolios/swensen-portfolio/?utm_source=gemini), [Pinwheel Portfolio](https://portfoliocharts.com/portfolios/pinwheel-portfolio/?utm_source=gemini), [Sandwich Portfolio](https://portfoliocharts.com/portfolios/sandwich-portfolio/?utm_source=gemini), [Richer Retirement Portfolio](https://portfoliocharts.com/portfolios/richer-retirement-portfolio/?utm_source=gemini)
* **수학적 특성**: 단순 $1/N$ 배분을 통한 추정 오차(Estimation Error) 제거 및 Shannon's Demon식 재조정 효과 극대화.

### Volatility Parity (수학적 변동성 패리티)

공분산 행렬 $\boldsymbol{\Sigma}$와 변동성 벡터 $\boldsymbol{\sigma}$를 명시적으로 고려하여, 고변동성 자산의 비중을 극도로 낮추고 저변동성 자산의 비중을 높여 한계 위험 기여도(MCR)를 일치시키는 모델입니다.

* **포트폴리오**: [Larry Portfolio](https://portfoliocharts.com/portfolios/larry-portfolio/?utm_source=gemini), [All Seasons Portfolio](https://portfoliocharts.com/portfolios/all-seasons-portfolio/?utm_source=gemini)
* **수학적 특성**: $w_i (\boldsymbol{\Sigma}\mathbf{w})_i \approx w_j (\boldsymbol{\Sigma}\mathbf{w})_j$, 채권 비중이 $70\sim 85\%$ 수준으로 극대화됨.

### Regime Parity (거시 경제 국면 패리티)

상관계수 시계열의 불안정성(Non-stationarity)을 전제하고, 거시경제 4대 국면(인플레이션, 디플레이션, 성장, 침체)에 직교(Orthogonal)하는 자산(주식, 장기국채, 금, 현금)을 대칭 배분하는 모델입니다.

* **포트폴리오**: [Permanent Portfolio](https://portfoliocharts.com/portfolios/permanent-portfolio/?utm_source=gemini), [Golden Butterfly Portfolio](https://portfoliocharts.com/portfolios/golden-butterfly-portfolio/?utm_source=gemini), [Golden Ratio Portfolio](https://portfoliocharts.com/portfolios/golden-ratio-portfolio/?utm_source=gemini), [Weird Portfolio](https://portfoliocharts.com/portfolios/weird-portfolio/?utm_source=gemini)
* **수학적 특성**: 주식-채권 외에 금(Gold) 등 대체 자산을 $20\%$ 이상 편입하여 체제 전환(Regime Shift) 시 테일 리스크 방어.

### Factor-Tilted (베타 / 팩터 틸트형)

단순 자산군 분산이 아니라 Fama-French의 규모(Size), 가치(Value) 프리미엄을 통해 기대수익률 벡터 $\boldsymbol{\mu}$를 상향시키면서 자산군 간 약한 상관계수를 활용하는 모델입니다.

* **포트폴리오**: [Coffeehouse Portfolio](https://portfoliocharts.com/portfolios/coffeehouse-portfolio/?utm_source=gemini), [Ultimate Buy and Hold Portfolio](https://portfoliocharts.com/portfolios/ultimate-buy-and-hold-portfolio/?utm_source=gemini), [No-Brainer Portfolio](https://portfoliocharts.com/portfolios/no-brainer-portfolio/?utm_source=gemini), [Ideal Index Portfolio](https://portfoliocharts.com/portfolios/ideal-index-portfolio/?utm_source=gemini)
* **수학적 특성**: $\boldsymbol{\mu}$의 팩터 프리미엄 극대화, 동일 자산군 내 세부 스타일 다변화.

---

## 21개 포트폴리오 매핑 요약

| 분류 레이블 | 설계 메커니즘 | 해당 포트폴리오 (21개 전수 매핑) |
| --- | --- | --- |
| **Value Weighted** | 시가총액 가중 및 주식/채권 단순 자본 분할 | Classic 60-40, Three-Fund, Total Stock Market, Global Market, Core Four |
| **Equal Weighted** | 다중 자산군(실물자산 포함) 균등/준균등 분할 | 7Twelve, Ivy, Swensen, Pinwheel, Sandwich, Richer Retirement |
| **Volatility Parity** | 공분산 및 변동성 기반 위험 기여도 균등화 | Larry, All Seasons |
| **Regime Parity** | 4대 거시 경제 국면(물가·성장) 직교 자산 배분 | Permanent, Golden Butterfly, Golden Ratio, Weird |
| **Factor-Tilted** | 소형주·가치주 등 스타일 팩터 프리미엄 결합 | Coffeehouse, Ultimate Buy and Hold, No-Brainer, Ideal Index |

