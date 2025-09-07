// Camera functionality for BananaVerse
let currentStream = null;
let capturedImageData = null;

function initializeCamera() {
    const cameraBtn = document.getElementById('camera-btn');
    const uploadBtn = document.getElementById('upload-btn');
    const cameraContainer = document.getElementById('camera-container');
    const uploadForm = document.getElementById('upload-form');
    const video = document.getElementById('camera-feed');
    const canvas = document.getElementById('photo-canvas');
    const captureBtn = document.getElementById('capture-btn');
    const retakeBtn = document.getElementById('retake-btn');
    const usePhotoBtn = document.getElementById('use-photo-btn');

    // Toggle between camera and upload
    cameraBtn.addEventListener('click', async () => {
        cameraContainer.classList.remove('hidden');
        uploadForm.classList.add('hidden');
        await startCamera();
    });

    uploadBtn.addEventListener('click', () => {
        uploadForm.classList.remove('hidden');
        cameraContainer.classList.add('hidden');
        stopCamera();
    });

    // Camera controls
    captureBtn.addEventListener('click', capturePhoto);
    retakeBtn.addEventListener('click', retakePhoto);
    usePhotoBtn.addEventListener('click', usePhoto);
}

async function startCamera() {
    try {
        const constraints = {
            video: {
                facingMode: 'user', // Front-facing camera preferred
                width: { ideal: 640 },
                height: { ideal: 480 }
            }
        };

        currentStream = await navigator.mediaDevices.getUserMedia(constraints);
        const video = document.getElementById('camera-feed');
        video.srcObject = currentStream;
        
        // Show capture button
        document.getElementById('capture-btn').classList.remove('hidden');
        
    } catch (error) {
        console.error('Camera access denied:', error);
        alert('Camera access denied. Please use the upload option instead.');
        document.getElementById('upload-form').classList.remove('hidden');
        document.getElementById('camera-container').classList.add('hidden');
    }
}

function stopCamera() {
    if (currentStream) {
        currentStream.getTracks().forEach(track => track.stop());
        currentStream = null;
    }
}

function capturePhoto() {
    const video = document.getElementById('camera-feed');
    const canvas = document.getElementById('photo-canvas');
    const ctx = canvas.getContext('2d');
    
    // Set canvas dimensions to match video
    canvas.width = video.videoWidth;
    canvas.height = video.videoHeight;
    
    // Draw the video frame to canvas
    ctx.drawImage(video, 0, 0);
    
    // Show canvas, hide video
    video.classList.add('hidden');
    canvas.classList.remove('hidden');
    
    // Update button states
    document.getElementById('capture-btn').classList.add('hidden');
    document.getElementById('retake-btn').classList.remove('hidden');
    document.getElementById('use-photo-btn').classList.remove('hidden');
    
    // Store the image data
    capturedImageData = canvas.toDataURL('image/jpeg', 0.8);
}

function retakePhoto() {
    const video = document.getElementById('camera-feed');
    const canvas = document.getElementById('photo-canvas');
    
    // Show video, hide canvas
    video.classList.remove('hidden');
    canvas.classList.add('hidden');
    
    // Update button states
    document.getElementById('capture-btn').classList.remove('hidden');
    document.getElementById('retake-btn').classList.add('hidden');
    document.getElementById('use-photo-btn').classList.add('hidden');
    
    // Clear captured data
    capturedImageData = null;
}

function usePhoto() {
    if (!capturedImageData) {
        alert('No photo captured');
        return;
    }
    
    // Convert dataURL to blob
    fetch(capturedImageData)
        .then(res => res.blob())
        .then(blob => {
            // Create form data
            const formData = new FormData();
            formData.append('photo', blob, 'camera-capture.jpg');
            
            // Show loading
            document.getElementById('loading-overlay').classList.remove('hidden');
            
            // Send to server
            fetch('/hx/figurine', {
                method: 'POST',
                body: formData
            })
            .then(response => response.text())
            .then(html => {
                document.getElementById('figurine-container').innerHTML = html;
                stopCamera();
                
                // Hide camera interface
                document.getElementById('camera-container').classList.add('hidden');
                
                // Auto-scroll to next step
                setTimeout(() => {
                    document.getElementById('scene-step').scrollIntoView({ 
                        behavior: 'smooth', 
                        block: 'start' 
                    });
                }, 500);
            })
            .catch(error => {
                console.error('Upload failed:', error);
                alert('Upload failed. Please try again.');
            })
            .finally(() => {
                document.getElementById('loading-overlay').classList.add('hidden');
            });
        });
}

// Helper function to check if camera is supported
function isCameraSupported() {
    return !!(navigator.mediaDevices && navigator.mediaDevices.getUserMedia);
}

// Initialize camera support check on load
document.addEventListener('DOMContentLoaded', () => {
    if (!isCameraSupported()) {
        const cameraBtn = document.getElementById('camera-btn');
        cameraBtn.disabled = true;
        cameraBtn.textContent = '📱 Camera Not Supported';
        cameraBtn.style.opacity = '0.5';
        
        // Show upload form by default
        document.getElementById('upload-form').classList.remove('hidden');
    }
});

// Clean up camera when page unloads
window.addEventListener('beforeunload', stopCamera);