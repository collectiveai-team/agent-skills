# Sections

This file defines all sections, their ordering, impact levels, and descriptions.
The section ID (in parentheses) is the filename prefix used to group rules.

---

## 1. Quality Gates (gates)

**Impact:** CRITICAL
**Description:** Mandatory gate order, blocking semantics, and validator reporting rules.

- [Block on Errors and Strict Warnings](gates-block-on-errors-and-strict-warnings.md) - CRITICAL
- [Run Quality Gates in Fixed Order](gates-run-in-order.md) - HIGH
- [Run SC Validators in KISS Purity SOLID Crash Order](gates-run-sc-validators-in-kiss-purity-solid-crash-order.md) - HIGH
- [Stop at the First Failing Gate](gates-stop-on-failure.md) - HIGH
- [Use Strict Principles Profile When Required](gates-use-strict-principles-profile-when-required.md) - HIGH
- [Use Validator Exit Codes in Automation](gates-use-validator-exit-codes-in-automation.md) - HIGH
- [Report Validator Summaries and Top Categories](gates-report-validator-summaries-and-top-categories.md) - MEDIUM
- [Use Flags in Supported Contexts](gates-use-flags-in-supported-contexts.md) - MEDIUM

## 2. SC Principles (sc)

**Impact:** CRITICAL
**Description:** Core simplicity, purity, SOLID, and fail-fast error-handling rules.

- [Never Raise Not Implemented Error in Overrides](sc-never-raise-notimplementederror-in-overrides.md) - CRITICAL
- [Never Use Bare Except](sc-let-it-crash-no-bare-except.md) - CRITICAL
- [Never Use Except Pass](sc-never-use-except-pass.md) - CRITICAL
- [Apply Enforcement Severity Correctly](sc-apply-enforcement-severity-correctly.md) - HIGH
- [Classify Core and Shell Layers Explicitly](sc-classify-core-and-shell-layers.md) - HIGH
- [Enforce KISS Thresholds](sc-enforce-kiss-thresholds.md) - HIGH
- [Handle Operational Failures at Boundaries](sc-handle-operational-failures-at-boundaries.md) - HIGH
- [Inject Infrastructure Through Abstractions](sc-inject-infrastructure-through-abstractions.md) - HIGH
- [Keep Async IO in Shell Layers](sc-keep-async-io-in-shell-layers.md) - HIGH
- [Keep Domain Logic Free of I/O](sc-functional-core-no-io.md) - HIGH
- [Keep Functions Small and Single-Purpose](sc-kiss-small-functions.md) - HIGH
- [Keep Logging and Print Out of Core Logic](sc-keep-logging-and-print-out-of-core.md) - HIGH
- [Pass Data Into Core and Return Values](sc-pass-data-into-core-and-return-values.md) - HIGH
- [Avoid Nested Try Fallback Ladders](sc-avoid-nested-try-fallback-ladders.md) - MEDIUM
- [Keep One Responsibility per Module](sc-solid-single-responsibility.md) - MEDIUM
- [Prefer Guard Clauses Over Nested Branches](sc-prefer-guard-clauses-over-nested-branches.md) - MEDIUM
- [Replace Type Branch Chains With Strategy](sc-replace-type-branch-chains-with-strategy.md) - MEDIUM
- [Run Decision Checklist Before Merge](sc-run-decision-checklist-before-merge.md) - MEDIUM
