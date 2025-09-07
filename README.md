# 🍌 BananaVerse - AI Toy Adventure Creator

**Transform your selfies into epic toy adventures with Google Gemini AI!**

[![Demo Video](https://img.shields.io/badge/▶️_Demo_Video-YouTube-red)](# "Add your video link here")
[![Live Demo](https://img.shields.io/badge/🚀_Live_Demo-Try_Now-brightgreen)](# "Add your deployed URL here")

## 🎯 What is BananaVerse?

BananaVerse turns your photos into personalized toy figurines and places them in AI-generated adventure scenes. Perfect for social media content, personalized gifts, or just having fun with AI creativity!

### ✨ Key Features

- 📸 **Smart Photo Capture** - Use camera or upload photos
- 🎭 **AI Figurine Transformation** - Turn selfies into collectible toy versions  
- 🎲 **Random Adventure Generator** - AI creates unique scenes every time
- 🎨 **Intelligent Composition** - Seamlessly merge figurine into backgrounds
- 💾 **Instant Download** - Get your adventure image immediately

## 🚀 Quick Start

### Prerequisites
- Go 1.19 or higher
- Google AI API Key ([Get one here](https://aistudio.google.com/app/apikey))

### Setup
```bash
# Clone the repository
git clone https://github.com/yourusername/bananaverse.git
cd bananaverse

# Install dependencies
go mod download

# Set your Google AI API key
export GOOGLE_AI_API_KEY="your-api-key-here"

# Run the application
go run main.go

# Open in browser
open http://localhost:8080
```

## 🤖 Gemini 2.5 Flash Integration

BananaVerse leverages **Gemini 2.5 Flash Image Preview** as the core engine for all AI functionality:

### Image Generation Features Used:
1. **Figurine Transformation** - Converts selfies into toy-style figurines with chibi proportions and glossy textures
2. **Scene Generation** - Creates cinematic adventure backgrounds from AI-generated themes  
3. **Image Composition** - Intelligently places figurines into scenes with proper scaling, lighting, and shadows
4. **Dynamic Content Creation** - Generates unique adventure scenarios on-demand

The application demonstrates Gemini's multimodal capabilities by seamlessly combining text prompts with image processing to create cohesive, engaging visual content. Each step relies on Gemini's understanding of visual composition, artistic styles, and narrative context to produce high-quality results.

## 🎬 How It Works

### 1. Photo to Figurine
Upload a selfie → **Gemini 2.5 Flash** analyzes facial features and creates a toy figurine version with:
- Chibi-style proportions
- Glossy plastic texture
- Preserved distinctive features

### 2. Adventure Generation  
Click any adventure button → **Gemini 2.5 Flash** generates:
- Themed backgrounds (cyberpunk alley, mystical forest, etc.)
- Appropriate lighting conditions
- Cinematic composition

### 3. Smart Composition
Hit "Merge" → **Gemini 2.5 Flash** combines:
- Your figurine with the background scene
- Realistic shadows and lighting
- Proper scaling and positioning

### 4. Download & Share
Get your personalized adventure image instantly!

## 🏗️ Technical Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Frontend      │    │   Go Backend     │    │  Gemini AI      │
│   HTML/JS       │◄──►│   HTTP Server    │◄──►│  2.5 Flash      │
│   HTMX          │    │   Image Storage  │    │  Image Preview  │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

**Backend**: Go with Google AI SDK  
**Frontend**: HTML + HTMX + Vanilla JavaScript  
**AI Engine**: Google Gemini 2.5 Flash Image Preview  
**Storage**: Local filesystem (production-ready for cloud storage)

## 🎮 Usage Examples

Perfect for:
- 🎂 **Personalized gifts** - Create custom toy adventures for friends
- 📱 **Social media content** - Unique, shareable images  
- 🎨 **Creative projects** - Story illustrations and character art
- 🎭 **Entertainment** - Fun family activity with AI

## 🔧 Development

### Project Structure
```
bananaverse/
├── main.go                 # Core application server
├── templates/index.html    # Main UI template  
├── static/
│   ├── css/style.css      # Styling
│   ├── js/camera.js       # Camera functionality
│   └── uploads/           # Generated images
└── README.md              # This file
```

### Key API Endpoints
- `GET /` - Main application interface
- `POST /hx/figurine` - Transform photo to figurine
- `GET /hx/random-adventures` - Generate 4 random adventures
- `POST /hx/scene` - Generate background scene
- `POST /hx/compose` - Merge figurine with background

## 🚀 Deployment

Ready for deployment on:
- **Google Cloud Run** (recommended)  
- **Railway**
- **Heroku**
- **Any Docker-compatible platform**

### Environment Variables
```bash
GOOGLE_AI_API_KEY=your_gemini_api_key
PORT=8080
```

## 🏆 Hackathon Highlights

**Innovation**: Novel application of AI image generation for personalized toy creation  
**Technical Execution**: Seamless integration of multiple Gemini 2.5 Flash features  
**Impact**: Makes AI art accessible and fun for everyone  
**Presentation**: Clean, intuitive interface with instant results

## 📄 License

MIT License - Build amazing things! 🚀

## 🔗 Links

- **Demo Video**: [Add your video link]
- **Live Demo**: [Add your deployed URL]  
- **Kaggle Competition**: https://www.kaggle.com/competitions/banana