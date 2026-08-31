
import { z } from "zod";
import { exractVideoId } from "../yt/id";
import { getTranscript } from "../yt/transcript";
import { findVideoById } from "../db/turso";
import { summarizeText } from "../summarize/llm";

export const summarizeInputSchema = z.object({
    url: z.string().describe("YouTube video URL or video ID"),
});

export type SummerizeInput = z.infer<typeof summarizeInputSchema>;

export async function summerizeYTVideo(input: SummerizeInput) {
    
    const videoId = exractVideoId(input.url);

    if (!videoId) {
        return {
            ok: false as const,
            error: "Invalid YT URL or video ID!",
        };
    }

    // Try to load metadata from Turso DB
    const saved = await findVideoById(videoId);

    // Fetch transcript
    let transcriptText: string;

    try {

        const transcript = await getTranscript(videoId);
        transcriptText = transcript.text;

    } catch (error) {
        
        const message = error instanceof Error ? error.message : String(error);

        return {
            ok: false as const,
            error: message,
            videoId,
            saved: saved
            ? {
                title: String(saved.title ?? ""),
                channel: String(saved.channel ?? ""),
                url: String(saved.url ?? ""),
            }
            : null,
        };
    }

    // Summerize
    const summary = await summarizeText({
        title:      saved ? String(saved.title ?? "") : undefined,
        channel:    saved ? String(saved.channel ?? "") : undefined,
        language:   saved ? String(saved.language ?? "") : undefined,
        transcript: transcriptText
    });

    return {
        ok: true as const,
        videoId,
        source:  saved ? "turso+transcript" : "transcript-only",
        title:   saved ? String(saved.title ?? "") : undefined,
        channel: saved ? String(saved.channel ?? "") : undefined,
        url:     saved ? String(saved.url ?? `https://www.youtube.com/watch?v=${videoId}`) : `https://www.youtube.com/watch?v=${videoId}`,
        summary,
    }


}