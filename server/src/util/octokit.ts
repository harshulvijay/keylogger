import { loadDotEnv } from "./env.ts";

// load environment variables (in development)
await loadDotEnv();

const personalToken = Buffer.from(
	process.env.GITHUB_PERSONAL_ACCESS_TOKEN as string,
	"base64"
).toString("utf8");

export default personalToken;
