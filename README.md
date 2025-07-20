# Image API

This is a blazing-fast image search microservice using [Serpapi](https://serpapi.com/).

It fetches high-quality, direct image URLs for a given search query using a Serpapi, parses the response, and returns JSON with usable image links.

## 🛠 How to Run
```bash
# 1. Clone the repo
git clone https://github.com/your-username/image-api.git
cd image-api

# 2. Create your .env file
cp .env.example .env
# Edit it to set your SerpAPI keys

# 3. Install dependencies
go mod tidy

# 4. Run the server
go run ./cmd/server/main.go
```

## 🔐 Environment Variables
You'll need a .env file in the root directory. Here's what's required:

```bash
PORT=3000
SERP_KEY1=your_serpapi_key_here
```

## 🧪 Example CURL Test
```bash
curl "http://localhost:3000/image?q=iphone+15&max=2"
```
