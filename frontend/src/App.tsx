import React, { useMemo, useState } from "react"

const API = (import.meta as any).env?.VITE_API_URL || "http://localhost:8080"

type Stats = {
  id: number
  token: string
  long_url: string
  clicks: number
  created_at: string
  updated_at: string
  last_accessed?: string | null
}

export default function App() {
  // --- state utama ---
  const [url, setUrl] = useState("")
  const [token, setToken] = useState("")
  const [stats, setStats] = useState<Stats | null>(null)

  // --- state bantuan untuk negative tests & utilities ---
  const [error, setError] = useState("")
  const [checkingToken, setCheckingToken] = useState("")
  const [busy, setBusy] = useState(false)

  const shortUrl = useMemo(() => (token ? `${window.location.origin}/${token}` : ""), [token])

  // === TC-1: Generate Short URL (Positive & Negative) ===
  async function onShorten(e: React.FormEvent) {
    e.preventDefault()
    setError("")
    setStats(null)
    setToken("")
    setBusy(true)
    try {
      const res = await fetch(`${API}/shorten`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ url }),
      })
      const data = await safeJson(res)
      if (res.ok) {
        setToken(data.token)
      } else {
        // tampilkan error di UI (400 / 422)
        setError(data?.error || `Failed with status ${res.status}`)
      }
    } catch (err: any) {
      setError(err?.message || "Network error")
    } finally {
      setBusy(false)
    }
  }

  // === TC-4: Load Stats (Positive & Negative) ===
  async function loadStats(forToken?: string) {
    const t = (forToken ?? token).trim()
    if (!t) {
      setError("Token is empty")
      return
    }
    setBusy(true)
    setError("")
    try {
      const res = await fetch(`${API}/stats/${t}`)
      const data = await safeJson(res)
      if (res.ok) setStats(data as Stats)
      else setError(data?.error || `Not found (status ${res.status})`)
    } catch (err: any) {
      setError(err?.message || "Network error")
    } finally {
      setBusy(false)
    }
  }

  // === TC-2: Redirect Test (cek status & Location tanpa auto-follow) ===
  async function testRedirect() {
    const t = checkingToken.trim()
    if (!t) {
      setError("Token to check is empty")
      return
    }
    setBusy(true)
    setError("")
    try {
      const res = await fetch(`${API}/${t}`, { redirect: "manual" as RequestRedirect })
      // Positive: 301/302 + Location header
      if (res.status === 301 || res.status === 302) {
        const loc = res.headers.get("Location")
        alert(`Redirect OK (${res.status}) → ${loc}`)
      } else if (res.status === 404) {
        // Negative: unknown token -> 404
        const data = await safeJson(res)
        alert(`Not Found (404). ${data?.error ?? ""}`.trim())
      } else {
        alert(`Unexpected status: ${res}`)
      }
    } catch (err: any) {
      setError(err?.message || "Network error")
    } finally {
      setBusy(false)
    }
  }

  // === Utility: copy URL ===
  async function copyShort() {
    if (!shortUrl) return
    try {
      await navigator.clipboard.writeText(shortUrl)
      toast("Copied!")
    } catch {
      toast("Copy failed")
    }
  }

  // === Utility helpers ===
  function toast(msg: string) {
    alert(msg)
  }
  async function safeJson(res: Response) {
    try {
      return await res.json()
    } catch {
      return null
    }
  }

  return (
    <div style={container}>
      <h1 style={{ marginBottom: 8 }}>URL Shortener</h1>
      <p style={{ marginTop: 0, opacity: 0.8 }}>
        Backend: <code>{API}</code>
      </p>

      {/* TC-1: Form shorten */}
      <form onSubmit={onShorten} style={row}>
        <input
          style={input}
          placeholder="https://example.com/very/long/url"
          value={url}
          onChange={(e) => setUrl(e.target.value)}
        />
        <button disabled={busy} style={button}>
          {busy ? "Processing..." : "Shorten"}
        </button>
      </form>

      {/* tampilkan error untuk negative cases (400, 422, 404, dll) */}
      {!!error && <p style={{ color: "crimson", marginTop: 8 }}>{error}</p>}

      {/* Hasil shorten + aksi lanjut */}
      {token && (
        <div style={{ ...card, marginTop: 16 }}>
          <p style={{ margin: 0 }}>
            <b>Short URL:</b>{" "}
            <a href={`/${token}`} target="_blank" rel="noreferrer">
              {shortUrl}
            </a>
          </p>
          <div style={{ display: "flex", gap: 8, marginTop: 8 }}>
            <button onClick={copyShort} style={button}>
              Copy
            </button>
            <button onClick={() => loadStats()} style={button}>
              Load Stats
            </button>
          </div>

          {/* TC-4: tampilkan stats jelas */}
          {stats && (
            <div style={{ marginTop: 12, lineHeight: 1.6 }}>
              <div>
                <b>Token:</b> {stats.token}
              </div>
              <div>
                <b>Long URL:</b>{" "}
                <a href={stats.long_url} target="_blank" rel="noreferrer">
                  {stats.long_url}
                </a>
              </div>
              <div>
                <b>Clicks:</b> {stats.clicks}
              </div>
              <div>
                <b>Created:</b> {date(stats.created_at)}
              </div>
              {stats.last_accessed && (
                <div>
                  <b>Last Accessed:</b> {date(stats.last_accessed)}
                </div>
              )}
              <details style={{ marginTop: 8 }}>
                <summary>Raw JSON</summary>
                <pre style={pre}>{JSON.stringify(stats, null, 2)}</pre>
              </details>
            </div>
          )}
        </div>
      )}

      {/* TC-2: Manual Redirect Test (tanpa auto-follow) */}
      <div style={{ ...card, marginTop: 20 }}>
        <h3 style={{ marginTop: 0 }}>Redirect to Long URL</h3>
        <div style={row}>
          <input
            style={input}
            placeholder="enter short token"
            value={checkingToken}
            onChange={(e) => setCheckingToken(e.target.value)}
          />
          <button onClick={testRedirect} style={button} disabled={busy}>
            Check
          </button>
          <button onClick={() => loadStats(checkingToken)} style={button} disabled={busy}>
            Load Stats by Token
          </button>
        </div>
        <p style={{ marginTop: 8, opacity: 0.8 }}>
          Expected: <b>301/302</b> with <b>Location</b> header (positive). Unknown token → <b>404</b> (negative).
        </p>
      </div>
    </div>
  )
}

/* ------- styles (inline sederhana biar tanpa dependency) ------- */
const container: React.CSSProperties = {
  maxWidth: 760,
  margin: "40px auto",
  fontFamily:
    'ui-sans-serif, system-ui, -apple-system, Segoe UI, Roboto, "Helvetica Neue", Arial, "Noto Sans", "Apple Color Emoji", "Segoe UI Emoji"',
  padding: "0 16px",
}
const row: React.CSSProperties = { display: "flex", gap: 8, alignItems: "center" }
const input: React.CSSProperties = {
  flex: 1,
  padding: "10px 12px",
  border: "1px solid #d0d0d0",
  borderRadius: 8,
  outline: "none",
}
const button: React.CSSProperties = {
  padding: "10px 14px",
  borderRadius: 8,
  border: "1px solid #d0d0d0",
  background: "#f8f8f8",
  cursor: "pointer",
}
const card: React.CSSProperties = {
  border: "1px solid #e5e5e5",
  borderRadius: 12,
  padding: 16,
  background: "#fff",
}
const pre: React.CSSProperties = {
  background: "#f7f7f7",
  borderRadius: 8,
  padding: 12,
  overflowX: "auto",
  fontSize: 13,
}

/* ------- small helper ------- */
function date(s?: string | null) {
  if (!s) return "-"
  try {
    const d = new Date(s)
    return isNaN(d.getTime()) ? s : d.toLocaleString()
  } catch {
    return s
  }
}
