# Discord Ranking Bot

A simple bot to track how often users posted images (actually attachments) to a specified channel

## Getting Started

1. Clone the repository
2. [Create a Discord bot](https://discordapp.com/developers/applications/ "Discord's Developer page")
3. Insert your credentials into the script
4. ```go run main.go```
5. Go in your specified channel and send images.
6. Go in any channel and write ```!ranking n```, where n is an optional parameter and specifies the top n ranks to show

### Prerequisites

This bot was created using the go programming language, if you don't know go: [Here you _go_, hehe](https://golang.org/doc/install, "Getting Started with go")

It uses 2 libraries, [DiscordGo](https://github.com/bwmarrin/discordgo "DiscordGo Github page") and [go-cache](https://github.com/patrickmn/go-cache "go-cache Github page")


## Contributing

Contributing is more than welcome! 
1. Fork the repository
2. Create a new branch 
3. Implement/fix whatever you want
4. Make sure to use gofmt!
5. Open a pull request:)

## Planned additions:
- Commands to start / stop / reset / change the tracking process
- Option to save the ranking data so you can restart the bot without losing everything
## Authors

* **Jonas Scholz** - [Code42Cate](https://github.com/Code42Cate)

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE.md](LICENSE.md) file for details

## Acknowledgments

* Hat tip to anyone whose code was used:D
