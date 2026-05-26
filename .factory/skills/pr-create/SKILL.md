# PR Creation — Couchbase Capella Terraform Provider

Use this skill to create a pull request for the current branch.

---

## Before you write anything

1. Run `git log main..HEAD --oneline` to see what's on the branch.
2. Run `git diff main..HEAD --stat -- '*.go' ':!internal/generated/api/openapi.gen.go'` to see what changed.
3. Read `.github/pull_request_template.md` for the required format.
4. Push the branch if not already pushed: `git push -u origin <branch>`.

---

## Writing the PR

Keep it short and human. Write like you're telling a teammate what you did and why.

- **Title:** `[JIRA-KEY] short description` — under 70 characters.
- **Description:** Say what changed and why. If bugs were found, name them and say whether they're fixed or documented.
- **No AI padding.** No "This PR introduces...", no bullet lists of every file touched, no restating the title.
- Only include a section if it has real content.

---

## PR template fields

### Jira
List every ticket this touches — the primary ticket and any bugs raised.

### Description
Two to four sentences covering:
- What was added or changed
- Any bugs fixed (with ticket reference)
- Any bugs documented/workaround (with ticket reference)

### Type of Change
Tick what applies. Usually bug fix + new feature for test branches.

### Testing
Tick "Acceptance tested" if tests were run. List which test groups passed under the Testing details block. Keep it brief — just test name pattern and PASS/FAIL.

### Further comments
Leave as "None." if there's nothing extra.

---

## Creating the PR

```bash
gh pr create \
  --title "[AV-XXXXX] short description" \
  --body "$(cat <<'EOF'
## Jira
* AV-XXXXX

## Description
<2-4 sentences>

## Type of Change
- [x] Bug fix
- [x] New feature

## Manual Testing Approach
### How was this change tested and do you have evidence? _**(REQUIRED: Select at least 1)**_
- [x] Acceptance tested

### Testing
<details open>
  <summary>Testing</summary>
<test results>
</details>

## Further comments
None.
EOF
)"
```

---

## After creating

Report back:
- PR URL
- Jira tickets linked
- Any uncommitted changes the user should be aware of
