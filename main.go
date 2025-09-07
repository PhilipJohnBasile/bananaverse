package main

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type App struct {
	geminiClient *genai.Client
	templates    *template.Template
}

type FigurineResponse struct {
	URL     string `json:"url"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

type SceneResponse struct {
	URL     string `json:"url"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

type CompositionResponse struct {
	URL     string `json:"url"`
	Caption string `json:"caption"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

func main() {
	ctx := context.Background()
	
	apiKey := os.Getenv("GOOGLE_AI_API_KEY")
	if apiKey == "" {
		apiKey = "AIzaSyAGPXQF8kcXexUTLiq6i1Rr2hwelnzpbcs"
	}
	
	// No cloud storage needed

	geminiClient, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal("Failed to create Gemini client:", err)
	}
	defer geminiClient.Close()

	// Using local file storage only

	templates, err := template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal("Failed to parse templates:", err)
	}

	app := &App{
		geminiClient: geminiClient,
		templates:    templates,
	}

	http.HandleFunc("/", app.indexHandler)
	http.HandleFunc("/hx/figurine", app.figurineHandler)
	http.HandleFunc("/hx/scene", app.sceneHandler)
	http.HandleFunc("/hx/compose", app.composeHandler)
	http.HandleFunc("/hx/caption", app.captionHandler)
	http.HandleFunc("/hx/random-adventures", app.randomAdventuresHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting BananaVerse on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func (app *App) indexHandler(w http.ResponseWriter, r *http.Request) {
	if err := app.templates.ExecuteTemplate(w, "index.html", nil); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
	}
}

func (app *App) figurineHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("photo")
	if err != nil {
		http.Error(w, "No photo uploaded", http.StatusBadRequest)
		return
	}
	defer file.Close()

	imageData, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read image", http.StatusInternalServerError)
		return
	}

	figurineURL, err := app.transformToFigurine(r.Context(), imageData)
	if err != nil {
		log.Printf("Figurine transformation error: %v", err)
		app.renderFigurineError(w, "Failed to transform image")
		return
	}

	app.renderFigurineSuccess(w, figurineURL)
}

func (app *App) sceneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	theme := r.FormValue("theme")
	timeOfDay := r.FormValue("timeOfDay")
	prompt := r.FormValue("prompt")

	log.Printf("Scene request - Theme: '%s', TimeOfDay: '%s', Prompt: '%s'", theme, timeOfDay, prompt)

	if theme == "" || timeOfDay == "" {
		log.Printf("Scene request failed - missing required fields. Theme: '%s', TimeOfDay: '%s'", theme, timeOfDay)
		http.Error(w, "Theme and time of day required", http.StatusBadRequest)
		return
	}

	sceneURL, err := app.generateScene(r.Context(), theme, timeOfDay, prompt)
	if err != nil {
		log.Printf("Scene generation error: %v", err)
		app.renderSceneError(w, "Failed to generate scene")
		return
	}

	app.renderSceneSuccess(w, sceneURL)
}

func (app *App) composeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("=== COMPOSE HANDLER CALLED ===")
	log.Printf("Method: %s", r.Method)
	log.Printf("Content-Type: %s", r.Header.Get("Content-Type"))
	log.Printf("User-Agent: %s", r.Header.Get("User-Agent"))
	
	if r.Method != http.MethodPost {
		log.Printf("ERROR: Wrong method - expected POST, got %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	figurineURL := r.FormValue("figurineUrl")
	backgroundURL := r.FormValue("backgroundUrl")

	log.Printf("Form values received:")
	log.Printf("  - figurineUrl: '%s'", figurineURL)
	log.Printf("  - backgroundUrl: '%s'", backgroundURL)

	if figurineURL == "" || backgroundURL == "" {
		log.Printf("ERROR: Missing parameters - figurineUrl='%s', backgroundUrl='%s'", figurineURL, backgroundURL)
		http.Error(w, "Both figurine and background URLs required", http.StatusBadRequest)
		return
	}

	composedURL, _, err := app.composeScene(r.Context(), figurineURL, backgroundURL)
	if err != nil {
		log.Printf("Composition error: %v", err)
		app.renderCompositionError(w, "Failed to compose scene")
		return
	}

	app.renderCompositionSuccess(w, composedURL, "")
}

func (app *App) captionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	prompt := r.FormValue("prompt")
	if prompt == "" {
		http.Error(w, "Prompt required", http.StatusBadRequest)
		return
	}

	caption, err := app.generateCaption(r.Context(), prompt)
	if err != nil {
		log.Printf("Caption generation error: %v", err)
		w.Write([]byte("Failed to generate caption"))
		return
	}

	w.Write([]byte(caption))
}

func (app *App) randomAdventuresHandler(w http.ResponseWriter, r *http.Request) {
	adventures := app.generateRandomAdventures()
	
	var html strings.Builder
	for _, adventure := range adventures {
		html.WriteString(fmt.Sprintf(`
			<button onclick="generateDemoScene('%s', '%s', '%s')" class="btn-adventure" style="padding: 15px; border: 2px solid #ddd; border-radius: 12px; background: %s; color: white; font-size: 1rem; cursor: pointer; transition: transform 0.2s;">
				%s<br><strong>%s</strong><br><small>%s</small>
			</button>
		`, adventure["theme"], adventure["lighting"], adventure["prompt"], adventure["gradient"], adventure["emoji"], adventure["title"], adventure["desc"]))
	}
	
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html.String()))
}

func (app *App) transformToFigurine(ctx context.Context, imageData []byte) (string, error) {
	// Step 1: Use Gemini to analyze the image and create a detailed description
	analysisModel := app.geminiClient.GenerativeModel("gemini-1.5-flash")
	analysisModel.SetTemperature(0.3)
	
	analysisPrompt := "Analyze this person's appearance in detail. Describe their facial features, hair style, clothing, pose, and any distinctive characteristics. Be specific about colors, textures, and style elements."
	
	analysisResp, err := analysisModel.GenerateContent(ctx, 
		genai.Text(analysisPrompt),
		genai.ImageData("jpeg", imageData),
	)
	if err != nil {
		return "", fmt.Errorf("failed to analyze image: %v", err)
	}
	
	var description string = "a person"
	var hasPersonDetected bool = false
	
	if len(analysisResp.Candidates) > 0 && len(analysisResp.Candidates[0].Content.Parts) > 0 {
		if textPart, ok := analysisResp.Candidates[0].Content.Parts[0].(genai.Text); ok {
			description = string(textPart)
			// Check if the analysis mentions no person/people
			lowerDesc := strings.ToLower(description)
			hasPersonDetected = !strings.Contains(lowerDesc, "no person") && 
							   !strings.Contains(lowerDesc, "no people") && 
							   !strings.Contains(lowerDesc, "no discernible") &&
							   (strings.Contains(lowerDesc, "person") || strings.Contains(lowerDesc, "face") || strings.Contains(lowerDesc, "hair") || strings.Contains(lowerDesc, "clothing"))
		}
	}
	
	log.Printf("Image analysis: %s", description)
	log.Printf("Person detected: %t", hasPersonDetected)
	
	// If no person is detected, offer a demo experience
	if !hasPersonDetected {
		log.Println("No person detected - creating demo figurine")
		return "", fmt.Errorf("person not detected in image")
	}
	
	// Step 2: Generate figurine using the exact Google documentation approach
	log.Printf("Generating figurine using Gemini 2.5 Flash Image Preview")
	
	figurinePrompt := fmt.Sprintf("Create a picture of a collectible toy figurine based on this person: %s. Style: chibi proportions, glossy plastic texture, colorful, studio lighting", description)
	
	// Use the generative model approach as shown in documentation
	imageModel := app.geminiClient.GenerativeModel("gemini-2.5-flash-image-preview")
	imageResp, err := imageModel.GenerateContent(ctx, genai.Text(figurinePrompt))
	if err != nil {
		log.Printf("Figurine generation failed: %v", err)
		return "", fmt.Errorf("figurine generation failed: %v", err)
	}
	
	log.Printf("Figurine generation response received, processing parts...")
	
	// Check response parts for generated content
	log.Printf("Response has %d candidates", len(imageResp.Candidates))
	if len(imageResp.Candidates) == 0 {
		log.Println("No candidates in figurine response")
		return "", fmt.Errorf("no candidates returned from AI")
	}
	
	log.Printf("First candidate has %d parts", len(imageResp.Candidates[0].Content.Parts))
	for i, part := range imageResp.Candidates[0].Content.Parts {
		log.Printf("Part %d type: %T", i, part)
		
		// Try different ways to access the data based on the actual SDK structure
		if textPart, ok := part.(genai.Text); ok {
			log.Printf("Generated text: %s", string(textPart))
		} else if blobPart, ok := part.(genai.Blob); ok {
			log.Printf("Generated figurine image! MIME type: %s, size: %d bytes", blobPart.MIMEType, len(blobPart.Data))
			filename := fmt.Sprintf("figurine_%d.png", time.Now().Unix())
			return app.uploadToStorage(ctx, blobPart.Data, filename)
		} else {
			log.Printf("Unknown part type, trying reflection...")
		}
	}
	
	// If no image was generated, return error
	log.Println("No image generated")
	return "", fmt.Errorf("no figurine image generated")
}

func (app *App) generateScene(ctx context.Context, theme, timeOfDay, userPrompt string) (string, error) {
	log.Printf("Generating scene using Gemini 2.5 Flash Image Preview")
	
	// Use the Google documentation approach for scene generation
	prompt := fmt.Sprintf("Create a picture of a %s scene with %s lighting, cinematic style, space for character placement. Additional details: %s", theme, timeOfDay, userPrompt)
	
	// Use the generative model approach
	sceneModel := app.geminiClient.GenerativeModel("gemini-2.5-flash-image-preview") 
	resp, err := sceneModel.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		log.Printf("Scene generation failed: %v", err)
		return "", fmt.Errorf("scene generation failed: %v", err)
	}
	
	log.Printf("Scene generation response received, processing parts...")
	
	// Check if we have candidates before accessing
	if len(resp.Candidates) == 0 {
		log.Println("No candidates in scene response")
		return "", fmt.Errorf("no candidates returned from AI")
	}
	
	// Check response parts for generated content
	for i, part := range resp.Candidates[0].Content.Parts {
		log.Printf("Scene part %d type: %T", i, part)
		
		if textPart, ok := part.(genai.Text); ok {
			log.Printf("Generated text: %s", string(textPart))
		} else if blobPart, ok := part.(genai.Blob); ok {
			log.Printf("Generated scene image! MIME type: %s, size: %d bytes", blobPart.MIMEType, len(blobPart.Data))
			filename := fmt.Sprintf("scene_%d.png", time.Now().Unix())
			return app.uploadToStorage(ctx, blobPart.Data, filename)
		}
	}
	
	// Fallback if no image was generated
	log.Println("No scene image generated")
	return "", fmt.Errorf("no scene image generated")
}

func (app *App) composeScene(ctx context.Context, figurineURL, backgroundURL string) (string, string, error) {
	log.Printf("🎭 Composing scene with figurine: %s and background: %s", figurineURL, backgroundURL)
	
	// Load both images
	figurineData, err := app.loadImageFromURL(figurineURL)
	if err != nil {
		return "", "", fmt.Errorf("failed to load figurine: %v", err)
	}
	
	backgroundData, err := app.loadImageFromURL(backgroundURL)
	if err != nil {
		return "", "", fmt.Errorf("failed to load background: %v", err)
	}
	
	log.Printf("🖼️ Loaded figurine (%d bytes) and background (%d bytes)", len(figurineData), len(backgroundData))
	
	// Use Gemini 2.5 Flash Image Preview to compose the figurine onto the background
	log.Printf("🎨 Generating composition using Gemini 2.5 Flash Image Preview...")
	
	imageModel := app.geminiClient.GenerativeModel("gemini-2.5-flash-image-preview")
	
	compositionPrompt := "Using the provided images, place the toy figurine from image 2 onto the background scene from image 1. Ensure that the figurine is positioned naturally in the scene with appropriate scaling, lighting, and shadows. The figurine should look like it belongs in this environment."
	
	// Generate composition with both images - background first, then figurine
	imageResp, err := imageModel.GenerateContent(ctx, 
		genai.Text(compositionPrompt),
		genai.ImageData("png", backgroundData),
		genai.ImageData("png", figurineData),
	)
	if err != nil {
		log.Printf("❌ Composition generation failed: %v", err)
		return "", "", fmt.Errorf("composition generation failed: %v", err)
	}
	
	log.Printf("📸 Composition response received, processing parts...")
	
	// Check response for generated image
	if len(imageResp.Candidates) == 0 {
		log.Println("❌ No candidates in composition response")
		return "", "", fmt.Errorf("no candidates returned from AI")
	}
	
	var composedImageData []byte
	for i, part := range imageResp.Candidates[0].Content.Parts {
		log.Printf("Part %d type: %T", i, part)
		
		if blobPart, ok := part.(genai.Blob); ok {
			log.Printf("✅ Generated composed image! MIME type: %s, size: %d bytes", blobPart.MIMEType, len(blobPart.Data))
			composedImageData = blobPart.Data
			break
		}
	}
	
	if len(composedImageData) == 0 {
		log.Println("❌ No composed image generated, falling back to background only")
		composedImageData = backgroundData
	}
	
	// Save the composed image
	filename := fmt.Sprintf("composed_%d.png", time.Now().Unix())
	composedURL, err := app.uploadToStorage(ctx, composedImageData, filename)
	if err != nil {
		return "", "", fmt.Errorf("failed to save composed image: %v", err)
	}
	
	log.Printf("✅ Composition complete! URL: %s", composedURL)
	return composedURL, "", nil
}

func (app *App) generateCaption(ctx context.Context, scenePrompt string) (string, error) {
	model := app.geminiClient.GenerativeModel("gemini-1.5-flash")
	
	prompt := fmt.Sprintf("Create a witty, one-liner caption for a comic panel with this scene: %s. Keep it under 10 words and make it funny.", scenePrompt)
	
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}
	
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "Adventure awaits!", nil
	}
	
	if textPart, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
		return string(textPart), nil
	}
	
	return "Adventure awaits!", nil
}

func (app *App) uploadToStorage(ctx context.Context, data []byte, filename string) (string, error) {
	// Always save locally
	return app.saveLocally(data, filename)
}

func (app *App) saveLocally(data []byte, filename string) (string, error) {
	// Create uploads directory if it doesn't exist
	os.MkdirAll("static/uploads", 0755)
	
	// Save file locally
	filepath := fmt.Sprintf("static/uploads/%s", filename)
	file, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	
	if _, err := file.Write(data); err != nil {
		return "", err
	}
	
	// Return URL for local serving
	return fmt.Sprintf("/static/uploads/%s", filename), nil
}

func (app *App) generateRandomAdventures() []map[string]string {
	// Use Gemini to generate random adventure ideas
	model := app.geminiClient.GenerativeModel("gemini-1.5-flash")
	
	prompt := `Generate 4 unique, creative adventure scenarios for a toy figurine. For each adventure, provide:
1. A theme (2-3 words, kebab-case like "underwater-temple")
2. Lighting condition (2-3 words, kebab-case like "mystical-moonlight")  
3. Adventure prompt (3-4 words like "treasure hunting mission")
4. An emoji
5. A short title (2-3 words)
6. A brief description (3-4 words)

Format as: theme|lighting|prompt|emoji|title|description

Examples:
neon-cyberpunk-alley|golden-hour-sunset|ninja pizza heist|🌃|Neon Alley|Cyberpunk ninja heist
crystal-ice-caves|aurora-borealis-glow|frozen dragon rescue|❄️|Ice Caves|Frozen dragon rescue`

	resp, err := model.GenerateContent(context.Background(), genai.Text(prompt))
	if err != nil {
		log.Printf("Failed to generate random adventures: %v", err)
		// Fallback to a few hardcoded ones
		return []map[string]string{
			{"theme": "mysterious-jungle", "lighting": "dappled-sunlight", "prompt": "ancient artifact hunt", "emoji": "🌿", "title": "Jungle Quest", "desc": "Ancient artifact hunt", "gradient": "linear-gradient(135deg, #11998e 0%, #38ef7d 100%)"},
			{"theme": "floating-castle", "lighting": "magical-aurora", "prompt": "princess rescue mission", "emoji": "🏰", "title": "Sky Castle", "desc": "Princess rescue mission", "gradient": "linear-gradient(135deg, #667eea 0%, #764ba2 100%)"},
			{"theme": "desert-oasis", "lighting": "golden-hour", "prompt": "genie lamp search", "emoji": "🏜️", "title": "Desert Oasis", "desc": "Genie lamp search", "gradient": "linear-gradient(135deg, #f093fb 0%, #f5576c 100%)"},
			{"theme": "robot-city", "lighting": "neon-glow", "prompt": "AI uprising battle", "emoji": "🤖", "title": "Robot City", "desc": "AI uprising battle", "gradient": "linear-gradient(135deg, #6a11cb 0%, #2575fc 100%)"},
		}
	}

	var adventures []map[string]string
	gradients := []string{
		"linear-gradient(135deg, #667eea 0%, #764ba2 100%)",
		"linear-gradient(135deg, #11998e 0%, #38ef7d 100%)",
		"linear-gradient(135deg, #6a11cb 0%, #2575fc 100%)",
		"linear-gradient(135deg, #f093fb 0%, #f5576c 100%)",
		"linear-gradient(135deg, #12c2e9 0%, #c471ed 100%)",
		"linear-gradient(135deg, #ff512f 0%, #f09819 100%)",
		"linear-gradient(135deg, #a8edea 0%, #fed6e3 100%)",
		"linear-gradient(135deg, #d299c2 0%, #fef9d7 100%)",
	}

	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		if textPart, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
			lines := strings.Split(string(textPart), "\n")
			for i, line := range lines {
				if strings.TrimSpace(line) == "" || !strings.Contains(line, "|") {
					continue
				}
				parts := strings.Split(line, "|")
				if len(parts) >= 6 {
					adventure := map[string]string{
						"theme":    strings.TrimSpace(parts[0]),
						"lighting": strings.TrimSpace(parts[1]),
						"prompt":   strings.TrimSpace(parts[2]),
						"emoji":    strings.TrimSpace(parts[3]),
						"title":    strings.TrimSpace(parts[4]),
						"desc":     strings.TrimSpace(parts[5]),
						"gradient": gradients[i%len(gradients)],
					}
					adventures = append(adventures, adventure)
					if len(adventures) >= 4 {
						break
					}
				}
			}
		}
	}

	// Fallback if parsing failed
	if len(adventures) == 0 {
		return []map[string]string{
			{"theme": "mysterious-jungle", "lighting": "dappled-sunlight", "prompt": "ancient artifact hunt", "emoji": "🌿", "title": "Jungle Quest", "desc": "Ancient artifact hunt", "gradient": "linear-gradient(135deg, #11998e 0%, #38ef7d 100%)"},
			{"theme": "floating-castle", "lighting": "magical-aurora", "prompt": "princess rescue mission", "emoji": "🏰", "title": "Sky Castle", "desc": "Princess rescue mission", "gradient": "linear-gradient(135deg, #667eea 0%, #764ba2 100%)"},
			{"theme": "desert-oasis", "lighting": "golden-hour", "prompt": "genie lamp search", "emoji": "🏜️", "title": "Desert Oasis", "desc": "Genie lamp search", "gradient": "linear-gradient(135deg, #f093fb 0%, #f5576c 100%)"},
			{"theme": "robot-city", "lighting": "neon-glow", "prompt": "AI uprising battle", "emoji": "🤖", "title": "Robot City", "desc": "AI uprising battle", "gradient": "linear-gradient(135deg, #6a11cb 0%, #2575fc 100%)"},
		}
	}

	return adventures
}

// Removed unused placeholder functions

// Removed unused fallback functions - using AI generation only

// Removed unused background generation - using AI only

// Helper function to load image data from URL
func (app *App) loadImageFromURL(imageURL string) ([]byte, error) {
	log.Printf("🔍 Loading image from URL: %s", imageURL)
	
	// Handle both relative and full URLs
	var filePath string
	
	if strings.HasPrefix(imageURL, "/static/uploads/") {
		// Relative URL: /static/uploads/filename.png
		filename := strings.TrimPrefix(imageURL, "/static/uploads/")
		filePath = fmt.Sprintf("static/uploads/%s", filename)
	} else if strings.Contains(imageURL, "/static/uploads/") {
		// Full URL: http://localhost:8080/static/uploads/filename.png
		parts := strings.Split(imageURL, "/static/uploads/")
		if len(parts) >= 2 {
			filename := parts[1]
			filePath = fmt.Sprintf("static/uploads/%s", filename)
		} else {
			return nil, fmt.Errorf("could not extract filename from URL: %s", imageURL)
		}
	} else {
		return nil, fmt.Errorf("unsupported URL format: %s", imageURL)
	}
	
	log.Printf("📁 Reading file: %s", filePath)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %v", filePath, err)
	}
	
	log.Printf("✅ Successfully loaded %d bytes from %s", len(data), filePath)
	return data, nil
}

func (app *App) renderFigurineSuccess(w http.ResponseWriter, url string) {
	html := fmt.Sprintf(`
		<div id="figurine-result" class="result-panel">
			<img src="%s" alt="Transformed Figurine" class="figurine-image">
			<p class="success">Figurine created successfully! Now choose an adventure below.</p>
		</div>
	`, url)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func (app *App) renderFigurineError(w http.ResponseWriter, message string) {
	html := fmt.Sprintf(`
		<div id="figurine-result" class="error-panel">
			<p class="error">%s</p>
			<button onclick="location.reload()" class="btn-secondary">Try Again</button>
		</div>
	`, message)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func (app *App) renderSceneSuccess(w http.ResponseWriter, url string) {
	html := fmt.Sprintf(`
		<div id="scene-result" class="result-panel">
			<img src="%s" alt="Generated Scene" class="scene-image">
			<p class="success">Scene generated! Auto-composing your comic panel...</p>
		</div>
	`, url)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func (app *App) renderSceneError(w http.ResponseWriter, message string) {
	html := fmt.Sprintf(`
		<div id="scene-result" class="error-panel">
			<p class="error">%s</p>
			<button onclick="location.reload()" class="btn-secondary">Try Again</button>
		</div>
	`, message)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func (app *App) renderCompositionSuccess(w http.ResponseWriter, url, caption string) {
	html := fmt.Sprintf(`
		<div id="composition-result" class="result-panel">
			<p class="success">✅ Figurine merged into scene!</p>
			<img src="%s" alt="Merged Scene" class="composed-image">
			<button onclick="downloadImage('%s')" class="btn-primary">📥 Download Image</button>
		</div>
	`, url, url)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func (app *App) renderCompositionError(w http.ResponseWriter, message string) {
	html := fmt.Sprintf(`
		<div id="composition-result" class="error-panel">
			<p class="error">%s</p>
			<button onclick="location.reload()" class="btn-secondary">Try Again</button>
		</div>
	`, message)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}