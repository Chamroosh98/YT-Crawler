
import { fetchTranscript } from "youtube-transcript-plus";

export type TranscriptResult = {
    text: string;
    language?: string;
};

/**
 * Fetch transcript for a YouTube video.
 */
export async function getTranscript(videoId: string): Promise<TranscriptResult> {
    try {
        const items = await fetchTranscript(videoId, {
        // optional: custom user agent helps in some environments
        userAgent:
            "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36",
        });

        if (!items || items.length === 0) {
        throw new Error("No transcript available for this video");
        }

        const text = items
        .map((item) => item.text)
        .join(" ")
        .replace(/\s+/g, " ")
        .trim();

        if (!text) {
        throw new Error("Transcript is empty");
        }

        return { text };
    } catch (err) {
        const message = err instanceof Error ? err.message : String(err);

        // Make common errors clearer
        if (/disabled|unavailable|not available|could not find/i.test(message)) {
        throw new Error("This video has no captions/transcript available");
        }

        throw new Error(`Failed to fetch transcript: ${message}`);
    }
}