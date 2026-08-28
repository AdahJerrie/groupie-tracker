# Groupie Trackers

A server-rendered web application, written in Go, that consumes the public [Groupie Trackers API](https://groupietrackers.herokuapp.com/api) and presents artist, concert-date, and location data as a browsable website.

Built with a deliberate focus on Go fundamentals, correct HTTP semantics, and idiomatic, production-minded backend code — no frontend framework, no external Go dependencies, standard library only.

---

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Data Model](#data-model)
- [How It Works](#how-it-works)
- [Getting Started](#getting-started)
- [Routes](#routes)
- [Error Handling](#error-handling)
- [Roadmap](#roadmap)
- [Author](#author)

---

## Overview

Groupie Trackers fetches four related resources from a public API — **artists**, **locations**, **dates**, and **relation** — merges them into a single browsable dataset, and serves that data as HTML pages rendered on the server. There is no client-side JavaScript and no database: all data is fetched once at startup and held in memory for the lifetime of the process.

The project intentionally avoids frameworks and third-party packages. Every part of the request/response cycle — routing, method checks, error handling, template rendering, and static file serving — is built directly on Go's standard library (`net/http`, `html/template`, `encoding/json`).

## Features

- Home page listing every artist as a card (image, name, year active).
- Individual artist detail page showing:
  - Founding year and first album date
  - Band members
  - Concert locations cross-referenced with dates (via the `relation` resource)
- Server-side HTML rendering with automatic XSS-escaping (`html/template`).
- Static asset serving (CSS) via `http.FileServer`.
- Defensive error handling: the server does not crash on bad input, missing artists, or upstream API failures.

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.22.2 |
| HTTP | `net/http` (standard library, no router framework) |
| Templating | `html/template` |
| Data format | JSON (`encoding/json`) |
| Frontend | Static HTML + CSS (no JS framework) |
| Data source | [Groupie Trackers public API](https://groupietrackers.herokuapp.com/api) |

No external Go modules are required — `go.mod` declares zero dependencies.

## Project Structure

```
groupie-tracker/
├── main.go              # Entry point: data fetching, routing, handlers
├── go.mod               # Module definition (Go 1.22.2, no external deps)
├── templates/
│   ├── index.html       # Home page — artist listing
│   └── artist.html      # Artist detail page
├── static/
│   └── style.css        # Site styling, served at /static/
├── parts/                # Standalone Go fundamentals exercises
│   ├── struct.go         # (structs — not wired into the running app)
│   ├── maps.go            # (maps)
│   └── pointer.go         # (pointers)
└── README.md
```

> **Note on `parts/`:** this package holds isolated learning exercises from earlier phases of the project (structs, maps, pointers) and is not imported by `main.go`. It is kept in the repo as a record of the learning process, not as part of the running application.

## Data Model

The API exposes four resources, joined by a shared numeric `id`:

| Resource | Key fields | Purpose |
|---|---|---|
| `artists` | `id`, `image`, `name`, `members []string`, `creationDate int`, `firstAlbum string` | Core artist/band info |
| `locations` | `id`, `locations []string` | Where each artist has played |
| `dates` | `id`, `dates []string` | When each artist has played |
| `relation` | `id`, `datesLocations map[string][]string` | Joins each location to its list of dates for that artist |

The `relation` resource is what makes the other three usable together — it's the map that answers "which dates did this artist play in this city?" for every location.

```go
type Artist struct {
    ID           int      `json:"id"`
    Image        string   `json:"image"`
    Name         string   `json:"name"`
    Members      []string `json:"members"`
    CreationDate int      `json:"creationDate"`
    FirstAlbum   string   `json:"firstAlbum"`
}

type Relation struct {
    ID             int                 `json:"id"`
    DatesLocations map[string][]string `json:"datesLocations"`
}
```

## How It Works

1. **Startup fetch:** on boot, `main()` calls `FetchArtists`, `FetchLocations`, `FetchDates`, and `FetchRelation`, each of which performs an HTTP `GET` against the live API, checks the status code, reads the body, and unmarshals it into the corresponding Go struct. If any of these fail, the server logs the error and exits (`log.Fatal`) rather than starting in a broken state.
2. **In-memory storage:** the fetched data is held in package-level variables (`artists`, `locations`, `dates`, `relations`) for the life of the process — there is no database and no re-fetching per request.
3. **Template parsing:** all HTML templates are parsed once at startup with `template.Must(template.ParseGlob("templates/*.html"))`, avoiding redundant disk reads on every request.
4. **Routing:** `http.NewServeMux` maps three routes — `/`, `/artist`, and `/static/` — to their handlers.
5. **Request handling:** each handler checks the HTTP method, validates any input (e.g. the `id` query parameter), looks up the requested data, and executes the matching template with `html/template` — which automatically escapes any data interpolated into HTML, preventing XSS.

## Getting Started

### Prerequisites

- [Go 1.22.2](https://go.dev/dl/) or later installed
- An internet connection (the server fetches live data from the Groupie Trackers API at startup)

### Run locally

```bash
# 1. Clone the repository
git clone https://github.com/AdahJerrie/groupie-tracker.git
cd groupie-tracker

# 2. Run the server
go run main.go
```

You should see:

```
starting server on :8080
```

Then open **http://localhost:8080** in your browser.

### Build a binary

```bash
go build -o groupie-tracker main.go
./groupie-tracker
```

## Routes

| Method | Path | Description |
|---|---|---|
| `GET` | `/` | Home page — lists all artists |
| `GET` | `/artist?id={id}` | Detail page for a single artist |
| `GET` | `/static/*` | Serves CSS and other static assets |

Any other method on these routes returns `405 Method Not Allowed`. An unknown path, a non-numeric `id`, or an `id` with no matching artist/relation returns `404 Not Found`.

## Error Handling

The project follows a "never crash at runtime" principle:

- All upstream API calls check the response status code before attempting to decode the body.
- All JSON decoding errors are wrapped with context using `fmt.Errorf("...: %w", err)`, preserving the original error for debugging.
- Handler-level errors (bad method, bad `id`, missing artist) return the correct HTTP status code instead of panicking.
- Fatal startup errors (e.g. the API is unreachable) fail fast with `log.Fatal` **before** the server starts accepting traffic, rather than serving a half-initialized app.

## Roadmap

This project is under active development as part of a structured learning path. Planned next steps:

- [ ] **Client-server event/action feature** — a user-triggered interaction that causes a fresh request/response round trip after initial page load (candidates under consideration: search by artist name, filter by creation year range, filter by member count, filter by concert location).
- [ ] Unit tests for the fetch and handler functions.
- [ ] Graceful handling of partial upstream API failures (currently a fatal fetch failure at startup stops the server entirely).

## Author

**Jerrie Adah**
GitHub: https://github.com/AdahJerrie

Built as part of the 01Edu project-based curriculum, under a self-directed "AI-native software engineering" mentorship approach emphasizing fundamentals-first learning, guided code review, and intentional AI usage.

## License

No license specified yet — all rights reserved by the author unless stated otherwise.
