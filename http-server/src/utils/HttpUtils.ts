export type HttpMethod = "GET" | "POST" | "PUT" | "DELETE";

export interface ParsedRequest {
  method: HttpMethod;
  url: string;
  headers: { [key: string]: string };
  body: string;
}

export function parseRequest(rawRequest: string): ParsedRequest {
  const [requestLine, ...rest] = rawRequest.split("\r\n");
  const [method, url] = requestLine.split(" ");

  // Validate the method
  if (!isValidHttpMethod(method)) {
    throw new Error(`Invalid HTTP method: ${method}`);
  }

  const headers: { [key: string]: string } = {};
  let bodyStart = rest.indexOf("");

  for (let i = 0; i < bodyStart; i++) {
    const [key, value] = rest[i].split(": ");
    headers[key.toLowerCase()] = value;
  }

  const body = rest.slice(bodyStart + 1).join("\r\n");

  return { method: method as HttpMethod, url, headers, body };
}

function isValidHttpMethod(method: string): method is HttpMethod {
  return ["GET", "POST", "PUT", "DELETE"].includes(method);
}
