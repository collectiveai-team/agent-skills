# Sections

This file defines all sections, their ordering, impact levels, and descriptions.
The section ID (in parentheses) is the filename prefix used to group rules.

---

## 1. SC Principles (sc)

**Impact:** CRITICAL
**Description:** Core simplicity, purity, SOLID, and fail-fast error-handling rules.

- [Never Raise Not Implemented Error in Overrides](sc-never-raise-notimplementederror-in-overrides.md) - CRITICAL
- [Never Use Bare Except](sc-let-it-crash-no-bare-except.md) - CRITICAL
- [Never Use Except Pass](sc-never-use-except-pass.md) - CRITICAL
- [Avoid Broad Except Without Contextual Re-raise](sc-avoid-broad-except-without-contextual-reraise.md) - HIGH
- [Classify Core and Shell Layers Explicitly](sc-classify-core-and-shell-layers.md) - HIGH
- [Inject Infrastructure Through Abstractions](sc-inject-infrastructure-through-abstractions.md) - HIGH
- [Keep Async IO in Shell Layers](sc-keep-async-io-in-shell-layers.md) - HIGH
- [Keep Domain Logic Free of I/O](sc-functional-core-no-io.md) - HIGH
- [Keep Functions Small and Single-Purpose](sc-kiss-small-functions.md) - HIGH
- [Keep Logging and Print Out of Core Logic](sc-keep-logging-and-print-out-of-core.md) - HIGH
- [Pass Data Into Core and Return Values](sc-pass-data-into-core-and-return-values.md) - HIGH
- [Avoid Defensive Null Guard Chains](sc-avoid-defensive-null-guard-chains.md) - MEDIUM
- [Avoid Nested Try Fallback Ladders](sc-avoid-nested-try-fallback-ladders.md) - MEDIUM
- [Keep One Responsibility per Module](sc-solid-single-responsibility.md) - MEDIUM
- [Prefer Guard Clauses Over Nested Branches](sc-prefer-guard-clauses-over-nested-branches.md) - MEDIUM
- [Replace Type Branch Chains With Strategy](sc-replace-type-branch-chains-with-strategy.md) - MEDIUM
