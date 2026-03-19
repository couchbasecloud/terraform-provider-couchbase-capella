---
applyTo: "examples/**"
---

# Examples Review Instructions

Guidelines for reviewing terraform provider capella examples directories in this repository.

## Writing Style

---

## Directory Structure

Each examples directory should be organisated as such:

```
examples/
  └── resource_name/
      ├── <crud>_<resource_name>.tf (optional, but must have at least one CRUD example)
      ├── list_<resource_name>.tf (optional, but at least one of get or list)
      ├── get_<resource_name>.tf (optional, but at least one of get or list)
      ├── main.tf
      ├── terraform.template.tfvars
      ├── variables.tf
      └── README.md
```

Where `<crud>` is one of `create`, `read`, `update`, or `delete`.

---

## Terraform Configuration Files

- All Terraform configuration files should be syntactically correct and follow best practices for Terraform code.
- Variable names should be descriptive and consistent across examples.
- The `main.tf` file should contain a complete example that can be applied end-to-end, while the CRUD files should focus on demonstrating specific operations.
- The `terraform.template.tfvars` file should include all necessary variables with placeholder values.

## README.md

All written text should be clear, concise, and free of spelling or grammatical errors. If it is possible to cut a word out, cut it out. The writing style should be consistent across all examples, and the tone should be professional yet approachable.

The README should include:

- A brief description of the example's purpose and functionality.
- Example walkthrough including sample outputs.
