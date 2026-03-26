# 🤖 Droid Runbook: The Fun Edition

### _"These aren't the data sources you're looking for..." — Actually, they are._

A delightfully irreverent (but accurate) guide to using AI droids to generate Terraform data sources for the Couchbase Capella provider.

---

## 🧠 WTF Is a Droid?

Think of a droid as your AI pair programmer that actually reads the docs.

You tell it what you want. It reads the OpenAPI spec. It writes the code. You review it. Ship it.

```
  You: "Generate a data source for Buckets"
  
  Droid: *generates 6 files, registers them, writes tests*
  
  You: 👁️👄👁️
```

> **"It's not stupid if it works."** — Every engineer who let an AI write their boilerplate

---

## 📋 Requirements

| Tool | Version | Vibe Check |
|---|---|---|
| Git | any | If you don't have this, we need to talk |
| Go | >= 1.21 | `go version` or go home |
| Terraform | >= 1.5.2 | The reason we're all here |
| AI Droid | with workspace access | Your new best friend |
| OpenAPI spec | `openapi.generated.yaml` | The sacred scroll |

---

## 🚀 The Workflow (a.k.a. "How to Vibe Code Responsibly")

### Step 1: Find Your Feature

```
   📖 OpenAPI Spec
       |
       |  grep -i "buckets" openapi.generated.yaml
       |
       v
   🎯 GET endpoint  → single resource data source
   🎯 LIST endpoint → list data source
```

No GET endpoint? Skip the single resource. No LIST endpoint? Skip the list. Both missing? Wrong spec, buddy.

> **"404: Feature Not Found"** — You, looking at the wrong section of the spec

---

### Step 2: Summon the Droid

Tell the droid:

> "Generate Terraform data sources for **[Feature]** using the `tf-datasource-gen` skill and the OpenAPI spec."

Then sit back and watch the magic happen. ✨

```
  ┌─────────────────────────────────┐
  │  Droid: "I have generated the   │
  │  following files..."            │
  │                                 │
  │  feature.go ✅                  │
  │  features.go ✅                 │
  │  feature_schema.go ✅           │
  │  features_schema.go ✅          │
  │  api structs ✅                 │
  │  provider registration ✅       │
  │  acceptance tests ✅            │
  │                                 │
  │  Me:                            │
  │  ┌───────────────────────────┐  │
  │  │  (ノ◕ヮ◕)ノ*:・゚✧        │  │
  │  └───────────────────────────┘  │
  └─────────────────────────────────┘
```

---

### Step 3: Review (Yes, You Still Have to Do This)

> **"Trust, but verify."** — Ronald Reagan, definitely talking about AI-generated Terraform providers

#### The Checklist of Doom™

**✅ Uses `ClientV1`?**
```go
// GOOD — This is the way
response, err := s.ClientV1.ExecuteWithRetry(...)

// BAD — We don't do that here
response, err := s.Client.Execute(...)
```

If the droid used an old client:

```
  ┌──────────────────────────────────────┐
  │                                      │
  │  "You were supposed to use ClientV1, │
  │   not destroy it!"                   │
  │                — Obi-Wan, probably   │
  │                                      │
  └──────────────────────────────────────┘
```

**✅ Struct embeds `*providerschema.Data`?**
```go
type Buckets struct {
    *providerschema.Data   // <- This. Always this.
}
```

**✅ Interface compliance assertions?**
```go
var (
    _ datasource.DataSource              = (*Buckets)(nil)
    _ datasource.DataSourceWithConfigure = (*Buckets)(nil)
)
```

> No assertions = no proof it works = 🔥

**✅ Schema validators on IDs?**
```go
capellaschema.AddAttr(attrs, "organization_id", builder, requiredStringWithValidator())
//                                                         ^^^^^^^^^^^^^^^^^^^^^^^^^^
//                          If this is missing, someone will pass an empty string
//                          and have a Very Bad Time™
```

**✅ Registered in `provider.go`?**

If your data source isn't in `DataSources()`, it literally does not exist to Terraform.

```
  Terraform: "I don't know her" 💅
```

---

### Step 4: The Build-Fix Loop

Also known as "The Five Stages of Grief (but only 5 retries)":

```bash
# 1. Format
goimports -w -local github.com/couchbasecloud/terraform-provider-couchbase-capella your_file.go

# 2. Vet
go vet ./internal/datasources/...

# 3. Build
VERSION=$(git describe --tags --abbrev=0)
go build -ldflags "-s -w -X 'github.com/couchbasecloud/terraform-provider-couchbase-capella/version.ProviderVersion=$VERSION'" -o ./bin/terraform-provider-couchbase-capella
```

```
  Attempt 1: ❌ errors
  Attempt 2: ❌ different errors
  Attempt 3: ❌ same errors but you're angrier
  Attempt 4: ❌ why
  Attempt 5: ✅ IT BUILDS 🎉
  
  OR
  
  Attempt 5: ❌ report it and walk away
```

> **"I don't always test my code, but when I do, I do it in production."**
> — Please don't. Run the tests. ⬇️

---

### Step 5: Acceptance Tests

```bash
make testacc
```

> ⚠️ **WARNING:** Acceptance tests create REAL resources that cost REAL money. 
> This is not a drill. This is your cloud bill.

```
  Your wallet:    📉📉📉
  Your coverage:  📈📈📈
```

Tests MUST use `resource.ParallelTest()`. Sequential tests are for people who enjoy watching paint dry.

---

## 📁 Where Things Go

```
internal/
├── api/                  # Structs go here (the "what does the API return" stuff)
├── datasources/          # The main event 🎭
│   ├── feature.go             # GET one thing
│   ├── features.go            # GET all the things  
│   ├── feature_schema.go      # What one thing looks like
│   └── features_schema.go     # What all the things look like
├── provider/
│   └── provider.go       # Register here or it doesn't exist
└── schema/               # Shared helpers

acceptance_tests/         # Prove it works (with real money)

.factory/
└── skills/
    └── tf-datasource-gen/
        └── SKILL.md      # The droid's instruction manual 📜
```

---

## 🔥 Common Fails & Fixes

| What Went Wrong | What To Do | Mood |
|---|---|---|
| Droid used old API client | Say "use ClientV1" | 😤 |
| Missing validators on IDs | Add `requiredStringWithValidator()` | 🤦 |
| Build fails repeatedly | `goimports` → `go vet` → build → repeat | 😮‍💨 |
| Data source invisible to TF | Register it in `provider.go` | 🫥 |
| Droid read `openapi.gen.go` | Tokens gone. Context destroyed. Start over. | 💀 |

> **NEVER let the droid read `internal/generated/api/openapi.gen.go`.**
> It will consume all your tokens and you will have nothing to show for it.
> It's the `node_modules` of this repo.

---

## 🧘 Pro Tips

1. **One feature per droid session.** Don't ask it to boil the ocean.
2. **Be specific.** "Generate data sources for Buckets using GET and LIST endpoints from the OpenAPI spec" > "do the thing"
3. **Always check `git diff`.** Trust no one. Not even yourself.
4. **Iterate.** First pass not perfect? Point at the problem. Droids love targeted feedback.
5. **Read the skill file.** `.factory/skills/tf-datasource-gen/SKILL.md` is the source of truth. When in doubt, RTFM.

---

## 🏁 TL;DR

```
1. Find endpoints in OpenAPI spec
2. Tell droid to generate data sources
3. Review: ClientV1, validators, registration
4. goimports → go vet → go build (repeat ≤ 5x)
5. make testacc (RIP wallet)
6. Ship it 🚢
```

```
  ┌─────────────────────────────────────────┐
  │                                         │
  │  "In the beginning there was main.go,   │
  │   and the droid said 'Let there be      │
  │   data sources,' and there were data    │
  │   sources. And the code review saw      │
  │   that it was good."                    │
  │                                         │
  │            — Genesis 1:1 (Go Edition)   │
  │                                         │
  └─────────────────────────────────────────┘
```

---

_Made with 🤖 and questionable humor._

