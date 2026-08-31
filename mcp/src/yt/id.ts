
/**
 * Extract YouTube video ID from various URL formats.
 * Supports:
 * - https://www.youtube.com/watch?v=VIDEO_ID
 * - https://youtu.be/VIDEO_ID
 * - https://www.youtube.com/shorts/VIDEO_ID
 * - https://www.youtube.com/embed/VIDEO_ID
 */

import { trim } from "zod";

export function exractVideoId(input: string): string | null {
    
    const trimmed = input.trim();

    // Already a bare video ID (11 chars typical)
    if (/^[a-zA-Z0-9_-]{11}$/.test(trimmed)) {
        return trimmed;
    }

    try {
        const url = new URL(trimmed);

        // youtu.be/VIDEO_ID
        if (url.hostname === "youtu.be") {
            const id = url.pathname.slice(1).split("/")[0];
            return id || null;
        }

        // youtube.com/watch?v=VIDEO_ID
        if (url.hostname.includes("youtube.com")) {
            const v = url.searchParams.get("v");
            if (v) return v;

            // /shorts/VIDEO_ID or /embed/VIDEO_ID
            const parts = url.pathname.split("/").filter(Boolean);
            if (parts[0] === "shorts" || parts[0] === "embed" || parts[0] === "v") {
                return parts[1] || null;
            }
        }
    } 
    catch {
        // not a valid url!
    }

    return null;
}
