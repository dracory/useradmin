# e2e

End-to-end tests that drive the useradmin panel in a real headless Chrome
via [chromedp](https://github.com/chromedp/chromedp) (CDP). Unlike go-rod,
chromedp has no `leakless.exe` helper, so it avoids the Windows Defender
false positive entirely.

## How it works

- `TestMain` mounts the admin on `httptest.NewServer` (in-process, random
  port) backed by a seeded in-memory SQLite userstore — no binary build
  needed. Roles and groups are enabled; two users are seeded as members
  of the first role and group.
- One shared browser is launched via `chromedp.NewExecAllocator` for the
  whole run; each test opens a fresh tab with a 30s deadline (`newTab`).
- Chrome binary resolution order: `CHROME_BIN` env → Windows Playwright
  `chromium_headless_shell` → WSL-side Windows shell → Linux Playwright
  shell → chromedp's default lookup (Chrome on PATH / install dirs).

## Run

```bash
go test -v ./e2e/     # or: task e2e
```

Note: on WSL, run tests with the Windows Go toolchain (`go.exe`) —
Windows-side Chrome binds 127.0.0.1, which WSL cannot reach (no
WSL→Windows localhost forwarding). The Linux Playwright shell also
requires system libs (libnspr4 etc.) that may be missing.

## Coverage

Every registered controller has at least one test:

- `user-manager` — seeded rows + create-modal flow
- `user-view`, `user-update`, `user-delete`, `user-impersonate`
- `role-manager` — seeded rows + create-modal flow
- `role-view`, `role-update`, `role-members`
- `group-manager`, `group-view`, `group-update`, `group-members`

## Gotchas learned

- Prefer `chromedp.WaitReady` over `chromedp.WaitVisible`: on some pages
  `WaitVisible` polls forever even though the DOM is fully loaded.
- Routing is query-param based: `?controller=user-manager|role-manager|...`
  (see `shared/consts.go`).
- Manager/members tables render via Vue after an AJAX fetch — wait for the
  second `tbody tr` before asserting (the raw v-for template row exists
  in the DOM before Vue mounts).
- `v-model` is compiled away by Vue — it does NOT appear in the rendered
  DOM. Match form inputs by `type` or position (`input[type="email"]`,
  `input[type="text"]:nth-of-type`), never by the `v-model` attribute.
- View pages (`user-view`, `role-view`, `group-view`) are server-rendered
  — `WaitReady("body")` is enough, no AJAX wait needed.
- The impersonate route is only registered when `OnUserImpersonate` is
  configured; setup provides a cookie-setting stub so it is testable.
- chromedp v0.23 rewrote the API (generic `Run[T]`); this suite pins the
  stable v0.14 line with matching `cdproto`/`go-json-experiment` versions.
