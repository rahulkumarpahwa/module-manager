/* ============================================================
   Module Manager Client - app.js
   Talks to the Module Manager Go backend using the endpoints
   defined in postman/module-manager.postman_collection.json:

     GET  /api/healthz
     POST /api/set-path            { "path": "..." }
     POST /api/main/scan?limit=&offset=
     POST /api/project/{id}        { "path": "..." }  (full project details)
     POST /api/project/infos       [ { "id": "...", "path": "..." }, ... ]  (bulk detail fetch)
     POST /api/project/open-vscode { "path": "..." }
     POST /api/project/open-winexplorer { "path": "..." }

   The backend responses have been evolving (the Postman examples
   show the intended shape). We defensively normalize BOTH the
   Postman shape and the current Go handler shape.
   ============================================================ */
"use strict";

(() => {
  const DEFAULT_BASE_URL = "http://localhost:8778";
  const DEFAULT_LIMIT = 20;
  const DEFAULT_OFFSET = 0;

  // Ordered list of fields shown in an expanded project row.
  // Mirrors the project object (types.Project) returned by POST /api/project/{id}.
  const PROJECT_FIELDS = [
    ["id", "id"],
    ["name", "name"],
    ["path", "path"],
    ["size", "size"],
    ["date_created", "date_created"],
    ["date_modified", "date_modified"],
    ["package_manager", "package_manager"],
    ["has_package_json", "has_package_json"],
    ["has_node_modules", "has_node_modules"],
    ["node_modules_size", "node_modules_size"],
    ["git_status", "git_status"],
    ["latest_commit", "latest_commit"],
    ["latest_commit_date", "latest_commit_date"],
    ["libraries", "libraries"],
    ["under_dir", "under_dir"],
  ];

  const els = {
    baseUrl: $("#baseUrl"),
    folderPath: $("#folderPath"),
    setPathBtn: $("#setPathBtn"),
    scanBtn: $("#scanBtn"),
    limit: $("#limit"),
    offset: $("#offset"),
    scannedPath: $("#scannedPath"),
    toast: $("#toast"),
    summary: $("#summary"),
    body: $("#projectsBody"),
    raw: $("#rawResponse"),
    toggleRawBtn: $("#toggleRawBtn"),
    healthDot: $("#healthDot"),
    healthText: $("#healthText"),
    healthBtn: $("#healthBtn"),
  };

  const projects = new Map(); // id -> normalized project (details merged on fetch)
  let lastRaw = null;         // last API response payload (for the Raw panel)
  let fetchToken = 0;         // guards against racing detail fetches

  // ---------- tiny DOM helpers ----------
  function $(sel) { return document.querySelector(sel); }

  function esc(v) {
    return String(v ?? "")
      .replaceAll("&", "&amp;")
      .replaceAll("<", "&lt;")
      .replaceAll(">", "&gt;")
      .replaceAll('"', "&quot;")
      .replaceAll("'", "&#39;");
  }

  function baseUrl() {
    return (els.baseUrl.value.trim() || DEFAULT_BASE_URL).replace(/\/+$/, "");
  }

  // CSS selector escaping for ids that may contain special characters.
  function sel(str) {
    return (typeof CSS !== "undefined" && CSS.escape)
      ? CSS.escape(str)
      : String(str).replace(/["\\]/g, "\\$&");
  }

  // ---------- formatting ----------
  function formatBytes(bytes) {
    if (bytes === null || bytes === undefined || bytes === "") return "—";
    const n = Number(bytes);
    if (!Number.isFinite(n)) return "—";
    if (Math.abs(n) < 1024) return `${n} B`;
    const units = ["KB", "MB", "GB", "TB", "PB"];
    let v = n;
    let u = -1;
    do { v /= 1024; u += 1; } while (Math.abs(v) >= 1024 && u < units.length - 1);
    return `${v.toFixed(2)} ${units[u]}`;
  }

  function formatDateTime(value) {
    if (!value) return "—";
    const d = new Date(value);
    if (Number.isNaN(d.getTime())) return esc(String(value));
    const pad = (n) => String(n).padStart(2, "0");
    const date = `${d.getUTCFullYear()}-${pad(d.getUTCMonth() + 1)}-${pad(d.getUTCDate())}`;
    const time = `${pad(d.getUTCHours())}:${pad(d.getUTCMinutes())}:${pad(d.getUTCSeconds())}`;
    return `${date} ${time}`;
  }

  function boolBadge(v) {
    if (v === true) return '<span class="badge yes">yes</span>';
    if (v === false) return '<span class="badge no">no</span>';
    return '<span class="badge unknown">—</span>';
  }

  // ---------- response normalization ----------
  // Scan response ("directory_details" = current Go handler shape,
  // top-level fields = Postman collection shape).
  function parseScan(json) {
    const payload = (json && json.data && json.data.directory_details) ? json.data.directory_details : json;

    const rawEntries = (payload && payload.projects) || (payload && payload.entries) || [];

    const summary = {
      path: payload.path || json.path || "",
      total_size: payload.total_size ?? payload.total_Size ?? json.total_size ?? null,
      total_projects: payload.total_projects ?? payload.projects_count ?? rawEntries.length,
    };

    const list = rawEntries.map((e, i) => ({
      id: e.id ?? "",
      name: e.name ?? "",
      path: e.path ?? "",
      is_dir: e.is_dir ?? null,
      total_size: e.total_size ?? e.size ?? 0,
      total_size_str: e.total_size_str ?? null,
      size: e.size ?? e.total_size ?? 0,
      date_created: e.date_created ?? null,
      date_modified: e.date_modified ?? e.modified_date ?? null,
      package_manager: e.package_manager ?? null,
      libraries: e.libraries ?? null,
      has_package_json: e.has_package_json ?? null,
      has_node_modules: e.has_node_modules ?? null,
      node_modules_size: e.node_modules_size ?? null,
      node_modules_size_str: e.node_modules_size_str ?? null,
      has_git: e.has_git ?? e.git_status ?? null,
      git_status: e.git_status ?? e.has_git ?? null,
      frameworks: e.frameworks ?? null,
      latest_commit: e.latest_commit ?? null,
      latest_commit_date: e.latest_commit_date ?? null,
      list_files: e.list_files ?? null,
      under_dir: e.under_dir ?? null,
      language: e["framework/language"] ?? e.framework_language ?? null,
      workspace: e.workspace ?? null,
      project: e.project ?? null,
      _detailLoading: false,
      _detailLoaded: false,
    }));

    return { summary, list };
  }

  // Project detail response from POST /api/project/{id}:
  //   "data": { "projects": [ { ...types.Project } ] }
  function parseProjectDetail(json) {
    if (!json) return null;
    if (json.data && Array.isArray(json.data.projects)) return json.data.projects[0] || null;
    if (json.data && json.data.project) return json.data.project;
    if (typeof json.project === "object" && json.project) return json.project;
    if (json.id || json.name) return json; // envelope is the project itself
    return null;
  }

  // Bulk detail response from POST /api/project/infos:
  //   "data": { "projects": [ ...types.Project ] }
  // A bare top-level array is tolerated as well.
  function parseProjectDetails(json) {
    if (!json) return [];
    if (Array.isArray(json)) return json;
    if (Array.isArray(json.projects)) return json.projects;
    if (json.data && Array.isArray(json.data.projects)) return json.data.projects;
    return [];
  }

  // ---------- rendering ----------
  function renderEmpty() {
    els.body.innerHTML = `
      <tr>
        <td colspan="9" class="empty">Nothing scanned yet. Set a path and click <b>Scan</b>.</td>
      </tr>`;
  }

  function renderSummary(summary) {
    const totalSize = formatBytes(summary.total_size);
    els.summary.innerHTML = summary.path
      ? `Path: <b>${esc(summary.path)}</b> &nbsp;·&nbsp; Projects: <b>${summary.total_projects}</b> &nbsp;·&nbsp; Size: <b>${totalSize}</b>
         <button class="btn btn-ghost" id="fetchAllBtn" type="button" title="Fetch POST /api/project/infos for every row (single batch request)">Fetch All Details</button>`
      : "";
    const fetchAll = $("#fetchAllBtn");
    if (fetchAll) fetchAll.addEventListener("click", fetchAllDetails);
  }

  function rowCells(p, index) {
    return `
      <td class="num">${index + 1}</td>
      <td class="mono">${esc(p.id)}</td>
      <td><b>${esc(p.name)}</b></td>
      <td class="mono path" title="${esc(p.path)}">${esc(p.path)}</td>
      <td class="mono">${formatDateTime(p.modified_date)}</td>
      <td>${boolBadge(p.workspace)}</td>
      <td>${boolBadge(p.project)}</td>
      <td>${esc(p.language || "—")}</td>
      <td class="num actions">
        <button class="btn btn-vscode open-vscode-btn" type="button" data-id="${esc(p.id)}" title="Open ${esc(p.path)} in VS Code">Open in VS Code</button>
        <button class="btn btn-explorer open-explorer-btn" type="button" data-id="${esc(p.id)}" title="Open ${esc(p.path)} in File Explorer">Open Explorer</button>
        <span class="caret">&#9654;</span>
      </td>`;
  }

  function renderProjects(list) {
    projects.clear();
    list.forEach((p) => projects.set(p.id, { ...p }));

    if (list.length === 0) {
      renderEmpty();
      return;
    }

    els.body.innerHTML = list
      .map((p, i) => `<tr class="project-row" data-id="${esc(p.id)}" data-index="${i}">${rowCells(p, i)}</tr>`)
      .join("");
    // Row clicks (expand/collapse) and the per-row "Open in VS Code" /
    // "Open Explorer" buttons are handled by a single delegated listener
    // on els.body (see events).
  }

  function detailGrid(p) {
    return PROJECT_FIELDS.map(([label, key]) => {
      let v = p[key];
      let html;
      if (v === null || v === undefined) {
        html = "—";
      } else if (typeof v === "boolean") {
        html = boolBadge(v);
      } else if (key === "under_dir") {
        // "under_dir": [ { name, size, path }, ... ]
        html = v.length
          ? `<div class="under-dir">${v.map((item) => {
              const sizeText = item.size ? formatBytes(item.size) : "";
              return `<div class="under-dir-row" title="${esc(item.path)}">
                        <span class="under-name">${esc(item.name)}</span>
                        ${sizeText ? `<span class="under-size">${sizeText}</span>` : ""}
                      </div>`;
            }).join("")}</div>`
          : "—";
      } else if (Array.isArray(v)) {
        html = v.length
          ? `<div class="chips">${v.map((item) => `<span class="chip">${esc(item)}</span>`).join("")}</div>`
          : "—";
      } else if (key === "size" || key === "node_modules_size") {
        html = esc(`${v} (${formatBytes(v)})`);
      } else if (key === "date_created" || key === "date_modified" || key === "latest_commit_date") {
        html = formatDateTime(v);
      } else {
        html = esc(String(v));
      }
      return `<div class="detail-item"><span class="k">${esc(label)}</span><span class="v">${html}</span></div>`;
    }).join("");
  }

  function expandDetails(p) {
    const extraFields = [];
    extraFields.push(`<div class="detail-item"><span class="k">workspace</span><span class="v">${boolBadge(p.workspace)}</span></div>`);
    extraFields.push(`<div class="detail-item"><span class="k">project</span><span class="v">${boolBadge(p.project)}</span></div>`);
    if (p.language) extraFields.push(`<div class="detail-item"><span class="k">framework/language</span><span class="v">${esc(p.language)}</span></div>`);
    extraFields.push(`<div class="detail-item"><span class="k">modified_date</span><span class="v">${formatDateTime(p.modified_date)}</span></div>`);

    return `
      <div class="detail-box">
        <div class="detail-grid">
          ${detailGrid(p)}
          ${extraFields.join("")}
        </div>
        <div class="detail-raw">Expand uses <code>POST /api/project/{id}</code> — status:
          ${p._detailLoading ? '<span class="loading-dot">loading…</span>' : (p._detailLoaded ? "loaded" : "not loaded yet")}
        </div>
      </div>`;
  }

  // ---------- table row expand / collapse ----------
  async function toggleRow(tr) {
    const id = tr.dataset.id;
    const p = projects.get(id);

    if (!p) return;

    const nextTd = tr.nextElementSibling;
    if (nextTd && nextTd.classList.contains("detail-row")) {
      // collapse
      nextTd.remove();
      tr.classList.remove("active");
      return;
    }

    tr.classList.add("active");

    const detailTr = document.createElement("tr");
    detailTr.className = "detail-row";
    detailTr.innerHTML = `<td colspan="9">${expandDetails(p)}</td>`;
    tr.after(detailTr);

    if (!p._detailLoading && !p._detailLoaded) {
      await fetchProjectDetails(id);
      // Re-render the expanded view with merged data (if still open)
      const updated = projects.get(id);
      if (tr.classList.contains("active") && tr.nextElementSibling === detailTr && updated) {
        detailTr.innerHTML = `<td colspan="9">${expandDetails(updated)}</td>`;
      }
    }
  }

  // ---------- API ----------
  async function apiRequest(path, options = {}) {
    const url = baseUrl() + path;
    const res = await fetch(url, options);
    const text = await res.text();
    let json = null;
    try { json = text ? JSON.parse(text) : null; } catch (_) { /* plain-text body */ }

    if (!res.ok) {
      const msg = json ? (json.message || JSON.stringify(json)) : (text || `HTTP ${res.status}`);
      const err = new Error(msg);
      err.status = res.status;
      err.json = json;
      throw err;
    }
    return { json, text, status: res.status };
  }

  async function setPath() {
    const pathValue = els.folderPath.value.trim();
    if (!pathValue) {
      showToast("path cannot be empty", "err");
      els.folderPath.focus();
      return;
    }

    setBusy(els.setPathBtn, true, "Setting…");
    try {
      const { json } = await apiRequest("/api/set-path", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ path: pathValue }),
      });
      lastRaw = json;
      renderRaw();
      showToast(`Set Path OK → ${json.data?.path || pathValue}`, "ok");
      els.scannedPath.textContent = `stored path: ${json.data?.path || pathValue}`;
    } catch (err) {
      showToast(`Set Path failed: ${err.message}`, "err");
    } finally {
      setBusy(els.setPathBtn, false);
    }
  }

  async function scan() {
    const pathValue = els.folderPath.value.trim();
    const limit = Number.parseInt(els.limit.value, 10) || DEFAULT_LIMIT;
    const offset = Number.parseInt(els.offset.value, 10) || DEFAULT_OFFSET;

    // Ensure a path is configured first (POST /api/set-path).
    if (pathValue) {
      try {
        await apiRequest("/api/set-path", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ path: pathValue }),
        });
      } catch (err) {
        showToast(`Set Path failed: ${err.message}`, "err");
        return;
      }
    }

    setBusy(els.scanBtn, true, "Scanning…");
    try {
      const { json } = await apiRequest(`/api/main/scan?limit=${limit}&offset=${offset}`, { method: "POST" });
      lastRaw = json;
      renderRaw();

      const parsed = parseScan(json);
      els.scannedPath.textContent = parsed.summary.path
        ? `scanned: ${parsed.summary.path}`
        : "scanned (no path echoed)";

      renderProjects(parsed.list);
      renderSummary({ ...parsed.summary, total_projects: parsed.list.length });

      const msg = `Scan OK — ${parsed.list.length} project(s), window offset=${offset}, limit=${limit}`;
      showToast(msg, "ok");
    } catch (err) {
      showToast(`Scan failed: ${err.message}`, "err");
    } finally {
      setBusy(els.scanBtn, false);
    }
  }

  function mergeDetail(project, detail) {
    for (const [, key] of PROJECT_FIELDS) {
      if (detail[key] !== undefined && detail[key] !== null) {
        project[key] = detail[key];
      }
    }
    if (detail.path) project.path = detail.path;
    if (detail.name) project.name = detail.name;
    project._detailLoaded = true;
    project._detailLoading = false;
  }

  async function fetchProjectDetails(id) {
    const p = projects.get(id);
    if (!p || p._detailLoading || p._detailLoaded) return;

    p._detailLoading = true;
    const token = ++fetchToken;

    try {
      const { json } = await apiRequest(`/api/project/${encodeURIComponent(id)}`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ path: p.path || "" }),
      });
      const detail = parseProjectDetail(json);
      if (detail) mergeDetail(projects.get(id), detail);
    } catch (err) {
      showToast(`Detail [${id}] failed: ${err.message}`, "err");
      if (projects.has(id)) projects.get(id)._detailLoading = false;
    } finally {
      if (token === fetchToken) patchSummaryColumn(id);
    }
  }

  // Fetch full details for every visible row in a single batch request.
  // Uses the bulk endpoint POST /api/project/infos, which takes an array of
  // { id, path } in the body, so the backend resolves all projects at once
  // instead of us firing one request per row.
  async function fetchAllDetails() {
    const list = [...projects.values()].filter((p) => !p._detailLoading && !p._detailLoaded);
    if (list.length === 0) {
      if (projects.size > 0) {
        showToast("All visible details are already loaded", "ok");
      }
      return;
    }

    const btn = $("#fetchAllBtn");
    if (btn) btn.disabled = true;

    let ok = 0;
    let failed = 0;
    let failReason = "";
    try {
      const { json } = await apiRequest("/api/project/infos", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(list.map((p) => ({ id: p.id, path: p.path }))),
      });
      lastRaw = json;
      renderRaw();

      // The batch worker pool resolves concurrently, so order is not
      // guaranteed — match each returned detail back to its row by id.
      const byId = new Map(
        parseProjectDetails(json)
          .filter((d) => d && d.id != null)
          .map((d) => [String(d.id), d])
      );

      for (const p of list) {
        const detail = byId.get(String(p.id));
        if (detail) {
          mergeDetail(projects.get(p.id), detail);
          ok += 1;
        } else {
          failed += 1;
        }
      }
    } catch (err) {
      failed = list.length;
      failReason = err.message;
    } finally {
      // Refresh visible rows so merged columns (pkg_json, git, …) update.
      rerenderRows();
      if (btn) {
        btn.disabled = false;
        btn.textContent = "Fetch All Details";
      }
    }

    if (failed) {
      showToast(`Details: ${ok} loaded, ${failed} failed${failReason ? ` — ${failReason}` : ""}`, "warn");
    } else {
      showToast(`Details: ${ok} loaded`, "ok");
    }
  }

  // Open a project folder via one of the "open" endpoints
  // (POST /api/project/open-vscode | POST /api/project/open-winexplorer).
  async function openProject(apiPath, id, btn, label) {
    const p = projects.get(id);
    if (!p) return;

    const projectPath = p.path;
    if (!projectPath) {
      showToast(`Project [${id}] has no path — cannot open ${label}`, "warn");
      return;
    }

    setBusy(btn, true, "Opening…");
    try {
      const { json } = await apiRequest(apiPath, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ path: projectPath }),
      });
      lastRaw = json;
      renderRaw();
      const openedPath = json?.data?.path || projectPath;
      showToast(`${label} opened → ${openedPath}`, "ok");
    } catch (err) {
      showToast(`Open ${label} failed: ${err.message}`, "err");
    } finally {
      if (btn && btn.isConnected) setBusy(btn, false);
    }
  }

  function openInVSCode(id, btn) {
    return openProject("/api/project/open-vscode", id, btn, "VS Code");
  }

  function openInWinExplorer(id, btn) {
    return openProject("/api/project/open-winexplorer", id, btn, "Explorer");
  }

  function rerenderRows() {
    const list = [...projects.values()];
    renderProjects(list);
    // Keep the summary (rebuilt only by scan) intact.
  }

  async function checkHealth() {
    els.healthDot.className = "health-dot busy";
    els.healthText.textContent = "Checking…";
    try {
      const { json } = await apiRequest("/api/healthz", { method: "GET" });
      els.healthDot.className = "health-dot ok";
      els.healthText.textContent = `Connected — ${json.message || "ok"}${json.current_time ? ` (${json.current_time})` : ""}`;
      return true;
    } catch (err) {
      els.healthDot.className = "health-dot err";
      els.healthText.textContent = `Unreachable — ${err.message}`;
      return false;
    }
  }

  // ---------- misc UI helpers ----------
  let toastTimer = null;

  function showToast(message, type = "ok", ms = 4000) {
    els.toast.textContent = message;
    els.toast.className = `toast ${type === "err" ? "err" : type === "warn" ? "warn" : ""}`;
    els.toast.hidden = false;
    clearTimeout(toastTimer);
    toastTimer = setTimeout(() => { els.toast.hidden = true; }, ms);
  }

  function setBusy(btn, busy, label) {
    if (!btn) return;
    if (busy) {
      btn.dataset.orig = btn.textContent;
      btn.disabled = true;
      btn.innerHTML = `<span class="spinner"></span>${label || "…"}`;
    } else {
      btn.disabled = false;
      btn.innerHTML = "";
      btn.textContent = btn.dataset.orig || "";
    }
  }

  function renderRaw() {
    if (lastRaw == null) {
      els.raw.hidden = true;
      els.toggleRawBtn.textContent = "show";
      return;
    }
    els.raw.textContent = JSON.stringify(lastRaw, null, 2);
  }

  function patchSummaryColumn(id) {
    const p = projects.get(id);
    if (!p) return;
    const tr = document.querySelector(`.project-row[data-id="${sel(id)}"]`);
    if (!tr) return;
    const index = Number(tr.dataset.index);
    tr.innerHTML = rowCells(p, index);
  }

  // ---------- events ----------
  els.setPathBtn.addEventListener("click", setPath);
  els.scanBtn.addEventListener("click", scan);
  els.healthBtn.addEventListener("click", () => checkHealth());
  els.folderPath.addEventListener("keydown", (e) => { if (e.key === "Enter") setPath(); });

  // Delegated clicks on the projects table: "Open in VS Code" / "Open Explorer"
  // buttons vs row expand.
  els.body.addEventListener("click", (e) => {
    const vscodeBtn = e.target.closest(".open-vscode-btn");
    if (vscodeBtn) {
      openInVSCode(vscodeBtn.dataset.id, vscodeBtn);
      return;
    }
    const explorerBtn = e.target.closest(".open-explorer-btn");
    if (explorerBtn) {
      openInWinExplorer(explorerBtn.dataset.id, explorerBtn);
      return;
    }
    const tr = e.target.closest(".project-row");
    if (tr) toggleRow(tr);
  });

  els.toggleRawBtn.addEventListener("click", () => {
    const show = els.raw.hidden;
    els.raw.hidden = !show;
    els.toggleRawBtn.textContent = show ? "hide" : "show";
    if (lastRaw == null) els.raw.textContent = "No response captured yet. Run Set Path / Scan first.";
  });

  // ---------- init ----------
  renderEmpty();
  renderRaw();
  checkHealth();
})();