import { HttpMethod, ParsedRequest } from "../utils/HttpUtils";

type RouteHandler = (request: ParsedRequest) => {
  statusCode: number;
  headers: { [key: string]: string };
  body: string;
};

export class Router {
  private routes: Map<string, Map<HttpMethod, RouteHandler>>;

  constructor() {
    this.routes = new Map();
  }

  public addRoute(
    method: HttpMethod,
    path: string,
    handler: RouteHandler,
  ): void {
    if (!this.routes.has(path)) {
      this.routes.set(path, new Map());
    }
    this.routes.get(path)!.set(method, handler);
  }

  public getHandler(method: HttpMethod, path: string): RouteHandler {
    const routeHandlers = this.routes.get(path);
    if (routeHandlers && routeHandlers.has(method)) {
      return routeHandlers.get(method)!;
    }

    return this.notFoundHandler;
  }

  private notFoundHandler(): {
    statusCode: number;
    headers: { [key: string]: string };
    body: string;
  } {
    return {
      statusCode: 404,
      headers: { "Content-Type": "text/plain" },
      body: "404 Not Found",
    };
  }
}
