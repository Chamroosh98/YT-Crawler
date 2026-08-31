
import { createClient, type Client } from "@libsql/client";

let client: Client | null = null;

export function getDB(): Client {

    if (client) return client;

    const url = process.env.TURSO_DATABASE_URL;
    const authToken = process.env.TURSO_AUTH_TOKEN;

    if (!url) {
        throw new Error("TURSO_DATABASE_URL is required!");
    }

    client = createClient ({
        url,
        authToken,
    });

    return client;
}

export async function findVideoById(videId: string) {
    const db = getDB();
    const result = await db.execute({
        sql: `SELECT id, title, channel, published_at, url, language
                FROM videos WHERE id = ? LIMIT 1`,
        args: [videId],
    });

    if (result.rows.length === 0) return null;

    return result.rows[0];
}


