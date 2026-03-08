import { useState, useRef, useEffect } from "react";
import * as pdfjsLib from 'pdfjs-dist';
pdfjsLib.GlobalWorkerOptions.workerSrc = new URL(
  'pdfjs-dist/build/pdf.worker.min.mjs',
  import.meta.url
).toString();

const API_BASE = "http://localhost:8080/api/v1";

const styles = `
  @import url('https://fonts.googleapis.com/css2?family=Space+Mono:wght@400;700&family=DM+Sans:wght@300;400;500;600&display=swap');

  *, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }

  :root {
    --bg: #0d0f14;
    --surface: #13161e;
    --surface2: #1a1e28;
    --border: #252a38;
    --accent: #5b8aff;
    --accent2: #a78bfa;
    --green: #34d399;
    --red: #f87171;
    --yellow: #fbbf24;
    --text: #e2e8f0;
    --muted: #64748b;
    --mono: 'Space Mono', monospace;
    --sans: 'DM Sans', sans-serif;
  }

  body { background: var(--bg); color: var(--text); font-family: var(--sans); min-height: 100vh; }

  .app { display: grid; grid-template-columns: 260px 1fr; min-height: 100vh; }

  /* SIDEBAR */
  .sidebar {
    background: var(--surface);
    border-right: 1px solid var(--border);
    display: flex;
    flex-direction: column;
    padding: 0;
    position: sticky;
    top: 0;
    height: 100vh;
    overflow-y: auto;
  }
  .sidebar-logo {
    padding: 20px 20px 16px;
    border-bottom: 1px solid var(--border);
  }
  .logo-title {
    font-family: var(--mono);
    font-size: 13px;
    font-weight: 700;
    color: var(--accent);
    letter-spacing: 2px;
    text-transform: uppercase;
  }
  .logo-sub {
    font-size: 11px;
    color: var(--muted);
    margin-top: 3px;
    font-family: var(--mono);
  }
  .nav-section { padding: 16px 12px 8px; }
  .nav-label {
    font-size: 10px;
    font-weight: 600;
    color: var(--muted);
    letter-spacing: 1.5px;
    text-transform: uppercase;
    padding: 0 8px;
    margin-bottom: 6px;
  }
  .nav-btn {
    display: flex;
    align-items: center;
    gap: 10px;
    width: 100%;
    padding: 9px 12px;
    background: none;
    border: none;
    border-radius: 8px;
    color: var(--muted);
    font-family: var(--sans);
    font-size: 13px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.15s;
    text-align: left;
  }
  .nav-btn:hover { background: var(--surface2); color: var(--text); }
  .nav-btn.active { background: rgba(91,138,255,0.12); color: var(--accent); }
  .nav-btn .icon { font-size: 15px; width: 18px; text-align: center; }

  /* COLLECTIONS LIST IN SIDEBAR */
  .collections-list { padding: 0 12px 12px; }
  .collection-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 7px 12px;
    border-radius: 6px;
    cursor: pointer;
    font-size: 12px;
    color: var(--muted);
    transition: all 0.15s;
    gap: 6px;
  }
  .collection-item:hover { background: var(--surface2); color: var(--text); }
  .collection-item.selected { background: rgba(91,138,255,0.1); color: var(--accent); }
  .collection-dot { width: 6px; height: 6px; border-radius: 50%; background: var(--accent2); flex-shrink: 0; }
  .collection-name { flex: 1; font-family: var(--mono); font-size: 11px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  .collection-del {
    background: none; border: none; cursor: pointer;
    color: transparent; font-size: 12px; padding: 2px;
    transition: color 0.15s; border-radius: 3px;
  }
  .collection-item:hover .collection-del { color: var(--red); }
  .no-collections { padding: 10px 12px; font-size: 11px; color: var(--muted); font-style: italic; }

  /* MAIN CONTENT */
  .main { display: flex; flex-direction: column; min-height: 100vh; overflow: hidden; }
  .topbar {
    padding: 16px 28px;
    border-bottom: 1px solid var(--border);
    display: flex;
    align-items: center;
    justify-content: space-between;
    background: var(--surface);
  }
  .page-title { font-size: 16px; font-weight: 600; color: var(--text); }
  .status-dot { width: 8px; height: 8px; border-radius: 50%; background: var(--green); display: inline-block; margin-right: 6px; animation: pulse 2s infinite; }
  @keyframes pulse { 0%,100%{opacity:1} 50%{opacity:0.4} }
  .status-text { font-size: 12px; color: var(--muted); display: flex; align-items: center; }

  /* CONTENT AREA */
  .content { flex: 1; padding: 28px; overflow-y: auto; }

  /* CARDS */
  .card {
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: 12px;
    padding: 22px;
    margin-bottom: 20px;
  }
  .card-title {
    font-size: 12px;
    font-weight: 600;
    color: var(--muted);
    letter-spacing: 1px;
    text-transform: uppercase;
    margin-bottom: 18px;
    display: flex;
    align-items: center;
    gap: 8px;
  }

  /* FORM ELEMENTS */
  .form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 14px; margin-bottom: 14px; }
  .form-row.single { grid-template-columns: 1fr; }
  .form-row.three { grid-template-columns: 1fr 1fr 1fr; }
  label { display: block; font-size: 11px; font-weight: 600; color: var(--muted); margin-bottom: 5px; letter-spacing: 0.5px; text-transform: uppercase; }
  input[type=text], input[type=number], select, textarea {
    width: 100%;
    background: var(--bg);
    border: 1px solid var(--border);
    border-radius: 8px;
    padding: 9px 12px;
    color: var(--text);
    font-family: var(--sans);
    font-size: 13px;
    outline: none;
    transition: border-color 0.15s;
  }
  input:focus, select:focus, textarea:focus { border-color: var(--accent); }
  textarea { resize: vertical; min-height: 90px; }
  select option { background: var(--surface2); }

  /* FILE UPLOAD */
  .file-drop {
    border: 2px dashed var(--border);
    border-radius: 10px;
    padding: 28px;
    text-align: center;
    cursor: pointer;
    transition: all 0.2s;
    background: var(--bg);
    position: relative;
  }
  .file-drop:hover, .file-drop.drag-over { border-color: var(--accent); background: rgba(91,138,255,0.04); }
  .file-drop input[type=file] { position: absolute; inset: 0; opacity: 0; cursor: pointer; width: 100%; height: 100%; }
  .file-drop-icon { font-size: 28px; margin-bottom: 8px; }
  .file-drop-text { font-size: 13px; color: var(--muted); }
  .file-drop-text span { color: var(--accent); }
  .file-selected { font-size: 12px; color: var(--green); margin-top: 6px; font-family: var(--mono); }
  .file-parsing { font-size: 12px; color: var(--yellow); margin-top: 8px; font-family: var(--mono); display: flex; align-items: center; gap: 6px; justify-content: center; }

  /* BUTTONS */
  .btn {
    padding: 9px 18px;
    border-radius: 8px;
    border: none;
    font-family: var(--sans);
    font-size: 13px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.15s;
    display: inline-flex;
    align-items: center;
    gap: 6px;
  }
  .btn-primary { background: var(--accent); color: #fff; }
  .btn-primary:hover { background: #4a79ef; }
  .btn-primary:disabled { opacity: 0.5; cursor: not-allowed; }
  .btn-ghost { background: var(--surface2); color: var(--text); border: 1px solid var(--border); }
  .btn-ghost:hover { border-color: var(--accent); color: var(--accent); }
  .btn-danger { background: rgba(248,113,113,0.1); color: var(--red); border: 1px solid rgba(248,113,113,0.2); }
  .btn-danger:hover { background: rgba(248,113,113,0.2); }

  /* TOGGLE */
  .toggle-row { display: flex; align-items: center; justify-content: space-between; padding: 10px 0; border-bottom: 1px solid var(--border); }
  .toggle-row:last-child { border-bottom: none; }
  .toggle-label { font-size: 13px; color: var(--text); }
  .toggle-sub { font-size: 11px; color: var(--muted); margin-top: 2px; }
  .toggle {
    position: relative; width: 38px; height: 20px; background: var(--border);
    border-radius: 10px; cursor: pointer; transition: background 0.2s;
  }
  .toggle.on { background: var(--accent); }
  .toggle::after {
    content: ''; position: absolute; top: 3px; left: 3px;
    width: 14px; height: 14px; border-radius: 50%; background: white;
    transition: left 0.2s;
  }
  .toggle.on::after { left: 21px; }

  /* RESPONSE AREA */
  .answer-box {
    background: var(--bg);
    border: 1px solid var(--border);
    border-radius: 10px;
    padding: 18px;
    font-size: 14px;
    line-height: 1.7;
    color: var(--text);
    white-space: pre-wrap;
    word-break: break-word;
    min-height: 60px;
  }
  .answer-box.empty { color: var(--muted); font-style: italic; font-size: 13px; }

  /* CHUNKS */
  .chunks-grid { display: grid; gap: 10px; margin-top: 16px; }
  .chunk-card {
    background: var(--bg);
    border: 1px solid var(--border);
    border-radius: 8px;
    padding: 14px;
    font-size: 12px;
  }
  .chunk-header { display: flex; align-items: center; gap: 8px; margin-bottom: 8px; flex-wrap: wrap; }
  .chunk-badge {
    padding: 2px 8px; border-radius: 4px; font-family: var(--mono);
    font-size: 10px; font-weight: 700; background: rgba(91,138,255,0.12); color: var(--accent);
  }
  .chunk-badge.score { background: rgba(52,211,153,0.12); color: var(--green); }
  .chunk-badge.section { background: rgba(167,139,250,0.12); color: var(--accent2); }
  .chunk-text { color: var(--muted); line-height: 1.6; font-size: 12px; }

  /* STATS GRID */
  .stats-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 14px; }
  .stat-card {
    background: var(--bg);
    border: 1px solid var(--border);
    border-radius: 8px;
    padding: 14px;
    text-align: center;
  }
  .stat-val { font-size: 22px; font-weight: 700; font-family: var(--mono); color: var(--accent); }
  .stat-label { font-size: 11px; color: var(--muted); margin-top: 3px; }

  /* ALERT */
  .alert {
    padding: 10px 14px; border-radius: 8px; font-size: 12px;
    margin-top: 12px; display: flex; align-items: center; gap: 8px;
  }
  .alert-success { background: rgba(52,211,153,0.1); border: 1px solid rgba(52,211,153,0.25); color: var(--green); }
  .alert-error { background: rgba(248,113,113,0.1); border: 1px solid rgba(248,113,113,0.25); color: var(--red); }
  .alert-info { background: rgba(91,138,255,0.1); border: 1px solid rgba(91,138,255,0.25); color: var(--accent); }

  /* LOADING */
  .spinner {
    width: 14px; height: 14px; border: 2px solid transparent;
    border-top-color: currentColor; border-radius: 50%;
    animation: spin 0.6s linear infinite; display: inline-block;
  }
  @keyframes spin { to { transform: rotate(360deg); } }

  /* PROCESSING TIME */
  .proc-time { font-family: var(--mono); font-size: 11px; color: var(--muted); margin-top: 8px; }

  /* SCROLLBAR */
  ::-webkit-scrollbar { width: 6px; }
  ::-webkit-scrollbar-track { background: transparent; }
  ::-webkit-scrollbar-thumb { background: var(--border); border-radius: 3px; }

  /* TABS */
  .tabs { display: flex; gap: 4px; margin-bottom: 20px; background: var(--surface2); padding: 4px; border-radius: 10px; }
  .tab {
    flex: 1; padding: 8px; border-radius: 7px; border: none;
    background: none; color: var(--muted); font-family: var(--sans);
    font-size: 12px; font-weight: 600; cursor: pointer; transition: all 0.15s;
  }
  .tab.active { background: var(--surface); color: var(--text); }

  /* RESPONSIVE */
  @media (max-width: 768px) {
    .app { grid-template-columns: 1fr; }
    .sidebar { display: none; }
    .form-row { grid-template-columns: 1fr; }
    .stats-grid { grid-template-columns: 1fr 1fr; }
  }
`;

const CHUNKING_STRATEGIES = [
  { value: "structural", label: "Structural" },
  { value: "semantic", label: "Semantic" },
  { value: "sentence_window", label: "Sentence Window" },
  { value: "parent_document", label: "Parent Document" },
  { value: "fixed_size", label: "Fixed Size" },
];

// ── PDF / file text extraction ────────────────────────────────────────────────
// Only change from the original: this function replaces the bare `uploadFile.text()` call.
// All request/response structures remain identical.
const extractFileContent = async (file) => {
  const isPDF = file.type === "application/pdf" || file.name.toLowerCase().endsWith(".pdf");

  if (isPDF) {
    const arrayBuffer = await file.arrayBuffer();
    const pdf = await pdfjsLib.getDocument({ data: arrayBuffer }).promise;
    let fullText = "";

    for (let pageNum = 1; pageNum <= pdf.numPages; pageNum++) {
      const page = await pdf.getPage(pageNum);
      const textContent = await page.getTextContent();

      // Preserve line breaks by tracking Y position of each text item
      let lastY = null;
      let pageText = "";
      for (const item of textContent.items) {
        if (lastY !== null && Math.abs(item.transform[5] - lastY) > 5) {
          pageText += "\n";
        }
        pageText += item.str;
        lastY = item.transform[5];
      }
      fullText += pageText + "\n\n";
    }

    return fullText.trim();
  }

  // For all other file types (.txt, .md, .json, .csv, .html)
  return await file.text();
};

export default function App() {
  const [page, setPage] = useState("query");
  const [collections, setCollections] = useState([]);
  const [selectedCollection, setSelectedCollection] = useState("");
  const [health, setHealth] = useState(null);
  const [loading, setLoading] = useState({});
  const [alerts, setAlerts] = useState({});

  // Collection form
  const [newCollName, setNewCollName] = useState("");

  // Document form
  const [docForm, setDocForm] = useState({
    collection_name: "",
    source: "",
    doc_type: "resume",
    content: "",
    strategy: "structural",
    fixed_size: 500,
    overlap: 50,
    sentence_window_size: 3,
    min_chunk_size: 100,
    max_chunk_size: 2000,
    preserve_paragraphs: true,
    extract_keywords: true,
  });
  const [uploadFile, setUploadFile] = useState(null);
  const [dragOver, setDragOver] = useState(false);
  const [inputMode, setInputMode] = useState("file"); // "file" | "text"
  const [parsingFile, setParsingFile] = useState(false); // shows spinner while PDF is being parsed

  // Query form
  const [queryForm, setQueryForm] = useState({
    collection_name: "",
    query: "",
    top_k: 5,
    reranker_enabled: true,
    include_parents: false,
    query_expansion: true,
    semantic_threshold: 0,
  });
  const [queryResult, setQueryResult] = useState(null);

  // Search form
  const [searchForm, setSearchForm] = useState({
    collection_name: "",
    query: "",
    top_k: 5,
  });
  const [searchResult, setSearchResult] = useState(null);

  // Collection stats
  const [collStats, setCollStats] = useState(null);

  const fileInputRef = useRef();

  const setLoad = (key, val) => setLoading(p => ({ ...p, [key]: val }));
  const showAlert = (key, type, msg) => {
    setAlerts(p => ({ ...p, [key]: { type, msg } }));
    setTimeout(() => setAlerts(p => { const n = { ...p }; delete n[key]; return n; }), 5000);
  };

  const req = async (method, path, body, isJson = true) => {
    const opts = { method, headers: {} };
    if (body && isJson) { opts.headers["Content-Type"] = "application/json"; opts.body = JSON.stringify(body); }
    const res = await fetch(`${API_BASE}${path}`, opts);
    const data = await res.json().catch(() => ({}));
    if (!res.ok) throw new Error(data.error || data.message || `HTTP ${res.status}`);
    return data;
  };

  const loadCollections = async () => {
    try {
      const data = await req("GET", "/collections");
      setCollections(data.collections || []);
    } catch (e) { console.error(e); }
  };

  useEffect(() => {
    loadCollections();
    fetch("http://localhost:8080/health").then(r => r.json()).then(() => setHealth("ok")).catch(() => setHealth("err"));
  }, []);

  // Sync selected collection to forms
  useEffect(() => {
    if (selectedCollection) {
      setDocForm(p => ({ ...p, collection_name: selectedCollection }));
      setQueryForm(p => ({ ...p, collection_name: selectedCollection }));
      setSearchForm(p => ({ ...p, collection_name: selectedCollection }));
    }
  }, [selectedCollection]);

  // CREATE COLLECTION
  const createCollection = async () => {
    if (!newCollName.trim()) return showAlert("coll", "error", "Collection name required");
    setLoad("coll", true);
    try {
      await req("POST", "/collections", { name: newCollName.trim() });
      showAlert("coll", "success", `Collection "${newCollName}" created`);
      setNewCollName("");
      await loadCollections();
    } catch (e) { showAlert("coll", "error", e.message); }
    setLoad("coll", false);
  };

  // DELETE COLLECTION
  const deleteCollection = async (name) => {
    if (!window.confirm(`Delete collection "${name}"?`)) return;
    try {
      await req("DELETE", `/collections/${name}`);
      if (selectedCollection === name) setSelectedCollection("");
      await loadCollections();
    } catch (e) { showAlert("coll", "error", e.message); }
  };

  // GET COLLECTION STATS
  const getCollStats = async (name) => {
    try {
      const data = await req("GET", `/collections/${name}`);
      setCollStats(data);
    } catch (e) { showAlert("stats", "error", e.message); }
  };

  // FILE DROP — just stores the File object, no reading yet
  const handleFileDrop = (e) => {
    e.preventDefault();
    setDragOver(false);
    const file = e.dataTransfer?.files?.[0] || e.target.files?.[0];
    if (file) setUploadFile(file);
  };

  // ADD DOCUMENT
  // ONLY change vs original: `uploadFile.text()` → `extractFileContent(uploadFile)`
  // + parsingFile spinner state for PDFs. Request body is identical.
  const addDocument = async () => {
    if (!docForm.collection_name) return showAlert("doc", "error", "Select a collection");
    setLoad("doc", true);
    try {
      let content = "";

      if (inputMode === "file" && uploadFile) {
        const isPDF = uploadFile.type === "application/pdf" || uploadFile.name.toLowerCase().endsWith(".pdf");
        if (isPDF) setParsingFile(true);
        content = await extractFileContent(uploadFile); // ← was: uploadFile.text()
        setParsingFile(false);
      } else if (inputMode === "text") {
        content = docForm.content;
      } else {
        throw new Error("Provide a file or paste content");
      }

      if (!content.trim()) throw new Error("Content is empty — PDF may have no extractable text");

      // Request body — UNCHANGED from original
      const body = {
        collection_name: docForm.collection_name,
        content,
        source: docForm.source || (uploadFile ? uploadFile.name : "direct-input"),
        doc_type: docForm.doc_type,
        chunking_config: {
          strategy: docForm.strategy,
          fixed_size: Number(docForm.fixed_size),
          overlap: Number(docForm.overlap),
          sentence_window_size: Number(docForm.sentence_window_size),
          min_chunk_size: Number(docForm.min_chunk_size),
          max_chunk_size: Number(docForm.max_chunk_size),
          preserve_paragraphs: docForm.preserve_paragraphs,
          extract_keywords: docForm.extract_keywords,
        },
      };

      await req("POST", "/documents", body);
      showAlert("doc", "success", "Document added and processed successfully!");
      setUploadFile(null);
      setDocForm(p => ({ ...p, content: "", source: "" }));
      await loadCollections();
    } catch (e) {
      setParsingFile(false);
      showAlert("doc", "error", e.message);
    }
    setLoad("doc", false);
  };

  // QUERY — request/response UNCHANGED
  const runQuery = async () => {
    if (!queryForm.collection_name) return showAlert("query", "error", "Select a collection");
    if (!queryForm.query.trim()) return showAlert("query", "error", "Enter a query");
    setLoad("query", true);
    setQueryResult(null);
    try {
      const body = {
        collection_name: queryForm.collection_name,
        query: queryForm.query,
        top_k: Number(queryForm.top_k),
        reranker_enabled: queryForm.reranker_enabled,
        include_parents: queryForm.include_parents,
        query_expansion: queryForm.query_expansion,
        semantic_threshold: Number(queryForm.semantic_threshold),
      };
      const data = await req("POST", "/query", body);
      setQueryResult(data);
    } catch (e) { showAlert("query", "error", e.message); }
    setLoad("query", false);
  };

  // SEARCH — request/response UNCHANGED
  const runSearch = async () => {
    if (!searchForm.collection_name) return showAlert("search", "error", "Select a collection");
    if (!searchForm.query.trim()) return showAlert("search", "error", "Enter a query");
    setLoad("search", true);
    setSearchResult(null);
    try {
      const body = {
        collection_name: searchForm.collection_name,
        query: searchForm.query,
        top_k: Number(searchForm.top_k),
      };
      const data = await req("POST", "/search", body);
      setSearchResult(data);
    } catch (e) { showAlert("search", "error", e.message); }
    setLoad("search", false);
  };

  const AlertBox = ({ id }) => alerts[id] ? (
    <div className={`alert alert-${alerts[id].type}`}>
      <span>{alerts[id].type === "success" ? "✓" : alerts[id].type === "error" ? "✕" : "ℹ"}</span>
      {alerts[id].msg}
    </div>
  ) : null;

  const Toggle = ({ val, onChange }) => (
    <div className={`toggle ${val ? "on" : ""}`} onClick={() => onChange(!val)} />
  );

  return (
    <>
      <style>{styles}</style>
      <div className="app">
        {/* SIDEBAR */}
        <aside className="sidebar">
          <div className="sidebar-logo">
            <div className="logo-title">RAG System</div>
            <div className="logo-sub">v1.0 · localhost:8080</div>
          </div>

          <div className="nav-section">
            <div className="nav-label">Navigation</div>
            {[
              { id: "query", icon: "⚡", label: "Query" },
              { id: "search", icon: "🔍", label: "Search" },
              { id: "documents", icon: "📄", label: "Add Document" },
              { id: "collections", icon: "🗂", label: "Collections" },
            ].map(n => (
              <button key={n.id} className={`nav-btn ${page === n.id ? "active" : ""}`} onClick={() => setPage(n.id)}>
                <span className="icon">{n.icon}</span>
                {n.label}
              </button>
            ))}
          </div>

          <div className="nav-section" style={{ flex: 1 }}>
            <div className="nav-label">Collections</div>
            <div className="collections-list">
              {collections.length === 0
                ? <div className="no-collections">No collections yet</div>
                : collections.map(c => {
                  const name = typeof c === "string" ? c : c.name || c;
                  return (
                    <div
                      key={name}
                      className={`collection-item ${selectedCollection === name ? "selected" : ""}`}
                      onClick={() => setSelectedCollection(name)}
                    >
                      <div className="collection-dot" />
                      <span className="collection-name">{name}</span>
                      <button className="collection-del" onClick={e => { e.stopPropagation(); deleteCollection(name); }}>✕</button>
                    </div>
                  );
                })
              }
            </div>
          </div>

          <div style={{ padding: "12px", borderTop: "1px solid var(--border)" }}>
            <div className="status-text">
              <span className="status-dot" style={{ background: health === "ok" ? "var(--green)" : health === "err" ? "var(--red)" : "var(--yellow)" }} />
              {health === "ok" ? "Server online" : health === "err" ? "Server offline" : "Checking..."}
            </div>
          </div>
        </aside>

        {/* MAIN */}
        <main className="main">
          <div className="topbar">
            <div className="page-title">
              {page === "query" && "⚡ Query Documents"}
              {page === "search" && "🔍 Semantic Search"}
              {page === "documents" && "📄 Add Document"}
              {page === "collections" && "🗂 Manage Collections"}
            </div>
            <button className="btn btn-ghost" style={{ fontSize: 11 }} onClick={loadCollections}>↻ Refresh</button>
          </div>

          <div className="content">

            {/* ── QUERY PAGE ── */}
            {page === "query" && (
              <>
                <div className="card">
                  <div className="card-title">⚡ Query Configuration</div>
                  <div className="form-row">
                    <div>
                      <label>Collection</label>
                      <select value={queryForm.collection_name} onChange={e => setQueryForm(p => ({ ...p, collection_name: e.target.value }))}>
                        <option value="">— Select Collection —</option>
                        {collections.map(c => { const n = typeof c === "string" ? c : c.name || c; return <option key={n} value={n}>{n}</option>; })}
                      </select>
                    </div>
                    <div>
                      <label>Top K Results</label>
                      <input type="number" min={1} max={20} value={queryForm.top_k} onChange={e => setQueryForm(p => ({ ...p, top_k: e.target.value }))} />
                    </div>
                  </div>
                  <div className="form-row single">
                    <div>
                      <label>Query</label>
                      <textarea value={queryForm.query} onChange={e => setQueryForm(p => ({ ...p, query: e.target.value }))}
                        placeholder="Ask anything about your documents..." rows={3}
                        onKeyDown={e => { if (e.key === "Enter" && e.ctrlKey) runQuery(); }} />
                    </div>
                  </div>
                  <div style={{ marginBottom: 16 }}>
                    {[
                      { key: "reranker_enabled", label: "Re-ranking", sub: "Improve result ordering" },
                      { key: "query_expansion", label: "Query Expansion", sub: "Add synonyms & related terms" },
                      { key: "include_parents", label: "Include Parents", sub: "Fetch parent chunks for context" },
                    ].map(opt => (
                      <div className="toggle-row" key={opt.key}>
                        <div>
                          <div className="toggle-label">{opt.label}</div>
                          <div className="toggle-sub">{opt.sub}</div>
                        </div>
                        <Toggle val={queryForm[opt.key]} onChange={v => setQueryForm(p => ({ ...p, [opt.key]: v }))} />
                      </div>
                    ))}
                    <div className="toggle-row">
                      <div>
                        <div className="toggle-label">Semantic Threshold</div>
                        <div className="toggle-sub">Min similarity score (0 = disabled)</div>
                      </div>
                      <input type="number" step={0.05} min={0} max={1}
                        value={queryForm.semantic_threshold}
                        onChange={e => setQueryForm(p => ({ ...p, semantic_threshold: e.target.value }))}
                        style={{ width: 80 }} />
                    </div>
                  </div>
                  <div style={{ display: "flex", gap: 10 }}>
                    <button className="btn btn-primary" onClick={runQuery} disabled={loading.query}>
                      {loading.query ? <><span className="spinner" /> Running…</> : "⚡ Run Query"}
                    </button>
                  </div>
                  <AlertBox id="query" />
                </div>

                {queryResult && (
                  <div className="card">
                    <div className="card-title">Answer</div>
                    <div className="answer-box">{queryResult.answer}</div>
                    {queryResult.processing_time && (
                      <div className="proc-time">⏱ {queryResult.processing_time.toFixed(3)}s · {queryResult.metadata_used ? "Filters applied" : "No filters"}</div>
                    )}

                    {queryResult.enhanced_chunks?.length > 0 && (
                      <>
                        <div className="card-title" style={{ marginTop: 20 }}>Retrieved Chunks ({queryResult.enhanced_chunks.length})</div>
                        <div className="chunks-grid">
                          {queryResult.enhanced_chunks.map((chunk, i) => (
                            <div className="chunk-card" key={chunk.id || i}>
                              <div className="chunk-header">
                                <span className="chunk-badge">#{i + 1}</span>
                                <span className="chunk-badge">{chunk.chunk_type}</span>
                                {chunk.section && <span className="chunk-badge section">{chunk.section}</span>}
                                {queryResult.similarity_scores?.[i] && (
                                  <span className="chunk-badge score">{(queryResult.similarity_scores[i] * 100).toFixed(1)}%</span>
                                )}
                                {queryResult.reranked_scores?.[i] && (
                                  <span className="chunk-badge score">re: {(queryResult.reranked_scores[i] * 100).toFixed(1)}%</span>
                                )}
                              </div>
                              <div className="chunk-text">{chunk.text?.slice(0, 300)}{chunk.text?.length > 300 ? "…" : ""}</div>
                              {chunk.keywords?.length > 0 && (
                                <div style={{ marginTop: 6, display: "flex", flexWrap: "wrap", gap: 4 }}>
                                  {chunk.keywords.slice(0, 6).map(k => (
                                    <span key={k} style={{ background: "rgba(167,139,250,0.08)", color: "var(--accent2)", padding: "1px 6px", borderRadius: 3, fontSize: 10, fontFamily: "var(--mono)" }}>{k}</span>
                                  ))}
                                </div>
                              )}
                            </div>
                          ))}
                        </div>
                      </>
                    )}
                  </div>
                )}
              </>
            )}

            {/* ── SEARCH PAGE ── */}
            {page === "search" && (
              <>
                <div className="card">
                  <div className="card-title">🔍 Semantic Search</div>
                  <div className="form-row">
                    <div>
                      <label>Collection</label>
                      <select value={searchForm.collection_name} onChange={e => setSearchForm(p => ({ ...p, collection_name: e.target.value }))}>
                        <option value="">— Select Collection —</option>
                        {collections.map(c => { const n = typeof c === "string" ? c : c.name || c; return <option key={n} value={n}>{n}</option>; })}
                      </select>
                    </div>
                    <div>
                      <label>Top K</label>
                      <input type="number" min={1} max={20} value={searchForm.top_k} onChange={e => setSearchForm(p => ({ ...p, top_k: e.target.value }))} />
                    </div>
                  </div>
                  <div className="form-row single">
                    <div>
                      <label>Search Query</label>
                      <textarea value={searchForm.query} onChange={e => setSearchForm(p => ({ ...p, query: e.target.value }))}
                        placeholder="Search for relevant chunks..." rows={2}
                        onKeyDown={e => { if (e.key === "Enter" && e.ctrlKey) runSearch(); }} />
                    </div>
                  </div>
                  <button className="btn btn-primary" onClick={runSearch} disabled={loading.search}>
                    {loading.search ? <><span className="spinner" /> Searching…</> : "🔍 Search"}
                  </button>
                  <AlertBox id="search" />
                </div>

                {searchResult && (
                  <div className="card">
                    <div className="card-title">Results</div>
                    {searchResult.answer && <div className="answer-box" style={{ marginBottom: 16 }}>{searchResult.answer}</div>}
                    {searchResult.enhanced_chunks?.length > 0 && (
                      <div className="chunks-grid">
                        {searchResult.enhanced_chunks.map((chunk, i) => (
                          <div className="chunk-card" key={chunk.id || i}>
                            <div className="chunk-header">
                              <span className="chunk-badge">#{i + 1}</span>
                              <span className="chunk-badge">{chunk.chunk_type}</span>
                              {chunk.section && <span className="chunk-badge section">{chunk.section}</span>}
                              {searchResult.similarity_scores?.[i] && (
                                <span className="chunk-badge score">{(searchResult.similarity_scores[i] * 100).toFixed(1)}%</span>
                              )}
                            </div>
                            <div className="chunk-text">{chunk.text?.slice(0, 300)}{chunk.text?.length > 300 ? "…" : ""}</div>
                          </div>
                        ))}
                      </div>
                    )}
                  </div>
                )}
              </>
            )}

            {/* ── ADD DOCUMENT PAGE ── */}
            {page === "documents" && (
              <div className="card">
                <div className="card-title">📄 Add Document</div>

                <div className="form-row">
                  <div>
                    <label>Collection *</label>
                    <select value={docForm.collection_name} onChange={e => setDocForm(p => ({ ...p, collection_name: e.target.value }))}>
                      <option value="">— Select Collection —</option>
                      {collections.map(c => { const n = typeof c === "string" ? c : c.name || c; return <option key={n} value={n}>{n}</option>; })}
                    </select>
                  </div>
                  <div>
                    <label>Document Type</label>
                    <select value={docForm.doc_type} onChange={e => setDocForm(p => ({ ...p, doc_type: e.target.value }))}>
                      {["resume", "article", "bible", "technical", "legal", "general"].map(t => <option key={t} value={t}>{t}</option>)}
                    </select>
                  </div>
                </div>

                <div className="form-row single" style={{ marginBottom: 14 }}>
                  <div>
                    <label>Source / Filename (optional)</label>
                    <input type="text" value={docForm.source} onChange={e => setDocForm(p => ({ ...p, source: e.target.value }))} placeholder="my-document.txt" />
                  </div>
                </div>

                {/* Input mode tabs */}
                <div className="tabs" style={{ marginBottom: 14 }}>
                  <button className={`tab ${inputMode === "file" ? "active" : ""}`} onClick={() => setInputMode("file")}>📁 Upload File</button>
                  <button className={`tab ${inputMode === "text" ? "active" : ""}`} onClick={() => setInputMode("text")}>✏️ Paste Text</button>
                </div>

                {inputMode === "file" ? (
                  <div
                    className={`file-drop ${dragOver ? "drag-over" : ""}`}
                    onDragOver={e => { e.preventDefault(); setDragOver(true); }}
                    onDragLeave={() => setDragOver(false)}
                    onDrop={handleFileDrop}
                  >
                    <input type="file" accept=".txt,.md,.pdf,.json,.csv,.html" onChange={handleFileDrop} ref={fileInputRef} />
                    <div className="file-drop-icon">📂</div>
                    <div className="file-drop-text">Drag & drop a file or <span>browse</span></div>
                    <div style={{ fontSize: 11, color: "var(--muted)", marginTop: 4 }}>.txt · .md · .pdf · .json · .csv · .html</div>
                    {/* Show filename when selected, spinner when parsing PDF */}
                    {uploadFile && !parsingFile && (
                      <div className="file-selected">
                        ✓ {uploadFile.name} ({(uploadFile.size / 1024).toFixed(1)} KB)
                        {(uploadFile.type === "application/pdf" || uploadFile.name.toLowerCase().endsWith(".pdf")) && (
                          <span style={{ color: "var(--accent)", marginLeft: 6 }}>· PDF (text extracted on upload)</span>
                        )}
                      </div>
                    )}
                    {parsingFile && (
                      <div className="file-parsing">
                        <span className="spinner" /> Extracting text from PDF…
                      </div>
                    )}
                  </div>
                ) : (
                  <div>
                    <label>Document Content *</label>
                    <textarea value={docForm.content} onChange={e => setDocForm(p => ({ ...p, content: e.target.value }))}
                      placeholder="Paste your document content here..." rows={8} />
                  </div>
                )}

                {/* Chunking config */}
                <div style={{ marginTop: 20 }}>
                  <div className="card-title">Chunking Strategy</div>
                  <div className="form-row three">
                    <div>
                      <label>Strategy</label>
                      <select value={docForm.strategy} onChange={e => setDocForm(p => ({ ...p, strategy: e.target.value }))}>
                        {CHUNKING_STRATEGIES.map(s => <option key={s.value} value={s.value}>{s.label}</option>)}
                      </select>
                    </div>
                    <div>
                      <label>Fixed Size</label>
                      <input type="number" value={docForm.fixed_size} onChange={e => setDocForm(p => ({ ...p, fixed_size: e.target.value }))} />
                    </div>
                    <div>
                      <label>Overlap</label>
                      <input type="number" value={docForm.overlap} onChange={e => setDocForm(p => ({ ...p, overlap: e.target.value }))} />
                    </div>
                  </div>
                  <div className="form-row three">
                    <div>
                      <label>Sentence Window</label>
                      <input type="number" value={docForm.sentence_window_size} onChange={e => setDocForm(p => ({ ...p, sentence_window_size: e.target.value }))} />
                    </div>
                    <div>
                      <label>Min Chunk Size</label>
                      <input type="number" value={docForm.min_chunk_size} onChange={e => setDocForm(p => ({ ...p, min_chunk_size: e.target.value }))} />
                    </div>
                    <div>
                      <label>Max Chunk Size</label>
                      <input type="number" value={docForm.max_chunk_size} onChange={e => setDocForm(p => ({ ...p, max_chunk_size: e.target.value }))} />
                    </div>
                  </div>
                  <div className="toggle-row">
                    <div>
                      <div className="toggle-label">Preserve Paragraphs</div>
                      <div className="toggle-sub">Keep paragraph boundaries intact</div>
                    </div>
                    <Toggle val={docForm.preserve_paragraphs} onChange={v => setDocForm(p => ({ ...p, preserve_paragraphs: v }))} />
                  </div>
                  <div className="toggle-row">
                    <div>
                      <div className="toggle-label">Extract Keywords</div>
                      <div className="toggle-sub">Auto-extract keywords from chunks</div>
                    </div>
                    <Toggle val={docForm.extract_keywords} onChange={v => setDocForm(p => ({ ...p, extract_keywords: v }))} />
                  </div>
                </div>

                <div style={{ marginTop: 18 }}>
                  <button className="btn btn-primary" onClick={addDocument} disabled={loading.doc || parsingFile}>
                    {loading.doc ? <><span className="spinner" /> Processing…</> : "📤 Add Document"}
                  </button>
                </div>
                <AlertBox id="doc" />
              </div>
            )}

            {/* ── COLLECTIONS PAGE ── */}
            {page === "collections" && (
              <>
                <div className="card">
                  <div className="card-title">➕ Create Collection</div>
                  <div className="form-row">
                    <div>
                      <label>Collection Name</label>
                      <input type="text" value={newCollName} onChange={e => setNewCollName(e.target.value)}
                        placeholder="e.g. resumes, knowledge-base"
                        onKeyDown={e => { if (e.key === "Enter") createCollection(); }} />
                    </div>
                    <div style={{ display: "flex", alignItems: "flex-end" }}>
                      <button className="btn btn-primary" onClick={createCollection} disabled={loading.coll}>
                        {loading.coll ? <><span className="spinner" /> Creating…</> : "Create"}
                      </button>
                    </div>
                  </div>
                  <AlertBox id="coll" />
                </div>

                <div className="card">
                  <div className="card-title">🗂 All Collections ({collections.length})</div>
                  {collections.length === 0 ? (
                    <div style={{ color: "var(--muted)", fontSize: 13, fontStyle: "italic" }}>No collections yet. Create one above.</div>
                  ) : (
                    <div style={{ display: "grid", gap: 10 }}>
                      {collections.map(c => {
                        const name = typeof c === "string" ? c : c.name || c;
                        return (
                          <div key={name} style={{ background: "var(--bg)", border: "1px solid var(--border)", borderRadius: 8, padding: "12px 16px", display: "flex", alignItems: "center", justifyContent: "space-between" }}>
                            <div style={{ display: "flex", alignItems: "center", gap: 10 }}>
                              <span style={{ fontSize: 16 }}>🗂</span>
                              <span style={{ fontFamily: "var(--mono)", fontSize: 13 }}>{name}</span>
                            </div>
                            <div style={{ display: "flex", gap: 8 }}>
                              <button className="btn btn-ghost" style={{ fontSize: 11, padding: "5px 12px" }} onClick={() => { getCollStats(name); setSelectedCollection(name); }}>Stats</button>
                              <button className="btn btn-danger" style={{ fontSize: 11, padding: "5px 12px" }} onClick={() => deleteCollection(name)}>Delete</button>
                            </div>
                          </div>
                        );
                      })}
                    </div>
                  )}
                </div>

                {collStats && (
                  <div className="card">
                    <div className="card-title">📊 Collection Stats — {selectedCollection}</div>
                    <div className="stats-grid">
                      <div className="stat-card">
                        <div className="stat-val">{collStats.total_documents ?? collStats.document_count ?? "—"}</div>
                        <div className="stat-label">Documents</div>
                      </div>
                      <div className="stat-card">
                        <div className="stat-val">{collStats.total_chunks ?? collStats.chunk_count ?? "—"}</div>
                        <div className="stat-label">Chunks</div>
                      </div>
                      <div className="stat-card">
                        <div className="stat-val">{collStats.name ?? selectedCollection}</div>
                        <div className="stat-label">Name</div>
                      </div>
                    </div>
                    <AlertBox id="stats" />
                  </div>
                )}
              </>
            )}

          </div>
        </main>
      </div>
    </>
  );
}