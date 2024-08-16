// src/data/content.ts

export const content = {
  home: {
    title: "Welcome to My Custom HTTP Web Server 🌐",
    description:
      "This is a minimalistic web server built from scratch using TypeScript and Node.js. 🦋",
    links: [
      { text: "About", url: "/about" },
      { text: "API Users", url: "/api/users" },
    ],
  },
  about: {
    title: "About This Project",
    content: `
      <h1>About My Custom HTTP Web Server 🌐</h1>

      <p>This project is a small, custom-built HTTP web server I made to learn how web servers work from the ground up.</p>

      <h2>Why I Built This 🛠️</h2>

      <p>I started this project because I was curious. I wanted to know what happens behind the scenes when you type a web address and hit enter. Here's what I aimed to learn:</p>

      <ul>
        <li>How HTTP really works - the language of the web</li>
        <li>How to use Node.js for networking at a low level</li>
        <li>How web servers handle different web addresses (URLs)</li>
        <li>How web servers deal with many users at once</li>
      </ul>

      <h2>How It Works ⚙️</h2>

      <p>This server uses Node.js's 'net' module to create a TCP server. It's like setting up a phone line for computers to talk to each other. Here's what happens when you visit a webpage on this server:</p>

      <ol>
        <li>The server waits for a computer (like your web browser) to call</li>
        <li>When a call comes in, the server listens to what the browser is asking for</li>
        <li>The server looks at the web address to figure out what information to send back</li>
        <li>It prepares the right information (like the text of a webpage)</li>
        <li>Finally, it sends this information back to your browser</li>
      </ol>

      <h2>The Parts of the Server 🧩</h2>

      <p>This project is made up of several files, each with a specific job:</p>

      <ul>
        <li><strong>index.ts</strong>: This is like the main control center. It starts everything up.</li>
        <li><strong>Server.ts</strong>: This file sets up the "phone line" for browsers to connect to.</li>
        <li><strong>RequestHandler.ts</strong>: This listens to what the browser is asking for.</li>
        <li><strong>ResponseBuilder.ts</strong>: This prepares the information to send back to the browser.</li>
        <li><strong>Router.ts</strong>: This decides what information to send based on the web address.</li>
        <li><strong>HttpUtils.ts</strong>: This has helper tools for understanding browser requests.</li>
      </ul>

      <h2>Why Not Use Existing Tools? ❓</h2>

      <p>You might wonder, "Why build a web server when there are already good ones out there?" It's a fair question! I'm not trying to replace those servers. Instead, I'm learning how they work by building a simple version myself.</p>

      <p>It's like learning to bake bread from scratch. You might not do it every day when you can buy bread at the store, but the process teaches you a lot about what goes into making bread.</p>


      <h2>What's Next? ⏭️</h2>

      <p>This project is just the beginning. There's still so much to learn and add:</p>

      <ul>
        <li>Handling different types of requests (like POST for sending data)</li>
        <li>Adding security features to protect against common web attacks</li>
        <li>Making the server faster by handling multiple requests at the same time</li>
        <li>Adding a way to serve files (like images or CSS) along with web pages</li>
      </ul>

      <p>While this server isn't ready for real-world websites, it's an excellent tool for learning. It's helped me understand what's really happening when I browse the web, and I hope it can help others learn too!</p>
    `,
  },
  users: [
    {
      id: 1,
      name: "Alice Johnson",
      email: "alice@example.com",
      role: "Developer",
    },
    { id: 2, name: "Bob Smith", email: "bob@example.com", role: "Designer" },
    {
      id: 3,
      name: "Charlie Brown",
      email: "charlie@example.com",
      role: "Manager",
    },
    { id: 4, name: "Diana Ross", email: "diana@example.com", role: "DevOps" },
    { id: 5, name: "Ethan Hunt", email: "ethan@example.com", role: "Tester" },
  ],
};

// Helper function to create an HTML page
export const createHtmlPage = (body: string) => `
  <!DOCTYPE html>
  <html lang="en">
  <head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Custom Web Server</title>
    <style>
      body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 800px; margin: 0 auto; padding: 20px; }
      h1 { color: #2c3e50; }
      a { color: #3498db; text-decoration: none; }
      a:hover { text-decoration: underline; }
    </style>
  </head>
  <body>
    ${body}
  </body>
  </html>
`;
