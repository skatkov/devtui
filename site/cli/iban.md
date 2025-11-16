---
title: iban
parent: CLI
---

## devtui iban

Generate test IBAN numbers

### Synopsis

Generate test IBAN numbers for a given country code using the banking library.

The country code is required and should be a valid ISO 3166-1 alpha-2 country code.
Use the --formatted flag to output the IBAN in paper format with spaces.

```bash
devtui iban <country-code> [flags]
```

### Examples

```bash
# Generate IBAN for Great Britain
devtui iban GB
# Generate formatted IBAN for Germany
devtui iban DE --format
devtui iban DE -f
# Generate IBAN for France
devtui iban FR
```

### Options

```
  -f, --format   output IBAN in paper format with spaces
  -h, --help     help for iban
```
