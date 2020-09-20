export class HTTPError extends Error {
  constructor(status, statusMessage, text, url) {
    const err = `HTTP ${status} ${statusMessage}: ${text} fetching ${url}`;
    super(err);
    this.status = status;
    this.statusMessage = statusMessage;
    this.text = text;
    this.url = url;
    this.err = err;
  }

  get message() {
    return this.err;
  }
}

export const fetchJson = async (url, args) => {
  const resp = await fetch(url, args);
  if (resp.status !== 200) {
    const text = await resp.text();
    throw new HTTPError(resp.status, resp.statusText, text, url);
  }
  return resp.json();
};
