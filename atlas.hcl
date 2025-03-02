data "template_dir" "schema" {
  path = "./schema/templates"
  vars = {
    key = "value"
    // Pass the --env value as a template variable.
    env  = atlas.env
  }
}

env "psql" {
  src = data.template_dir.schema.url
  dev = "docker://postgres/17/dev?search_path=public"
  migration {
    dir = "file://schema/migrations"
  }
}