import { debug } from "./debug";

export interface StringMap<T> {
  [key: string]: T
}

interface RequestOption {
  method: string;
  headers: StringMap<string>;
  body: string;
  mode: RequestMode;
}

export interface FetchOptions {
  url: string;
  options: RequestOption;
}


interface FetchResponse {
  status: number;
  headers: string[][];
  text: string;
  finalurl: string;
}

export interface FetchResult {
  response: FetchResponse;
  error_stack: string;
}

type AsyncExecutor<T,K> = (c: T) => Promise<K>

export let FetchExecutor: AsyncExecutor<
  FetchOptions,
  FetchResult
> = async (c: FetchOptions) => {
  var r: FetchResult = {
    error_stack: "",
    response: {
      status: 0,
      headers: [],
      text: "",
      finalurl: "",
    },
  };
  try {
    debug("pre req data: ", c);
    debug("request ", c.options.method, c.options.mode, c.url);

    let rx = {} as Response;
    if (c.options.method == "get" || c.options.method == "GET") {
      rx = await fetch(c.url, {
        body: null,
        headers: c.options.headers,
        method: c.options.method,
        mode: c.options.mode,
      });
    } else {
      rx = await fetch(c.url, c.options);
    }
    debug(rx);
    r.response.finalurl = rx.url;
    r.response.status = rx.status;
    r.response.text = await rx.text();
    rx.headers.forEach((v, k) => {
      debug(k, v);
      r.response.headers.push([k, v]);
    });
  } catch (e: any) {
    r.error_stack = e.stack;
  } finally {
    return r;
  }
};
