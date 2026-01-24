---
title: uuidgenerate
parent: CLI
---

## devtui uuidgenerate

Generate a UUID

### Synopsis

Generate a UUID of a specified version.

By default, generates a version 4 UUID. Versions 3 and 5 accept a namespace value.

```bash
devtui uuidgenerate [flags]
```

### Examples

```bash
# Generate a default UUID (v4)
devtui uuidgenerate
# Generate a UUID v7
devtui uuidgenerate --uuid-version 7
# Generate a UUID v3 with namespace
devtui uuidgenerate --uuid-version 3 --namespace example.com
```

### Options

```
  -h, --help               help for uuidgenerate
  -n, --namespace string   namespace for UUID v3/v5 generation
  -v, --uuid-version int   UUID version to generate (1-7) (default 4)
```
