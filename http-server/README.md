# 🚀 My Custom HTTP Web Server

## 🧐 About This Project

This is a lightweight, custom-built HTTP web server created for learning purposes. It's designed to help understand the fundamentals of how web servers work under the hood.

I built this server because I was curious about what happens behind the scenes when you type a web address and hit enter.

## 🛠️ How It Works

This server uses Node.js's 'net' module to create a TCP server. Here's a simple breakdown of the process:

1. 🎧 The server waits for a computer (like your web browser) to connect
2. 📥 When a connection comes in, the server listens to what the browser is asking for
3. 🔍 The server looks at the web address to figure out what information to send back
4. 📝 It prepares the right information (like the text of a webpage)
5. 📤 Finally, it sends this information back to your browser

![Design](https://github.com/user-attachments/assets/95989302-3c75-482d-b269-483e3d182830)

## 📁 Project Structure

- `index.ts`: The main control center. It starts everything up.
- `Server.ts`: Sets up the "phone line" for browsers to connect to.
- `RequestHandler.ts`: Listens to what the browser is asking for.
- `ResponseBuilder.ts`: Prepares the information to send back to the browser.
- `Router.ts`: Decides what information to send based on the web address.
- `HttpUtils.ts`: Helper tools for understanding browser requests.

## 🚀 Usage

1. Clone the repository:
   ```
   git clone https://github.com/mutasim77/lets-build-x/http-server.git
   ```

2. Install dependencies:
   ```
   yarn
   ```

3. Run the server:
   ```
   yarn dev
   ```

4. Visit `http://localhost:3000` in your web browser. 🌐


## 🔮 Future Enhancements
- 📮 Handle different types of requests (like POST for sending data)
- 🔒 Add security features to protect against common web attacks
- ⚡ Make the server faster by handling multiple requests at the same time
- 🖼️ Add a way to serve files (like images or CSS) along with web pages

## 🤝 Collaboration
I'm always open to learning from others! If you have ideas for improvements or want to contribute, please feel free to open an issue or submit a pull request.

## 📜 License
This project is open source and available under the [MIT License](./LICENSE).

---
Thank you for checking out my custom HTTP web server project. I hope it helps you learn about web server architecture as much as it helped me.

Happy coding, and enjoy exploring how the web works under the hood! 🎉🖥️
