module roersla.no/askeladden

require (
	github.com/bwmarrin/discordgo v0.29.0
	gopkg.in/yaml.v3 v3.0.1
)

go 1.24.1

replace roersla.no/askeladden => ./

require (
	github.com/go-sql-driver/mysql v1.9.3
	github.com/google/uuid v1.6.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	golang.org/x/crypto v0.0.0-20210421170649-83a5a9bb288b // indirect
	golang.org/x/sys v0.0.0-20201119102817-f84b799fce68 // indirect
)
