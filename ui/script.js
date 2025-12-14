document.addEventListener("DOMContentLoaded", () => {
    const generateForm = document.getElementById("generate-form");
    const generateResponse = document.getElementById("generate-response");

    generateForm.addEventListener("submit", async (event) => {
        event.preventDefault();
        const formData = new FormData(generateForm);
        const prompt = formData.get("prompt");

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

    const chatHistory = document.getElementById("chat-history");
    const chatForm = document.getElementById("chat-form");

    chatForm.addEventListener("submit", async (event) => {
        event.preventDefault();
        const formData = new FormData(chatForm);
        const prompt = formData.get("prompt");

        // Display user message
        const userMessage = document.createElement("div");
        userMessage.textContent = `You: ${prompt}`;
        chatHistory.appendChild(userMessage);

        const response = await fetch("/chat", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ prompt }),
        });

        const result = await response.json();

        if (result.taskID) {
            const botMessage = document.createElement("div");
            botMessage.textContent = `Bot: Task created with ID: ${result.taskID}`;
            chatHistory.appendChild(botMessage);
            pollTaskStatus(result.taskID, chatHistory, true);
        } else {
            const botMessage = document.createElement("div");
            botMessage.innerHTML = `Bot: <pre>${JSON.stringify(result, null, 2)}</pre>`;
            chatHistory.appendChild(botMessage);
        }
    });

    const multimodalForm = document.getElementById("multimodal-form");
    const multimodalResponse = document.getElementById("multimodal-response");

    multimodalForm.addEventListener("submit", async (event) => {
        event.preventDefault();
        const formData = new FormData(multimodalForm);
        const prompt = formData.get("prompt");
        const imageFile = formData.get("image");

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

            const { taskID } = await response.json();
            pollTaskStatus(taskID, multimodalResponse);
        };
    });

    const steganographyForm = document.getElementById("steganography-form");
    const steganographyResponse = document.getElementById("steganography-response");

    steganographyForm.addEventListener("submit", async (event) => {
        event.preventDefault();
        const formData = new FormData(steganographyForm);
        const prompt = formData.get("prompt");
        const imageFile = formData.get("image");

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

            const { taskID } = await response.json();
            pollTaskStatus(taskID, steganographyResponse);
        };
    });

    const summarizeForm = document.getElementById("summarize-form");
    const summarizeResponse = document.getElementById("summarize-response");

    summarizeForm.addEventListener("submit", async (event) => {
        event.preventDefault();
        const formData = new FormData(summarizeForm);
        const data = formData.get("data");

        const response = await fetch("/summarize", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ data }),
        });

        const { taskID } = await response.json();
        pollTaskStatus(taskID, summarizeResponse);
    });

    async function pollTaskStatus(taskID, responseElement, isChat = false) {
        const interval = setInterval(async () => {
            const response = await fetch(`/task/${taskID}`);
            const task = await response.json();

            if (task.status === "completed") {
                clearInterval(interval);
                if (isChat) {
                    const botMessage = document.createElement("div");
                    botMessage.innerHTML = `Bot: <pre>${JSON.stringify(task.result, null, 2)}</pre>`;
                    responseElement.appendChild(botMessage);
                } else {
                    responseElement.innerHTML = `<pre>${JSON.stringify(task.result, null, 2)}</pre>`;
                }
            } else if (task.status === "failed") {
                clearInterval(interval);
                if (isChat) {
                    const botMessage = document.createElement("div");
                    botMessage.innerHTML = `Bot: <pre>Error: ${task.error}</pre>`;
                    responseElement.appendChild(botMessage);
                } else {
                    responseElement.innerHTML = `<pre>Error: ${task.error}</pre>`;
                }
            }
        }, 2000);
    }
});
