---
title: license
parent: CLI
---

## devtui license

License management commands

### Synopsis

Commands for activating, validating, and deactivating licenses

### Options

```
  -h, --help   help for license
```

## devtui license activate

Activate a license

### Synopsis

Activate a license

```bash
devtui license activate [flags]
```

### Examples

```bash
devtui license activate --key=YOUR_LICENSE_KEY
```

### Options

```
  -h, --help         help for activate
      --key string   License key
```

## devtui license deactivate

Deactivate a license

### Synopsis

Deactivate a license

```bash
devtui license deactivate [flags]
```

### Examples

```bash
devtui license deactivate
```

### Options

```
  -h, --help   help for deactivate
```

## devtui license validate

Validate a license

### Synopsis

Reads a license and validates it

```bash
devtui license validate [flags]
```

### Examples

```bash
devtui license validate
```

### Options

```
  -h, --help   help for validate
```
