def method:
  . == "get" or . == "put" or . == "post" or . == "delete" or
  . == "patch" or . == "options" or . == "head" or . == "trace";

def component_ref:
  type == "string" and startswith("#/components/");

def refs:
  [.. | objects | .["$ref"]? // empty | select(component_ref)] | unique;

def component_path($ref):
  $ref | ltrimstr("#/") | split("/");

def expand_refs($doc):
  (. + [
    .[] as $ref
    | (component_path($ref)) as $path
    | ($doc | getpath($path)? | refs[])
  ]) | unique;

def reachable_refs($doc; $roots):
  ($roots | refs)
  | until((expand_refs($doc) - .) | length == 0; expand_refs($doc));

def prune_component_map($section; $refs):
  with_entries(
    .key as $key
    | select($refs | index("#/components/" + $section + "/" + $key))
  );

def prune_components($refs):
  .components.schemas = ((.components.schemas // {}) | prune_component_map("schemas"; $refs))
  | .components.parameters = ((.components.parameters // {}) | prune_component_map("parameters"; $refs))
  | .components.requestBodies = ((.components.requestBodies // {}) | prune_component_map("requestBodies"; $refs))
  | .components.responses = ((.components.responses // {}) | prune_component_map("responses"; $refs))
  | .components.headers = ((.components.headers // {}) | prune_component_map("headers"; $refs))
  | .components.examples = ((.components.examples // {}) | prune_component_map("examples"; $refs))
  | .components.links = ((.components.links // {}) | prune_component_map("links"; $refs))
  | .components.callbacks = ((.components.callbacks // {}) | prune_component_map("callbacks"; $refs));

def filter_paths($tag):
  with_entries(
    .value |= with_entries(
      select((.key | method | not) or ((.value.tags // []) | index($tag)))
    )
    | select(.value | keys | map(select(method)) | length > 0)
  );

. as $full
| (.paths | filter_paths($tag)) as $paths
| ($full | .paths = $paths) as $filtered
| reachable_refs($filtered; $paths) as $reachable
| $filtered
| prune_components($reachable)
