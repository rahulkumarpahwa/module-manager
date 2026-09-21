# Client

A small, dependency-free browser client for the Module Manager Go backend.
It calls the exact endpoints from [`postman/module-manager.postman_collection.json`](../postman/module-manager.postman_collection.json):

| Endpoint | Action in the UI |
| --- | --- |
| `GET /api/healthz` | Live health pill in the header |
| `POST /api/set-path` | **Set Path** button (+ auto-run before **Scan** when the field is filled) |
| `POST /api/main/scan?limit=&offset=` | **Scan** button; renders one project per row |
| `POST /api/project/{id}` | Clicking a row loads the full project structure for that project (body: `{ "path": "<project path>" }`) |
| `POST /api/project/infos` | **Fetch All Details** loads every visible row in one request (body: `[ { "id": "...", "path": "..." }, ... ]`) |
| `POST /api/project/open-vscode` | **Open in VS Code** button at the end of every project row |
| `POST /api/project/open-winexplorer` | **Open Explorer** button at the end of every project row |

## Files

- `index.html` – page structure
- `style.css` – styling
- `app.js` – API calls + rendering
- `server.js` – optional dev server (static files + `/api/*` proxy)

## How to run

1. Start the backend first (reads its port from `.env`, default `8778`):

   ```bash
   go run ./cmd/api
   ```

2. Serve the client. **Recommended** – use the bundled proxy server
   (avoids browser CORS restrictions, since the backend does not send
   CORS headers yet):

   ```bash
   cd client
   node server.js
   # client:  http://localhost:8080
   # /api/* -> http://localhost:8778
   ```

3. Open `http://localhost:8080` and:
   - type a folder path (e.g. `C:\Programming\projects`) → **Set Path**
   - click **Scan** → each detected project appears as a row showing the scan-entry
     fields: `id`, `name`, `path`, `modified_date`, `workspace`, `project`, `framework/language`
   - click a row to expand its full details, loaded via `POST /api/project/{id}`
   - click **Fetch All Details** in the summary bar to load every visible row's
     full details in a single batch request to `POST /api/project/infos`
   - click **Open in VS Code** at the end of a row to open that project's folder
     in VS Code (`POST /api/project/open-vscode`)
   - click **Open Explorer** at the end of a row to open that project's folder
     in Windows File Explorer (`POST /api/project/open-winexplorer`)

### Alternative: open the files directly

You can open `client/index.html` straight from disk, but the calls to
`http://localhost:8778` will only succeed if the backend later sends the
appropriate `Access-Control-Allow-Origin` headers (no CORS middleware exists
today). If you change the backend port, update the **Base URL** field in the UI.

## Notes

- The Postman examples show the *intended* response shapes; the Go handlers
  behind some endpoints are still evolving, so `app.js` normalizes both shapes
  (e.g. scan `data.directory_details.entries` **and** top-level `projects`).
- `POST /api/project/{id}` is the project-detail endpoint (single project). It
  expects a JSON body of `{ "path": "<project path>" }` and returns the full
  `types.Project` (envelope `data.projects[0]`) including `under_dir`
  (top-level entries), sizes, package-manager info and git info.
- `POST /api/project/infos` is the bulk variant. It expects a JSON array body of
  `[ { "id": "<project id>", "path": "<project path>" }, ... ]` and returns all
  details at once (`data.projects` is an array of `types.Project`, order not
  guaranteed — the client matches each result back to its row by `id`). The
  **Fetch All Details** button uses this endpoint so the whole table is resolved
  in a single request.
- **Open in VS Code** calls `POST /api/project/open-vscode`, which runs the
  `code` CLI on the backend machine — make sure VS Code's `code` command is
  available in the backend's `PATH` (`code --version`).
- **Open Explorer** calls `POST /api/project/open-winexplorer`, which launches
  `explorer.exe` — it only works when the backend runs on Windows.