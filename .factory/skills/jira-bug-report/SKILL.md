# Jira Bug Report — Couchbase Capella Terraform Provider

Use this skill to raise a bug ticket in the AV project on
`jira.issues.couchbase.com`.

---

## Two modes — know which one applies before you write anything

### Mode A — Backend bug (API/server, fixed by backend dev)

- No branch, no PR, no test names, no test file paths.
- Write as a human operator who hit the bug manually.
- Only include what the dev needs: what you did, what happened, what should happen, and the HCL to reproduce it.

### Mode B — Provider bug fixed by us

- Include branch, the file and the change needed.
- No test names or file paths.
- Write the opening as a user who hit the bug, then give the root cause and fix in plain terms.

### Mode C — Provider bug fixed by backend dev

- Same as Mode A. No branch, no test references.

---

## Decision checklist

```
Is the bug in the Capella API/server?  →  Mode A
Is the bug in the Terraform provider?
  Will we fix it?    →  Mode B
  Will dev fix it?   →  Mode C (same as Mode A)
```

---

## Fixed Jira coordinates

| Field | Value |
|---|---|
| Cloud ID | `7fa05bac-b453-4b39-9ec3-830a6365e08a` |
| Project key | `AV` |
| Issue type | `Bug` |

### Environment Impacted (`customfield_10106`) — required

| Option | ID |
|---|---|
| Sandbox | `10100` |
| Dev | `10101` |
| Stage | `10102` |
| Prod | `10103` |

Bugs found during sandbox validation use **Sandbox** (`10100`) unless otherwise known.

---

## Summary format

```
<terraform_resource_or_datasource> <what is wrong>
```

Examples:
- `couchbase-capella_collection accepts negative max_ttl without validation error`
- `couchbase-capella_buckets datasource crashes with Value Conversion Error`

Under 100 characters. No trailing punctuation.

---

## Writing the description — general rules

Keep it short and human. Write like a developer who hit the bug and is telling a colleague.
- One sentence opening: what you did and what went wrong.
- No verbose explanations of things the reader can see in the error or code.
- Only include a section if it adds information not already obvious from the others.
- Never pad with "This means operators cannot..." or "This is actively incorrect because..." — state the impact once, plainly.

---

## Description format — Mode A / Mode C

```markdown
<One sentence: what you did and what happened.>

**Actual:**
<Error or state — in a code block.>

**Expected:**
<What should happen.>

**Impact:**
<One line: what breaks for the operator.>

**Steps to reproduce**

<Minimal HCL — in an hcl code block.>

<Command and output that shows the bug.>
```

---

## Description format — Mode B

```markdown
<One sentence: what the user sees.>

**Actual:**
<Error — in a code block.>

**Expected:**
<What should happen.>

**Impact:**
<One line.>

**Root cause:**
<File and the specific problem — one or two sentences max.>

**Steps to reproduce**

<Minimal HCL — in an hcl code block.>

<Error output.>

**Fix:**
<What to change and where — one or two sentences.>

**Branch:** `<branch-name>`
```

---

## Step-by-step procedure

### 1. Determine the mode (A, B, or C)

### 2. Collect inputs

| Input | Mode A/C | Mode B |
|---|---|---|
| What you did and what happened | ✓ | ✓ |
| Actual error/state | ✓ | ✓ |
| Expected behaviour | ✓ | ✓ |
| Impact | ✓ | ✓ |
| Minimal HCL to reproduce | ✓ | ✓ |
| Root cause (file + reason) | — | ✓ |
| Fix (file + change) | — | ✓ |
| Branch name | — | ✓ |
| Test name / file path | ✗ never | ✗ never |

### 3. Create the issue

```json
{
  "cloudId": "7fa05bac-b453-4b39-9ec3-830a6365e08a",
  "projectKey": "AV",
  "issueTypeName": "Bug",
  "summary": "<summary>",
  "contentFormat": "markdown",
  "description": "<formatted description>",
  "additional_fields": {
    "customfield_10106": { "id": "<environment_id>" },
    "labels": ["terraform-provider"]
  }
}
```

### 4. Report back

Output:
- Jira key and URL
- One-line bug summary
- Which mode was used and why

---

## Optional: assign the ticket

```
mcp__claude_ai_Atlassian__lookupJiraAccountId(searchString: "<name or email>")
```

Then pass `assignee_account_id` to `createJiraIssue` or `editJiraIssue`.

---

## Acceptance test update rule (separate from ticket content)

When a test with `ExpectError` finds that no error is returned:
1. Remove `ExpectError`.
2. Add a `Check` block asserting the actual (broken) state.
3. Add a single-line comment above the function: `// <TICKET>: <what the bug is>; restore ExpectError once fixed.`

When a provider crash makes the whole step fail before checks run:
1. Add `ExpectError` matching the crash message.
2. Add a single-line comment above the function with the same format.

These test updates are internal — never copy them into the Jira ticket.

---

## Real examples

### AV-132307 — Backend bug (Mode A)

**Bug:** API accepts `max_ttl = -1` on a collection with no error.
**Opening:** "While creating a collection with `max_ttl` set to `-1`, the Capella API accepted the request without returning a validation error."

### AV-132308 — Provider bug fixed by us (Mode B)

**Bug:** `couchbase-capella_buckets` datasource crashes with a Value Conversion Error.
**Opening:** "When using the `couchbase-capella_buckets` datasource to list buckets for a cluster, Terraform fails immediately with a Value Conversion Error."
**Root cause (short):** "`OneBucket` struct has a `vbuckets` tfsdk field that is missing from `BucketsSchema()` in `internal/datasources/buckets_schema.go`."
**Fix (short):** "Add `vbuckets` as a computed int64 to `BucketsSchema()`."
