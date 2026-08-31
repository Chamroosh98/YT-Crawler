
import { summerizeYTVideo } from "./tools/summarize";

const url = process.argv[2];
if (!url) {
  console.error("Usage: bun run src/test-summarize.ts <youtube-url>");
  process.exit(1);
}

const result = await summerizeYTVideo({ url });
console.log(JSON.stringify(result, null, 2));