import * as net from "net";
import { RequestHandler } from "./RequestHandler";
import { Router } from "../router/Router";

export class Server {
  private server: net.Server;

  constructor(
    private port: number,
    private router: Router,
  ) {
    // Create a TCP server
    this.server = net.createServer((socket) => this.handleConnection(socket));
  }

  // Method to start the server
  public start(): void {
    this.server.listen(this.port, () => {
      console.log(`Server listening on port ${this.port}`);
    });
  }

  // Method to handle incoming connections
  private handleConnection(socket: net.Socket): void {
    // Create a RequestHandler for each connection
    const requestHandler = new RequestHandler(socket, this.router);

    // Handle incoming data
    socket.on("data", (data) => {
      requestHandler.handleRequest(data);
    });

    // Handle connection errors
    socket.on("error", (error) => {
      console.error("Socket error:", error);
    });
  }
}
