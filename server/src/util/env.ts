/**
 * Checks if `NODE_ENV` is set to `production`
 * 
 * @returns {boolean}
 */
export function isProduction(): boolean {
  return process.env.NODE_ENV === "production";
}

/**
 * Loads environment variables from the `.env` file in development.
 */
export async function loadDotEnv() {
	// load environment variables from `.env` only in development
	if (!isProduction()) {
		// conditionally import dotenv
		// `await` is used to make sure the environment variables load before the
		// application starts
		await import("dotenv").then(({ config }) => config({
      path: ""
    }));
	}
}
