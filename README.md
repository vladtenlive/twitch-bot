Made by ChatGPT

This is a simple program for retrieving the viewer count for a Twitch user using the Twitch API. It was developed live on a Twitch stream using the Go programming language and includes a number of functions for interacting with the Twitch API, including retrieving user information and refreshing access tokens.

Dependencies:
- Go 1.16 or later
- github.com/joho/godotenv
- net/http
- encoding/json

To use the program, clone the repository and set the following environment variables in a .env file in the project directory:
- TWITCH_CLIENT_ID: Your Twitch app client ID
- TWITCH_CLIENT_SECRET: Your Twitch app client secret

Then, run the program using the 'go run' command:
$ go run main.go
