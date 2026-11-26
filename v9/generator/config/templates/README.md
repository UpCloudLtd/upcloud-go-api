# oapi-codegen Template Patches

`client.tmpl` and `client-with-responses.tmpl` are rendered files. Their source
is the pinned upstream template in the oapi-codegen module, plus the local
patches in `patches/`. Do not edit the rendered files by hand.

## Upgrade oapi-codegen

From `v9`:

1. Set the new version: `go get github.com/oapi-codegen/oapi-codegen/v2@vX.Y.Z`.
2. Render the templates: `bash generator/scripts/render-templates.sh`.
3. If a patch fails with a reject, the script prints the reject path. Edit
   `patches/<name>.patch` to apply to the new upstream file. Then run step 2
   again.
4. Update the `refs/tags/` comments in `generator/config/*.yaml`.
5. Regenerate: `bash generator/scripts/generate.sh`, then
   `make generate-request-adapters`. Review the public API diff.

`make vet` runs `render-templates.sh check`. It fails when a rendered file
differs from upstream plus its patch.
