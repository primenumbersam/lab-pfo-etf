(function () {
  document.querySelectorAll(".highlight").forEach(block => {
    const btn = document.createElement("button");
    btn.className = "copy-code-btn";
    btn.textContent = "Copy";
    btn.setAttribute("aria-label", "Copy code");

    btn.addEventListener("click", () => {
      const code = block.querySelector("pre code") || block.querySelector("pre");
      if (!code) return;

      navigator.clipboard.writeText(code.innerText.trim()).then(() => {
        btn.textContent = "Copied!";
        setTimeout(() => { btn.textContent = "Copy"; }, 1500);
      });
    });

    block.appendChild(btn);
  });
})();
