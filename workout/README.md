To re-generate swagger doc, run

```bash
go install github.com/swaggo/swag/cmd/swag@latest # if not already installed
swag init --parseDependency --parseInternal
```
