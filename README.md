# yanalytics

yanalytics is a cookie-less, self-hosted, open-source alternative to Google Analytics

### Self Host the yanalytics Server

Starts analytics server

```sh
PORT=80 HOST=your-yanalitics-server.com go run main.go
```

### Add yanalitics to your website

```html
...
    ...
        ...
        <div>bla bla bla</div>
        <script async src="http://your-yanalytics-server.com/y.js"></script>
    </body>
</html>
```

### How it works

yanalytics is composed of a javascript tracker file and a http endpoint where the js tracker send http requests
The js file is returned on-demand by the yanalytics server which forces the browser to cache it automatically

### Requesting the JavaScript tracker file

This happens every time the yanalytics script tag is added to any website. It calls `your-yanalytics-server.com/y.js`

```sh
1. browser GET /y.js - - - - - - - - - - - - - - - - - - - - - - - - - - - > your-yanalytics-server.com
2. browser < - - - - - - - - - - - - - - -  "Last-Modified: Thu, Jan 1 1970" your-yanalytics-server.com
?. Few days later...
3. browser GET /y.js "If-Modified-Since: Thu, Jan 1 1970"  - - - - - - - - > your-yanalytics-server.com
4. browser < - - - - - - - - - - - - - - - - - - - - - - - -  "Not Modified" your-yanalytics-server.com
5. browser gets /y.js from local cache
```

At step 2 the analytics server generates a unique visitor id which is returned with the js tracker file.
Hopefully the browser has cached the file and every time the js tracker gets executed will always have the same visitor id. Effectively the visitor id acts like a cookie

### Tracking Page Views

Once the JavaScript tracker is downloaded it get executed and send an `XMLHttpRequest` to `your-yanalytics-server.com/y` which decoded the payload and tracks the page view event.
Example of Analytics Request

```sh
➜  yanalytics git:(master) ✗ PORT=8000 HOST=http://localhost:8000 go run main.go
storing main.AnalyticsRequest{URL:"http://localhost:8080/", UserAgent:"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.0 Safari/605.1.15", Source:"", Referrer:"", WindowWidth:(*int)(0xc0001420f8), UserID:"00a23b0f-ed16-4798-9b4f-472fc351950e"}
storing main.AnalyticsRequest{URL:"http://localhost:8080/", UserAgent:"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.0 Safari/605.1.15", Source:"", Referrer:"", WindowWidth:(*int)(0xc00010e118), UserID:"00a23b0f-ed16-4798-9b4f-472fc351950e"}
```