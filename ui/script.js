document.addEventListener("DOMContentLoaded", () => {
    // Chat functionality
    const chatHistory = document.getElementById("chat-history");
    const chatForm = document.getElementById("chat-form");

    // Add message to chat history
    function addMessage(content, isUser = false) {
        const messageDiv = document.createElement("div");
        messageDiv.className = isUser ? "message user-message" : "message bot-message";
        
        const timestamp = new Date().toLocaleTimeString();
        messageDiv.innerHTML = `
            <div class="message-header">
                <strong>${isUser ? "👤 You" : "🤖 AI"}</strong>
                <span class="timestamp">${timestamp}</span>
            </div>
            <div class="message-content">${formatContent(content)}</div>
        `;
        
        chatHistory.appendChild(messageDiv);
        chatHistory.scrollTop = chatHistory.scrollHeight;
    }

    // Format content (handle newlines and basic markdown)
    function formatContent(content) {
        if (!content) return "";
        // Convert newlines to breaks
        let formatted = content.replace(/\n/g, "<br>");
        // Basic bold formatting
        formatted = formatted.replace(/\*\*(.*?)\*\*/g, "<strong>$1</strong>");
        // Basic bullet points
        formatted = formatted.replace(/^• /gm, "&bull; ");
        return formatted;
    }

    // Handle chat form submission
    chatForm.addEventListener("submit", async (event) => {
        event.preventDefault();
        const formData = new FormData(chatForm);
        const prompt = formData.get("prompt").trim();
        
        if (!prompt) return;

        // Display user message
        addMessage(prompt, true);
        
        // Clear input
        chatForm.querySelector("input[name='prompt']").value = "";

        // Show typing indicator
        const typingDiv = document.createElement("div");
        typingDiv.className = "message bot-message typing";
        typingDiv.innerHTML = "<em>AI is thinking...</em>";
        chatHistory.appendChild(typingDiv);
        chatHistory.scrollTop = chatHistory.scrollHeight;

        try {
            const response = await fetch("/chat", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ prompt }),
            });

            const result = await response.json();
            
            // Remove typing indicator
            chatHistory.removeChild(typingDiv);

            if (result.response) {
                addMessage(result.response, false);
            } else if (result.taskID) {
                // Poll for async task
                pollTaskStatus(result.taskID, chatHistory);
            } else {
                addMessage("Received response: " + JSON.stringify(result), false);
            }
        } catch (error) {
            chatHistory.removeChild(typingDiv);
            addMessage("Error: " + error.message, false);
        }
    });

    // Quick action buttons
    document.querySelectorAll(".action-btn").forEach(btn => {
        btn.addEventListener("click", () => {
            const prompt = btn.getAttribute("data-prompt");
            if (prompt) {
                chatForm.querySelector("input[name='prompt']").value = prompt;
                chatForm.dispatchEvent(new Event("submit"));
            }
        });
    });

    // System status
    async function loadSystemStatus() {
        const statusDiv = document.getElementById("system-status");
        try {
            const response = await fetch("/health");
            const state = await response.json();
            
            statusDiv.innerHTML = `
                <div class="status-item">
                    <span class="status-label">Vector Store Chunks:</span>
                    <span class="status-value">${state.vector_store_chunks || 0}</span>
                </div>
                <div class="status-item">
                    <span class="status-label">Topological Coefficients:</span>
                    <span class="status-value">${state.topological_coeffs || 0}</span>
                </div>
                <div class="status-item">
                    <span class="status-label">E8 Roots:</span>
                    <span class="status-value">${state.e8_roots || 0}</span>
                </div>
                <div class="status-item">
                    <span class="status-label">Phase Lock State:</span>
                    <span class="status-value ${state.phase_lock_state ? 'active' : ''}">${state.phase_lock_state ? "✓ Active" : "○ Inactive"}</span>
                </div>
                <div class="status-item">
                    <span class="status-label">Working Memory Size:</span>
                    <span class="status-value">${state.working_memory_size || 0}</span>
                </div>
            `;
        } catch (error) {
            statusDiv.innerHTML = `<p class="error">Unable to load status: ${error.message}</p>`;
        }
    }

    document.getElementById("refresh-status").addEventListener("click", loadSystemStatus);
    
    // Load status on page load
    loadSystemStatus();

    // Generate Content
    const generateForm = document.getElementById("generate-form");
    const generateResponse = document.getElementById("generate-response");

    generateForm.addEventListener("submit", async (event) => {
        event.preventDefault();
        const formData = new FormData(generateForm);
        const prompt = formData.get("prompt");

        generateResponse.innerHTML = '<em>Generating...</em>';

        const response = await fetch("/generate", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ prompt }),
        });

        const { taskID } = await response.json();
        pollTaskStatus(taskID, generateResponse);
    });

    // Summarize
    const summarizeForm = document.getElementById("summarize-form");
    const summarizeResponse = document.getElementById("summarize-response");

    summarizeForm.addEventListener("submit", async (event) => {
        event.preventDefault();
        const formData = new FormData(summarizeForm);
        const data = formData.get("data");

        summarizeResponse.innerHTML = '<em>Summarizing...</em>';

        const response = await fetch("/summarize", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ data }),
        });

        const result = await response.json();
        summarizeResponse.innerHTML = result.summary ? `<pre>${formatContent(result.summary)}</pre>` : '<em>No summary generated</em>';
    });

    // Multimodal
    const multimodalForm = document.getElementById("multimodal-form");
    const multimodalResponse = document.getElementById("multimodal-response");

    multimodalForm.addEventListener("submit", async (event) => {
        event.preventDefault();
        const formData = new FormData(multimodalForm);
        const prompt = formData.get("prompt");
        const imageFile = formData.get("image");

        if (!imageFile || imageFile.size === 0) {
            multimodalResponse.innerHTML = '<em>Please select an image file</em>';
            return;
        }

        multimodalResponse.innerHTML = '<em>Analyzing image...</em>';

        const reader = new FileReader();
        reader.readAsDataURL(imageFile);
        reader.onload = async () => {
            const image = reader.result.split(",")[1];

            const response = await fetch("/multimodal", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ prompt, image }),
            });

            const result = await response.json();
            multimodalResponse.innerHTML = result.response ? `<pre>${formatContent(result.response)}</pre>` : '<em>No response generated</em>';
        };
    });

    // Steganography
    const steganographyForm = document.getElementById("steganography-form");
    const steganographyResponse = document.getElementById("steganography-response");

    steganographyForm.addEventListener("submit", async (event) => {
        event.preventDefault();
        const formData = new FormData(steganographyForm);
        const prompt = formData.get("prompt");
        const imageFile = formData.get("image");

        if (!imageFile || imageFile.size === 0) {
            steganographyResponse.innerHTML = '<em>Please select an image file</em>';
            return;
        }

        steganographyResponse.innerHTML = '<em>Processing...</em>';

        const reader = new FileReader();
        reader.readAsDataURL(imageFile);
        reader.onload = async () => {
            const image = reader.result.split(",")[1];

            const response = await fetch("/steganography", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ prompt, image }),
            });

            const result = await response.json();
            steganographyResponse.innerHTML = result.response ? `<pre>${formatContent(result.response)}</pre>` : '<em>No response generated</em>';
        };
    });

    // Poll task status for async operations
    async function pollTaskStatus(taskID, responseElement) {
        const interval = setInterval(async () => {
            try {
                const response = await fetch(`/task/${taskID}`);
                const task = await response.json();

                if (task.status === "completed") {
                    clearInterval(interval);
                    if (task.result && task.result.response) {
                        const messageDiv = document.createElement("div");
                        messageDiv.className = "message bot-message";
                        messageDiv.innerHTML = `
                            <div class="message-header">
                                <strong>🤖 AI</strong>
                                <span class="timestamp">${new Date().toLocaleTimeString()}</span>
                            </div>
                            <div class="message-content">${formatContent(task.result.response)}</div>
                        `;
                        responseElement.appendChild(messageDiv);
                        responseElement.scrollTop = responseElement.scrollHeight;
                    } else {
                        responseElement.innerHTML += `<pre>${JSON.stringify(task.result, null, 2)}</pre>`;
                    }
                } else if (task.status === "failed") {
                    clearInterval(interval);
                    responseElement.innerHTML += `<pre class="error">Error: ${task.error}</pre>`;
                }
            } catch (error) {
                clearInterval(interval);
                responseElement.innerHTML += `<pre class="error">Error polling task: ${error.message}</pre>`;
            }
        }, 2000);
    }
});
