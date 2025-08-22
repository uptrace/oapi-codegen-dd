# Migrating from v2

There are some incompatible changes that were introduced in v3 of the codegen.<br/>

## Extensions:

The following extensions are no longer supported:
- `x-order`
- `x-oapi-codegen-only-honour-go-name`

## User templates

HTTP path is not supported

## Custom name normalizer

Not supported

## Server code generation

Not supported

## Configuration changes

```yaml
package: ✅
generate: ✅
    iris-server: ❌
    chi-server: ❌
    fiber-server: ❌
    echo-server: ❌
    gin-server: ❌
    gorilla-server: ❌
    std-http-server: ❌
    strict-server: ❌
    client: ✅
      🆕🐣new properties:
        name: string
        timeout: time.duration
    models: ❌ always generated
    embedded-spec: ❌
    server-urls: ❌
  🆕🐣new properties:
    omit-description: bool
    default-int-type: "int64"
compatibility: ❌
output-options: ➡️renamed to output
    skip-fmt: ❌
    skip-prune: ➡ moved to config root
    include-tags: ➡ moved to filter include
    exclude-tags: ➡ moved to filter.exclude
    include-operation-ids: ➡ moved to filter.include
    exclude-operation-ids: ➡ moved to filter.exclude
    user-templates: ➡ moved to the config root
    exclude-schemas: ❌moved to filter.exclude
    response-type-suffix: ❌
    client-type-name: ➡ moved to generate.client.name
    initialism-overrides: ❌
    additional-initialisms: ❌
    nullable-type: ❌
    disable-type-aliases-for-type: ❌
    name-normalizer: ❌
    overlay: ❌there is spec merge functionality
    yaml-tags: ❌
    client-response-bytes-function: ❌
    prefer-skip-optional-pointer: ❌
    prefer-skip-optional-pointer-with-omitzero: ❌
    prefer-skip-optional-pointer-on-container-types: ❌

  🆕🐣new properties:
    use-single-file: bool
    directory: string
    filename: string
import-mapping: ❌
additional-imports: ✅
```
