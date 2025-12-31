data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "ariga.io/atlas-provider-gorm",
    "load",
    "--path", "./internal/database/models",
    "--dialect", "postgres",
  ]
}

data "composite_schema" "app" {
  # Load enum types first.
  schema "public" {
    url = "file://internal/database/enum-schema.sql"
  }
  # Then, load the GORM models.
  schema "public" {
    url = data.external_schema.gorm.url
  }
}

env "gorm" {
  src = data.composite_schema.app.url
  dev = "docker://postgres/15/dev"
  migration {
    dir = "file://internal/database/migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}

env "local" {
  url = "postgresql://postgres:postgres@localhost:5432/statify?search_path=public&sslmode=disable"
  migration {
    dir = "file://internal/database/migrations"
  }
}