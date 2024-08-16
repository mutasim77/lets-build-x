// src/index.ts

import { Server } from "./server/Server";
import { Router } from "./router/Router";
import { ParsedRequest } from "./utils/HttpUtils";
import { content, createHtmlPage } from "./content";

// Create a router
const router = new Router();

// Home route
router.addRoute("GET", "/", (request: ParsedRequest) => {
  const { title, description, links } = content.home;
  const linksHtml = links
    .map((link) => `<li><a href="${link.url}">${link.text}</a></li>`)
    .join("");
  const body = `
    <h1>${title}</h1>
    <p>${description}</p>
    <ul>${linksHtml}</ul>
  `;
  return {
    statusCode: 200,
    headers: { "Content-Type": "text/html" },
    body: createHtmlPage(body),
  };
});

// About route
router.addRoute("GET", "/about", (request: ParsedRequest) => {
  return {
    statusCode: 200,
    headers: { "Content-Type": "text/html" },
    body: createHtmlPage(content.about.content),
  };
});

// API Users route
router.addRoute("GET", "/api/users", (request: ParsedRequest) => {
  return {
    statusCode: 200,
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(content.users),
  };
});

// Create an instance of the Server class, listening on port 3000 and using our router
const server = new Server(3000, router);

// Start the server
server.start();

console.log("Server is running on http://localhost:3000");
