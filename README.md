# ProxyInBrowser

Here is Proxy in browser Project, a project that aims to provide a simple way in impersonate any web request with fetch API
and create a simple http proxy service at local.

You can make request via http proxy and that request will send to the browser which executed payloads and let the browser request on behalf of you.

## How to

The "ProxyInBrowser" project will established an HTTP proxy through a browser to execute web requests. It leverages the fetch API to allow a victim's browser to make customized requests as per the attacker's parameters, enabling the attacker to receive responses from the victim's browser. 

A typical use case for this project is in XSS (Cross-Site Scripting) attacks where, after injecting a generated malicious script, the payload from this project is automatically loaded, and JavaScript is executed to establish a WebSocket connection back to the attacker. The WebSocket is for command and control communication, which can bypass some CSP but will not automatically rebuilt unless the XSS trigger is reactivated. Also it will persist by using a specific client trace id inside localstorage, which allows controller backend knows which client is.

The main security measure against such exploits is a well-configured Content Security Policy (CSP) that can prevent XSS and block tools like ProxyInBrowser. 

The primary technical challenge involves stripping browser HTTPS requests by using a methodology similar to Burp Suite to create an HTTP proxy and performing MITM attacks with self-signed CA certificates. This setup allows manipulation of Fetch API calls and CORS responses to bypass security measures in browsers, considering ongoing updates and security enhancements.

## Installation

### Pre-requisites

```bash
rlwrap # for better readline support
pbcopy # copy payload
```

### start

```bash
make
```

## Usage

```bash
./server 
```

paste payload on website or developer kit console.

### Example

```bash
Console> help
```

### Usage Demo

[![ProxyInBrowser Usage Demo](https://markdown-videos-api.jorgenkh.no/url?url=https%3A%2F%2Fyoutu.be%2FoJyczopfzrc)](https://youtu.be/oJyczopfzrc)

## Known Issue

no-cors mode fetch command will let chrome broswer ban javascript get response from some where. It will happen when cross site CDN js/image resource is included in that website.

So Fetch can't impersonate any request that browser does. :( 

## Sponsor

[Patreon](https://patreon.com/Skyworshiper?utm_medium=unknown&utm_source=join_link&utm_campaign=creatorshare_creator&utm_content=copyLink)