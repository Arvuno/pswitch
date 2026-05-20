const messages = {
  en: {
    "brand.title": "Operations Hub",
    "brand.subtitle": "A clean admin dashboard for usage, provider health, and runtime configuration.",
    "common.language": "Language",
    "common.updated": "Updated",
    "common.refresh": "Refresh",
    "common.inShort": "In",
    "common.outShort": "Out",
    "common.reqShort": "Req",
    "common.failShort": "Fail",
    "common.share": "Share",
    "auth.title": "Admin token required",
    "auth.subtitle": "Enter",
    "auth.suffix": "to continue.",
    "auth.unlock": "Unlock",
    "auth.placeholder": "Admin token",
    "tabs.overview": "Overview",
    "tabs.usage": "Usage",
    "tabs.providers": "Providers",
    "tabs.config": "Config",
    "topbar.kicker": "Monitoring first",
    "topbar.overview": "Service overview",
    "topbar.overviewDesc": "Focus on request volume, token consumption, and provider health.",
    "topbar.usage": "Usage windows",
    "topbar.usageDesc": "Compare 24h and 7d token activity without leaving the current page.",
    "topbar.providers": "Provider analytics",
    "topbar.providersDesc": "See which upstreams carry traffic, consume tokens, or produce errors.",
    "topbar.config": "Runtime config",
    "topbar.configDesc": "Keep runtime settings, routes, and providers in a compact editing workspace.",
    "overview.kicker": "Key metrics",
    "overview.heading": "What matters right now",
    "overview.totalRequests": "Total requests",
    "overview.totalTokens": "Total tokens",
    "overview.totalFailures": "Total failures",
    "overview.healthyProviders": "Healthy providers",
    "overview.tokens24h": "24h tokens",
    "overview.tokens7d": "7d tokens",
    "overview.snapshotKicker": "Runtime snapshot",
    "overview.snapshotHeading": "Current service state",
    "runtime.listen": "Listen",
    "runtime.mode": "Mode",
    "runtime.routes": "Routes",
    "runtime.providers": "Providers",
    "providers.spotlightKicker": "Provider spotlight",
    "providers.spotlightHeading": "Top providers by token usage",
    "providers.spotlightEmpty": "No provider data yet.",
    "providers.summaryKicker": "Provider summary",
    "providers.summaryHeading": "Provider estate at a glance",
    "providers.totalProviders": "Configured providers",
    "providers.healthRatio": "Healthy ratio",
    "providers.topProvider": "Top provider",
    "providers.requests24h": "24h requests",
    "models.kicker": "Model usage",
    "models.heading": "Usage by model",
    "models.empty": "No model usage data yet.",
    "models.tokens24h": "24h",
    "models.tokens7d": "7d",
    "usage.kicker": "Time windows",
    "usage.heading": "24h and 7d usage",
    "usage.trendKicker": "Trends",
    "usage.trendHeading": "Token activity",
    "usage.trendMeta": "Each bar shows total tokens per bucket.",
    "sections.hourly24h": "Hourly 24h",
    "sections.daily7d": "Daily 7d",
    "sections.providersKicker": "Provider analytics",
    "sections.providersTitle": "Calls, tokens, and failures",
    "sections.healthKicker": "Health",
    "sections.healthTitle": "Provider state",
    "table.provider": "Provider",
    "table.status": "Status",
    "table.requests": "Requests",
    "table.tokens": "Tokens",
    "table.failures": "Failures",
    "table.tokens24h": "24h Tokens",
    "table.tokens7d": "7d Tokens",
    "config.kicker": "Configuration",
    "config.title": "Runtime config",
    "config.save": "Save config",
    "config.listen": "Listen",
    "config.mode": "Mode",
    "config.modeGuideTitle": "Mode behavior",
    "config.modeRoundRobinLabel": "Round robin",
    "config.modeRoundRobinDesc": "Rotate requests across healthy providers to spread traffic evenly.",
    "config.modeSequentialLabel": "Sequential",
    "config.modeSequentialDesc": "Always try providers in list order and fail over only when the current one breaks.",
    "config.modeLeastFailuresLabel": "Least failures",
    "config.modeLeastFailuresDesc": "Prefer healthy providers that have accumulated the fewest upstream failures.",
    "config.failureThreshold": "Failure threshold",
    "config.cooldown": "Cooldown",
    "config.healthCheckInterval": "Health check interval",
    "config.healthCheckTimeout": "Health check timeout",
    "config.routes": "Routes",
    "config.providers": "Providers",
    "config.addRoute": "Add route",
    "config.addProvider": "Add provider",
    "config.generalTab": "General",
    "config.routesTab": "Routes",
    "config.providersTab": "Providers",
    "config.remove": "Remove",
    "config.enabled": "Enabled",
    "config.routePrefixPlaceholder": "/codex",
    "config.routeModelPlaceholder": "advertised model",
    "config.routeUpstreamPlaceholder": "upstream model",
    "config.providerNamePlaceholder": "provider name",
    "config.providerBaseURLPlaceholder": "https://provider.example/v1",
    "config.providerAPIKeyPlaceholder": "sk-...",
    "window.last24h": "Last 24 hours",
    "window.last7d": "Last 7 days",
    "window.requests": "Requests",
    "window.failures": "Failures",
    "window.inputTokens": "Input",
    "window.outputTokens": "Output",
    "window.totalTokens": "Total",
    "status.healthy": "healthy",
    "status.broken": "broken",
    "status.failures": "Consecutive failures",
    "status.nextProbe": "Next probe",
    "status.requests": "Requests",
    "status.tokens": "Tokens",
    "status.na": "n/a",
    "message.adminTokenRequired": "Admin token required.",
    "message.savedReloaded": "Config saved and hot reloaded.",
    "message.savedRestart": "Config saved. Some changes require a restart."
  },
  "zh-CN": {
    "brand.title": "运营控制台",
    "brand.subtitle": "以简洁、清晰为核心，集中展示用量、Provider 健康状态和运行配置。",
    "common.language": "语言",
    "common.updated": "更新时间",
    "common.refresh": "刷新",
    "common.inShort": "入",
    "common.outShort": "出",
    "common.reqShort": "请求",
    "common.failShort": "失败",
    "common.share": "占比",
    "auth.title": "需要管理员令牌",
    "auth.subtitle": "请输入",
    "auth.suffix": "后继续。",
    "auth.unlock": "解锁",
    "auth.placeholder": "管理员令牌",
    "tabs.overview": "总览",
    "tabs.usage": "使用情况",
    "tabs.providers": "Providers",
    "tabs.config": "配置",
    "topbar.kicker": "监控优先",
    "topbar.overview": "服务总览",
    "topbar.overviewDesc": "优先关注请求量、Token 消耗和 Provider 健康状态。",
    "topbar.usage": "使用窗口",
    "topbar.usageDesc": "在同一工作区内对比 24h 与 7d 的 Token 使用情况。",
    "topbar.providers": "Provider 分析",
    "topbar.providersDesc": "快速查看哪些上游承担流量、消耗 Token 或产生错误。",
    "topbar.config": "运行配置",
    "topbar.configDesc": "将运行设置、路由和 Provider 编辑集中到更紧凑的工作区。",
    "overview.kicker": "关键指标",
    "overview.heading": "当前最值得关注的数据",
    "overview.totalRequests": "累计调用",
    "overview.totalTokens": "累计 Token",
    "overview.totalFailures": "累计失败",
    "overview.healthyProviders": "健康 Provider",
    "overview.tokens24h": "24h Token",
    "overview.tokens7d": "7d Token",
    "overview.snapshotKicker": "运行快照",
    "overview.snapshotHeading": "当前服务状态",
    "runtime.listen": "监听地址",
    "runtime.mode": "模式",
    "runtime.routes": "路由数",
    "runtime.providers": "Provider 数",
    "providers.spotlightKicker": "Provider 聚焦",
    "providers.spotlightHeading": "按 Token 使用排序的重点 Provider",
    "providers.spotlightEmpty": "暂无 Provider 数据。",
    "providers.summaryKicker": "Provider 概览",
    "providers.summaryHeading": "快速查看整体 Provider 状态",
    "providers.totalProviders": "已配置 Provider",
    "providers.healthRatio": "健康占比",
    "providers.topProvider": "Token 最高 Provider",
    "providers.requests24h": "24h 请求",
    "models.kicker": "模型使用量",
    "models.heading": "按模型统计使用量",
    "models.empty": "暂无模型使用数据。",
    "models.tokens24h": "24h",
    "models.tokens7d": "7d",
    "usage.kicker": "时间窗口",
    "usage.heading": "24h 与 7d 使用情况",
    "usage.trendKicker": "趋势",
    "usage.trendHeading": "Token 活动",
    "usage.trendMeta": "每个柱条代表一个时间桶内的 Token 总量。",
    "sections.hourly24h": "按小时 24h",
    "sections.daily7d": "按天 7d",
    "sections.providersKicker": "Provider 统计",
    "sections.providersTitle": "调用、Token 与失败",
    "sections.healthKicker": "健康状态",
    "sections.healthTitle": "Provider 状态",
    "table.provider": "Provider",
    "table.status": "状态",
    "table.requests": "调用次数",
    "table.tokens": "Token",
    "table.failures": "失败次数",
    "table.tokens24h": "24h Token",
    "table.tokens7d": "7d Token",
    "config.kicker": "配置",
    "config.title": "运行配置",
    "config.save": "保存配置",
    "config.listen": "监听地址",
    "config.mode": "模式",
    "config.modeGuideTitle": "模式说明",
    "config.modeRoundRobinLabel": "轮询",
    "config.modeRoundRobinDesc": "在健康的 Provider 之间轮流分配请求，尽量均摊流量。",
    "config.modeSequentialLabel": "顺序优先",
    "config.modeSequentialDesc": "始终按 Provider 列表顺序尝试，只有当前上游异常时才切到下一个。",
    "config.modeLeastFailuresLabel": "最少错误优先",
    "config.modeLeastFailuresDesc": "优先选择累计上游失败次数更少的健康 Provider。",
    "config.failureThreshold": "失败阈值",
    "config.cooldown": "冷却时间",
    "config.healthCheckInterval": "健康检查间隔",
    "config.healthCheckTimeout": "健康检查超时",
    "config.routes": "路由",
    "config.providers": "Providers",
    "config.addRoute": "新增路由",
    "config.addProvider": "新增 Provider",
    "config.generalTab": "基础设置",
    "config.routesTab": "路由",
    "config.providersTab": "Providers",
    "config.remove": "移除",
    "config.enabled": "启用",
    "config.routePrefixPlaceholder": "/codex",
    "config.routeModelPlaceholder": "对外展示模型",
    "config.routeUpstreamPlaceholder": "上游模型",
    "config.providerNamePlaceholder": "provider 名称",
    "config.providerBaseURLPlaceholder": "https://provider.example/v1",
    "config.providerAPIKeyPlaceholder": "sk-...",
    "window.last24h": "最近 24 小时",
    "window.last7d": "最近 7 天",
    "window.requests": "请求",
    "window.failures": "失败",
    "window.inputTokens": "输入",
    "window.outputTokens": "输出",
    "window.totalTokens": "总计",
    "status.healthy": "健康",
    "status.broken": "异常",
    "status.failures": "连续失败",
    "status.nextProbe": "下次探测",
    "status.requests": "请求",
    "status.tokens": "Token",
    "status.na": "无",
    "message.adminTokenRequired": "需要管理员令牌。",
    "message.savedReloaded": "配置已保存并热重载。",
    "message.savedRestart": "配置已保存，但部分变更需要重启生效。"
  }
};

const state = {
  token: localStorage.getItem("pswitch_admin_token") || "",
  locale: localStorage.getItem("pswitch_admin_locale") || inferLocale(),
  activeTab: "overview",
  activeConfigTab: "general",
  isEditingConfig: false,
  meta: null,
  current: null,
  stats: null,
  view: {
    kpiCards: [],
    runtimeRows: [],
    overviewWindowCards: [],
    usagePanels: [],
    providerSummaryCards: [],
    spotlightCards: [],
    modelCards: [],
    overviewBars: [],
    hourlyBars: [],
    dailyBars: [],
    providerRows: new Map(),
    providerStatusCards: new Map()
  }
};

const authCard = document.getElementById("auth-card");
const authForm = document.getElementById("auth-form");
const tokenInput = document.getElementById("token-input");
const messageEl = document.getElementById("message");
const formEl = document.getElementById("config-form");
const routesEl = document.getElementById("routes");
const providerEditorEl = document.getElementById("provider-editor");
const routeTemplate = document.getElementById("route-row-template");
const providerTemplate = document.getElementById("provider-row-template");
const languageSwitch = document.getElementById("language-switch");
const lastUpdatedEl = document.getElementById("last-updated");
const panelTitleEl = document.getElementById("panel-title");
const panelDescriptionEl = document.getElementById("panel-description");
const topbarHealthyEl = document.getElementById("topbar-healthy");
const topbarRequestsEl = document.getElementById("topbar-requests");
const topbarTokensEl = document.getElementById("topbar-tokens");
const kpiGridEl = document.getElementById("kpi-grid");
const runtimeSummaryEl = document.getElementById("runtime-summary");
const overviewWindowGridEl = document.getElementById("overview-window-grid");
const providerSpotlightEl = document.getElementById("provider-spotlight");
const modelUsageGridEl = document.getElementById("model-usage-grid");
const overviewHourlyChartEl = document.getElementById("overview-hourly-chart");
const usagePanelsEl = document.getElementById("usage-panels");
const hourlyChartEl = document.getElementById("hourly-chart");
const dailyChartEl = document.getElementById("daily-chart");
const providerSummaryGridEl = document.getElementById("provider-summary-grid");
const providerMetricsBodyEl = document.getElementById("provider-metrics-body");
const providerStatusGridEl = document.getElementById("provider-status-grid");
const refreshButtonEl = document.getElementById("refresh-button");

async function boot() {
  languageSwitch.value = state.locale;
  applyTranslations();
  initializeStaticView();
  updatePanelCopy();
  bindEvents();

  state.meta = await request("./api/meta");
  if (state.meta.token_required && !state.token) {
    authCard.classList.remove("hidden");
    return;
  }

  authCard.classList.add("hidden");
  await loadAll();
  window.setInterval(() => {
    if (!authCard.classList.contains("hidden")) {
      return;
    }
    void loadAll(true);
  }, 15000);
}

function bindEvents() {
  document.querySelectorAll(".nav-button").forEach((button) => {
    button.addEventListener("click", () => {
      state.activeTab = button.dataset.tab;
      document.querySelectorAll(".nav-button").forEach((node) => node.classList.remove("active"));
      document.querySelectorAll(".screen").forEach((node) => node.classList.remove("active"));
      button.classList.add("active");
      document.getElementById(`panel-${state.activeTab}`).classList.add("active");
      state.isEditingConfig = state.activeTab === "config" && state.isEditingConfig;
      updatePanelCopy();
    });
  });

  document.querySelectorAll(".config-tab-button").forEach((button) => {
    button.addEventListener("click", () => {
      state.activeConfigTab = button.dataset.configTab;
      document.querySelectorAll(".config-tab-button").forEach((node) => node.classList.remove("active"));
      document.querySelectorAll(".config-panel").forEach((node) => node.classList.remove("active"));
      button.classList.add("active");
      document.getElementById(`config-panel-${state.activeConfigTab}`).classList.add("active");
    });
  });

  document.getElementById("add-route").addEventListener("click", () => {
    routesEl.appendChild(buildRouteRow());
    state.isEditingConfig = true;
  });

  document.getElementById("add-provider").addEventListener("click", () => {
    providerEditorEl.appendChild(buildProviderRow());
    state.isEditingConfig = true;
  });

  authForm.addEventListener("submit", async (event) => {
    event.preventDefault();
    state.token = tokenInput.value.trim();
    localStorage.setItem("pswitch_admin_token", state.token);
    authCard.classList.add("hidden");
    await loadAll();
  });

  document.getElementById("save-button").addEventListener("click", async () => {
    const payload = collectConfig();
    try {
      const result = await request("./api/config", {
        method: "PUT",
        body: JSON.stringify(payload)
      });
      state.current = result.state;
      state.isEditingConfig = false;
      await loadAll(true);
      const notices = result.messages || [];
      if (result.requires_restart) {
        notices.unshift(t("message.savedRestart"));
      } else {
        notices.unshift(t("message.savedReloaded"));
      }
      showMessage("success", notices.join(" "));
    } catch (error) {
      showMessage("error", error.message);
    }
  });

  languageSwitch.addEventListener("change", () => {
    state.locale = languageSwitch.value;
    localStorage.setItem("pswitch_admin_locale", state.locale);
    applyTranslations();
    if (state.current && state.stats) {
      renderState({ preserveConfig: state.isEditingConfig });
    } else {
      updatePanelCopy();
      updateEditorRowTranslations();
    }
  });

  formEl.addEventListener("input", () => {
    if (state.activeTab === "config") {
      state.isEditingConfig = true;
    }
  });

  refreshButtonEl.addEventListener("click", async () => {
    await loadAll(true);
  });
}

function initializeStaticView() {
  initializeKpiCards();
  initializeRuntimeRows();
  initializeOverviewWindows();
  initializeUsagePanels();
  initializeProviderSummaryCards();
}

function initializeKpiCards() {
  const labels = [
    "overview.totalRequests",
    "overview.totalTokens",
    "overview.tokens24h",
    "overview.tokens7d",
    "overview.totalFailures",
    "overview.healthyProviders"
  ];

  kpiGridEl.innerHTML = "";
  state.view.kpiCards = labels.map((key) => {
    const card = document.createElement("div");
    const label = document.createElement("span");
    const value = document.createElement("strong");
    const note = document.createElement("p");
    card.className = "kpi-card";
    label.className = "kpi-label";
    value.className = "kpi-value";
    note.className = "kpi-note";
    label.dataset.i18nDynamic = key;
    card.append(label, value, note);
    kpiGridEl.appendChild(card);
    return { key, label, value, note };
  });
}

function initializeRuntimeRows() {
  const keys = ["runtime.listen", "runtime.mode", "runtime.routes", "runtime.providers"];
  runtimeSummaryEl.innerHTML = "";
  state.view.runtimeRows = keys.map((key) => {
    const row = document.createElement("div");
    const label = document.createElement("span");
    const value = document.createElement("strong");
    row.className = "runtime-row";
    label.dataset.i18nDynamic = key;
    row.append(label, value);
    runtimeSummaryEl.appendChild(row);
    return { key, label, value };
  });
}

function initializeOverviewWindows() {
  const windows = ["window.last24h", "window.last7d"];
  overviewWindowGridEl.innerHTML = "";
  state.view.overviewWindowCards = windows.map((key) => {
    const card = document.createElement("div");
    const title = document.createElement("div");
    const value = document.createElement("strong");
    const detail = document.createElement("p");
    card.className = "window-brief-card";
    title.className = "window-title";
    value.className = "window-value";
    detail.className = "window-detail";
    title.dataset.i18nDynamic = key;
    card.append(title, value, detail);
    overviewWindowGridEl.appendChild(card);
    return { key, title, value, detail };
  });
}

function initializeUsagePanels() {
  const windows = ["window.last24h", "window.last7d"];
  usagePanelsEl.innerHTML = "";
  state.view.usagePanels = windows.map((key) => {
    const card = document.createElement("div");
    const title = document.createElement("div");
    const metrics = document.createElement("div");
    card.className = "window-card";
    title.className = "window-title";
    metrics.className = "window-metrics";
    title.dataset.i18nDynamic = key;

    const metricKeys = [
      "window.requests",
      "window.failures",
      "window.totalTokens",
      "window.inputTokens",
      "window.outputTokens"
    ].map((metricKey) => {
      const item = document.createElement("div");
      const label = document.createElement("span");
      const value = document.createElement("strong");
      label.dataset.i18nDynamic = metricKey;
      item.append(label, value);
      metrics.appendChild(item);
      return { key: metricKey, label, value };
    });

    card.append(title, metrics);
    usagePanelsEl.appendChild(card);
    return { key, title, metricKeys };
  });
}

function initializeProviderSummaryCards() {
  const items = [
    "providers.totalProviders",
    "providers.healthRatio",
    "providers.topProvider",
    "providers.requests24h"
  ];

  providerSummaryGridEl.innerHTML = "";
  state.view.providerSummaryCards = items.map((key) => {
    const card = document.createElement("div");
    const label = document.createElement("span");
    const value = document.createElement("strong");
    const note = document.createElement("p");
    card.className = "summary-card";
    label.className = "summary-label";
    value.className = "summary-value";
    note.className = "summary-note";
    label.dataset.i18nDynamic = key;
    card.append(label, value, note);
    providerSummaryGridEl.appendChild(card);
    return { key, label, value, note };
  });
}

async function loadAll(silent = false) {
  try {
    const [current, stats] = await Promise.all([
      request("./api/state"),
      request("./api/stats")
    ]);
    state.current = current;
    state.stats = stats;
    renderState({ preserveConfig: state.isEditingConfig });
    if (!silent) {
      showMessage("", "");
    }
  } catch (error) {
    if (error.status === 401) {
      authCard.classList.remove("hidden");
      showMessage("error", t("message.adminTokenRequired"));
      return;
    }
    showMessage("error", error.message);
  }
}

function renderState(options = {}) {
  if (!options.preserveConfig) {
    renderConfig(state.current.config);
  } else {
    updateEditorRowTranslations();
  }

  updateDynamicTranslations();
  updatePanelCopy();
  renderTopbar(state.stats.overview);
  renderOverview(state.current.config, state.stats);
  renderUsage(state.stats);
  renderProviders(state.stats);
  lastUpdatedEl.textContent = formatDateTime(state.stats.server_time || state.current.server_time);
}

function renderTopbar(overview) {
  topbarHealthyEl.textContent = `${formatNumber(overview.healthy_providers)} / ${formatNumber(overview.providers_count)}`;
  topbarRequestsEl.textContent = formatNumber(overview.total_requests);
  topbarTokensEl.textContent = formatTokenCount(overview.total_tokens);
  topbarTokensEl.title = formatNumber(overview.total_tokens);
}

function renderOverview(config, stats) {
  updateKpis(stats.overview, stats.windows);
  updateRuntimeSummary(config);
  updateOverviewWindows(stats.windows);
  updateProviderSpotlight(stats.providers || []);
  updateModelUsage(stats.models || []);
  updateBars(overviewHourlyChartEl, state.view.overviewBars, (stats.series.hourly_24h || []).slice(-12));
}

function renderUsage(stats) {
  updateUsagePanels(stats.windows);
  updateBars(hourlyChartEl, state.view.hourlyBars, stats.series.hourly_24h || []);
  updateBars(dailyChartEl, state.view.dailyBars, stats.series.daily_7d || []);
}

function renderProviders(stats) {
  const providers = stats.providers || [];
  updateProviderSummary(stats.overview, stats.windows, providers);
  updateProviderTable(providers);
  updateProviderStatus(providers);
}

function renderConfig(cfg) {
  formEl.elements.listen.value = cfg.listen;
  formEl.elements.mode.value = cfg.mode;
  formEl.elements.failure_threshold.value = cfg.failure_threshold;
  formEl.elements.cooldown.value = cfg.cooldown;
  formEl.elements.health_check_interval.value = cfg.health_check_interval;
  formEl.elements.health_check_timeout.value = cfg.health_check_timeout;

  routesEl.innerHTML = "";
  cfg.routes.forEach((route) => routesEl.appendChild(buildRouteRow(route)));

  providerEditorEl.innerHTML = "";
  cfg.providers.forEach((provider) => providerEditorEl.appendChild(buildProviderRow(provider)));
}

function updateKpis(overview, windows) {
  const failureRate = overview.total_requests > 0
    ? `${Math.round((overview.total_failures / overview.total_requests) * 100)}%`
    : "0%";

  const values = new Map([
    [
      "overview.totalRequests",
      {
        value: formatNumber(overview.total_requests),
        note: joinMetrics([
          [t("window.last24h"), windows.last_24h.request_count],
          [t("window.last7d"), windows.last_7d.request_count]
        ])
      }
    ],
    [
      "overview.totalTokens",
      {
        value: formatTokenCount(overview.total_tokens),
        note: joinMetrics([
          [t("common.inShort"), overview.total_input_tokens],
          [t("common.outShort"), overview.total_output_tokens]
        ])
      }
    ],
    [
      "overview.tokens24h",
      {
        value: formatTokenCount(windows.last_24h.total_tokens),
        note: joinMetrics([
          [t("common.inShort"), windows.last_24h.input_tokens],
          [t("common.outShort"), windows.last_24h.output_tokens]
        ])
      }
    ],
    [
      "overview.tokens7d",
      {
        value: formatTokenCount(windows.last_7d.total_tokens),
        note: joinMetrics([
          [t("common.reqShort"), windows.last_7d.request_count],
          [t("common.failShort"), windows.last_7d.failure_count]
        ])
      }
    ],
    [
      "overview.totalFailures",
      {
        value: formatNumber(overview.total_failures),
        note: `${t("common.failShort")} ${failureRate}`
      }
    ],
    [
      "overview.healthyProviders",
      {
        value: `${formatNumber(overview.healthy_providers)} / ${formatNumber(overview.providers_count)}`,
        note: `${formatNumber(overview.providers_count)} ${t("runtime.providers").toLowerCase()}`
      }
    ]
  ]);

  state.view.kpiCards.forEach((card) => {
    const entry = values.get(card.key);
    card.label.textContent = t(card.key);
    card.value.textContent = entry ? entry.value : "-";
    card.value.title = entry && isTokenMetricKey(card.key) ? formatNumber(resolveTokenMetricValue(card.key, overview, windows)) : "";
    card.note.textContent = entry ? entry.note : "";
  });
}

function updateRuntimeSummary(config) {
  const values = new Map([
    ["runtime.listen", config.listen],
    ["runtime.mode", formatModeLabel(config.mode)],
    ["runtime.routes", String(config.routes.length)],
    ["runtime.providers", String(config.providers.length)]
  ]);

  state.view.runtimeRows.forEach((row) => {
    row.label.textContent = t(row.key);
    row.value.textContent = values.get(row.key) || "-";
  });
}

function updateOverviewWindows(windows) {
  const snapshots = [windows.last_24h, windows.last_7d];
  state.view.overviewWindowCards.forEach((card, index) => {
    const snapshot = snapshots[index];
    card.title.textContent = t(card.key);
    card.value.textContent = formatTokenCount(snapshot.total_tokens);
    card.value.title = formatNumber(snapshot.total_tokens);
    card.detail.textContent = joinMetrics([
      [t("common.reqShort"), snapshot.request_count],
      [t("common.failShort"), snapshot.failure_count]
    ]);
  });
}

function updateProviderSpotlight(providers) {
  const ordered = [...providers].sort((a, b) => (b.total_tokens || 0) - (a.total_tokens || 0));
  const totalTokens = ordered.reduce((sum, item) => sum + Number(item.total_tokens || 0), 0);
  const topProviders = ordered.slice(0, 4);

  syncCollection(providerSpotlightEl, state.view.spotlightCards, Math.max(topProviders.length, 1), createSpotlightCard);

  if (topProviders.length === 0) {
    const card = state.view.spotlightCards[0];
    card.root.classList.add("empty");
    card.title.textContent = t("providers.spotlightEmpty");
    card.badge.textContent = "";
    card.meta.textContent = "";
    card.share.textContent = "";
    card.tokens.textContent = "";
    card.bar.style.width = "0%";
    return;
  }

  topProviders.forEach((provider, index) => {
    const share = totalTokens > 0 ? Math.round((provider.total_tokens / totalTokens) * 100) : 0;
    const card = state.view.spotlightCards[index];
    card.root.classList.remove("empty");
    card.title.textContent = provider.name;
    card.badge.textContent = provider.healthy ? t("status.healthy") : t("status.broken");
    card.badge.className = `badge ${provider.healthy ? "healthy" : "broken"}`;
    card.meta.textContent = joinMetrics([
      [t("common.reqShort"), provider.request_count],
      [t("common.failShort"), provider.failure_count]
    ]);
    card.share.textContent = `${t("common.share")} ${share}%`;
    card.tokens.textContent = formatTokenCount(provider.total_tokens);
    card.tokens.title = formatNumber(provider.total_tokens);
    card.bar.style.width = `${share}%`;
  });
}

function createSpotlightCard() {
  const root = document.createElement("article");
  const header = document.createElement("div");
  const title = document.createElement("strong");
  const badge = document.createElement("span");
  const meta = document.createElement("p");
  const progress = document.createElement("div");
  const bar = document.createElement("span");
  const footer = document.createElement("div");
  const share = document.createElement("span");
  const tokens = document.createElement("strong");

  root.className = "spotlight-card";
  header.className = "spotlight-header";
  badge.className = "badge healthy";
  meta.className = "spotlight-meta";
  progress.className = "spotlight-progress";
  footer.className = "spotlight-footer";

  progress.appendChild(bar);
  header.append(title, badge);
  footer.append(share, tokens);
  root.append(header, meta, progress, footer);

  return { root, title, badge, meta, bar, share, tokens };
}

function updateModelUsage(models) {
  const ordered = [...models].sort((a, b) => (b.total_tokens || 0) - (a.total_tokens || 0)).slice(0, 6);
  syncCollection(modelUsageGridEl, state.view.modelCards, Math.max(ordered.length, 1), createModelCard);

  if (ordered.length === 0) {
    const card = state.view.modelCards[0];
    card.root.classList.add("empty");
    card.name.textContent = t("models.empty");
    card.share.textContent = "";
    card.total.textContent = "";
    card.meta.textContent = "";
    return;
  }

  const totalTokens = ordered.reduce((sum, item) => sum + Number(item.total_tokens || 0), 0);
  ordered.forEach((model, index) => {
    const card = state.view.modelCards[index];
    const share = totalTokens > 0 ? Math.round((model.total_tokens / totalTokens) * 100) : 0;
    card.root.classList.remove("empty");
    card.name.textContent = model.name;
    card.share.textContent = `${t("common.share")} ${share}%`;
    card.total.textContent = formatTokenCount(model.total_tokens);
    card.total.title = formatNumber(model.total_tokens);
    card.meta.textContent = [
      `${t("window.requests")} ${formatCompact(model.request_count)}`,
      `${t("window.failures")} ${formatCompact(model.failure_count)}`,
      `${t("models.tokens24h")} ${formatCompact(model.last_24h.total_tokens)}`,
      `${t("models.tokens7d")} ${formatCompact(model.last_7d.total_tokens)}`
    ].join(" · ");
  });
}

function createModelCard() {
  const root = document.createElement("article");
  const header = document.createElement("div");
  const name = document.createElement("strong");
  const share = document.createElement("span");
  const total = document.createElement("strong");
  const meta = document.createElement("p");

  root.className = "model-card";
  header.className = "model-card-header";
  total.className = "model-card-total";
  meta.className = "model-card-meta";
  header.append(name, share);
  root.append(header, total, meta);

  return { root, name, share, total, meta };
}

function updateUsagePanels(windows) {
  const snapshots = [windows.last_24h, windows.last_7d];
  state.view.usagePanels.forEach((panel, index) => {
    const snapshot = snapshots[index];
    panel.title.textContent = t(panel.key);
    const values = new Map([
      ["window.requests", snapshot.request_count],
      ["window.failures", snapshot.failure_count],
      ["window.totalTokens", snapshot.total_tokens],
      ["window.inputTokens", snapshot.input_tokens],
      ["window.outputTokens", snapshot.output_tokens]
    ]);
    panel.metricKeys.forEach((metric) => {
      metric.label.textContent = t(metric.key);
      const rawValue = values.get(metric.key);
      metric.value.textContent = metric.key === "window.totalTokens" || metric.key === "window.inputTokens" || metric.key === "window.outputTokens"
        ? formatTokenCount(rawValue)
        : formatNumber(rawValue);
      metric.value.title = metric.key === "window.totalTokens" || metric.key === "window.inputTokens" || metric.key === "window.outputTokens"
        ? formatNumber(rawValue)
        : "";
    });
  });
}

function updateProviderSummary(overview, windows, providers) {
  const topProvider = [...providers].sort((a, b) => (b.total_tokens || 0) - (a.total_tokens || 0))[0];
  const summary = new Map([
    [
      "providers.totalProviders",
      {
        value: formatNumber(overview.providers_count),
        note: `${formatNumber(overview.healthy_providers)} ${t("status.healthy")}`
      }
    ],
    [
      "providers.healthRatio",
      {
        value: overview.providers_count > 0
          ? `${Math.round((overview.healthy_providers / overview.providers_count) * 100)}%`
          : "0%",
        note: `${formatNumber(overview.healthy_providers)} / ${formatNumber(overview.providers_count)}`
      }
    ],
    [
      "providers.topProvider",
      {
        value: topProvider ? topProvider.name : "-",
        note: topProvider ? `${formatTokenCount(topProvider.total_tokens)} token` : "-"
      }
    ],
    [
      "providers.requests24h",
      {
        value: formatNumber(windows.last_24h.request_count),
        note: `${t("window.last7d")} ${formatNumber(windows.last_7d.request_count)}`
      }
    ]
  ]);

  state.view.providerSummaryCards.forEach((card) => {
    const entry = summary.get(card.key);
    card.label.textContent = t(card.key);
    card.value.textContent = entry ? entry.value : "-";
    card.note.textContent = entry ? entry.note : "";
  });
}

function updateBars(root, cache, points) {
  syncCollection(root, cache, points.length, createBarItem);
  const max = Math.max(1, ...points.map((item) => item.total_tokens || 0));
  points.forEach((point, index) => {
    const item = cache[index];
    const height = Math.max(8, Math.round(((point.total_tokens || 0) / max) * 100));
    item.bar.style.height = `${height}%`;
    item.value.textContent = formatCompact(point.total_tokens || 0);
    item.label.textContent = point.label;
    item.root.title = `${point.label}: ${formatNumber(point.total_tokens || 0)}`;
  });
}

function createBarItem() {
  const root = document.createElement("div");
  const bar = document.createElement("div");
  const value = document.createElement("strong");
  const label = document.createElement("span");
  root.className = "bar-item";
  bar.className = "bar";
  root.append(bar, value, label);
  return { root, bar, value, label };
}

function updateProviderTable(providers) {
  const seen = new Set();
  providers.forEach((provider) => {
    let row = state.view.providerRows.get(provider.name);
    if (!row) {
      row = createProviderTableRow();
      state.view.providerRows.set(provider.name, row);
      providerMetricsBodyEl.appendChild(row.root);
    }

    seen.add(provider.name);
    row.name.textContent = provider.name;
    row.url.textContent = provider.base_url;
    row.badge.textContent = provider.healthy ? t("status.healthy") : t("status.broken");
    row.badge.className = `badge ${provider.healthy ? "healthy" : "broken"}`;
    row.requests.textContent = formatNumber(provider.request_count);
    row.tokensTotal.textContent = formatTokenCount(provider.total_tokens);
    row.tokensTotal.title = formatNumber(provider.total_tokens);
    row.tokensDetail.textContent = joinMetrics([
      [t("common.inShort"), provider.input_tokens],
      [t("common.outShort"), provider.output_tokens]
    ]);
    row.failures.textContent = formatNumber(provider.failure_count);
    row.tokens24h.textContent = formatTokenCount(provider.last_24h.total_tokens);
    row.tokens24h.title = formatNumber(provider.last_24h.total_tokens);
    row.tokens7d.textContent = formatTokenCount(provider.last_7d.total_tokens);
    row.tokens7d.title = formatNumber(provider.last_7d.total_tokens);
  });
  pruneMapEntries(state.view.providerRows, providerMetricsBodyEl, seen);
}

function createProviderTableRow() {
  const row = document.createElement("tr");
  const nameCell = document.createElement("td");
  const nameWrap = document.createElement("div");
  const name = document.createElement("strong");
  const url = document.createElement("span");
  const statusCell = document.createElement("td");
  const badge = document.createElement("span");
  const requests = document.createElement("td");
  const tokensCell = document.createElement("td");
  const tokensWrap = document.createElement("div");
  const tokensTotal = document.createElement("strong");
  const tokensDetail = document.createElement("span");
  const failures = document.createElement("td");
  const tokens24h = document.createElement("td");
  const tokens7d = document.createElement("td");

  nameWrap.className = "table-cell-stack";
  tokensWrap.className = "table-cell-stack";
  nameWrap.append(name, url);
  tokensWrap.append(tokensTotal, tokensDetail);
  nameCell.appendChild(nameWrap);
  tokensCell.appendChild(tokensWrap);
  statusCell.appendChild(badge);
  row.append(nameCell, statusCell, requests, tokensCell, failures, tokens24h, tokens7d);

  return { root: row, name, url, badge, requests, tokensTotal, tokensDetail, failures, tokens24h, tokens7d };
}

function updateProviderStatus(providers) {
  const seen = new Set();
  providers.forEach((provider) => {
    let card = state.view.providerStatusCards.get(provider.name);
    if (!card) {
      card = createProviderStatusCard();
      state.view.providerStatusCards.set(provider.name, card);
      providerStatusGridEl.appendChild(card.root);
    }

    seen.add(provider.name);
    card.title.textContent = provider.name;
    card.badge.textContent = provider.healthy ? t("status.healthy") : t("status.broken");
    card.badge.className = `badge ${provider.healthy ? "healthy" : "broken"}`;
    card.baseURL.textContent = provider.base_url;
    card.requestLabel.textContent = t("status.requests");
    card.requestValue.textContent = formatNumber(provider.request_count);
    card.tokenLabel.textContent = t("status.tokens");
    card.tokenValue.textContent = formatTokenCount(provider.total_tokens);
    card.tokenValue.title = formatNumber(provider.total_tokens);
    card.failureLabel.textContent = t("status.failures");
    card.failureValue.textContent = formatNumber(provider.consecutive_failures);
    card.probeLabel.textContent = t("status.nextProbe");
    card.probeValue.textContent = provider.next_probe_at && provider.next_probe_at !== "0001-01-01T00:00:00Z"
      ? formatDateTime(provider.next_probe_at)
      : t("status.na");
  });
  pruneMapEntries(state.view.providerStatusCards, providerStatusGridEl, seen);
}

function createProviderStatusCard() {
  const root = document.createElement("article");
  const header = document.createElement("div");
  const title = document.createElement("strong");
  const badge = document.createElement("span");
  const baseURL = document.createElement("p");
  const metrics = document.createElement("div");

  const request = createStatusMetric();
  const token = createStatusMetric();
  const failure = createStatusMetric();
  const probe = createStatusMetric();

  root.className = "status-card";
  header.className = "status-card-header";
  metrics.className = "status-metrics";
  header.append(title, badge);
  metrics.append(request.root, token.root, failure.root, probe.root);
  root.append(header, baseURL, metrics);

  return {
    root,
    title,
    badge,
    baseURL,
    requestLabel: request.label,
    requestValue: request.value,
    tokenLabel: token.label,
    tokenValue: token.value,
    failureLabel: failure.label,
    failureValue: failure.value,
    probeLabel: probe.label,
    probeValue: probe.value
  };
}

function createStatusMetric() {
  const root = document.createElement("div");
  const label = document.createElement("span");
  const value = document.createElement("strong");
  root.className = "status-metric";
  root.append(label, value);
  return { root, label, value };
}

function buildRouteRow(route = { prefix: "", type: "openai", model: "", upstream_model: "", enabled: true }) {
  const node = routeTemplate.content.firstElementChild.cloneNode(true);
  node.querySelector('[data-field="prefix"]').value = route.prefix || "";
  node.querySelector('[data-field="type"]').value = route.type || "openai";
  node.querySelector('[data-field="model"]').value = route.model || "";
  node.querySelector('[data-field="upstream_model"]').value = route.upstream_model || "";
  node.querySelector('[data-field="enabled"]').checked = route.enabled !== false;
  node.querySelector('[data-action="remove"]').addEventListener("click", () => {
    node.remove();
    state.isEditingConfig = true;
  });
  hydrateRouteRow(node);
  return node;
}

function buildProviderRow(provider = { name: "", base_url: "", api_key: "", enabled: true }) {
  const node = providerTemplate.content.firstElementChild.cloneNode(true);
  node.querySelector('[data-field="name"]').value = provider.name || "";
  node.querySelector('[data-field="base_url"]').value = provider.base_url || "";
  node.querySelector('[data-field="api_key"]').value = provider.api_key || "";
  node.querySelector('[data-field="enabled"]').checked = provider.enabled !== false;
  node.querySelector('[data-action="remove"]').addEventListener("click", () => {
    node.remove();
    state.isEditingConfig = true;
  });
  hydrateProviderRow(node);
  return node;
}

function hydrateRouteRow(node) {
  node.querySelector('[data-role="enabled-text"]').textContent = t("config.enabled");
  node.querySelector('[data-action="remove"]').textContent = t("config.remove");
  node.querySelector('[data-field="prefix"]').placeholder = t("config.routePrefixPlaceholder");
  node.querySelector('[data-field="model"]').placeholder = t("config.routeModelPlaceholder");
  node.querySelector('[data-field="upstream_model"]').placeholder = t("config.routeUpstreamPlaceholder");
}

function hydrateProviderRow(node) {
  node.querySelector('[data-role="enabled-text"]').textContent = t("config.enabled");
  node.querySelector('[data-action="remove"]').textContent = t("config.remove");
  node.querySelector('[data-field="name"]').placeholder = t("config.providerNamePlaceholder");
  node.querySelector('[data-field="base_url"]').placeholder = t("config.providerBaseURLPlaceholder");
  node.querySelector('[data-field="api_key"]').placeholder = t("config.providerAPIKeyPlaceholder");
}

function updateEditorRowTranslations() {
  tokenInput.placeholder = t("auth.placeholder");
  routesEl.querySelectorAll(".route-row").forEach((node) => hydrateRouteRow(node));
  providerEditorEl.querySelectorAll(".provider-row").forEach((node) => hydrateProviderRow(node));
}

function updatePanelCopy() {
  panelTitleEl.textContent = t(`topbar.${state.activeTab}`);
  panelDescriptionEl.textContent = t(`topbar.${state.activeTab}Desc`);
}

function updateDynamicTranslations() {
  document.querySelectorAll("[data-i18nDynamic]").forEach((node) => {
    node.textContent = t(node.dataset.i18nDynamic);
  });
  updateEditorRowTranslations();
}

function collectConfig() {
  return {
    listen: formEl.elements.listen.value.trim(),
    mode: formEl.elements.mode.value,
    failure_threshold: Number(formEl.elements.failure_threshold.value),
    cooldown: formEl.elements.cooldown.value.trim(),
    health_check_interval: formEl.elements.health_check_interval.value.trim(),
    health_check_timeout: formEl.elements.health_check_timeout.value.trim(),
    routes: Array.from(routesEl.querySelectorAll(".route-row")).map((row) => ({
      prefix: row.querySelector('[data-field="prefix"]').value.trim(),
      type: row.querySelector('[data-field="type"]').value,
      model: row.querySelector('[data-field="model"]').value.trim(),
      upstream_model: row.querySelector('[data-field="upstream_model"]').value.trim(),
      enabled: row.querySelector('[data-field="enabled"]').checked
    })),
    providers: Array.from(providerEditorEl.querySelectorAll(".provider-row")).map((row) => ({
      name: row.querySelector('[data-field="name"]').value.trim(),
      base_url: row.querySelector('[data-field="base_url"]').value.trim(),
      api_key: row.querySelector('[data-field="api_key"]').value.trim(),
      enabled: row.querySelector('[data-field="enabled"]').checked
    }))
  };
}

function syncCollection(root, cache, desiredLength, factory) {
  while (cache.length < desiredLength) {
    const item = factory();
    cache.push(item);
    root.appendChild(item.root);
  }
  while (cache.length > desiredLength) {
    const item = cache.pop();
    item.root.remove();
  }
}

function pruneMapEntries(map, root, seen) {
  for (const [key, item] of map.entries()) {
    if (seen.has(key)) {
      continue;
    }
    item.root.remove();
    map.delete(key);
  }
}

async function request(url, options = {}) {
  const headers = new Headers(options.headers || {});
  if (!headers.has("Content-Type") && options.body) {
    headers.set("Content-Type", "application/json");
  }
  if (state.token) {
    headers.set("X-PSwitch-Admin-Token", state.token);
  }
  const response = await fetch(url, { ...options, headers });
  const data = await response.json().catch(() => ({}));
  if (!response.ok) {
    const error = new Error(data.error || `Request failed with ${response.status}`);
    error.status = response.status;
    throw error;
  }
  return data;
}

function showMessage(kind, text) {
  if (!text) {
    messageEl.className = "message-card hidden";
    messageEl.textContent = "";
    return;
  }
  messageEl.className = `message-card ${kind}`;
  messageEl.textContent = text;
}

function t(key) {
  const localeMessages = messages[state.locale] || messages.en;
  return localeMessages[key] || messages.en[key] || key;
}

function applyTranslations() {
  document.documentElement.lang = state.locale;
  document.querySelectorAll("[data-i18n]").forEach((node) => {
    node.textContent = t(node.dataset.i18n);
  });
  tokenInput.placeholder = t("auth.placeholder");
}

function inferLocale() {
  return (navigator.language || "").toLowerCase().startsWith("zh") ? "zh-CN" : "en";
}

function formatNumber(value) {
  return new Intl.NumberFormat(state.locale).format(Number(value || 0));
}

function formatTokenCount(value) {
  const number = Number(value || 0);
  if (Math.abs(number) < 10000) {
    return formatNumber(number);
  }
  return new Intl.NumberFormat(state.locale, {
    notation: "compact",
    maximumFractionDigits: number >= 100000000 ? 0 : 1
  }).format(number);
}

function formatCompact(value) {
  return new Intl.NumberFormat(state.locale, {
    notation: "compact",
    maximumFractionDigits: 1
  }).format(Number(value || 0));
}

function formatDateTime(value) {
  if (!value) {
    return "-";
  }
  return new Intl.DateTimeFormat(state.locale, {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit"
  }).format(new Date(value));
}

function joinMetrics(items) {
  return items.map(([label, value]) => `${label} ${formatCompact(value)}`).join(" · ");
}

function formatModeLabel(mode) {
  switch (mode) {
    case "round_robin":
      return t("config.modeRoundRobinLabel");
    case "sequential":
      return t("config.modeSequentialLabel");
    case "least_failures":
      return t("config.modeLeastFailuresLabel");
    default:
      return mode || "-";
  }
}

function isTokenMetricKey(key) {
  return key === "overview.totalTokens" || key === "overview.tokens24h" || key === "overview.tokens7d";
}

function resolveTokenMetricValue(key, overview, windows) {
  switch (key) {
    case "overview.totalTokens":
      return overview.total_tokens;
    case "overview.tokens24h":
      return windows.last_24h.total_tokens;
    case "overview.tokens7d":
      return windows.last_7d.total_tokens;
    default:
      return 0;
  }
}

boot();
