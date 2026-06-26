<p align="center">
<img src="docs/auditory.png" alt="yatas-logo" width="30%">
<p align="center">

# YATAS — Plugin Template

Yet Another Testing &amp; Auditing Solution

This repository is a **template** for building your own [YATAS](https://github.com/padok-team/yatas) plugin. A plugin audits the resources of a provider (a cloud account, a SaaS, an internal API, ...) against a set of checks and reports `OK`/`FAIL` results back to YATAS.

Each plugin audits **one kind** of thing (a Kubernetes plugin, an ArgoCD plugin, a cloud plugin, ...) — just like [yatas-aws](https://github.com/padok-team/yatas-aws) audits AWS. It ships with a runnable example category so you can build, run and test it immediately, then replace the example with your own checks.

## Getting started

1. **Rename the plugin.** Replace `yatas-template` / `template` with your plugin name in:
   - `go.mod` (the module path)
   - `Makefile` (the binary name and install path)
   - `internal/constants.go` (`PluginName`)
   - `package.json` (`name`, `repository`, ...)
2. **Build it:**
   ```bash
   make build
   ```
3. **Install it locally** to test against a real YATAS install:
   ```bash
   make install
   ```
   Then in your `.yatas.yml`, set the plugin `version` to `local` to use the build you just installed (see [How to test?](#how-to-test)).

## Architecture

A plugin audits a list of **targets** (the things to check — clusters, instances, accounts, ...). For each target it builds a **session** (the authenticated client), then runs every **check**, grouped into **categories**. This is the same shape as yatas-aws.

Work fans out so nothing blocks anything else:

```
Run (main.go)
└── run                  one goroutine per TARGET    (.yatas.yml accounts)
    └── initTest         one goroutine per CATEGORY  (checks/<category>/)
        └── RunChecks    one goroutine per CHECK
            └── checkIf...  emits OK/FAIL Results
```

- **`internal/`** — `Target` is your typed connection config (rename it; ≈ yatas-aws's `AWS_Account`) and `Session` is the authenticated client (≈ `aws.Config`). `UnmarshalConfig` reads `.yatas.yml` into `[]Target`.
- **`main.go`** — `run` fans out per target; `initSession` builds the client; `initTest` lists every category (one `CheckMacroTest` line each), exactly like yatas-aws's `initTest`.
- **`checks/<category>/`** — one folder per category (≈ `aws/s3`, `aws/ec2`, ...), with this file convention:
  | File | Responsibility |
  | --- | --- |
  | `getter.go` | All API / IO. Fetches the resources to audit, using the `Session`. |
  | `<category>.go` | `RunChecks`: fetch once, fan out the checks, collect results. |
  | `<check>.go` | One check: turns fetched data into `OK`/`FAIL` results. Never calls an API. |
  | `<check>_test.go` | Unit test feeding hand-built data and asserting the status. |

To make this a real plugin: rename `Target`/`Session` and fill in your connection fields (`internal/`), build the client in `initSession` (`main.go`), then replace the `example` category with your own.

## How to add a new check

To an existing category (e.g. `checks/example`):

1. Create `checks/example/<myCheck>.go` with a `checkIf...` function (see `exampleCompliant.go`).
2. Give it a new, never-reused ID constant, e.g. `TEMPLATE_EXAMPLE_002`.
3. If it needs new data, add the fetch to `getter.go`.
4. Add one line to `RunChecks`:
   ```go
   go commons.CheckTest(checkConfig.Wg, c, MyCheckID, checkIfMyCondition)(checkConfig, resources, MyCheckID)
   ```
5. Add a `_test.go` with a passing and a failing case.

To add a new **category**:

1. Create `checks/<category>/` with a `getter.go` and a `<category>.go` exporting `RunChecks(wg, s, c, queue)`.
2. Add one line in `initTest` (`main.go`):
   ```go
   go commons.CheckMacroTest(&wg, c, mycategory.RunChecks)(&wg, s, c, queue)
   ```

`CheckTest` and `CheckMacroTest` are generic wrappers from `yatas/plugins/commons`: they register the check/category on the wait group and honor the user's `include`/`exclude` config — a disabled check becomes a no-op, no extra code needed.

## Configuration

Users enable and configure your plugin in their `.yatas.yml`. See `.yatas.example.yml` for a minimal example:

```yaml
plugins:
  - name: "template"
    enabled: true
    source: "github.com/padok-team/yatas-template"
    version: "local"

pluginsConfiguration:
  - pluginName: "template"
    accounts:
      - name: "my-target"
        # server: "https://api.example.com"
        # token: "${API_TOKEN}"
```

Each entry under `accounts` is one target to audit. `name` labels it in the report; every other key is read by `internal/UnmarshalConfig` into your `Target` struct — add a `case` there for each field you define.

### Exclude / include checks

Users can select which of your checks run, by ID:

```yaml
plugins:
  - name: "template"
    enabled: true
    exclude:
      - TEMPLATE_EXAMPLE_001
    # or, to run ONLY specific checks:
    include:
      - TEMPLATE_EXAMPLE_001
```

### Ignore known results

Users can silence individual results (optionally with a regex):

```yaml
ignore:
  - id: "TEMPLATE_EXAMPLE_001"
    regex: true
    values:
      - "Resource .* is not compliant"
```

## How to test?

Run the unit tests:

```bash
make test
```

To test end-to-end, run `make install` and set `version: "local"` for your plugin in `.yatas.yml` — YATAS will load the binary you just installed instead of downloading a release.

### Debug logs

```bash
export YATAS_LOG=debug
```

Levels: `debug`, `info`, `warn`, `error`, `fatal`, `panic`, and `off` (default).

## How to deploy?

The provided GitHub Actions workflow (`.github/workflows/release.yml`) releases your plugin with [GoReleaser](https://goreleaser.com/) on tag push. The released binary must start with `yatas-` and end with the name of your plugin so YATAS can resolve it.

## PR labels

`.github/workflows/label.yml` automatically labels pull requests based on which files changed, using the rules in `.github/labeler.yml`. The template ships with generic rules (`dependencies`, `main`, `makefile`, `readme`, `github-actions`).

As your plugin grows, add one rule per category so reviewers can see at a glance what a PR touches. Each rule maps a label to a glob:

```yaml
example:
  - changed-files:
      - any-glob-to-any-file: 'checks/example/**/*'
```

Add the matching label in your repository settings (Issues &amp; PRs → Labels) for it to be applied.
