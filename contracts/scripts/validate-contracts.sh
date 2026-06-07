#!/usr/bin/env bash
set -euo pipefail

root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

python3 - <<'PY' "$root_dir"
import json
import sys
from pathlib import Path

root = Path(sys.argv[1])
errors = []


def load_json(path):
    try:
        return json.loads(path.read_text())
    except json.JSONDecodeError as exc:
        errors.append(f"{path.relative_to(root)}: invalid JSON: {exc}")
        return None


def check_value(path, location, schema, value):
    if "const" in schema and value != schema["const"]:
        errors.append(f"{path.relative_to(root)}: {location} must equal {schema['const']!r}")

    expected_type = schema.get("type")
    if expected_type == "object" and not isinstance(value, dict):
        errors.append(f"{path.relative_to(root)}: {location} must be object")
        return
    if expected_type == "string" and not isinstance(value, str):
        errors.append(f"{path.relative_to(root)}: {location} must be string")
        return
    if expected_type == "integer" and not isinstance(value, int):
        errors.append(f"{path.relative_to(root)}: {location} must be integer")
        return

    if isinstance(value, str) and len(value) < schema.get("minLength", 0):
        errors.append(f"{path.relative_to(root)}: {location} is shorter than minLength")

    enum = schema.get("enum")
    if enum and value not in enum:
        errors.append(f"{path.relative_to(root)}: {location} must be one of {enum}")

    if expected_type == "object":
        required = schema.get("required", [])
        for key in required:
            if key not in value:
                errors.append(f"{path.relative_to(root)}: {location}.{key} is required")
        for key, child_schema in schema.get("properties", {}).items():
            if key in value:
                check_value(path, f"{location}.{key}", child_schema, value[key])


schemas = {}
for path in sorted(root.glob("events/*.schema.json")):
    schema = load_json(path)
    if schema is not None:
        schemas[path.name.replace(".schema", "")] = schema

for path in sorted(root.glob("samples/*.json")):
    sample = load_json(path)
    schema = schemas.get(path.name)
    if sample is not None and schema is None:
        errors.append(f"{path.relative_to(root)}: missing matching event schema")
    elif sample is not None:
        check_value(path, "$", schema, sample)

for path in sorted(root.glob("openapi/*.yaml")):
    text = path.read_text()
    for token in ("openapi: 3.0.3", "securitySchemes:", "serviceJwt:"):
        if token not in text:
            errors.append(f"{path.relative_to(root)}: missing {token}")

    operation = None
    operation_headers = {}
    for line in text.splitlines():
        stripped = line.strip()
        if stripped == "operationId: getUserProfile":
            operation = "getUserProfile"
            operation_headers[operation] = set()
        elif stripped.startswith("operationId: "):
            operation = stripped.removeprefix("operationId: ")
            operation_headers[operation] = set()
        elif operation and "#/components/parameters/IdempotencyKey" in stripped:
            operation_headers[operation].add("IdempotencyKey")
        elif operation and "#/components/parameters/CorrelationId" in stripped:
            operation_headers[operation].add("CorrelationId")

    for operation_id, headers in operation_headers.items():
        if "CorrelationId" not in headers:
            errors.append(f"{path.relative_to(root)}: {operation_id} missing X-Correlation-ID")
        if operation_id != "getUserProfile" and "IdempotencyKey" not in headers:
            errors.append(f"{path.relative_to(root)}: {operation_id} missing Idempotency-Key")

if errors:
    for error in errors:
        print(error, file=sys.stderr)
    sys.exit(1)
PY
