Go’s environment is **very different from Python’s `venv`** — Go does **not create per-project virtual environments** by default.

Instead, Go uses:

✅ A global module cache
✅ A global build cache
✅ Environment variables

Let’s go step-by-step for **Windows**.

---

# 🧹 1️⃣ Clear Go Build Cache (compiled binaries)

This removes cached compiled packages:

```bash
go clean -cache
```

✔ Deletes build artifacts
✔ Safe
✔ No effect on source code

---

# 📦 2️⃣ Clear Module Download Cache

This removes downloaded dependencies:

```bash
go clean -modcache
```

✔ Deletes:

```
C:\Users\<you>\go\pkg\mod
```

✔ Forces Go to re-download modules next time

---

# 🗑 3️⃣ Clear Test Cache

```bash
go clean -testcache
```

---

# 💣 4️⃣ Clear EVERYTHING (most aggressive)

```bash
go clean -cache -modcache -testcache
```

---

# 📍 Where caches live on Windows

Typically:

```
Build Cache:
C:\Users\<you>\AppData\Local\go-build

Module Cache:
C:\Users\<you>\go\pkg\mod
```

You can manually delete them too.

---

# 🧰 5️⃣ Manual Deletion (Windows Explorer / CMD)

### Open Run → type:

```
%LocalAppData%\go-build
```

Delete folder.

---

### Module cache:

```
%UserProfile%\go\pkg\mod
```

Delete folder.

---

# 🌍 6️⃣ Reset Go Environment Variables (optional)

Check current Go env:

```bash
go env
```

Important variables:

* `GOPATH`
* `GOROOT`
* `GOMODCACHE`
* `GOCACHE`

---

### Remove custom env vars (PowerShell)

```powershell
[Environment]::SetEnvironmentVariable("GOPATH", $null, "User")
```

Or via:

```
System Properties → Environment Variables
```

---

# 🧨 7️⃣ Completely Remove Go (like deleting Python venv)

If you want a **full reset**:

### ✔ Uninstall Go

Control Panel → Programs → Uninstall Go

---

### ✔ Delete remaining folders

Usually:

```
C:\Program Files\Go
C:\Users\<you>\go
C:\Users\<you>\AppData\Local\go-build
```

---

### ✔ Remove PATH entry

Environment Variables → PATH → Remove Go

---

# ⚠️ Important Difference vs Python venv

| Python venv              | Go                 |
| ------------------------ | ------------------ |
| Per-project isolated env | ❌ Not default      |
| Local dependency copies  | ❌ Global mod cache |
| Activation required      | ❌ No activation    |
| Multiple envs common     | Rare in Go         |

---

# ✅ If you WANT project-level isolation in Go

You can simulate isolation:

```bash
set GOPATH=C:\myproject\gopath
set GOMODCACHE=C:\myproject\modcache
```

Or PowerShell:

```powershell
$env:GOMODCACHE="C:\myproject\modcache"
```

---

# 🎯 Most common cleanup command

For 99% of cases:

```bash
go clean -cache -modcache -testcache
```