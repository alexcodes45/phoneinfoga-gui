# Dev Notes

- Use `wails init` to generate the UI shell and `wails.json`. Then delete `cmd/desktop/main.go` stub and replace with Wails bootstrap code.
- Implement PhoneInfoga serve launcher and CLI fallback in `internal/phoneinfoga`.
- Ent: add edges between Scan and Result, etc., as needed, then run `make generate`.
