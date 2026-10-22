<p align="center">
  <img src="static/logo.svg" width="96" height="96" alt="Ledger logo">
</p>

<h1 align="center">Ledger</h1>

<p align="center">A small self-hosted household budget, built to replace a yearly budget spreadsheet.</p>

Ledger is one Go binary with three dependencies (MySQL driver, go-oidc, x/oauth2) and no frontend build step.

## Features

- A yearly budget grid with income, expenses, the disposable amount and allocation, by month, with totals that update as you type
- Charts beside the grid: year totals, income vs. expenses per month, largest expenses and allocation
- Spreadsheet-style keyboard use: arrow keys or vim keys (`hjkl`, `gg`/`G`, `{`/`}`, `x`, `u`, `y`/`p`), Enter saves; press `?` for help
- New years start empty or as a copy of the previous year
- Statistics across years: totals, change from the year before and every line over time, for a chosen range of years
- Export a year to an OpenDocument spreadsheet (`.ods`) with the totals as formulas
- Login with OIDC, e.g. passkeys via [Pocket ID](https://pocket-id.org), or local mode without login for trying it out
- English and Danish, chosen per user or taken from the browser; new languages are JSON files
- Per-user settings for language, money format (`1.234,56` or `1,234.56`) and light/dark theme
- Room for plugins (pension, investments, …), see `plugin.go`

Planned: bank CSV import and budget vs. actual.

## Try it

Start a throwaway MySQL, fill it with a fictional demo household and run Ledger in local mode (no login, only reachable from your own machine):

```sh
docker run -d --name ledger-demo-db -p 127.0.0.1:3306:3306 \
  -e MYSQL_DATABASE=ledger -e MYSQL_USER=ledger -e MYSQL_PASSWORD=ledger -e MYSQL_RANDOM_ROOT_PASSWORD=yes mysql:8.4

go build
export LEDGER_DB_DSN='ledger:ledger@tcp(127.0.0.1:3306)/ledger'
./ledger demo        # or "./ledger demo da" for Danish line names; give MySQL ~20 s to start first
./ledger             # then open http://localhost:8080
```

`ledger demo` adds three years of budgets, ending with the current year. It refuses to run on a database that already has budgets, so it cannot mix demo data into real data. With Docker Compose: `docker compose run --rm ledger demo`.

## Configuration

| Variable | Required | Description |
|---|---|---|
| `LEDGER_DB_DSN` | yes | e.g. `ledger:secret@tcp(localhost:3306)/ledger` |
| `OIDC_ISSUER` | no | e.g. `https://id.example.com`. Without it the app runs in local mode: no login, loopback only |
| `OIDC_CLIENT_ID` | with OIDC | |
| `OIDC_CLIENT_SECRET` | no | Leave empty for a public client (PKCE is always used) |
| `LEDGER_BASE_URL` | no | Public URL, default `http://localhost:8080` |
| `LEDGER_ADDR` | no | Listen address, default `:8080` (`127.0.0.1:8080` in local mode) |
| `LEDGER_LOCALES_DIR` | no | Extra or overriding language files |

The database schema is created and migrated automatically on startup.

### Pocket ID

Create an OIDC client with the callback URL `<LEDGER_BASE_URL>/auth/callback`.
Everyone who can log in to the client has access, so restrict it to your household with *Allowed user groups*.

## Running

```sh
docker compose up -d   # see compose.yaml
# or
go build && LEDGER_DB_DSN=... OIDC_ISSUER=... OIDC_CLIENT_ID=... ./ledger
```

## Languages

Each language is a `locales/<code>.json` file. To add one, copy `locales/en.json` to e.g. `de.json`, translate it,
and either put it in `locales/` (compiled in) or in `LEDGER_LOCALES_DIR`. Missing keys fall back to English.
`_decimal` and `_thousands` control number formatting.

## Plugins

A plugin implements the `Plugin` interface in `plugin.go` and calls `register` from `init()`.
It can add migrations (tracked per plugin), routes behind login and menu entries.

## License

MIT, see [LICENSE](LICENSE).
