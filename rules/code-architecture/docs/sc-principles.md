# SC Principles

## Purpose
Define non-negotiable architecture and code-health principles for Python backend work.

These principles are mandatory for design, implementation, and refactoring decisions:

1. **KISS** - Keep behavior simple, readable, and easy to change.
2. **Functional Core, Imperative Shell** - Keep business logic pure; isolate I/O at edges.
3. **SOLID** - Maintain clear responsibilities and dependency direction.
4. **Let It Crash** - Fail fast on programming errors, handle operational failures at boundaries.

## Scope
Apply this rule set to production Python backend code, especially:

- Domain and business logic
- Services and orchestration
- API, CLI, and adapter boundaries
- Refactors and new feature implementation

## Principle Interaction Model

These principles are designed to work together:

- **KISS + SOLID**: keep units small and responsibilities clear.
- **Purity + Let It Crash**: isolate side effects and avoid hiding failures.
- **SOLID + Purity**: dependency inversion keeps business rules independent from infrastructure.

When trade-offs appear, prioritize clarity of intent and boundary correctness.

## Principle 1: KISS (Keep It Simple)

Prefer simple, explicit flows over cleverness.

### Enforcement thresholds

| Metric | Warning | Error | Why it matters |
|---|---:|---:|---|
| Cyclomatic complexity | > 7 | > 10 | Too many branches increase defect risk |
| Cognitive complexity | - | > 15 | Deeply nested control flow hurts readability |
| Function length | - | > 50 lines | Long functions mix concerns |
| Nesting depth | - | > 4 levels | Excess nesting hides intent |
| Parameter count | > 5 | - | Too many inputs suggests split responsibilities |

Function length is counted inclusively. Nesting includes `if`, `for`, `while`, `with`, and `try` blocks.

### Cyclomatic vs cognitive complexity

- **Cyclomatic** complexity counts decision points.
- **Cognitive** complexity penalizes nesting depth more heavily.
- A function can have moderate cyclomatic score but still be hard to read due to deep nesting.

### KISS guidance

- Extract helper functions when branches and side paths grow.
- Prefer guard clauses over nested `if/else` pyramids.
- Name intermediate conditions to explain intent.
- Keep one main reason to change per function.

### KISS anti-patterns

- Large orchestration and decision logic mixed in one function.
- Deeply nested control structures hiding primary flow.
- Utility functions with too many inputs and mode flags.

## Principle 2: Functional Core, Imperative Shell

Business decisions belong in pure functions. I/O belongs at boundaries.

### Layer model

| Layer | Typical paths | I/O allowed | Policy |
|---|---|---|---|
| Core | `*/domain/*`, `*/logic/*`, `*/services/*`, `*/utils/*`, `*/core/*` | No | I/O in core is an error |
| Shell | `*/handlers/*`, `*/adapters/*`, `*/api/*`, `*/cli/*`, `*/scripts/*`, `*/tests/*` | Yes | Keep effects localized and explicit |

Treat these as shell context even when outside shell paths:

- `*/archive/*`
- `*/examples/*`
- `*/benchmarks/*`
- `conftest.py`
- `setup.py`
- `*__main__.py`

### Side effects to keep out of core

- File I/O (`open`, `read`, `write`, `Path.read_text`, `Path.write_text`)
- Network (`requests`, `httpx`, `urllib`, `socket`)
- Database (`execute`, `query`, `session.add`, cursor operations)
- Subprocess (`subprocess.run`, `Popen`, `os.system`)
- Global state (`global`, `nonlocal`)
- Side effects (`print`, `logging.*`, `logger.*`)
- Async I/O (`async def`, `await`, `async for`, `async with`)

### Purity guidance

- Pass data into core functions; do not fetch it there.
- Return values from core; perform persistence/output in adapters.
- Keep async and transport concerns in shell layers.

### Purity anti-patterns

- Calling repository/query code from domain transformation functions.
- Logging and printing inside core calculation logic.
- Mixing validation rules with transport/protocol concerns.

## Principle 3: SOLID

Use SOLID as a design pressure test, not ceremony.

### Detection heuristics

| Principle | Common smell | Severity | Typical correction |
|---|---|---|---|
| SRP | File >300 lines, class with >5 public methods | Warning | Split by responsibility |
| OCP | `if/elif` type chains, repeated `isinstance` branching | Warning | Strategy/registry dispatch |
| LSP | Override that raises `NotImplementedError` | Error | Fix contract or inheritance model |
| ISP | Interface/protocol with >7 methods | Warning | Split into smaller focused interfaces |
| DIP | Business logic instantiates concrete infrastructure directly | Warning | Depend on abstractions and inject adapters |

### SOLID smell matrix

| Smell | Principle | Recommended fix |
|---|---|---|
| File >300 lines | SRP | Extract modules by responsibility |
| Type-driven branch chains | OCP | Strategy pattern or registry dispatch |
| Override that throws | LSP | Honor contract or avoid inheritance |
| Interface with 10+ methods | ISP | Split into focused interfaces |
| `new Service()` in core logic | DIP | Inject dependencies via abstraction |

### SOLID guidance

- Keep modules narrowly focused and composable.
- Prefer extension points over modifying central branching logic.
- Ensure subtypes satisfy base contracts in behavior, not only signatures.
- Keep interfaces small and role-focused.
- Point dependencies inward: policy depends on abstractions, not frameworks.

### SOLID anti-patterns

- Feature modules becoming catch-all files over time.
- Type-switch logic growing each time a new variant is added.
- Interfaces that force callers to implement irrelevant methods.
- Core services constructing infrastructure clients directly.

## Principle 4: Let It Crash

Do not hide programming errors. Surface them early and clearly.

### Prohibited patterns

| Pattern | Severity | Why |
|---|---|---|
| Bare `except:` | Error | Catches everything, including interrupts |
| `except: pass` | Error | Silent failure obscures defects |
| Broad `except Exception` without re-raise or boundary translation | Warning | Swallows root cause |
| Defensive `if not x: return` chains masking bugs | Warning | Hides contract violations |
| Nested try/except fallback ladders | Warning | Creates opaque error paths |

### When to let it crash

- Validation and programming errors in internal logic.
- Contract violations and impossible states.
- Internal operations where failure indicates a defect.

### Where to handle errors explicitly

- API/CLI boundaries (user-safe messages)
- Adapter integrations (network, persistence, external systems)
- Resource cleanup and transactional boundaries
- External API failures requiring retry/fallback/degradation
- Data persistence paths where data loss must be prevented

### Let-it-crash guidance

- Validate inputs at boundaries; trust internal contracts.
- Catch specific exception types only when you can recover or translate.
- Re-raise with context when preserving stack information.

### Not flagged as anti-patterns

- Boundary handling in `*/adapters/*`, `*/api/*`, and `*/cli/*` contexts.
- Swallowing only when explicit logging and justification are present.
- Catch-then-re-raise with contextual logging.

### Let-it-crash anti-patterns

- Catch-all error handling to keep code "working" in unknown states.
- Defensive null-guard chains replacing explicit contracts.
- Nested fallback trees that make true failure causes opaque.

## Operational Guidelines

- Validate at boundaries, not deep in domain internals.
- Keep dependencies flowing inward (`core <- domain <- adapters`).
- Prefer explicit contracts over implicit behavior.
- Replace broad exception handling with typed, intentional handling.

## Enforcement Semantics

- **Error** means design risk is high enough to block integration.
- **Warning** means remediation is expected, but may be deferred with justification.
- **Strict mode** promotes warnings to blocking outcomes.

Execution and script details are defined in `rules/code-architecture/docs/quality-gates.md`.

## Refactoring Patterns

- **Complexity issues**: extract methods, apply guard clauses, decompose conditionals.
- **Purity issues**: inject dependencies, move I/O to adapters, return values instead of side effects.
- **SOLID issues**: split responsibilities, introduce strategy/registry, isolate interfaces.
- **Crash-handling issues**: remove catch-all handlers, handle recoverable failures at boundaries.

## Decision Checklist

Use this checklist before merging backend code:

1. Is core business logic free of file/network/database/subprocess side effects?
2. Are functions readable without deep nesting and mode-flag branching?
3. Do dependencies point inward toward abstractions?
4. Are exceptions handled intentionally at boundaries, not swallowed internally?
5. Can another engineer identify module responsibility in one sentence?

## External References

- https://en.wikipedia.org/wiki/KISS_principle
- https://martinfowler.com/bliki/FunctionalCoreImperativeShell.html
- https://en.wikipedia.org/wiki/SOLID
- https://docs.python.org/3/tutorial/errors.html
