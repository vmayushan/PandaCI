data "template_dir" "schema" {
  path = "./schema/templates"
  vars = {
    key = "value"
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
