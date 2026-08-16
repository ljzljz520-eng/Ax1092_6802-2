const sections = {
  news: "新闻",
  interview: "访谈",
  editorial: "社论",
  event: "活动"
};

const statuses = {
  draft: "草稿",
  pending_review: "待审核",
  published: "已发布",
  returned: "退回",
  archived: "已归档",
  completed: "已完成"
};

const transitionActions = {
  draft: { target: "pending_review", label: "提交审核" },
  returned: { target: "pending_review", label: "重新提交" },
  pending_review: { target: "published", label: "审核发布" },
  published: { target: "archived", label: "归档" },
  archived: { target: "completed", label: "提交为已完成" }
};

const state = {
  articles: [],
  allArticles: [],
  selected: null,
  section: "news",
  status: "",
  search: "",
  queue: "",
  view: "edit"
};

const elements = {
  articleList: document.querySelector("#article-list"),
  editorEmpty: document.querySelector("#editor-empty"),
  editorContent: document.querySelector("#editor-content"),
  activeSectionLabel: document.querySelector("#active-section-label"),
  resultCount: document.querySelector("#result-count"),
  searchInput: document.querySelector("#search-input"),
  editorStatus: document.querySelector("#editor-status"),
  editorSection: document.querySelector("#editor-section"),
  editorEdition: document.querySelector("#editor-edition"),
  editorAuthor: document.querySelector("#editor-author"),
  titleInput: document.querySelector("#title-input"),
  summaryInput: document.querySelector("#summary-input"),
  bodyInput: document.querySelector("#body-input"),
  editorForm: document.querySelector("#editor-form"),
  previewPane: document.querySelector("#preview-pane"),
  previewKicker: document.querySelector("#preview-kicker"),
  previewTitle: document.querySelector("#preview-title"),
  previewSummary: document.querySelector("#preview-summary"),
  previewByline: document.querySelector("#preview-byline"),
  previewBody: document.querySelector("#preview-body"),
  wordCount: document.querySelector("#word-count"),
  saveButton: document.querySelector("#save-button"),
  transitionButton: document.querySelector("#transition-button"),
  saveState: document.querySelector("#save-state"),
  completedCount: document.querySelector("#completed-count"),
  progressLabel: document.querySelector("#issue-progress-label"),
  progressBar: document.querySelector("#issue-progress-bar"),
  toast: document.querySelector("#toast")
};

async function request(path, options = {}) {
  const response = await fetch(path, {
    headers: { "Content-Type": "application/json", ...options.headers },
    ...options
  });
  const payload = await response.json();
  if (!response.ok) {
    throw new Error(payload.error || "请求未完成");
  }
  return payload;
}

async function loadWorkspace() {
  try {
    const [all, completed] = await Promise.all([
      request("/api/articles"),
      request("/api/queues/completed")
    ]);
    state.allArticles = all.articles;
    elements.completedCount.textContent = completed.articles.length;
    updateCounts();
    await loadArticles();
  } catch (error) {
    showToast(error.message, true);
  }
}

async function loadArticles(preserveSelection = false) {
  try {
    let path;
    if (state.queue === "completed") {
      path = "/api/queues/completed";
    } else {
      const query = new URLSearchParams({ section: state.section });
      if (state.status) query.set("status", state.status);
      path = `/api/articles?${query}`;
    }
    const payload = await request(path);
    state.articles = payload.articles;
    renderArticleList();

    const selectedStillVisible = preserveSelection && state.selected && state.articles.some(article => article.id === state.selected.id);
    if (selectedStillVisible) {
      await selectArticle(state.selected.id);
    } else if (state.articles.length > 0) {
      await selectArticle(state.articles[0].id);
    } else {
      state.selected = null;
      renderEditor();
    }
  } catch (error) {
    showToast(error.message, true);
  }
}

function filteredArticles() {
  const query = state.search.trim().toLocaleLowerCase("zh-CN");
  if (!query) return state.articles;
  return state.articles.filter(article => `${article.title} ${article.author}`.toLocaleLowerCase("zh-CN").includes(query));
}

function renderArticleList() {
  const articles = filteredArticles();
  elements.resultCount.textContent = `${articles.length} 篇`;
  if (articles.length === 0) {
    elements.articleList.innerHTML = '<div class="empty-list">当前筛选下暂无稿件</div>';
    return;
  }

  elements.articleList.innerHTML = articles.map(article => `
    <button class="article-row ${state.selected?.id === article.id ? "active" : ""}" type="button" data-article-id="${escapeHTML(article.id)}">
      <div class="row-top">
        <span class="status-badge status-${article.status}">${statuses[article.status]}</span>
        <span class="result-count">${escapeHTML(article.edition)}</span>
      </div>
      <h3>${escapeHTML(article.title)}</h3>
      <p>${escapeHTML(article.summary)}</p>
      <div class="row-footer"><span>${escapeHTML(article.author)}</span><span>${escapeHTML(article.updatedLabel)}</span></div>
    </button>
  `).join("");
}

async function selectArticle(id) {
  try {
    state.selected = await request(`/api/articles/${encodeURIComponent(id)}`);
    renderArticleList();
    renderEditor();
  } catch (error) {
    showToast(error.message, true);
  }
}

function renderEditor() {
  const article = state.selected;
  elements.editorEmpty.classList.toggle("hidden", Boolean(article));
  elements.editorContent.classList.toggle("hidden", !article);
  if (!article) return;

  elements.editorStatus.className = `status-badge status-${article.status}`;
  elements.editorStatus.textContent = statuses[article.status];
  elements.editorSection.textContent = sections[article.section];
  elements.editorEdition.textContent = article.edition;
  elements.editorAuthor.textContent = `作者：${article.author}`;
  elements.titleInput.value = article.title;
  elements.summaryInput.value = article.summary;
  elements.bodyInput.value = article.body;
  elements.saveState.textContent = "内容已同步";

  const action = transitionActions[article.status];
  elements.transitionButton.hidden = !action;
  elements.transitionButton.disabled = !action;
  elements.transitionButton.textContent = action?.label || "流程已结束";
  updateWordCount();
  renderPreview();
  switchView(state.view);
}

function renderPreview() {
  if (!state.selected) return;
  elements.previewKicker.textContent = `${sections[state.selected.section]} · ${state.selected.edition}`;
  elements.previewTitle.textContent = elements.titleInput.value;
  elements.previewSummary.textContent = elements.summaryInput.value;
  elements.previewByline.textContent = `青岚校报记者 ${state.selected.author}`;
  elements.previewBody.textContent = elements.bodyInput.value;
}

function switchView(view) {
  state.view = view;
  document.querySelectorAll("[data-view]").forEach(button => button.classList.toggle("active", button.dataset.view === view));
  elements.editorForm.classList.toggle("hidden", view !== "edit");
  elements.previewPane.classList.toggle("hidden", view !== "preview");
  if (view === "preview") renderPreview();
}

async function saveArticle() {
  if (!state.selected) return;
  elements.saveButton.disabled = true;
  elements.saveState.textContent = "正在保存";
  try {
    const article = await request(`/api/articles/${encodeURIComponent(state.selected.id)}`, {
      method: "PATCH",
      body: JSON.stringify({
        title: elements.titleInput.value,
        summary: elements.summaryInput.value,
        body: elements.bodyInput.value
      })
    });
    state.selected = article;
    await refreshData();
    elements.saveState.textContent = "内容已同步";
    showToast("稿件已保存");
  } catch (error) {
    elements.saveState.textContent = "保存失败";
    showToast(error.message, true);
  } finally {
    elements.saveButton.disabled = false;
  }
}

async function transitionArticle() {
  if (!state.selected) return;
  const action = transitionActions[state.selected.status];
  if (!action) return;
  elements.transitionButton.disabled = true;
  try {
    await request(`/api/articles/${encodeURIComponent(state.selected.id)}/transitions`, {
      method: "POST",
      body: JSON.stringify({ status: action.target })
    });
    showToast(`稿件已更新为${statuses[action.target]}`);
    await refreshData();
  } catch (error) {
    showToast(error.message, true);
  } finally {
    elements.transitionButton.disabled = false;
  }
}

async function refreshData() {
  const selectedID = state.selected?.id;
  const [all, completed] = await Promise.all([
    request("/api/articles"),
    request("/api/queues/completed")
  ]);
  state.allArticles = all.articles;
  elements.completedCount.textContent = completed.articles.length;
  updateCounts();
  await loadArticles(true);
  if (selectedID && !state.selected) await selectArticle(selectedID);
}

function updateCounts() {
  Object.keys(sections).forEach(section => {
    const count = state.allArticles.filter(article => article.section === section).length;
    const target = document.querySelector(`[data-count="${section}"]`);
    if (target) target.textContent = count;
  });
  const currentIssue = state.allArticles.filter(article => article.edition === "第 128 期");
  const advanced = currentIssue.filter(article => ["pending_review", "published", "completed"].includes(article.status)).length;
  elements.progressLabel.textContent = `${advanced} / ${currentIssue.length}`;
  elements.progressBar.style.width = currentIssue.length ? `${Math.round((advanced / currentIssue.length) * 100)}%` : "0%";
}

function selectSection(section) {
  state.section = section;
  state.queue = "";
  state.status = "";
  state.search = "";
  elements.searchInput.value = "";
  elements.activeSectionLabel.textContent = sections[section];
  document.querySelectorAll("[data-section]").forEach(button => button.classList.toggle("active", button.dataset.section === section));
  document.querySelectorAll("[data-queue]").forEach(button => button.classList.remove("active"));
  document.querySelectorAll("[data-status]").forEach(button => button.classList.toggle("active", button.dataset.status === ""));
  document.querySelector("#status-tabs").classList.remove("hidden");
  loadArticles();
}

function selectCompletedQueue() {
  state.queue = "completed";
  state.search = "";
  elements.searchInput.value = "";
  elements.activeSectionLabel.textContent = "完成清单";
  document.querySelectorAll("[data-section]").forEach(button => button.classList.remove("active"));
  document.querySelectorAll("[data-queue]").forEach(button => button.classList.add("active"));
  document.querySelector("#status-tabs").classList.add("hidden");
  loadArticles();
}

function updateWordCount() {
  const count = Array.from(elements.bodyInput.value.replace(/\s/g, "")).length;
  elements.wordCount.textContent = `${count} 字`;
  elements.saveState.textContent = "有未保存修改";
}

function showToast(message, isError = false) {
  elements.toast.textContent = message;
  elements.toast.classList.toggle("error", isError);
  elements.toast.classList.add("show");
  window.clearTimeout(showToast.timer);
  showToast.timer = window.setTimeout(() => elements.toast.classList.remove("show"), 2400);
}

function escapeHTML(value) {
  return String(value)
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#039;");
}

document.querySelector("#section-list").addEventListener("click", event => {
  const button = event.target.closest("[data-section]");
  if (button) selectSection(button.dataset.section);
});

document.querySelector("[data-queue='completed']").addEventListener("click", selectCompletedQueue);

document.querySelector("#status-tabs").addEventListener("click", event => {
  const button = event.target.closest("[data-status]");
  if (!button) return;
  state.status = button.dataset.status;
  document.querySelectorAll("[data-status]").forEach(item => item.classList.toggle("active", item === button));
  loadArticles();
});

elements.articleList.addEventListener("click", event => {
  const row = event.target.closest("[data-article-id]");
  if (row) selectArticle(row.dataset.articleId);
});

elements.searchInput.addEventListener("input", event => {
  state.search = event.target.value;
  renderArticleList();
});

document.querySelectorAll("[data-view]").forEach(button => {
  button.addEventListener("click", () => switchView(button.dataset.view));
});

[elements.titleInput, elements.summaryInput, elements.bodyInput].forEach(input => {
  input.addEventListener("input", () => {
    updateWordCount();
    renderPreview();
  });
});

elements.saveButton.addEventListener("click", saveArticle);
elements.transitionButton.addEventListener("click", transitionArticle);
document.querySelector("#new-article").addEventListener("click", () => showToast("新稿件入口已预留"));

loadWorkspace();
