---
name: tf-bugfix
description: Diagnose and fix bugs in the Terraform Capella provider with acceptance tests.
---

# Terraform Bug Fix

## Instructions

0. Read the Jira ticket or GitHub issue provided by the user.
   Summarize: what is the bug, what is the expected behavior,
   and what is the actual behavior.

1. Search the codebase for the function, file, or code path
   mentioned in the bug report. Read the relevant source files
   to understand the current behavior.

2. Identify the root cause. Explain why the current code
   produces the bug (e.g., missing bounds check, nil pointer,
   incorrect logic).

3. Propose a robust fix that addresses the root cause.

4. Add an acceptance test in `acceptance_tests/` to
   validate the fix end-to-end:
   - Follow existing patterns: use `resource.ParallelTest()`,
     `globalProtoV6ProviderFactory`, and helper functions like
     `randomStringWithPrefix`.
   - Name the test in this format `TestAcc<Feature>_AV-XXXXX` e.g. `TestAccProject_AV-12345`.

5. Do not run the acceptance test.