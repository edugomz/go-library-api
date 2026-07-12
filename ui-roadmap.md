# UI Roadmap

Server-rendered frontend using Go `html/template` (Gin) + htmx for interactivity.
No build step, no JS framework.
Extends the existing `internal/handlers/web` / `internal/views` pattern.

---

## ✅ Built (this pass)

* [x] Shared nav partial (`internal/views/nav.html`) included on every page
* [x] Single stylesheet (`static/style.css`), served via `r.Static("/static", "./static")`
* [x] htmx (`unpkg.com/htmx.org@2.0.4` CDN `<script>` tag), used for login/register form submits
* [x] Book list (`/books`) and book detail (`/books/:id`, shows author link)
* [x] Author list (`/authors`) and author detail (`/authors/:id`, shows their books)
* [x] Register page/form (`/register`) wired to `AuthService.Register`
* [x] Login page/form (`/login`) wired to `AuthService.Login`
* [x] Logout (`POST /logout`)

---

## 🍪 Auth storage tradeoff

The JWT is stored in an `HttpOnly` cookie (`token`), set directly by the web handler after calling `AuthService.Login`, not by round-tripping through `/api/v1/auth/login` from client JS.

* Chosen because: no XSS-exposed token in `localStorage`, and it reuses the same server-rendered form-post pattern already used elsewhere in this codebase (no client-side JS needed beyond htmx).
* Tradeoff: the cookie is not marked `Secure` yet (dev runs over plain HTTP), and there's no CSRF protection on the state-changing web routes (`/login`, `/register`, `/logout`). Also, the cookie is only used to toggle nav display (`isLoggedIn`); no web route actually checks/requires it yet, since `internal/handlers/web` calls services directly and bypasses the JWT-protected `/api/v1/*` routes entirely.
* Revisit before any of this ships behind a real domain: add `Secure`, add CSRF tokens on mutating forms, and decide whether web routes should actually enforce login.

---

## 🔜 Deferred (pick up in future tasks)

* [ ] Reviews UI (list/add reviews on a book detail page) - `ReviewHandler`/`ReviewService` already exist on the API side
* [ ] Reading lists UI - no service/repository exists yet for this; models only
* [ ] Pagination and filtering on book/author lists (services currently return full `GetAll()` result sets)
* [ ] Create-book / create-author forms in the web UI (API endpoints exist, no web form yet)
* [ ] Client-side validation (currently relies on HTML5 `required`/`minlength` attributes only)
* [ ] Styling polish (currently a single utilitarian stylesheet, no responsive/mobile layout pass)
* [ ] Session/token hardening: `Secure` cookie flag, CSRF protection, token refresh/expiry handling in the UI (currently a stale token just fails silently until re-login)
* [ ] Protecting web routes behind login (currently anyone can browse `/books` and `/authors` without a session, since the web handler talks to services directly)
* [ ] Flash messages / toast pattern for cross-request feedback (e.g. "registered successfully") instead of relying on redirect + htmx inline errors only
