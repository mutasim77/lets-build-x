import * as net from "net";
import { Router } from "../router/Router";
import { ResponseBuilder } from "./ResponseBuilder";
import { parseRequest, ParsedRequest } from "../utils/HttpUtils";

export class RequestHandler {
  constructor(
    private socket: net.Socket,
    private router: Router,
  ) {}

  public handleRequest(data: Buffer): void {
    try {
      // Parse the raw request data
      const request: ParsedRequest = parseRequest(data.toString());

      // Get the appropriate handler from the router
      const handler = this.router.getHandler(request.method, request.url);

      // Execute the handler and get the response
      const responseData = handler(request);

      // Build the HTTP response
      const response = new ResponseBuilder()
        .setStatusCode(responseData.statusCode)
        .setHeaders(responseData.headers)
        .setBody(responseData.body)
        .build();

      // Send the response and end the connection
      this.socket.write(response);
      this.socket.end();
    } catch (error) {
      console.error("Error handling request:", error);

      // Send a 500 Internal Server Error response
      const errorResponse = new ResponseBuilder()
        .setStatusCode(500)
        .setHeaders({ "Content-Type": "text/plain" })
        .setBody("500 Internal Server Error")
        .build();

      this.socket.write(errorResponse);
      this.socket.end();
    }
  }
}
