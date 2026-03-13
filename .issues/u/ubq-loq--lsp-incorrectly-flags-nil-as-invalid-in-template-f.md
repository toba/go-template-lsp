---
# ubq-loq
title: LSP incorrectly flags nil as invalid in template function arguments
status: completed
type: bug
priority: normal
created_at: 2026-03-13T16:24:25Z
updated_at: 2026-03-13T16:36:21Z
sync:
    clickup:
        synced_at: "2026-03-13T16:36:56Z"
        task_id: 868hw2vdu
---

## Description

When calling a template function that accepts variadic `...any` parameters, the LSP flags `nil` as an invalid argument with a diagnostic error (red squiggly underline).

## Example

```html
{{template "page-link" linkTo .Page "portfolio-list" nil false false}}
```

The function signature is `func linkTo(args ...any) any`, so `nil` is a perfectly valid argument — it satisfies the `any` type. The LSP should accept `nil` as a valid value for any parameter type.

## Steps to Reproduce

1. Define a template function with a variadic `...any` parameter (e.g. `linkTo`)
2. Call it in a template with `nil` as one of the arguments
3. Observe the red squiggly underline on `nil`

## Expected Behavior

`nil` should be accepted as a valid argument for `any`-typed (and other nillable) parameters without diagnostic errors.

## Actual Behavior

`nil` is flagged as an invalid argument.


## Summary of Changes

Registered `nil` as a builtin identifier in `getBuiltinFunctionDefinition()` (`internal/template/analyzer/analyzer.go`), giving it a zero-arg signature returning `types.Typ[types.UntypedNil]`. This mirrors how `true` and `false` are already registered. When `nil` appears in a template expression, the analyzer now resolves it through the builtin function map instead of flagging it as an undefined function.

### Files Changed
- `internal/template/analyzer/analyzer.go` — added `nil` to builtin signatures
- `cmd/go-template-lsp/main_test.go` — added `TestNilAsTemplateArgument` with 5 test cases
