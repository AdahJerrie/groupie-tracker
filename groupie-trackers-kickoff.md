# Groupie Trackers — Project Kickoff

*Catch-up kickoff, produced after Phase 0 (HTTP fundamentals), Phase 1 (Go fundamentals: structs, slices, maps, pointers), and the start of Phase 2 (JSON modeling) were already underway. Written now, with real API data already verified, rather than at zero — so it reflects ground truth, not guesses.*

---

## 1. Project Brief

**What we're building:** A Go HTTP server that fetches data from the public Groupie Trackers API (four linked resources — `artists`, `locations`, `dates`, `relation`), merges them into a unified per-artist view, and serves that data as browsable HTML pages. The project also requires at least one live client-server interaction feature — something a user triggers *after* the initial page load that causes a fresh request/response round-trip (feature TBD, to be designed once the server and templates exist).

**Why this project, for this learner:** Jerrie has prior Go exposure but wants the foundation rebuilt deliberately and completely — "every edge, nook and cranny" — covering HTTP fundamentals, Go data modeling, JSON, server-side rendering, and one full client-server interaction loop, with AI used as a mentor/reviewer rather than a code generator.

**Success criteria:**
- All four API resources are correctly modeled and merged without data loss or silent mismatches.
- The server never crashes — all failure modes (bad data, network issues, missing pages) degrade gracefully with proper status codes and a real fallback UI.
- The code is idiomatic Go: proper error handling, exported/unexported field usage, separation of concerns (data logic separate from request handling).
- At least one genuine client-server event/action feature exists and is understood, not copy-pasted.
- The learner can explain *why* each piece of code is shaped the way it is — not just that it works.

**Known constraints (from the brief and the constitution):**
- Backend: Go only.
- No crashes, ever — including on upstream API failure.
- Errors must be handled, not ignored or silently swallowed.
- Good practices expected: naming, structure, idiomatic patterns.
- Unit tests recommended.

---

## 2. Engineering Breakdown

**The four real resources** (verified directly against the live API):

| Resource | Shape | Notes |
|---|---|---|
| `artists` | `id`, `image`, `name`, `members []string`, `creationDate int`, `firstAlbum string`, plus links to that artist's `locations`, `concertDates`, `relations` | Confirmed live — e.g. Queen, `id:1`. |
| `locations` | per-artist list of place strings, e.g. `"dunedin-new_zealand"` | Lowercase, hyphenated city-country format. |
| `dates` | per-artist list of date strings, e.g. `"10-02-2020"` | DD-MM-YYYY format. |
| `relation` | `id`, `datesLocations map[string][]string` | **The critical join.** Maps each location string directly to the list of dates the artist played there. Verified live for artist 1. |

**Architectural layers we'll build, in order:**
1. **Data layer** — Go structs matching real JSON; fetching all four resources; merging them into one `FullArtist`-style structure per artist.
2. **Server layer** — `net/http` routing and handlers; the "waiter" that never touches raw data logic directly.
3. **Presentation layer** — `html/template` + CSS; turning merged Go data into actual pages.
4. **Interaction layer** — the mandatory event/action feature; a real, scoped request/response triggered by user action.
5. **Cross-cutting concerns** — error handling, logging, testing, project structure — applied throughout, not bolted on at the end.

**Edge cases and failure modes to design around (not an afterthought — design for these from the start):**
- Upstream API is slow, down, or returns malformed JSON.
- A `relation` location key doesn't perfectly match a `locations` entry (data inconsistency between resources).
- Missing/empty fields — nil vs. empty slice distinction matters here (Phase 1 lesson directly applies).
- A user requests a non-existent artist ID or page — must 404 gracefully, not panic.
- If we ever cache fetched data in memory, concurrent requests reading/writing it must be race-safe.
- `html/template` auto-escapes by design — we must never bypass this for "convenience," since this is also a security boundary (basic XSS protection), not just a formatting detail.

---

## 3. AI Opportunity Map

*This section was the actual gap in our process — flagged and corrected here. The point isn't "where does Claude help teach" (that's the whole mentoring relationship) — it's specifically: where should AI-as-a-tool plug into the **engineering workflow** of this project, and where would using it actively hurt the learning goal?*

**Good fits — low learning value, high time cost, safe to automate or accelerate with AI:**
- Generating repetitive HTML/CSS boilerplate *after* the first card/page layout has been designed and understood by hand (e.g., replicating a card layout across artist entries).
- Drafting documentation comments, README content, and commit messages once the underlying code is understood and written by the learner.
- Scaffolding table-driven test cases (Go idiom) once the *concept* of what's being tested is understood — AI can generate the repetitive table rows, the learner defines what correctness means.
- Generating sample/mock JSON fixtures for testing edge cases (malformed data, missing fields) — useful for testing error handling without depending on the live API being in a specific state.
- Linting/formatting automation (`gofmt`, `go vet`) — already mechanical, zero learning value in doing it by hand.

**Poor fits — must NOT be handed to AI, because the learning value is the entire point:**
- The actual data-modeling and merge logic (structs, JSON tags, combining the four resources) — this is core engineering judgment being built right now.
- The error-handling decisions themselves (what should happen when X fails) — AI can review these decisions, but the learner must make them first.
- Designing the event/action feature (Phase 7) — this requires understanding tradeoffs the learner hasn't practiced yet; designing it via AI would skip the exact skill the project exists to build.
- Debugging compiler errors by just pasting them into AI for a fix without first reasoning about *why* — we've been correctly avoiding this; should continue.

**Middle ground — the actual current practice, and it's a legitimate AI-native pattern:** using AI as a **reviewer, rubber duck, and edge-case prober** (what we're doing in every code review so far) is itself a real, industry-relevant AI engineering skill — not a shortcut. Continue this explicitly, and treat it as a named skill being practiced, not just "how Claude happens to help."

---

## 4. Industry Perspective

How a professional or AI-native team would likely diverge from a pure learning-exercise approach on this same project:

- **Project structure:** rather than one `main.go`, a real team would use Go's conventional layout — e.g. `/internal/models` (structs), `/internal/fetch` (API client logic), `/internal/handlers` (HTTP layer), `/web/templates`, `/web/static` — enforcing the separation-of-concerns we've already discussed conceptually, but as actual folder boundaries.
- **Not fetching live on every request:** Since this project depends on a *third-party* API we don't control, a production team would not call it fresh on every page load. They'd fetch once at startup (or on an interval) and serve from an in-memory cache — directly relevant to our "never crash" requirement, since it means a temporary upstream outage doesn't take our site down with it.
- **Structured logging:** real services favor structured logging (e.g. Go's `log/slog`, stdlib since Go 1.21) over plain `fmt.Println`/`log.Println`, so logs are filterable and machine-readable in production. Worth introducing once we're past basic `if err != nil` fluency.
- **Testing strategy:** the merge/join logic (artists + relation → unified view) is the highest-risk, highest-value code to unit test — exactly where table-driven tests shine, since there are many small input variations (missing location, empty dates, mismatched keys) to verify cheaply.
- **CI automation:** even a small project benefits from a CI pipeline (e.g. GitHub Actions) running `go build`, `go vet`, and `go test` on every push — itself a legitimate "good AI opportunity" item, since scaffolding a CI YAML file is low-learning-value, high-leverage automation.
- **Graceful degradation as a first-class concern**, not a final cleanup step — directly enforced by this project's own "never crash" rule, and worth treating with the seriousness a real production team would.

---

## 5. Implementation Plan

Builds on the roadmap already in motion — formalized here as the official plan, with an explicit AI-Opportunity-Discovery checkpoint folded in before each remaining phase from now on (per the constitution's Phase 2/3, rather than skipped as it was earlier).

| Phase | Status | Focus |
|---|---|---|
| 0 | ✅ Done | HTTP request/response fundamentals |
| 1 | ✅ Done | Go fundamentals: structs, slices, maps, pointers |
| 2 | 🔄 In progress | JSON: struct tags, `Unmarshal`, error handling idiom |
| 3 | Next | Live API fetching (`net/http` client) + merging the four resources |
| 4 | Upcoming | Go HTTP server: routing, handlers (the "waiter" layer) |
| 5–6 | Upcoming | `html/template` + CSS: rendering merged data as real pages |
| 7 | Upcoming | Designing and building the event/action feature (brainstormed together, not pre-decided) |
| 8 | Upcoming | Crash-proofing: panics, recovery, graceful error pages |
| 9 | Upcoming | Unit testing, especially the merge logic |
| 10 | Upcoming | Project structure & security review pass |
| Closing | Upcoming | Senior Engineer Review + AI-Native Retrospective + AI-Native Score (1–10) + next-concept recommendations |

**Process change starting now:** before each phase above, we'll briefly run the constitution's AI Opportunity Discovery + AI-native discussion (lightweight, not necessarily a new written document each time) — e.g., before Phase 3, explicitly discussing whether/how AI tooling could assist live API fetching and where it shouldn't.
