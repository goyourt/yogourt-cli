Yogourt-cli is made to init and use Yogourt projects easily.

Project initialization

```shell
go install github.com/goyourt/yogourt/cli@latest
yogourt init <project-name>
docker-compose up -d
yogourt migrate
go run main.go
```

Your yogourt app is now running at http://localhost:8080

You can configure your app from the config file `config.yaml`