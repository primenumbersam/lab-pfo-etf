---
title: "Glossary"
weight: 40
bookToc: true
cascade:
  bookToc: true
---

# Glossary

US ETF Portfolio Comparison 을 위한 Definition, Methodology & Data Dictionary

---

## 1. 4대 핵심 평가 메트릭 (Core Metrics)

### (1) Geometric Return ($g_p$, 복리 실질 수익률)

10년 롤링 윈도우(Rolling 10Y Window)를 기준으로 측정한 장기 보수적 복리 실질 수익률이다. 단일 연도의 일시적 급등락에 따른 행운을 배제하고 투자자가 실제로 체감할 수 있는 하방 보수치를 측정하기 위해, 10년 롤링 기하평균 수익률 분포의 15번째 백분위수(15th Percentile)를 기본 대표값으로 채택한다.

$$ g_p = \left( \prod_{t=1}^{T} (1 + R_{p,t}^{\text{real}}) \right)^{\frac{1}{T}} - 1 $$

여기서 $R_{p,t}^{\text{real}}$은 연도 $t$의 물가조정(CPI 차감) 연간 실질 총수익률을 나타낸다.

### (2) Omega Ratio ($\Omega$, 비대칭 위험수익비)

전통적인 샤프 지수(Sharpe Ratio)가 전제하는 정규분포 가정을 배제하고, 수익률 분포의 왜도(Skewness)와 첨도(Kurtosis) 등 고차 모멘트를 직접 반영하는 비대칭 위험 지표다. 무위험 단기 자산(T-Bill 실질 수익률)을 임계 수익률 $\tau$로 설정하여, 임계치 대비 손실 기대값에 대한 초과 이익 기대값의 누적 확률 면적 비율을 측정한다.

$$ \Omega(\tau) = \frac{\int_{\tau}^{\infty} (1 - F(r)) \, dr}{\int_{-\infty}^{\tau} F(r) \, dr} = \frac{\mathbb{E}[\max(R - \tau, 0)]}{\mathbb{E}[\max(\tau - R, 0)]} $$

$\Omega > 1.0$은 임계값 이상의 이익 기대치가 손실 위험을 상회함을 의미하며, 값이 클수록 하방 위험 대비 상방 수익 잠재력이 우수하다.

### (3) RMSE of DD (Root Mean Squared Error of Drawdown, 누적 하락 고통)

고점 대비 하락률(Drawdown)의 최고점 낙폭(MDD)뿐만 아니라, 하락 구간에 머무는 체류 기간(Duration)에 2차 페널티(Squared Penalty)를 부여하는 누적 하락 지표다. 피터 마틴(Peter Martin)의 얼서 지수(Ulcer Index)와 수학적으로 동등하며, 장기 침체 구간이 길어질수록 지표값이 가파르게 증가한다.

$$ \text{RMSE}_{\text{DD}} = \sqrt{ \frac{1}{T} \sum_{t=1}^{T} D_t^2 } $$

$$ D_t = \frac{\max_{0 \le s \le t} W_s - W_t}{\max_{0 \le s \le t} W_s} \times 100 (\%) $$

여기서 $W_t$는 시점 $t$의 포트폴리오 누적 자산 가치를 뜻한다. 수치가 낮을수록 하락 충격이 얕고 회복 탄력성이 빠르다.

### (4) Start Date Sensitivity (진입 시점 민감도)

투자 시작 연도에 따라 최종 누적 성과가 엇갈리는 경로 의존성(Path Dependency) 위험을 수치화한 척도다. 1995년부터 2025년까지 임의의 연도에 투자를 시작해 10년을 보유했을 때 얻게 되는 연환산 복리 수익률들의 표준편차를 정규화하여 측정한다.

$$ \sigma_{\text{start}} = \sqrt{ \frac{1}{K} \sum_{s=1}^{K} (g_s - \bar{g})^2 } $$

민감도 수치가 낮을수록 시장 진입 타이밍에 구애받지 않고 언제 투자를 시작하든 균일한 성과를 보장한다.

### (5) Composite Score & Rank (종합 점수 및 순위)

메인 대시보드에서 21개 포트폴리오의 상대적 우열을 가리기 위해 4대 지표를 0~100점 범위로 Min-Max 정규화하고 가중합산(Weighted Sum)하여 산출한다.

- 수익 지표($g_p, \Omega$): $S_i = \frac{X_i - \min(X)}{\max(X) - \min(X)} \times 100$
- 위험 지표($\text{RMSE}_{\text{DD}}, \sigma_{\text{start}}$): $S_i = \frac{\max(X) - X_i}{\max(X) - \min(X)} \times 100$ (역방향 점수화)

$$ \text{Composite Score} = \sum_{k=1}^{4} w_k \cdot S_k, \quad \sum_{k=1}^{4} w_k = 1.0 $$

---

## 2. 포트폴리오 전략 분류 체계 (5 Categories)

본 플랫폼에 수록된 21개 자산배분 모델은 자산 선정 및 비중 결정 원리에 따라 5대 카테고리로 분류된다.

| 카테고리 | 핵심 철학 및 비중 결정 원리 | 소속 포트폴리오 (총 21개) |
| :--- | :--- | :--- |
| Value-Weighted | 자본시장 시가총액 가중 및 전통적 주식·채권 단순 결합 | Classic 60-40, Three-Fund, Total Stock Market, Global Market, Ideal Index |
| Equal-Weighted | 편입된 4~12개 자산군을 동일 비율(1/N)로 기계적 분산 | 7Twelve, Coffeehouse, No-Brainer, Sandwich, Ultimate Buy and Hold, Pinwheel |
| Volatility Parity | 자산별 과거 변동성에 반비례하여 리스크 기여도를 균등 배분 | All Seasons, Golden Ratio |
| Regime Parity | 번영·침체·인플레이션·디플레이션 4대 경제 국면에 대응하는 자산 배치 | Golden Butterfly, Permanent, Weird |
| Factor-Tilted | 규모(Size), 가치(Value), 배당 등 학술적 초과수익 팩터 가중 | Core Four, Ivy, Larry, Richer Retirement, Swensen |

---

## 3. 시각화 컴포넌트 구조 및 해석법

### (1) 시작 연도 × 종료 연도 실질 수익률 히트맵 (Return Heat Map)
- **축 구성**: 가로축은 투자 종료 연도(1995~2025), 세로축은 투자 시작 연도(1995~2025)를 나타낸다.
- **기하학적 형태**: 종료 연도는 항상 시작 연도 이상($\text{End Year} \ge \text{Start Year}$)이므로, 주대각선 상단에 데이터가 집중되는 31×31 우상향 직각삼각형(Upper Right Triangle)으로 렌더링된다. 모든 행의 우측 끝은 2025년으로 일직선 정렬된다.
- **정방형(Square) 격자**: 가로 셀 너비와 세로 셀 높이를 1:1 동일한 픽셀 단위로 고정하여 시각적 왜곡을 방지한다.
- **색상 스케일**: 연환산 복리 실질 수익률(CAGR)을 기준으로 -10% 이하(진한 적색)부터 0%(중립 회색 `#334155`), +10% 이상(진한 녹색)까지 대칭 스펙트럼으로 매핑한다.

### (2) 투자 기간별 수익률 수렴 깔때기 (CAGR Distribution along N)
- **구조**: 보유 기간 $N$이 1년부터 30년까지 증가함에 따라 연환산 실질 복리 수익률의 변동폭이 어떻게 수렴하는지를 시각화한다.
- **분위수 밴드**: 최댓값(Max), 상위 85%, 중앙값(Median), 하위 15%, 최솟값(Min) 5개 분위수 영역을 퍼널(Funnel) 형태로 배치하여, 보유 기간이 길어질수록 원금 손실 확률이 소멸하고 장기 기대값으로 좁혀지는 패턴을 보여준다.

### (3) 10년 롤링 복리 실질 수익률 (Rolling 10Y CAGR)
- **구조**: 매 10년간의 이동 구간(예: 1995~2004, 1996~2005, ..., 2016~2025)에 걸쳐 기록된 연환산 복리 실질 수익률 시계열 추세선이다.
- **해석**: 시장 사이클(닷컴 버블, 금융위기, 인플레이션 충격)을 통과하는 포트폴리오의 실질 방어력을 추적한다.

### (4) 연도별 실질 수익률 (Annual Real Returns)
- **구조**: 1995년부터 2025년까지 각 연도의 물가조정 실질 총수익률을 바 차트로 표시한다.
- **지표**: 전체 연도 중 플러스 수익 연도 수, 마이너스 수익 연도 수, 손실 연도 비율을 상단 배지에 표기한다.

---

## 4. 데이터 파이프라인 및 백테스팅 아키텍처

### (1) 실측 데이터셋 구축 (1995~2025, 31년)
- **모의 데이터 배제 원칙**: 백테스팅 데이터의 신뢰성을 담보하기 위해 가공된 추정 데이터나 모의 시계열을 배제하고, 실측 기반 통계가 검증된 1995년부터 2025년까지의 31년 시계열을 분석 범위로 통일했다.
- **CPI 차감 실질 수익률**: 세인트루이스 연방준비은행(FRED)의 미국 도시 소비자물가지수(`CPIAUCNS`) 연말 수치를 바탕으로 연간 인플레이션율을 도출하여 모든 명목 수익률에서 물가상승분을 차감했다.
- **자산군 원천 데이터**:
  - 미국 주식 팩터: Fama-French Data Library의 시가총액 가중 연간 수익률 (소형 가치, 소형 성장, 대형 가치, 대형 성장, 소형 블렌드).
  - 미국 국채: FRED 만기별 국채 금리(`TB3MS`, `GS1`, `GS10`, `GS20`) 기반 듀레이션 자본손익 및 쿠폰 수익률.
  - 글로벌 주식 및 대체 자산: Portfolio Charts 벤치마크 및 Yahoo Finance 수정종가, LBMA 금 현물가.

### (2) ETF 유니버스 동기화 (`tools/fetcher`)
- **실시간 API 연동**: Yahoo Finance v8/v10 엔드포인트(`quoteSummary`, `chart`)를 비동기 호출하여 37개 대표 ETF의 AUM, 운용보수, 상위 편입 종목(Top Holdings)을 정규화한다.
- **상장일 자동 추출**: Yahoo Finance Chart 메타데이터의 `firstTradeDate` 타임스탬프를 `YYYY-MM-DD`로 자동 변환하여 상장일(`inception_date`)을 취득한다.
- **공식 지수 매핑**: `goquery` 기반 웹 프로필 스크래핑과 공시 벤치마크 지수 카탈로그(`etfIndexCatalog`)를 결합하여 정확한 추종 지수(`index`)를 연동한다.

### (3) 포트폴리오 연산 엔진 (`tools/calculator`)
- **연간 리밸런싱**: 매년 12월 말 자산 비중을 원래 목표치로 재조정(Annual Rebalancing)하는 바이앤홀드(Buy-and-Hold) 모델을 적용한다.
- **배당 및 이자 전액 재투자**: 모든 현금흐름은 발생 즉시 재투자되는 총수익(Gross Total Return) 기준이다.
- **산출 결과물**: `data/scores.json`(21개 포트폴리오 지표 요약) 및 `data/series/<slug>.json`(연도별 수익률, 10년 롤링, 퍼널 통계, 31×31 히트맵 데이터)으로 자동 적재된다.

---

## 5. 면책 조항 (Disclaimer)

본 사이트의 모든 데이터, 수치, 차트 및 분석 결과는 비상업적 학술 연구 및 ETF 자산배분 전략 학습 목적으로 제공된다. 금융투자상품에 대한 매수 또는 매도 권유가 아니며, 과거 성과가 미래 수익을 보장하지 않는다.
