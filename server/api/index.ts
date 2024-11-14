import express, { json } from "express";
import apiRouter from "../src/routes/api.ts";
import { loadDotEnv } from "../src/util/env.ts";

// load environment variables (in development)
await loadDotEnv();

// create a new Express app
const app = express();

// use the JSON middleware to parse body content as JSON
app.use(json());

// setup the router
app.use("/api", apiRouter);

app.listen(process.env.APPLICATION_PORT, () => {
	console.log(`[info] listening on port ${process.env.APPLICATION_PORT}`);
});

export default app;
