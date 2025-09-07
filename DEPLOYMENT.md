# 🚀 BananaVerse Deployment Guide

## Quick Deploy Options

### Option 1: Railway ✅ (Successfully Deployed!)  
**Live Demo**: https://web-production-17b25.up.railway.app

**Steps taken:**
1. ✅ Pushed code to GitHub: https://github.com/PhilipJohnBasile/bananaverse
2. ✅ Connected Railway to GitHub repo
3. ✅ Added `railway.toml` for Docker build configuration
4. ✅ Fixed Railway settings:
   - Pre-deploy Command: *(empty)*
   - Custom Start Command: `./bananaverse`  
5. ✅ Set environment variable: `GOOGLE_AI_API_KEY`
6. ✅ Successfully deployed! 🚀

### Option 2: Google Cloud Run
```bash
# Build and deploy
gcloud builds submit --tag gcr.io/PROJECT_ID/bananaverse
gcloud run deploy bananaverse \
  --image gcr.io/PROJECT_ID/bananaverse \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars GOOGLE_AI_API_KEY=your-key
```

### Option 3: Heroku
```bash
# Add Heroku remote and deploy
heroku create your-app-name
heroku config:set GOOGLE_AI_API_KEY=your-key
git push heroku main
```

## Pre-Deployment Checklist

### ✅ Code Ready
- [x] All features working locally
- [x] AI generation functioning  
- [x] Image composition working
- [x] Random adventures generating
- [x] Download functionality working

### ✅ Documentation
- [x] README.md with clear setup instructions
- [x] GEMINI_INTEGRATION.md (200 words)
- [x] Demo script for video
- [x] Environment variables documented

### ✅ Required Files
- [x] go.mod and go.sum
- [x] main.go (core application)
- [x] templates/index.html
- [x] static/ directory with CSS/JS
- [x] Dockerfile (if using containerized deployment)

## Environment Variables
```bash
GOOGLE_AI_API_KEY=your_gemini_api_key
PORT=8080  # Optional, defaults to 8080
```

## Testing Deployment
1. ✅ Home page loads
2. ✅ Photo upload works
3. ✅ Figurine transformation works
4. ✅ Random adventures generate
5. ✅ Scene generation works  
6. ✅ Image composition works
7. ✅ Download functionality works

## Submission Checklist

### 1. 📹 Video Demo (2 min max)
- [ ] Record following DEMO_SCRIPT.md
- [ ] Upload to YouTube/Twitter
- [ ] Make publicly viewable
- [ ] Get shareable link

### 2. 🔗 Public Project Link  
- [x] ✅ **Deployed to Railway**: https://web-production-17b25.up.railway.app
- [x] ✅ **Tested deployed version** - all features working
- [x] ✅ **No login required** - publicly accessible
- [x] ✅ **Public URL obtained** and ready for submission

### 3. 📝 Gemini Integration Writeup
- [x] GEMINI_INTEGRATION.md (200 words max)
- [x] Describes Gemini 2.5 Flash features used
- [x] Explains why they're central to app

### 4. 🏆 Kaggle Submission
- [ ] Go to https://www.kaggle.com/competitions/banana
- [ ] Submit with all 3 components:
  - Video demo URL
  - Public project link  
  - Gemini integration text

## Tips for Success

🎯 **Innovation (40%)**: Highlight unique AI figurine creation  
⚙️ **Technical (30%)**: Emphasize seamless Gemini integration  
🌟 **Impact (20%)**: Show social media/gift potential  
🎬 **Presentation (10%)**: Professional video and clean UI

Good luck! 🍌✨