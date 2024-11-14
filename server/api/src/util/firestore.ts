import admin from "firebase-admin";
import { loadDotEnv } from "./env.ts";

// load environment variables (in development)
await loadDotEnv();

const privateKey = Buffer.from(
	process.env.FIREBASE_CREDS_PRIVATE_KEY as string,
	"base64"
).toString("utf8");

// initialize the app
admin.initializeApp({
	credential: admin.credential.cert({
    clientEmail: process.env.FIREBASE_CREDS_CLIENT_EMAIL,
    privateKey,
    projectId: process.env.FIREBASE_CREDS_PROJECT_ID,
  }),
});

const db = admin.firestore();

export default db;
