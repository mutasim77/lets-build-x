# 🚀 Redis Clone : Learning by Building

## 🧠 What is this project?
This is a simple Redis clone built in Go. I created it because I wanted to understand how Redis works internally.
It's not about making a better Redis, but about learning how things work under the hood.

## 🏗️ Why I made this
I've always been curious about how databases like Redis actually work. Instead of just reading about it, I decided to build my own version. This project helped me learn:

- How in-memory databases store and retrieve data
- How Redis handles different data types
- How the Redis protocol (RESP) works
- How to make data persist even after the server restarts

## 🛠️ What can it do?
This Redis clone can:
- Store and retrieve key-value pairs
- Handle basic Redis commands (GET, SET, DEL, EXISTS, INCR)
- Talk to Redis clients using the Redis protocol
- Save data to disk so it's not lost when the server stops

## 🚀 Getting Started
1. Clone the repository:
```bash
git clone https://github.com/mutasim77/lets-build-x.git
```

2. Go to the project folder:
```bash
cd lets-build-x/redis
```

3. Start the server
```bash
make run-redis
```

4. You can now connect to it using any Redis client!

## 🔧 Command Arsenal
My Redis comes packed with a set of powerful commands:

- `PING`: The classic "Are you there?" command.
- `SET key value`: Store a key-value pair.
- `GET key`: Retrieve a value by its key.
- `DEL key [key ...]`: Delete one or more keys.
- `EXISTS key [key ...]`: Check if one or more keys exist.
- `INCR key`: Increment the integer value of a key.

## 🧪 Testing Your Metal
I believe in the power of testing! Run test suite to ensure everything's working smoothly:
```bash
make test
```

## 🎨 Project Structure
Here's a quick tour of projects' codebase:
- `main.go`: Starts the server
- `server/`: Handles client connections
- `command/`: Implements Redis commands
- `storage/`: Manages data storage
- `resp/`: Handles the Redis protocol
- `aof/`: Manages saving data to disk

## 🤝 Contributing
Got ideas? Found a bug? Want to add a cool feature? I'm all ears! Feel free to open issues, submit pull requests, or just share your thoughts.

## 📜 License
This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.


Remember, the goal here isn't to replace Redis, but to learn by doing.

Happy coding! ❤️
