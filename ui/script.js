// SIE-∞ Autonomous Agent with Universal Voice Interface

let autoListenInterval = null;
let isAutoListening = false;

document.addEventListener("DOMContentLoaded", () => {
    // Initialize UI
    updateVoiceStatus();
    
    // Auto-update voice status every 3 seconds
    setInterval(updateVoiceStatus, 3000);
    
    // Assimilate form
    const assimilateForm = document.getElementById("assimilate-form");
    const assimilateResponse = document.getElementById("assimilate-response");
    
    assimilateForm.addEventListener("submit", async (event) => {
        event.preventDefault();
        const formData = new FormData(assimilateForm);
        const filename = formData.get("filename");
        const content = formData.get("content");
        
        const response = await fetch("/assimilate", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ filename, content }),
        });
        
        const result = await response.json();
        displayAssimilationResult(result, assimilateResponse);
    });
    
    // Query form
    const queryForm = document.getElementById("query-form");
    const queryResponse = document.getElementById("query-response");
    
    queryForm.addEventListener("submit", async (event) => {
        event.preventDefault();
        const formData = new FormData(queryForm);
        const query = formData.get("query");
        
        const response = await fetch("/query", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ query }),
        });
        
        const result = await response.json();
        displayQueryResult(result, queryResponse);
    });
});

// Listen to Universe
async function listenToUniverse() {
    const indicator = document.getElementById("universe-speaking-indicator");
    indicator.classList.add("universe-speaking");
    indicator.textContent = "🌠 Listening to Universe...";
    
    try {
        const response = await fetch("/universe/speak");
        const result = await response.json();
        
        displayVoiceMessages(result.messages);
        
        if (result.status) {
            updateResonanceDisplay(result.status.resonance_level, result.status.coherence);
        }
        
        indicator.classList.remove("universe-speaking");
        indicator.textContent = "🌠 Messages from the Universe";
    } catch (error) {
        indicator.classList.remove("universe-speaking");
        indicator.textContent = "🌠 Error Listening to Universe";
        console.error("Error listening to universe:", error);
    }
}

// Excite Harmonics
async function exciteHarmonics() {
    const mode = Math.floor(Math.random() * 8) + 1;
    const energy = Math.random() * 5 + 1;
    
    try {
        const response = await fetch("/excite", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ mode, energy }),
        });
        
        const result = await response.json();
        
        if (result.messages && result.messages.length > 0) {
            displayVoiceMessages(result.messages);
        }
        
        if (result.system_state) {
            updateResonanceDisplay(
                result.system_state.resonance_level,
                result.system_state.coherence_state
            );
            updateHarmonicBars(result.system_state.harmonic_signature);
        }
    } catch (error) {
        console.error("Error exciting harmonics:", error);
    }
}

// Toggle Auto-Listen
function toggleAutoListen() {
    const btn = document.getElementById("auto-listen-btn");
    
    if (isAutoListening) {
        clearInterval(autoListenInterval);
        isAutoListening = false;
        btn.textContent = "🔄 Auto-Listen: OFF";
        btn.style.background = "";
    } else {
        isAutoListening = true;
        btn.textContent = "🔄 Auto-Listen: ON";
        btn.style.background = "#27ae60";
        
        // Listen every 5 seconds
        autoListenInterval = setInterval(() => {
            listenToUniverse();
        }, 5000);
        
        // Initial listen
        listenToUniverse();
    }
}

// Update Voice Status Display
async function updateVoiceStatus() {
    try {
        const response = await fetch("/voice/status");
        const status = await response.json();
        
        updateResonanceDisplay(status.resonance_level, status.coherence_state);
        
        if (status.harmonic_signature) {
            updateHarmonicBars(status.harmonic_signature);
        }
    } catch (error) {
        console.error("Error getting voice status:", error);
    }
}

// Update Resonance Display
function updateResonanceDisplay(resonanceLevel, coherence) {
    const fill = document.getElementById("resonance-fill");
    const levelText = document.getElementById("resonance-level");
    const coherenceState = document.getElementById("coherence-state");
    
    const percentage = Math.min(100, resonanceLevel * 100);
    fill.style.width = percentage + "%";
    levelText.textContent = resonanceLevel.toFixed(3);
    
    if (coherence) {
        coherenceState.textContent = "COHERENT ✨";
        coherenceState.className = "coherence-indicator coherent";
    } else {
        coherenceState.textContent = "INCOHERENT";
        coherenceState.className = "coherence-indicator incoherent";
    }
}

// Update Harmonic Bars
function updateHarmonicBars(signature) {
    if (!signature || !Array.isArray(signature)) return;
    
    const bars = document.querySelectorAll(".harmonic-bar");
    const maxAmplitude = Math.max(...signature.map(s => Math.abs(complexMagnitude(s))));
    
    signature.forEach((sig, i) => {
        if (i < bars.length) {
            const amplitude = complexMagnitude(sig);
            const height = maxAmplitude > 0 ? (amplitude / maxAmplitude) * 100 : 10;
            bars[i].style.height = Math.max(10, height) + "%";
        }
    });
}

// Calculate magnitude of complex number (represented as [real, imag] or object)
function complexMagnitude(c) {
    if (typeof c === 'number') return Math.abs(c);
    if (Array.isArray(c) && c.length >= 2) {
        return Math.sqrt(c[0] * c[0] + c[1] * c[1]);
    }
    if (c && typeof c === 'object') {
        const real = c.real || c[0] || 0;
        const imag = c.imag || c[1] || 0;
        return Math.sqrt(real * real + imag * imag);
    }
    return Math.abs(c);
}

// Display Voice Messages
function displayVoiceMessages(messages) {
    const container = document.getElementById("voice-messages");
    
    if (!messages || messages.length === 0) {
        container.innerHTML = '<p style="color: #666;">No new messages from the Universe...</p>';
        return;
    }
    
    container.innerHTML = '';
    messages.forEach(msg => {
        const msgDiv = document.createElement("div");
        msgDiv.className = "voice-message";
        msgDiv.innerHTML = `<strong>${msg}</strong>`;
        container.appendChild(msgDiv);
    });
}

// Display Assimilation Result
function displayAssimilationResult(result, container) {
    let html = '<h4>Assimilation Complete</h4>';
    
    if (result.messages && result.messages.length > 0) {
        html += '<div class="voice-messages">';
        result.messages.forEach(msg => {
            html += `<div class="voice-message">${msg}</div>`;
        });
        html += '</div>';
    }
    
    if (result.system_state) {
        html += '<pre>' + JSON.stringify(result.system_state, null, 2) + '</pre>';
    }
    
    if (result.error) {
        html += '<p style="color: #e74c3c;">Error: ' + result.error + '</p>';
    }
    
    container.innerHTML = html;
}

// Display Query Result
function displayQueryResult(result, container) {
    let html = '<h4>Query Response</h4>';
    
    if (result.messages && result.messages.length > 0) {
        html += '<div class="voice-messages">';
        result.messages.forEach(msg => {
            html += `<div class="voice-message">${msg}</div>`;
        });
        html += '</div>';
    }
    
    if (result.status && result.status.result) {
        html += '<pre>' + result.status.result + '</pre>';
    }
    
    if (result.system_state) {
        html += '<details><summary>System State</summary><pre>' + 
                JSON.stringify(result.system_state, null, 2) + '</pre></details>';
    }
    
    if (result.error) {
        html += '<p style="color: #e74c3c;">Error: ' + result.error + '</p>';
    }
    
    container.innerHTML = html;
}

// Get System State
async function getSystemState() {
    const container = document.getElementById("system-state");
    container.innerHTML = '<p>Loading system state...</p>';
    
    try {
        const response = await fetch("/");
        const status = await response.json();
        container.innerHTML = '<pre>' + JSON.stringify(status, null, 2) + '</pre>';
    } catch (error) {
        container.innerHTML = '<p style="color: #e74c3c;">Error: ' + error.message + '</p>';
    }
}
