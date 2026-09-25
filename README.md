# US ETF Portfolio Comparison & Analytics

Go 언어 데이터 파이프라인과 `Hugo Extended` (v0.158.0+) 기반으로 구축된 포트폴리오 분석 및 시각화 플랫폼.  
장기 실질 총수익률(Real Total Return, 1995~2025 실측치)과 비대칭 위험 측정치(Omega Ratio, RMSE of Drawdown, Start Date Sensitivity)를 기반으로 21개 자산배분전략과 37개 대표 ETF Universe를 interactive하게 탐색한다.

---

## 프로젝트 개요

- **목적**: 1995년부터 2025년까지 최근 31년간의 실측 시계열 데이터를 바탕으로 ETF를 이용한 여러 자산배분 전략들의 실질 복리 성과와 경로 의존성(Path Dependency)을 다차원 분석.
- **데이터 주기**: 연간/월간 단위 갱신, `data/*.json` 형태로 정적 적재 후 정적 사이트 빌드.
- **기술 스택**:
  - **수집/연산 백엔드 (Data Pipeline)**: Go (`tools/fetcher`, `tools/calculator`, `goquery` 스크래핑, Yahoo Finance API 연동)
  - **정적 사이트 빌더 (Front-end SSG)**: Hugo Extended (`hugo-book` 단일 테마)
  - **시각화 (Interactive Charts)**: Apache ECharts (ETF Universe Heatmap, 31×31 Return Heatmap, 10Y Rolling, CAGR Funnel)
  - **수식 렌더링**: Goldmark Passthrough + KaTeX

---

## ETF 유니버스 및 자산군 분류

`data/portfolios.json`의 21개 포트폴리오 전략 구성에 포함된 18개 고유 자산 유형(`asset_type`)을 4대 자산군(`asset_class`)으로 분류한 37개 대표 ETF.

### 1. US Equities (6종)
- Stocks-US-Large: SPY, VOO, QQQ, DIA
- Stocks-US-Small-Value: VBR
- Stocks-US-Small: VB, IJR
- Stocks-US-Large-Value: SCHD, VIG, SPYV
- Stocks-US-Large-Growth: SPYG, VUG, IWF
- Stocks-US-Small-Growth: VBK

### 2. Global Equities (4종)
- Stocks-Developed-Large: VEA, SCHF, VT, URTH
- Stocks-Emerging-Large: VWO, IEMG
- Stocks-Developed-Small: VSS
- Stocks-Developed-Large-Value: EFV

### 3. Fixed Income & Cash (5종)
- Bond-US-Intermediate: VGIT, GOVT
- Bond-US-Bills: SGOV, VBIL, BIL
- Bond-US-Long: VGLT, TLT
- Bond-US-Short: VGSH, SHY
- Bond-Developed-Intermediate: IGOV

### 4. Commodities & Real Estate (3종)
- REITs-US: VNQ, SCHH
- Gold-Global: IAU, GLD
- Commodities-Global: GSG

---

## 사이트 정보 구조 및 페이지별 사양

### 1. Main 대시보드 (`/`)
- **ETF Portfolio Score**: 21개 자산배분전략의 횡단면적 순위 시각화.
- **성과평가 지표 (metrics)**:
  1. Geometric Return ($g_p$): 10년 롤링 윈도우 기준 보수적 복리 실질 수익률 (15th percentile)
  2. Omega Ratio ($\Omega$): T-Bill 실질 금리를 임계값($\tau$)으로 둔 비대칭 수익 지표
  3. RMSE of DD: 하락 폭과 체류 기간을 2차 페널티로 반영한 누적 고통 지표 (Ulcer Index)
  4. Start Date Sensitivity: 10년 보유 시 진입 시점에 따른 성과 분산도
- 사용자 가중치($w_k$) 실시간 슬라이더 조정 및 종합 점수(Composite Score) 실시간 재계산 인터랙션 제공.

### 2. Portfolios 페이지 (`/portfolios/`)
- **포트폴리오 카드 그리드 (`{{< portfolio_grid >}}`)**: `data/portfolios.json` 및 `data/scores.json`을 기반으로 동적 렌더링. 5개 카테고리(Value-Weighted, Equal-Weighted, Volatility Parity, Regime Parity, Factor-Tilted) 실시간 필터 버튼 탑재.
- **포트폴리오 상세 분석 (`{{< portfolio_detail >}}`)**:
  - Author 정보 및 세부 설명·자산군 유의사항(`details` 태그 내 `note` 배치).
  - 자산 배분 비중 테이블.
  - ECharts 기반 4개 interactive charts:
    - Col-2 (상단): 10년 롤링 복리 실질 수익률(Rolling 10Y CAGR) & 투자 기간별 분위수 수렴(CAGR along N Funnel, 1~30년)
    - Full-Width (중단): 1995~2025 연도별 실질 수익률 바 차트 (손실 연도 통계 배지)
    - Full-Width (하단): 31×31 시작 연도 $\times$ 종료 연도 실질 복리수익률 히트맵(Return Heat Map, 1:1 정방형 격자, 우상향 직각삼각형 배치)

### 3. ETFs 유니버스 페이지 (`/etfs/`)
- **ETF Universe Heatmap (`{{< etf_treemap >}}`)**: 최신 AUM(면적)과 기간별 수익률 지표(색상: 3M Real, 3M Nominal, YTD, 1Y, 1M 다중 전환) 기반 반응형 트리맵.
- **동적 정렬 ETF 테이블 (`{{< etf_table >}}`)**: Ticker, ETF Name, Index, Inception Date, Expense Ratio, AUM, Asset Class, Asset Type 다중 컬럼 오름차순/내림차순 바닐라 JS 정렬 지원.

### 4. Glossary 페이지 (`/glossary/`)
- 메트릭 연산 공식(LaTeX), 5개 portfolio categories 정의 표, 시각화 차트 해석법, 1995~2025 데이터 파이프라인(FRED, Fama-French, Portfolio Charts, Yahoo Finance) 및 백테스팅 설계 수록.

---

## 주요 시스템 아키텍처

1. **데이터 수집 파이프라인 (`tools/fetcher`)**
   - Yahoo Finance API (`quoteSummary`, `chart`) 비동기 호출로 37개 ETF의 AUM, 운용보수, Holdings, 실시간 기간별 수익률 수집.
   - Chart API의 `firstTradeDate`로부터 상장일(`inception_date`) 자동 변환.
   - `goquery` 기반 웹 프로필 스크래핑과 공식 벤치마크 지수 카탈로그(`etfIndexCatalog`)를 결합한 하이브리드 추종 지수(`index`) 동기화.

2. **포트폴리오 연산 엔진 (`tools/calculator`)**
   - `data/series/asset_returns.json`의 31년간(1995~2025) 실측 자산군 실질 수익률과 21개 포트폴리오 비중 결합.
   - 연간 리밸런싱 가정 하에 10년 롤링 기하평균, Omega, RMSE of DD, Start Date Sensitivity, 1~30년 퍼널 통계, 31×31 삼각형 히트맵 매트릭스를 일괄 연산하여 `data/scores.json` 및 `data/series/<slug>.json` 생성.

3. **단일 테마 및 컴포넌트 숏코드**
   - `hugo-book` 단일 테마 기반 경량화 아키텍처.
   - `layouts/shortcodes/` 내 `portfolio_grid.html`, `portfolio_detail.html`, `etf_table.html`, `etf_treemap.html`, `katex.html`로 모듈화.

4. **수식 및 인터랙티브 시각화**
   - Goldmark `passthrough` 확장 및 KaTeX 수식 자동 렌더링.
   - Apache ECharts CDN 연동 및 다크모드/라이트모드 자동 감지 테마 반응형 렌더링.

---

## 프로젝트 구조

```text
lab-pfo-etf/
├── assets/
│   └── styles/custom.css       # Hugo Pipes 커스텀 CSS (다크모드/Callout/코드복사)
├── static/                     # 정적 미디어 에셋 (웹 루트 1:1 복사)
│   └── images/                 # frame.jpg, character-samsoon.jpg 등 정적 이미지
├── content/                    # 마크다운 콘텐츠
│   ├── _index.md               # 메인 랜딩 (Portfolio Score 인터랙티브 대시보드)
│   ├── portfolios/             # 포트폴리오 목록 (_index.md) 및 21개 상세 페이지
│   ├── etfs/                   # ETF Universe 트리맵 및 정렬 테이블
│   └── glossary/               # 계량 지표 정의, 방법론 및 데이터 사전 (_index.md)
├── data/                       # 정적 JSON 데이터
│   ├── etfs.json               # 37개 ETF 스냅샷 (AUM, 수수료, Index, Inception Date, Holdings)
│   ├── portfolios.json         # 21개 포트폴리오 구성 비중, Author, Note
│   ├── scores.json             # 21개 포트폴리오 4대 핵심 메트릭 및 랭킹 요약
│   └── series/                 # 자산군 및 포트폴리오별 장기 시계열 (1995~2025)
│       ├── asset_returns.json  # 18개 자산군 실질 총수익률
│       └── *.json              # 21개 포트폴리오별 31x31 히트맵 및 퍼널 시계열
├── layouts/
│   └── shortcodes/             # 사이트 커스텀 숏코드
│       ├── portfolio_grid.html # 포트폴리오 카드 그리드 및 필터
│       ├── portfolio_detail.html # 상세 분석 (Author/Note details, 4대 ECharts)
│       ├── etf_table.html      # 바닐라 JS 정렬 테이블
│       ├── etf_treemap.html    # ETF 유니버스 트리맵
│       └── katex.html          # KaTeX 수식 블록
├── tools/                      # Go 기반 데이터 수집/연산 엔진
│   ├── fetcher/main.go         # Yahoo Finance + goquery ETF 메타데이터 수집기
│   ├── calculator/main.go      # 포트폴리오 백테스팅 메트릭 연산기
│   └── goquery/main.go         # goquery 웹 스크래핑 예제
├── themes/
│   └── hugo-book/              # 핵심 북 테마 엔진
├── hugo.yaml                   # 사이트 설정
└── README.md
```

---

## 시작하기

### 요구 사양
- [Hugo Extended](https://github.com/gohugoio/hugo/releases) v0.158.0 이상 (권장: v0.165.0+)
- Go 1.25+ (데이터 수집기 및 계산기 실행 시)

### 실행 방법

```bash
# 1. ETF 메타데이터 및 실시간 수익률 동기화
go run ./tools/fetcher

# 2. 포트폴리오 백테스팅 메트릭 및 시계열 연산
go run ./tools/calculator

# 3. 로컬 개발 서버 실행
hugo serve
```

브라우저에서 `http://localhost:1313/lab-pfo-etf/` 접속.
