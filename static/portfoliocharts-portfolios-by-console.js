(async () => {
  const slugs = [
    "7twelve-portfolio", "all-seasons-portfolio", "classic-60-40-portfolio",
    "coffeehouse-portfolio", "core-four-portfolio", "global-market-portfolio",
    "golden-butterfly-portfolio", "golden-ratio-portfolio", "ideal-index-portfolio",
    "ivy-portfolio", "larry-portfolio", "no-brainer-portfolio",
    "permanent-portfolio", "pinwheel-portfolio", "richer-retirement-portfolio",
    "sandwich-portfolio", "swensen-portfolio", "three-fund-portfolio",
    "total-stock-market-portfolio", "ultimate-buy-and-hold-portfolio",
    "weird-portfolio"
  ];

  const database = {};

  for (const slug of slugs) {
    const url = `https://portfoliocharts.com/portfolios/${slug}/`;
    try {
      const html = await (await fetch(url)).text();
      const doc = new DOMParser().parseFromString(html, "text/html");
      const container = doc.querySelector(".entry-content, article, main") || doc.body;

      const name = doc.querySelector("h1")?.innerText.trim() || slug;
      const composition = [];
      const authorLines = [];
      const overviewLines = [];

      // 1. Asset Allocation 테이블 파싱
      const headings = Array.from(doc.querySelectorAll("h1, h2, h3"));
      const allocHeading = headings.find(h => h.innerText.trim().toLowerCase() === "asset allocation");
      if (allocHeading) {
        let sibling = allocHeading.nextElementSibling;
        let targetTable = null;
        while (sibling) {
          if (sibling.tagName === "TABLE") {
            targetTable = sibling;
            break;
          }
          const nested = sibling.querySelector("table");
          if (nested) {
            targetTable = nested;
            break;
          }
          sibling = sibling.nextElementSibling;
        }

        if (targetTable) {
          targetTable.querySelectorAll("tbody tr, tr").forEach(tr => {
            const cells = tr.querySelectorAll("td");
            if (cells.length >= 2) {
              const rawWeight = cells[0].innerText.trim().replace("%", "");
              const elementText = cells[1].innerText.trim();
              const weightVal = parseFloat(rawWeight);
              if (!isNaN(weightVal) && elementText.length > 0 && !elementText.toLowerCase().includes("asset class")) {
                composition.push({
                  element: elementText,
                  weight: Math.round((weightVal / 100) * 1000) / 1000
                });
              }
            }
          });
        }
      }

      // 2. Asset Notes 파싱 (CoBlocks 아코디언 블록 타겟팅)
      let assetNotes = "";
      const accordionTitles = doc.querySelectorAll(".wp-block-coblocks-accordion-item__title");
      for (const titleEl of accordionTitles) {
        if (titleEl.innerText.trim().toLowerCase().includes("asset notes")) {
          const itemWrapper = titleEl.closest(".wp-block-coblocks-accordion-item") || titleEl.parentElement;
          const contentEl = itemWrapper.querySelector(".wp-block-coblocks-accordion-item__content");
          if (contentEl) {
            assetNotes = contentEl.innerText.trim();
          } else {
            const clone = itemWrapper.cloneNode(true);
            const cloneTitle = clone.querySelector(".wp-block-coblocks-accordion-item__title");
            if (cloneTitle) cloneTitle.remove();
            assetNotes = clone.innerText.trim();
          }
          break;
        }
      }

      // 3. Author 및 Overview 파싱 (이전 정상 작동했던 State Machine 방식 적용)
      let currentSection = "";
      const elements = container.querySelectorAll("h1, h2, h3, p, ol, ul");

      elements.forEach(el => {
        const text = el.innerText.trim();
        if (!text) return;

        const lowerText = text.toLowerCase();

        if (el.tagName.startsWith("H")) {
          if (lowerText === "author") {
            currentSection = "author";
            return;
          } else if (lowerText === "overview") {
            currentSection = "overview";
            return;
          } else if (["performance", "charts", "comparisons", "alternatives", "articles", "discussion", "asset allocation"].includes(lowerText)) {
            currentSection = "other";
            return;
          }
        }

        if (currentSection === "author") {
          if (el.tagName === "P" || el.tagName === "H2") {
            authorLines.push(text);
          }
        } else if (currentSection === "overview") {
          if (text.includes("Featured Discussion")) {
            currentSection = "other";
            return;
          }
          if (el.tagName === "P" || el.tagName === "OL" || el.tagName === "UL") {
            overviewLines.push(text);
          }
        }
      });

      database[slug] = {
        name,
        composition,
        asset_notes: assetNotes,
        author: authorLines.join("\n\n").trim(),
        overview: overviewLines.join("\n\n").trim()
      };

      console.log(`Parsed: ${slug} | Comp: ${composition.length} | Notes: ${assetNotes.length > 0 ? "O" : "X"} | Author: ${authorLines.length > 0 ? "O" : "X"} | Overview: ${overviewLines.length > 0 ? "O" : "X"}`);
    } catch (e) {
      console.error(`Error parsing ${slug}:`, e);
    }
  }

  const blob = new Blob([JSON.stringify(database, null, 2)], { type: "application/json" });
  const a = document.createElement("a");
  a.href = URL.createObjectURL(blob);
  a.download = "portfoliocharts.json";
  a.click();
})();