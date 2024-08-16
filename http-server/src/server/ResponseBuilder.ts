export class ResponseBuilder {
  private statusCode: number = 200;
  private headers: { [key: string]: string } = {};
  private body: string = "";

  // Set the HTTP status code
  public setStatusCode(statusCode: number): ResponseBuilder {
    this.statusCode = statusCode;
    return this;
  }

  // Set response headers
  public setHeaders(headers: { [key: string]: string }): ResponseBuilder {
    this.headers = { ...this.headers, ...headers };
    return this;
  }

  // Set response body
  public setBody(body: string): ResponseBuilder {
    this.body = body;
    return this;
  }

  // Build the complete HTTP response
  public build(): string {
    const statusText = this.getStatusText(this.statusCode);
    let response = `HTTP/1.1 ${this.statusCode} ${statusText}\r\n`;

    // Add Content-Length header
    this.headers["Content-Length"] = this.body.length.toString();

    // Add headers to the response
    for (const [key, value] of Object.entries(this.headers)) {
      response += `${key}: ${value}\r\n`;
    }

    // Add an empty line to separate headers from body
    response += "\r\n";

    // Add the body
    response += this.body;

    return response;
  }

  // Helper method to get status text for a given status code
  private getStatusText(statusCode: number): string {
    const statusTexts: { [key: number]: string } = {
      200: "OK",
      404: "Not Found",
      500: "Internal Server Error",
    };
    return statusTexts[statusCode] || "Unknown Status";
  }
}
