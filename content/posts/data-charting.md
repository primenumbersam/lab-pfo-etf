---
date: 2023-09-23
tags: [tools]
---

# Data Charts and Fetch

## Chart.js & ECharts 3D

{{< chartjs data="sample_mock_2026" key="sales_bar" title="2026년 월별 매출 및 순이익" >}}

{{< echarts data="sample_mock_2026" key="workload_3d" title="주간/시간대별 작업 부하 3D Bar" height="500px" >}}

## Web Scrapping

Go vs Python 

Go 생태계에도 목적과 사용 패턴에 따라 잘 검증된 도구들이 마련되어 있습니다. 파이썬의 BeautifulSoup이나 Scrapy처럼 Go에서도 CSS 셀렉터 기반 추출부터 고성능 크롤러 프레임워크까지 지원합니다.

| 도구 | 특징 및 용도 | Python 대응 |
| --- | --- | --- |
| **goquery** | CSS 셀렉터 기반 DOM 파싱 (jQuery 문법 유사). 단순 URL 요청 후 파싱. | BeautifulSoup |
| **colly** | 비동기, 병렬 처리, 자동 재시도, Rate Limiting을 내장한 웹 크롤링 도구. | Scrapy |
| **chromedp** | Headless Chrome 제어. JavaScript 렌더링(SPA, 동적 데이터)이 필요한 사이트용. | Selenium, Playwright |
| **net/html** | Go 공식 서브패키지(`golang.org/x/net/html`). 외부 의존성 없이 저수준 노드 트리 탐색. | html.parser |
