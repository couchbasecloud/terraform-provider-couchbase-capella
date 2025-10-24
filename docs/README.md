# Templates Directory

This directory contains custom documentation templates for the Terraform provider that are processed by `tfplugindocs`.

## How It Works

When you run `make build-docs`, the `tfplugindocs` tool:
1. Generates documentation from Go schema definitions
2. Copies any files from `templates/` to `docs/` 
3. Preserves custom content that would otherwise be overwritten

## Directory Structure

```
templates/
  ├── guides/              # Custom upgrade guides and other guides
  │   ├── 1.1.0-upgrade-guide.md
  │   ├── 1.2.0-upgrade-guide.md
  │   └── ...
  ├── index.md.tmpl        # Custom provider homepage
  └── README.md            # This file
```

## Adding New Upgrade Guides

To add a new upgrade guide:

1. Create the file in `templates/guides/`:
   ```bash
   touch templates/guides/1.6.0-upgrade-guide.md
   ```

2. Write your upgrade guide content

3. Run `make build-docs` to generate the documentation
   - The guide will appear in `docs/guides/`

## Important Notes

⚠️ **DO NOT edit files in `docs/` directly** - they will be overwritten!

- ✅ **DO** edit files in `templates/` 
- ✅ **DO** edit Go schema files (e.g., `internal/resources/*_schema.go`)
- ✅ **DO** run `make build-docs` after making changes
- ❌ **DON'T** manually edit generated files in `docs/`

## Custom Templates

For more advanced customization, you can create `.tmpl` files:
- `index.md.tmpl` - Provider homepage
- `resources/<resource>.md.tmpl` - Custom resource documentation
- `data-sources/<datasource>.md.tmpl` - Custom data source documentation

See [tfplugindocs documentation](https://github.com/hashicorp/terraform-plugin-docs) for template syntax.

