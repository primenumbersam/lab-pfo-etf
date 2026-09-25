---
title: "US ETF Portfolio Score"
---

# US ETF Portfolios' Annual Score

US ETFs를 이용한 21개의 자산배분 전략을 4개의 계량 평가 기준으로 종합 평가해 보는 대화형 순위 현황판.

<video controls preload="metadata" playsinline poster="images/video-frame.jpg" style="width: 100%; max-width: 820px; border-radius: 8px; margin: 1rem 0; box-shadow: 0 4px 12px rgba(0,0,0,0.15); display: block;">
  <source src="video.mp4" type="video/mp4">
  브라우저가 HTML5 비디오 재생을 지원하지 않는다.
</video>


<style>
.metric-weights-panel {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(210px, 1fr));
  gap: 1rem;
  padding: 1.25rem;
  border: 1px solid rgba(128,128,128,0.25);
  border-radius: 8px;
  background: var(--theme-box-bg, rgba(255,255,255,0.02));
  margin: 1.5rem 0;
}
.weight-control {
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
}
.weight-control label {
  font-size: 0.82rem;
  font-weight: 600;
  display: flex;
  justify-content: space-between;
}
.weight-control input[type="range"] {
  width: 100%;
  cursor: pointer;
}
.cat-btn-group {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin: 1rem 0;
}
.cat-btn {
  padding: 0.35rem 0.75rem;
  border-radius: 9999px;
  border: 1px solid var(--theme-border, #4c566a);
  background: transparent;
  color: inherit;
  cursor: pointer;
  font-size: 0.82rem;
  transition: all 0.2s ease;
}
.cat-btn.active, .cat-btn:hover {
  background: #3b82f6;
  border-color: #3b82f6;
  color: #fff;
}
.ranking-table {
  display: table !important;
  width: 100% !important;
  min-width: 100%;
  border-collapse: collapse;
  font-size: 0.88rem;
  margin: 0.5rem 0 1.5rem 0;
}
.ranking-table th {
  padding: 0.75rem 0.5rem;
  border-bottom: 2px solid var(--theme-border, #4c566a);
  text-align: left;
  font-weight: 600;
  font-size: 0.82rem;
  color: #888;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}
.ranking-table td {
  padding: 0.75rem 0.5rem;
  border-bottom: 1px solid rgba(128,128,128,0.18);
}
.ranking-row {
  cursor: pointer;
  transition: background 0.15s ease;
}
.ranking-row:hover {
  background: rgba(59, 130, 246, 0.08);
}
.rank-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 26px;
  height: 26px;
  border-radius: 50%;
  font-weight: 700;
  font-size: 0.82rem;
  background: rgba(128,128,128,0.15);
}
.rank-1 { background: #fbbf24; color: #000; }
.rank-2 { background: #94a3b8; color: #000; }
.rank-3 { background: #d97706; color: #fff; }
.score-pill {
  font-weight: 700;
  color: #3b82f6;
  font-size: 0.95rem;
}
.cat-tag {
  display: inline-block;
  font-size: 0.72rem;
  padding: 0.15rem 0.45rem;
  border-radius: 4px;
  background: rgba(128,128,128,0.15);
}
</style>

## Metric Weights (Interactive)

평가지표 가중치를 변경하면 실시간으로 21개 ETF Portfolio Strategies의 종합 점수와 순위가 계산됨.

<div class="metric-weights-panel">
  <div class="weight-control">
    <label for="w-geom"><span>수익률 ({{< katex >}}g_p{{< /katex >}})</span> <span id="val-geom">25%</span></label>
    <input type="range" id="w-geom" min="0" max="100" value="25" oninput="onWeightChange()">
  </div>
  <div class="weight-control">
    <label for="w-omega"><span>비대칭도 ({{< katex >}}\Omega{{< /katex >}})</span> <span id="val-omega">25%</span></label>
    <input type="range" id="w-omega" min="0" max="100" value="25" oninput="onWeightChange()">
  </div>
  <div class="weight-control">
    <label for="w-dd"><span>낙폭 통제 (RMSE DD)</span> <span id="val-dd">25%</span></label>
    <input type="range" id="w-dd" min="0" max="100" value="25" oninput="onWeightChange()">
  </div>
  <div class="weight-control">
    <label for="w-sens"><span>타이밍 무관성 (Sensitivity)</span> <span id="val-sens">25%</span></label>
    <input type="range" id="w-sens" min="0" max="100" value="25" oninput="onWeightChange()">
  </div>
</div>

---

## ETF Portfolio Rankings by Weighted Metric

(Real Rate = CPI adjusted 10 year Rolling Window 15th percentile)

<div class="cat-btn-group">
  <button class="cat-btn active" data-cat="all" onclick="setCategoryFilter('all')">All (21)</button>
  <button class="cat-btn" data-cat="value-weighted" onclick="setCategoryFilter('value-weighted')">Value-Weighted</button>
  <button class="cat-btn" data-cat="equal-weighted" onclick="setCategoryFilter('equal-weighted')">Equal-Weighted</button>
  <button class="cat-btn" data-cat="volatility-parity" onclick="setCategoryFilter('volatility-parity')">Volatility Parity</button>
  <button class="cat-btn" data-cat="regime-parity" onclick="setCategoryFilter('regime-parity')">Regime Parity</button>
  <button class="cat-btn" data-cat="factor-tilted" onclick="setCategoryFilter('factor-tilted')">Factor-Tilted</button>
</div>

<div style="overflow-x: auto;">
  <table class="ranking-table" id="ranking-table">
    <thead>
      <tr>
        <th style="width: 60px;">Rank</th>
        <th>Portfolio Name</th>
        <th>Category</th>
        <th style="text-align: right;">Real Rate ({{< katex >}}g_p{{< /katex >}})</th>
        <th style="text-align: right;">Omega ({{< katex >}}\Omega{{< /katex >}})</th>
        <th style="text-align: right;">RMSE DD</th>
        <th style="text-align: right;">Sensitivity</th>
        <th style="text-align: right; width: 100px;">Weighted Score</th>
      </tr>
    </thead>
    <tbody id="ranking-body">
      <!-- JavaScript populated -->
    </tbody>
  </table>
</div>

{{< ranking_scores >}}

<script>
let rawScores = [];
try {
  const sData = document.getElementById('scores-raw-data');
  if (sData) {
    rawScores = JSON.parse(sData.textContent);
    if (typeof rawScores === 'string') rawScores = JSON.parse(rawScores);
  }
} catch(e) {
  console.error('Failed to parse scores data:', e);
}
let currentCategory = 'all';

function normalize(val, min, max, invert=false) {
  if (max === min) return 50;
  let norm = ((val - min) / (max - min)) * 100;
  return invert ? 100 - norm : norm;
}

function calculateAndRender() {
  const w1 = parseFloat(document.getElementById('w-geom').value);
  const w2 = parseFloat(document.getElementById('w-omega').value);
  const w3 = parseFloat(document.getElementById('w-dd').value);
  const w4 = parseFloat(document.getElementById('w-sens').value);

  const totalW = (w1 + w2 + w3 + w4) || 1;

  const geomVals = rawScores.map(d => d.geometric_return);
  const omegaVals = rawScores.map(d => d.omega_ratio);
  const ddVals = rawScores.map(d => d.rmse_dd);
  const sensVals = rawScores.map(d => d.start_date_sensitivity);

  const minG = Math.min(...geomVals), maxG = Math.max(...geomVals);
  const minO = Math.min(...omegaVals), maxO = Math.max(...omegaVals);
  const minD = Math.min(...ddVals), maxD = Math.max(...ddVals);
  const minS = Math.min(...sensVals), maxS = Math.max(...sensVals);

  const evaluated = rawScores.map(item => {
    const s1 = normalize(item.geometric_return, minG, maxG);
    const s2 = normalize(item.omega_ratio, minO, maxO);
    const s3 = normalize(item.rmse_dd, minD, maxD, true);
    const s4 = normalize(item.start_date_sensitivity, minS, maxS, true);

    const finalScore = ((s1 * w1) + (s2 * w2) + (s3 * w3) + (s4 * w4)) / totalW;
    return {
      ...item,
      computedScore: Math.round(finalScore * 10) / 10
    };
  });

  evaluated.sort((a, b) => b.computedScore - a.computedScore);

  const tbody = document.getElementById('ranking-body');
  if (!tbody) return;
  tbody.innerHTML = '';

  let visibleRank = 1;
  evaluated.forEach(item => {
    const isVisible = currentCategory === 'all' || item.category_slug === currentCategory;
    if (!isVisible) return;

    const tr = document.createElement('tr');
    tr.className = 'ranking-row';
    tr.onclick = () => window.location.href = './portfolios/' + item.slug + '/';

    let rankBadgeClass = '';
    if (visibleRank === 1) rankBadgeClass = 'rank-1';
    else if (visibleRank === 2) rankBadgeClass = 'rank-2';
    else if (visibleRank === 3) rankBadgeClass = 'rank-3';

    tr.innerHTML = `
      <td><span class="rank-badge ${rankBadgeClass}">${visibleRank}</span></td>
      <td style="font-weight: 600;"><a href="./portfolios/${item.slug}/" style="color: inherit; text-decoration: none;">${item.name}</a></td>
      <td><span class="cat-tag">${item.category}</span></td>
      <td style="text-align: right; color: #10b981; font-weight: 500;">${item.geometric_return}%</td>
      <td style="text-align: right; color: #3b82f6; font-weight: 500;">${item.omega_ratio}</td>
      <td style="text-align: right; color: #f59e0b; font-weight: 500;">${item.rmse_dd}%</td>
      <td style="text-align: right; color: #ec4899; font-weight: 500;">${item.start_date_sensitivity}</td>
      <td style="text-align: right;"><span class="score-pill">${item.computedScore.toFixed(1)}</span></td>
    `;
    tbody.appendChild(tr);
    visibleRank++;
  });
}

function updateWeightsAndRender() {
  const wGeom = parseFloat(document.getElementById('w-geom').value) || 0;
  const wOmega = parseFloat(document.getElementById('w-omega').value) || 0;
  const wDd = parseFloat(document.getElementById('w-dd').value) || 0;
  const wSens = parseFloat(document.getElementById('w-sens').value) || 0;

  const weights = [wGeom, wOmega, wDd, wSens];
  const total = weights.reduce((a, b) => a + b, 0);

  let pGeom = 0, pOmega = 0, pDd = 0, pSens = 0;
  if (total > 0) {
    const exact = weights.map(w => (w / total) * 100);
    const floored = exact.map(v => Math.floor(v));
    const diff = 100 - floored.reduce((a, b) => a + b, 0);
    const remainders = exact.map((v, i) => ({ index: i, rem: v - floored[i] }))
                            .sort((a, b) => b.rem - a.rem);
    for (let i = 0; i < diff; i++) {
      floored[remainders[i].index]++;
    }
    [pGeom, pOmega, pDd, pSens] = floored;
  }

  document.getElementById('val-geom').textContent = pGeom + '%';
  document.getElementById('val-omega').textContent = pOmega + '%';
  document.getElementById('val-dd').textContent = pDd + '%';
  document.getElementById('val-sens').textContent = pSens + '%';

  calculateAndRender();
}

function onWeightChange() {
  updateWeightsAndRender();
}

function setCategoryFilter(cat) {
  currentCategory = cat;
  document.querySelectorAll('.cat-btn').forEach(btn => {
    btn.classList.toggle('active', btn.dataset.cat === cat);
  });
  calculateAndRender();
}

document.addEventListener('DOMContentLoaded', updateWeightsAndRender);
if (document.readyState === 'interactive' || document.readyState === 'complete') {
  updateWeightsAndRender();
}
</script>
