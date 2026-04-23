# koanf

Provider置き場

## API

### Providers

Install with `go get -u github.com/naycoma/koanf/providers/$provider`

| Package | Provider                                                               | Description    |
| ------- | ---------------------------------------------------------------------- | -------------- |
| fetch   | `fetch.Provider(url string)`                                           | URL Fetching   |
| gist    | `gist.Provider(user string, id string, file string)`                   | GitHub Gist    |
| mongo   | `mongo.Provider(uri, database, collection string, filter ...bson.M)`   | MongoDB (BSON) |
| walk    | `walk.Provider(parser koanf.Parser, root string, fileNames ...string)` | WalkDir        |

### Parsers

| Package | Provider                        | Description           |
| ------- | ------------------------------- | --------------------- |
| expand  | `expand.Parser(p koanf.Parser)` | wrap `os.ExpandEnv()` |
