# driftwatch

A CLI tool that detects configuration drift between deployed services and their declared infrastructure state.

---

## Installation

```bash
go install github.com/yourusername/driftwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/driftwatch.git
cd driftwatch
go build -o driftwatch .
```

---

## Usage

Run a drift check against your declared infrastructure state:

```bash
driftwatch check --config ./infra/state.yaml --env production
```

Example output:

```
[✓] api-service       — in sync
[✗] worker-service    — drift detected: replicas declared=3, actual=1
[✓] cache-service     — in sync

1 drift(s) found across 3 services.
```

### Common Flags

| Flag | Description |
|------|-------------|
| `--config` | Path to the declared state file |
| `--env` | Target environment to check |
| `--output` | Output format: `text`, `json`, `yaml` |
| `--fail-on-drift` | Exit with code 1 if drift is detected |

---

## Configuration

`driftwatch` reads a `state.yaml` file describing your expected service configuration. See [`docs/config.md`](docs/config.md) for the full schema reference.

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

---

## License

[MIT](LICENSE)