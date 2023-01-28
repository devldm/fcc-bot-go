# **FCC-Bot**

`ver. 2.3.3`

## **Description**

FCC-Bot is a discord bot used on the FCC Seoul Discord server.
It is written in Go.

Its primary library is [discordgo](https://github.com/bwmarrin/discordgo).

---

## **Getting Started**


### **Installation**

1. Clone the repository.
2. Run `go get` to install the packages needed.


### **Running**

#### **Environment Variables**

You will need to create .env files before running the bot.

| **env variable**            | **how to obtain**                                                                   |
|---------------------------|------------------------------------------------------------------------------------|
| BOT_TOKEN                 | The private key Discord gives you for a bot                                        |
| BOT_ID                    | The bot's user id: right click the bot in the server                               |
| GUILD_ID                  | The guild's id: right click the server's name in discord                           |
| LOG_CHANNEL               | The log channel's id: create a channel, and right click it                         |
| INTRO_CHANNEL             | The introduction channel's id: create a channel, and right click it                |
| LEARNING_RESOURCE_CHANNEL | The learning resource channel's id: create a channel, and right click it           |
| RFR_POST                  | The post to listen for react-for-role reactions: create a post, and right click it |
| ROLE_VERIFIED             | The role users receive when validated: create a role, and right click it           |

#### **Commands**

| **Command**                   | use                                                 |
|-------------------------------|-----------------------------------------------------|
|`go run ./app/*.go`             | Run the app using dev envs                          |
| `go run ./app/*.go -p`          | Run the app using production envs                   |
| `go build -o fccbot ./app/*.go` | Build the app to a file named 'fccbot'              |
| `make run`                      | Makefile command for running the app in dev mode    |
| `make build`                    | Makefile command for building the app               |
| `make deploy`                 | Makefile command to run the app as a process in PM2 |
| `make stop`                    | Makefile command to stop the app with PM2           |

>**_NOTE:_** The above commands assume you are running the bot on Mac or Linux, and have PM2 installed globally.