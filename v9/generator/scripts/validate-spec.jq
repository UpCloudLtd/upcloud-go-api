def operations:
  .paths[] | to_entries[] |
  select(.key | IN("get", "put", "post", "delete", "options", "head", "patch", "trace")) |
  .value;

def valid_input:
type == "object" and
(.openapi | type == "string" and test("\\A3\\.[01]\\.[0-9]+\\z")) and
(.info | type == "object") and
(.info.title | type == "string" and length > 0) and
(.info.version | type == "string" and length > 0) and
(.paths | type == "object" and length > 0) and
(.tags | type == "array" and length > 0) and
all(.tags[]; .name | type == "string" and test("\\A[A-Za-z][A-Za-z0-9_-]*\\z")) and
([.tags[].name] as $names | ($names | unique | length) == ($names | length)) and
([operations] | length > 0) and
([.tags[].name] as $names |
  all(operations;
    (.tags | type == "array" and length == 1) and
    all(.tags[]; . as $tag | $names | index($tag) != null)));

  length == 1 and (.[0] | valid_input)