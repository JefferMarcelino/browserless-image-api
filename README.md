# Browserless Image API

This is a blazing-fast image search microservice using Yandex Image Search + [Browserless.io](https://docs.browserless.io) (headless Chrome as a service).

It fetches high-quality, direct image URLs for a given search query using a headless browser, parses the response, and returns JSON with usable image links.

## 🛠 How to Run
```bash
# 1. Clone the repo
git clone https://github.com/your-username/browserless-image-api.git
cd browserless-image-api

# 2. Create your .env file
cp .env.example .env
# Edit it to set your Browserless token and url

# 3. Install dependencies
go mod tidy

# 4. Run the server
go run ./cmd/server/main.go
```

## 🔐 Environment Variables
You'll need a .env file in the root directory. Here's what's required:

```bash
PORT=3000
BROWSERLESS_URL=https://chrome.browserless.io
BROWSERLESS_TOKEN=your_browserless_token_here
```

## 🧪 Example CURL Test
```bash
curl "http://localhost:3000/image?q=iphone+15&max=2"
```

### Response:
```json
{
  "query": "iphone 15",
  "images": [
    "https://example.com/image1.jpg",
    "https://example.com/image2.jpg"
  ]
}
```
