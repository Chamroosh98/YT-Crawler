
import OpenAI from "openai";

const MAX_CHARS = 12000;     // Keep prompt reasonable!

function getClient() {
    const apiKey = process.env.OPENAI_API_KEY || process.env.GROQ_API_KEY;

    if (!apiKey) {
        throw new Error("Agent AI apiKey is required!");
    }

    const isGrok = Boolean(process.env.GROQ_API_KEY) && ! process.env.OPENAI_API_KEY;

    return new OpenAI ({
        apiKey,
        baseURL: isGrok ? "https://api.groq.com/openai/v1" : undefined,
    });
}

function getModel() {
    if (process.env.GROQ_API_KEY && !process.env.OPENAI_API_KEY) {
        return process.env.GROQ_MODEL || "llama-3.3-70b-versatile";
    }

    return process.env.OPENAI_MODEL || "gpt-4o-mini";
}

/**
    Summarize transcript text into a concise Persian/English summary.
*/

export async function summarizeText(params: {
    title?: string;
    channel?: string;
    transcript: string;
    language?: string;
}): Promise<string> {
    
    const client = getClient();
    const model  = getModel();

    let transcript = params.transcript;
    if (transcript.length > MAX_CHARS) {
        transcript = transcript.slice(0, MAX_CHARS) + "\n\n[... transcript truncated ...]";
    }

    const system = `You are a helpful assistant that summarizes YouTube videos.
Write a clear, structured summary.
Prefer the same language as the transcript when possible.
Include:
- Main topic (1-2 sentences)
- Key points (bullet list)
- Conclusion / takeaway (1-2 sentences)
Keep it concise.`;

    const user = [
        params.title ? `Title : ${params.title}` : null,
        params.channel? `Channel : ${params.channel}` : null,
        params.language ? `Language hint : ${params.language}` : null,
        "",
        "Transcript :",
        transcript, 
        ]
        .filter(Boolean)
        .join("\n");

    const response = await client.chat.completions.create({
        model,
        messages: [
            { role: "system", content: system },
            { role: "user",   content: user   },
        ],
        temperature: 0.3,
    });

    const content = response.choices[0]?.message?.content?.trim();

    if (!content) {
        throw new Error("Empty response from LLM!");
    }

    return content;
}




