// Export functionality for BananaVerse comics
let comicPanels = [];

function downloadCurrentPanel() {
    const composedImage = document.querySelector('.composed-image');
    const caption = document.querySelector('.caption');
    
    if (!composedImage) {
        alert('No panel to download');
        return;
    }
    
    // Create a canvas for the comic panel
    const canvas = document.createElement('canvas');
    const ctx = canvas.getContext('2d');
    
    // Set canvas size (comic panel dimensions)
    const panelWidth = 800;
    const panelHeight = 600;
    canvas.width = panelWidth;
    canvas.height = panelHeight;
    
    // Create image object
    const img = new Image();
    img.crossOrigin = 'anonymous';
    
    img.onload = function() {
        // Fill background
        ctx.fillStyle = '#ffffff';
        ctx.fillRect(0, 0, panelWidth, panelHeight);
        
        // Draw border
        ctx.strokeStyle = '#333333';
        ctx.lineWidth = 8;
        ctx.strokeRect(4, 4, panelWidth - 8, panelHeight - 8);
        
        // Calculate image dimensions to fit in panel
        const padding = 20;
        const maxWidth = panelWidth - padding * 2;
        const maxHeight = panelHeight - padding * 2 - 80; // Space for caption
        
        let imgWidth = img.width;
        let imgHeight = img.height;
        
        // Scale image to fit
        const scale = Math.min(maxWidth / imgWidth, maxHeight / imgHeight);
        imgWidth *= scale;
        imgHeight *= scale;
        
        // Center image
        const x = (panelWidth - imgWidth) / 2;
        const y = (panelHeight - imgHeight - 80) / 2 + padding;
        
        // Draw image
        ctx.drawImage(img, x, y, imgWidth, imgHeight);
        
        // Draw caption if exists
        if (caption) {
            const captionText = caption.textContent;
            drawCaption(ctx, captionText, panelWidth, panelHeight);
        }
        
        // Download the canvas
        const link = document.createElement('a');
        link.download = `bananaverse-panel-${Date.now()}.png`;
        link.href = canvas.toDataURL();
        link.click();
    };
    
    img.src = composedImage.src;
}

function drawCaption(ctx, text, panelWidth, panelHeight) {
    const captionHeight = 60;
    const captionY = panelHeight - captionHeight - 20;
    const captionPadding = 15;
    
    // Caption background
    ctx.fillStyle = 'rgba(255, 255, 255, 0.95)';
    ctx.fillRect(captionPadding, captionY, panelWidth - captionPadding * 2, captionHeight);
    
    // Caption border
    ctx.strokeStyle = '#333333';
    ctx.lineWidth = 3;
    ctx.strokeRect(captionPadding, captionY, panelWidth - captionPadding * 2, captionHeight);
    
    // Caption text
    ctx.fillStyle = '#333333';
    ctx.font = 'bold 18px Arial, sans-serif';
    ctx.textAlign = 'center';
    ctx.textBaseline = 'middle';
    
    // Word wrap the text
    const maxWidth = panelWidth - captionPadding * 4;
    const lines = wrapText(ctx, text, maxWidth);
    
    const lineHeight = 22;
    const startY = captionY + captionHeight / 2 - (lines.length - 1) * lineHeight / 2;
    
    lines.forEach((line, index) => {
        ctx.fillText(line, panelWidth / 2, startY + index * lineHeight);
    });
}

function wrapText(ctx, text, maxWidth) {
    const words = text.split(' ');
    const lines = [];
    let currentLine = words[0];

    for (let i = 1; i < words.length; i++) {
        const word = words[i];
        const width = ctx.measureText(currentLine + ' ' + word).width;
        if (width < maxWidth) {
            currentLine += ' ' + word;
        } else {
            lines.push(currentLine);
            currentLine = word;
        }
    }
    lines.push(currentLine);
    return lines;
}

function startComicCreation() {
    const currentPanel = document.querySelector('.comic-panel');
    if (!currentPanel) {
        alert('Create a panel first!');
        return;
    }
    
    // Clone the current panel and add to collection
    const panelData = {
        imageUrl: currentPanel.querySelector('.composed-image').src,
        caption: currentPanel.querySelector('.caption').textContent,
        timestamp: Date.now()
    };
    
    comicPanels.push(panelData);
    
    // Show comic builder
    const comicBuilder = document.getElementById('comic-builder');
    comicBuilder.classList.remove('hidden');
    
    updateComicBuilder();
    
    // Scroll to comic builder
    comicBuilder.scrollIntoView({ behavior: 'smooth', block: 'start' });
}

function updateComicBuilder() {
    const panelsContainer = document.getElementById('comic-panels');
    panelsContainer.innerHTML = '';
    
    comicPanels.forEach((panel, index) => {
        const panelElement = document.createElement('div');
        panelElement.className = 'comic-panel-preview';
        panelElement.innerHTML = `
            <div class="panel-header">
                <span>Panel ${index + 1}</span>
                <button onclick="removePanel(${index})" class="remove-btn">❌</button>
            </div>
            <div class="mini-comic-panel">
                <img src="${panel.imageUrl}" alt="Panel ${index + 1}">
                <div class="mini-caption">${panel.caption}</div>
            </div>
        `;
        panelsContainer.appendChild(panelElement);
    });
    
    // Add "Add Panel" button
    const addButton = document.createElement('div');
    addButton.className = 'add-panel-btn';
    addButton.innerHTML = `
        <div class="add-panel-placeholder" onclick="addNewPanel()">
            <span class="plus-icon">+</span>
            <span>Add Panel</span>
        </div>
    `;
    panelsContainer.appendChild(addButton);
}

function removePanel(index) {
    comicPanels.splice(index, 1);
    updateComicBuilder();
}

function addNewPanel() {
    // Scroll back to photo step to create another panel
    document.getElementById('photo-step').scrollIntoView({ 
        behavior: 'smooth', 
        block: 'start' 
    });
    
    // Reset the form
    resetWorkflow();
}

function resetWorkflow() {
    // Clear previous results
    document.getElementById('figurine-container').innerHTML = '';
    document.getElementById('scene-container').innerHTML = '';
    document.getElementById('composition-container').innerHTML = '<p class="instruction">Complete the previous steps to compose your scene!</p>';
    
    // Reset forms
    document.querySelector('.scene-form').reset();
    document.getElementById('photo-input').value = '';
    document.getElementById('photo-preview').innerHTML = '';
    
    // Hide camera interface
    document.getElementById('camera-container').classList.add('hidden');
    document.getElementById('upload-form').classList.add('hidden');
}

function exportComic() {
    if (comicPanels.length === 0) {
        alert('Add some panels first!');
        return;
    }
    
    // Create canvas for comic strip
    const canvas = document.createElement('canvas');
    const ctx = canvas.getContext('2d');
    
    // Calculate dimensions
    const panelWidth = 400;
    const panelHeight = 300;
    const cols = Math.min(comicPanels.length, 2); // Max 2 columns
    const rows = Math.ceil(comicPanels.length / cols);
    const gap = 20;
    
    canvas.width = cols * panelWidth + (cols + 1) * gap;
    canvas.height = rows * panelHeight + (rows + 1) * gap;
    
    // Fill background
    ctx.fillStyle = '#f0f0f0';
    ctx.fillRect(0, 0, canvas.width, canvas.height);
    
    // Draw title
    ctx.fillStyle = '#333333';
    ctx.font = 'bold 24px Arial, sans-serif';
    ctx.textAlign = 'center';
    ctx.fillText('🍌 BananaVerse Comic', canvas.width / 2, 30);
    
    let loadedImages = 0;
    const totalImages = comicPanels.length;
    
    // Draw each panel
    comicPanels.forEach((panel, index) => {
        const img = new Image();
        img.crossOrigin = 'anonymous';
        
        img.onload = function() {
            const col = index % cols;
            const row = Math.floor(index / cols);
            
            const x = gap + col * (panelWidth + gap);
            const y = gap + 50 + row * (panelHeight + gap); // 50px for title
            
            // Draw panel border
            ctx.strokeStyle = '#333333';
            ctx.lineWidth = 4;
            ctx.strokeRect(x, y, panelWidth, panelHeight);
            
            // Draw panel background
            ctx.fillStyle = '#ffffff';
            ctx.fillRect(x + 2, y + 2, panelWidth - 4, panelHeight - 4);
            
            // Draw image
            const imgPadding = 10;
            const captionSpace = 40;
            const availableHeight = panelHeight - imgPadding * 2 - captionSpace;
            const availableWidth = panelWidth - imgPadding * 2;
            
            const scale = Math.min(availableWidth / img.width, availableHeight / img.height);
            const scaledWidth = img.width * scale;
            const scaledHeight = img.height * scale;
            
            const imgX = x + (panelWidth - scaledWidth) / 2;
            const imgY = y + imgPadding;
            
            ctx.drawImage(img, imgX, imgY, scaledWidth, scaledHeight);
            
            // Draw caption
            const captionY = y + panelHeight - captionSpace + 5;
            ctx.fillStyle = 'rgba(255, 255, 255, 0.9)';
            ctx.fillRect(x + 10, captionY, panelWidth - 20, 30);
            
            ctx.strokeStyle = '#333333';
            ctx.lineWidth = 2;
            ctx.strokeRect(x + 10, captionY, panelWidth - 20, 30);
            
            ctx.fillStyle = '#333333';
            ctx.font = '14px Arial, sans-serif';
            ctx.textAlign = 'center';
            ctx.textBaseline = 'middle';
            ctx.fillText(panel.caption, x + panelWidth / 2, captionY + 15);
            
            loadedImages++;
            
            // When all images are loaded, download the comic
            if (loadedImages === totalImages) {
                const link = document.createElement('a');
                link.download = `bananaverse-comic-${Date.now()}.png`;
                link.href = canvas.toDataURL();
                link.click();
            }
        };
        
        img.src = panel.imageUrl;
    });
}

// Add some CSS for the comic builder
const comicBuilderStyles = `
<style>
.comic-panel-preview {
    border: 2px solid #e2e8f0;
    border-radius: 10px;
    padding: 15px;
    background: white;
    position: relative;
}

.panel-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 10px;
    font-weight: bold;
    color: #4a5568;
}

.remove-btn {
    background: none;
    border: none;
    cursor: pointer;
    font-size: 16px;
}

.mini-comic-panel {
    position: relative;
    border: 2px solid #333;
    border-radius: 8px;
    background: white;
    padding: 5px;
}

.mini-comic-panel img {
    width: 100%;
    height: auto;
    border-radius: 4px;
}

.mini-caption {
    position: absolute;
    bottom: 8px;
    left: 8px;
    right: 8px;
    background: rgba(255, 255, 255, 0.95);
    border: 1px solid #333;
    border-radius: 12px;
    padding: 4px 8px;
    font-size: 10px;
    font-weight: bold;
    text-align: center;
}

.add-panel-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 200px;
}

.add-panel-placeholder {
    border: 2px dashed #cbd5e0;
    border-radius: 10px;
    padding: 40px;
    text-align: center;
    cursor: pointer;
    transition: all 0.3s ease;
    background: #f8f9fa;
    width: 100%;
}

.add-panel-placeholder:hover {
    border-color: #667eea;
    background: #f0f4ff;
}

.plus-icon {
    font-size: 2rem;
    display: block;
    margin-bottom: 10px;
    color: #667eea;
}
</style>`;

// Inject styles
document.head.insertAdjacentHTML('beforeend', comicBuilderStyles);