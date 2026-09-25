---
title: "ETFs"
weight: 30
---

# ETF Universe

{{% details title="ETF 데이터 (data/etfs.json) 관리 파이프라인" open=false %}}
`data/etfs.json` 파일은 `tools/fetcher`를 통해 Yahoo Finance API에서 자동 추출한 37개 ETF의 최신 AUM, 수수료율, 상위 편입 종목(Top Holdings)을 정규화하여 관리한다.
{{% /details %}}

<h2 style="margin-top: 1.25rem; margin-bottom: 0.25rem;">ETF Heatmap</h2>
<p style="margin-top: 0; margin-bottom: 0.4rem; font-size: 0.85rem; color: #888;">
  ETFs categorized by classes and types. Size represents recent AUM, and Heat represents return performance (3M Real, 3M Nominal, YTD, 1Y, 1M).
</p>

{{< etf_treemap >}}

## ETF List
총자산(AUM) 순으로 기본 정렬됨. 각 열 헤더를 클릭하면 해당 기준(Asset Type, AUM, 수수료율 등)으로 오름차순/내림차순 정렬됨.

{{< etf_table >}}
